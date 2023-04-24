package chaincode

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"log"
	"strings"
	"time"
)

type SmartContract struct {
}

type Book struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
	Borrower    string `json:"borrower"`
	Publisher   string `json:"publisher"`
	BookKey     string `json:"bookKey"`
}

type Record struct {
	BookID      string `json:"bookID"`
	Borrower    string `json:"borrower"`
	LendingTime int64  `json:"lendingTime"`
	ReturnTime  int64  `json:"returnTime"`
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	books := []Book{
		{ID: "B1", Name: "Book1", Author: "Author1", ISBN: "111-1111111111", Description: "This is book 1", Publisher: "p1", Available: true},
		{ID: "B2", Name: "Book2", Author: "Author2", ISBN: "222-2222222222", Description: "This is book 2", Publisher: "P1", Available: true},
		{ID: "B3", Name: "Book3", Author: "Author3", ISBN: "333-3333333333", Description: "This is book 3", Publisher: "p1", Available: true},
		{ID: "B4", Name: "Book4", Author: "Author4", ISBN: "444-4444444444", Description: "This is book 4", Publisher: "p2", Available: true},
		{ID: "B5", Name: "Book5", Author: "Author5", ISBN: "555-5555555555", Description: "This is book 5", Publisher: "p2", Available: true},
	}

	for _, book := range books {
		bookKey := s.generateBookKey(&book)
		book.BookKey = bookKey
		bookJSON, err := json.Marshal(book)
		if err != nil {
			return shim.Error("Marshal failed")
		}

		err = stub.PutState(book.ID, bookJSON)
		if err != nil {
			return shim.Error("failed to put to world state. %v")
		}
	}
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "borrowBook" {
		// 借书方法
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2: book ID and borrower")
		}
		err := s.borrowBook(stub, args[0], args[1])
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	} else if function == "returnBook" {
		// 还书方法
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1: book ID")
		}
		err := s.returnBook(stub, args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	} else if function == "addBook" {
		// 添加书籍方法
		if len(args) != 5 {
			return shim.Error("Incorrect number of arguments. Expecting 5: book name, author, publisher, ISBN, description")
		}
		err := s.addBook(stub, args[0], args[1], args[2], args[3], args[4])
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	} else if function == "QueryBooksByPattern" {
		// 模糊查询方法
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1: book pattern")
		}
		books, err := s.QueryBooksByPattern(stub, args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		bookJSON, err := json.Marshal(books)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(bookJSON)
	} else {
		return shim.Error("Invalid function name.")
	}
}

// 借书方法
func (s *SmartContract) borrowBook(stub shim.ChaincodeStubInterface, bookID string, borrower string) error {
	// 检查调用合约的账户是否有权限
	err := s.checkCallerAuthorization(stub)
	if err != nil {
		return err
	}
	book, err := s.GetBook(stub, bookID)
	if err != nil {
		return fmt.Errorf("failed to get book %s: %v", bookID, err)
	}
	if book.Borrower != "" {
		return fmt.Errorf("book %s is already borrowed", bookID)
	}

	book.Borrower = borrower
	book.Available = false
	record := Record{BookID: bookID, Borrower: borrower, LendingTime: getCurrentTime(), ReturnTime: 0}
	if err := s.RecordTransaction(stub, record); err != nil {
		return err
	}
	err = s.UpdateBook(stub, book)
	if err != nil {
		return fmt.Errorf("failed to update book %s: %v", bookID, err)
	}

	return nil
}

// 还书方法
func (s *SmartContract) returnBook(stub shim.ChaincodeStubInterface, bookID string) error {
	// 检查调用合约的账户是否有权限
	err := s.checkCallerAuthorization(stub)
	if err != nil {
		return err
	}
	book, err := s.GetBook(stub, bookID)
	if err != nil {
		return fmt.Errorf("failed to get book %s: %v", bookID, err)
	}

	if book.Borrower == "" {
		return fmt.Errorf("book %s is not borrowed", bookID)
	}

	book.Borrower = ""
	book.Available = true
	record := Record{BookID: bookID, ReturnTime: getCurrentTime()}
	if err := s.RecordTransaction(stub, record); err != nil {
		return err
	}
	err = s.UpdateBook(stub, book)
	if err != nil {
		return err
	}

	return nil
}

