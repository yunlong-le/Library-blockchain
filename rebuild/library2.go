package rebuild

/*
import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
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
}

// InitLedger adds a base set of books to the ledger
func (s *SmartContract) InitLedger(stub shim.ChaincodeStubInterface) error {
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

		err = stub.PutState(book.ID, bookJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// Invoke processes a transaction
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "borrowBook" {
		return s.borrowBook(stub, args)
	} else if function == "returnBook" {
		return s.returnBook(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// 借书方法
func (l *SmartContract) borrowBook(stub shim.ChaincodeStubInterface, bookID string, borrower string) error {
	bookBytes, err := stub.GetState(bookID)
	if err != nil {
		return fmt.Errorf("failed to get book %s: %v", bookID, err)
	}
	if bookBytes == nil {
		return fmt.Errorf("book %s does not exist", bookID)
	}

	var book Book
	err = json.Unmarshal(bookBytes, &book)
	if err != nil {
		return err
	}

	if book.Borrower != "" {
		return fmt.Errorf("book %s is already borrowed", bookID)
	}

	book.Borrower = borrower
	book.Available = false
	bookBytes, err = json.Marshal(book)
	if err != nil {
		return err
	}

	err = stub.PutState(bookID, bookBytes)
	if err != nil {
		return fmt.Errorf("failed to update book %s: %v", bookID, err)
	}

	return nil
}

// 还书方法
func (l *SmartContract) returnBook(stub shim.ChaincodeStubInterface, bookID string) error {
	bookBytes, err := stub.GetState(bookID)
	if err != nil {
		return fmt.Errorf("failed to get book %s: %v", bookID, err)
	}
	if bookBytes == nil {
		return fmt.Errorf("book %s does not exist", bookID)
	}

	var book Book
	err = json.Unmarshal(bookBytes, &book)
	if err != nil {
		return err
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

func (s *SmartContract) checkCallerAuthorization(stub shim.ChaincodeStubInterface) error {
	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return fmt.Errorf("failed to get creator of transaction: %v", err)
	}

	cert, err := cid.GetX509Certificate(creatorBytes)
	if err != nil {
		return fmt.Errorf("failed to get X509 certificate of client: %v", err)
	}

	if cert.Subject.CommonName != "admin" {
		return fmt.Errorf("caller is not authorized")
	}

	return nil
}

func (s *SmartContract) generateBookKey(book *Book) string {
	data := []byte(fmt.Sprintf("%s|%s|%s|%s", book.Name, book.Author, book.Publisher, book.ISBN))
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
*/
