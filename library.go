package Library_blockchain

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

// Book describes basic details of what a book object contains
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

// SmartContract provides functions for managing a library
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of books to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
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
			return err
		}

		err = ctx.GetStub().PutState(book.ID, bookJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// 借书方法
func (l *SmartContract) borrowBook(ctx contractapi.TransactionContextInterface, bookID string, borrower string) error {
	book, err := l.GetBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get book %s: %v", bookID, err)
	}

	if book.Borrower != "" {
		return fmt.Errorf("book %s is already borrowed", bookID)
	}

	book.Borrower = borrower
	book.Available = false
	err = l.UpdateBook(ctx, book)
	if err != nil {
		return err
	}

	return nil
}

// 还书方法
func (l *SmartContract) returnBook(ctx contractapi.TransactionContextInterface, bookID string) error {
	book, err := l.GetBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get book %s: %v", bookID, err)
	}

	if book.Borrower == "" {
		return fmt.Errorf("book %s is not borrowed", bookID)
	}

	book.Borrower = ""
	book.Available = true
	err = l.UpdateBook(ctx, book)
	if err != nil {
		return err
	}

	return nil
}

// 根据书名、作者、出版社、ISBN等信息增加书籍
func (s *SmartContract) addBook(ctx contractapi.TransactionContextInterface, bookName string, author string, publisher string, isbn string, Description string) error {
	// 检查调用合约的账户是否有添加图书的权限
	err := s.checkCallerAuthorization(ctx)
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

	// 检查图书是否已经存在
	books, err := s.QueryBooksByPattern(ctx, bookKey)

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
	err = ctx.GetStub().PutState(id, bookJSON)
	if err != nil {
		return fmt.Errorf("failed to put book to world state: %v", err)
	}

	return nil
}

func (s *SmartContract) QueryBooksByPattern(ctx contractapi.TransactionContextInterface, pattern string) ([]*Book, error) {
	// Create an empty slice to store the query results
	var results []*Book

	// Get an iterator over all books in the state database
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	//defer iterator.Close()

	// Loop through all books and add any that match the pattern to the results slice
	for iterator.HasNext() {
		bookBytes, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		// Deserialize the book object from the state database
		var book Book
		err = json.Unmarshal(bookBytes.Value, &book)
		if err != nil {
			return nil, err
		}

		// Check if any book field matches the pattern
		if strings.Contains(book.Name, pattern) ||
			strings.Contains(book.Author, pattern) ||
			strings.Contains(book.Publisher, pattern) ||
			strings.Contains(book.ISBN, pattern) ||
			strings.Contains(book.ID, pattern) ||
			strings.Contains(book.BookKey, pattern) {

			// Add the book to the results slice
			results = append(results, &book)
		}
	}

	return results, nil
}

func (s *SmartContract) GetBook(ctx contractapi.TransactionContextInterface, bookID string) (*Book, error) {
	bookBytes, err := ctx.GetStub().GetState(bookID)
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

func (s *SmartContract) UpdateBook(ctx contractapi.TransactionContextInterface, book *Book) error {
	existingBook, err := s.GetBook(ctx, book.ID)
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

	err = ctx.GetStub().PutState(book.ID, bookBytes)
	if err != nil {
		return fmt.Errorf("failed to update book %s: %v", book.ID, err)
	}

	return nil
}

func (s *SmartContract) checkCallerAuthorization(ctx contractapi.TransactionContextInterface) error {

	err := ctx.GetClientIdentity().AssertAttributeValue("username", "admin")
	if err != nil {
		return fmt.Errorf("failed to get MSP ID of client: %v", err)
	}

	return nil
}

func (s *SmartContract) generateBookKey(book *Book) string {
	data := []byte(fmt.Sprintf("%s|%s|%s|%s", book.Name, book.Author, book.Publisher, book.ISBN))
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
