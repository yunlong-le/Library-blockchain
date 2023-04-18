package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Book struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
	Borrower    string `json:"borrower"`
	Publisher   string `json:"publisher"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	books := []Book{
		{ID: "B1", Name: "Book1", Author: "Author1", ISBN: "111-1111111111", Description: "This is book 1", Publisher: "p1", Available: true, Borrower: ""},
		{ID: "B2", Name: "Book2", Author: "Author2", ISBN: "222-2222222222", Description: "This is book 2", Publisher: "P1", Available: true, Borrower: ""},
		{ID: "B3", Name: "Book3", Author: "Author3", ISBN: "333-3333333333", Description: "This is book 3", Publisher: "p1", Available: true, Borrower: ""},
		{ID: "B4", Name: "Book4", Author: "Author4", ISBN: "444-4444444444", Description: "This is book 4", Publisher: "p2", Available: true, Borrower: ""},
		{ID: "B5", Name: "Book5", Author: "Author5", ISBN: "555-5555555555", Description: "This is book 5", Publisher: "p2", Available: true, Borrower: ""},
	}

	for _, book := range books {
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

func (s *SmartContract) CreateBook(ctx contractapi.TransactionContextInterface, id string, bookName string, author string, publisher string, isbn string, description string) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the book %s already exists", id)
	}

	// 创建图书对象
	book := &Book{
		Name:        bookName,
		Author:      author,
		Publisher:   publisher,
		ISBN:        isbn,
		Borrower:    "",
		Available:   true,
		Description: description,
	}
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookJSON)
}

// ReadBook returns the book stored in the world state with given id.
func (s *SmartContract) ReadBook(ctx contractapi.TransactionContextInterface, id string) (*Book, error) {
	bookJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bookJSON == nil {
		return nil, fmt.Errorf("the book %s does not exist", id)
	}

	var book Book
	err = json.Unmarshal(bookJSON, &book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

// UpdateBook updates an existing book in the world state with provided parameters.
func (s *SmartContract) UpdateBook(ctx contractapi.TransactionContextInterface, id string, bookName string, author string, publisher string, isbn string, description string, borrower string, available bool) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the book %s does not exist", id)
	}

	// overwriting original book with new book
	// 创建图书对象
	book := &Book{
		ID:          id,
		Name:        bookName,
		Author:      author,
		Publisher:   publisher,
		ISBN:        isbn,
		Borrower:    borrower,
		Available:   available,
		Description: description,
	}
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookJSON)
}

// DeleteBook deletes a given book from the world state.
func (s *SmartContract) DeleteBook(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the book %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// BookExists returns true when book with given ID exists in world state
func (s *SmartContract) BookExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	bookJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return bookJSON != nil, nil
}

// TransferBook updates the owner field of book with given id in world state.
func (s *SmartContract) borrowBook(ctx contractapi.TransactionContextInterface, id string, borrower string) error {
	book, err := s.ReadBook(ctx, id)
	if err != nil {
		return err
	}

	book.Borrower = borrower

	if borrower != "" {
		book.Available = false
	} else {
		book.Available = true
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookJSON)
}

func (s *SmartContract) returnBook(ctx contractapi.TransactionContextInterface, id string) error {
	book, err := s.ReadBook(ctx, id)
	if err != nil {
		return err
	}

	book.Borrower = ""
	book.Available = true

	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookJSON)
}

// GetAllBooks returns all books found in world state
func (s *SmartContract) GetAllBooks(ctx contractapi.TransactionContextInterface) ([]*Book, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all books in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var books []*Book
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var book Book
		err = json.Unmarshal(queryResponse.Value, &book)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}