// 根据书名、作者、出版社、ISBN等信息增加书籍
func (s *SmartContract) addBook(stub shim.ChaincodeStubInterface, bookName string, author string, publisher string, isbn string, Description string) error {
	// 检查调用合约的账户是否有权限
	err := s.checkCallerAuthorization(stub)
	if err != nil {
		return err
	}

	// 创建图书对象
	book := &Book{
		Name:      bookName,
		Author:    author,
		Publisher: publisher,
		ISBN:      isbn,
	}

	bookKey := s.generateBookKey(book)

	// 根据bookKey检查图书是否已经存在
	books, err := s.QueryBooksByPattern(stub, bookKey)

	if len(books) != 0 {
		return fmt.Errorf("the book already exists with book key: %s", bookKey)
	}

	book.BookKey = bookKey
	book.Borrower = ""
	book.Available = true
	book.Description = Description
	id := uuid.New().String()
	book.ID = id

	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}
	err = stub.PutState(book.ID, bookJSON)
	if err != nil {
		return fmt.Errorf("failed to put book to world state: %v", err)
	}

	return nil
}

func (s *SmartContract) QueryBooksByPattern(stub shim.ChaincodeStubInterface, pattern string) ([]*Book, error) {

	var results []*Book

	iterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := iterator.Close(); err != nil {
			log.Printf("Failed to close iterator: %s", err)
		}
	}()

	for iterator.HasNext() {
		bookBytes, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var book Book
		err = json.Unmarshal(bookBytes.Value, &book)
		if err != nil {
			return nil, err
		}

		if strings.Contains(book.Name, pattern) ||
			strings.Contains(book.Author, pattern) ||
			strings.Contains(book.Publisher, pattern) ||
			strings.Contains(book.ISBN, pattern) ||
			strings.Contains(book.ID, pattern) ||
			strings.Contains(book.BookKey, pattern) {

			results = append(results, &book)
		}
	}

	return results, nil
}

// 记录借还书信息
func (s *SmartContract) RecordTransaction(stub shim.ChaincodeStubInterface, record Record) error {
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}

	key := "record-" + record.BookID
	if record.ReturnTime == 0 {
		// Record lending of the book
		if err := stub.PutState(key, recordBytes); err != nil {
			return fmt.Errorf("failed to put record state: %v", err)
		}
	} else {
		// Record return of the book
		recordValue, err := stub.GetState(key)
		if err != nil {
			return fmt.Errorf("failed to get record state: %v", err)
		}
		if recordValue == nil {
			return fmt.Errorf("record not found for book ID: %s", record.BookID)
		}

		var existingRecord Record
		if err := json.Unmarshal(recordValue, &existingRecord); err != nil {
			return fmt.Errorf("failed to unmarshal existing record: %v", err)
		}

		existingRecord.ReturnTime = record.ReturnTime

		existingRecordBytes, err := json.Marshal(existingRecord)
		if err != nil {
			return fmt.Errorf("failed to marshal existing record: %v", err)
		}

		if err := stub.PutState(key, existingRecordBytes); err != nil {
			return fmt.Errorf("failed to put existing record state: %v", err)
		}
	}

	return nil
}

// 根据id获取图书
func (s *SmartContract) GetBook(stub shim.ChaincodeStubInterface, bookID string) (*Book, error) {
	bookBytes, err := stub.GetState(bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bookBytes == nil {
		return nil, fmt.Errorf("book %s does not exist", bookID)
	}

	var book Book
	err = json.Unmarshal(bookBytes, &book)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal book: %v", err)
	}
	return &book, nil
}

// 根据一个book实例更新图书
func (s *SmartContract) UpdateBook(stub shim.ChaincodeStubInterface, book *Book) error {
	existingBook, err := s.GetBook(stub, book.ID)
	if err != nil {
		return err
	}

	existingBook.Available = book.Available
	existingBook.Name = book.Name
	existingBook.Author = book.Author
	existingBook.ISBN = book.ISBN
	existingBook.Publisher = book.Publisher
	existingBook.Description = book.Description
	existingBook.BookKey = book.BookKey
	existingBook.Borrower = book.Borrower

	bookBytes, err := json.Marshal(existingBook)
	if err != nil {
		return err
	}

	err = stub.PutState(book.ID, bookBytes)
	if err != nil {
		return fmt.Errorf("failed to update book %s: %v", book.ID, err)
	}

	return nil
}

// 检查提交者用户名是否为admin
func (s *SmartContract) checkCallerAuthorization(stub shim.ChaincodeStubInterface) error {
	certificate, err := stub.GetCreator()
	if err != nil {
		return fmt.Errorf("failed to get certificate of client: %v", err)
	}

	cert, err := x509.ParseCertificate(certificate)
	if err != nil {
		return fmt.Errorf("failed to parse client certificate: %v", err)
	}

	if cert.Subject.CommonName != "admin" {
		return fmt.Errorf("caller is not authorized")
	}

	return nil
}

// 根据book实例的书名、作者、出版社、ISBN号生成一个BookKey
func (s *SmartContract) generateBookKey(book *Book) string {
	data := []byte(fmt.Sprintf("%s|%s|%s|%s", book.Name, book.Author, book.Publisher, book.ISBN))
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}
