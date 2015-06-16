package books

import (
	"errors"
	"fmt"
	"github.com/sogko/slumber/domain"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Book collection name
const BooksCollection string = "books"

func NewBookRepository(db domain.IDatabase) *BookRepository {
	return &BookRepository{db}
}

// Book repository
type BookRepository struct {
	DB domain.IDatabase
}

func (repo *BookRepository) CreateBook(book *Book) error {
	book.ID = bson.NewObjectId()
	book.CreatedDate = time.Now()
	book.LastModifiedDate = time.Now()
	return repo.DB.Insert(BooksCollection, book)
}

func (repo *BookRepository) GetBooks() Books {
	books := Books{}
	err := repo.DB.FindAll(BooksCollection, nil, &books, 50, "")
	if err != nil {
		return Books{}
	}
	return books
}

func (repo *BookRepository) CountBooks(field string, query string) int {
	q := domain.Query{}
	if query != "" {
		if field != "" {
			q[field] = domain.Query{
				"$regex":   fmt.Sprintf("^%v.*", query),
				"$options": "i",
			}
		} else {
			// if not field is specified, we do a text search on pre-defined text index
			q["$text"] = domain.Query{
				"$search": query,
			}
		}
	}

	count, err := repo.DB.Count(BooksCollection, q)
	if err != nil {
		return 0
	}
	return count
}

func (repo *BookRepository) DeleteBooks(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var objectIds []bson.ObjectId
	for _, id := range ids {
		if bson.IsObjectIdHex(id) {
			objectIds = append(objectIds, bson.ObjectIdHex(id))
		}
	}
	if len(objectIds) == 0 {
		return nil
	}
	err := repo.DB.RemoveAll(BooksCollection, domain.Query{"_id": bson.M{"$in": objectIds}})
	return err
}

// DeleteAllBooks Delete all books
func (repo *BookRepository) DeleteAllBooks() error {
	err := repo.DB.DropCollection(BooksCollection)
	return err
}

// GetBook Get book specified by the id
func (repo *BookRepository) GetBookById(id string) (*Book, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var book Book
	err := repo.DB.FindOne(BooksCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &book)
	return &book, err
}

// UpdateBookById Update book specified by the id
func (repo *BookRepository) UpdateBook(id string, inBook *ChangeBook) (*Book, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	// serialize to a sub-set of allowed Book fields to update
	update := domain.Query{
		"lastModifiedDate": time.Now(),
	}
	if inBook.Author != "" {
		update["author"] = inBook.Author
	}
	if inBook.Name != "" {
		update["name"] = inBook.Name
	}
	if inBook.Description != "" {
		update["description"] = inBook.Description
	}
	if inBook.ISBN != "" {
		update["isbn"] = inBook.ISBN
	}
	if inBook.Format != "" {
		update["format"] = inBook.Format
	}
	query := domain.Query{"_id": bson.ObjectIdHex(id)}
	change := domain.Change{
		Update:    domain.Query{"$set": update},
		ReturnNew: true,
	}
	var changedBook Book
	err := repo.DB.Update(BooksCollection, query, change, &changedBook)
	return &changedBook, err
}

// DeleteBook deletes book specified by the id
func (repo *BookRepository) DeleteBook(id string) error {

	if !bson.IsObjectIdHex(id) {
		return errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	err := repo.DB.RemoveOne(BooksCollection, domain.Query{"_id": bson.ObjectIdHex(id)})
	return err
}
