package repositories

import (
	"errors"
	"fmt"
	"github.com/sogko/slumber-books-api-example/domain"
	serverDomain "github.com/sogko/slumber/domain"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Book collection name
const BooksCollection string = "books"

// Book repository
type BookRepository struct {
	DB serverDomain.IDatabase
}

func (repo *BookRepository) CreateBook(book *domain.Book) error {
	book.ID = bson.NewObjectId()
	book.CreatedDate = time.Now()
	book.LastModifiedDate = time.Now()
	return repo.DB.Insert(BooksCollection, book)
}

func (repo *BookRepository) GetBooks() domain.Books {
	books := domain.Books{}
	err := repo.DB.FindAll(BooksCollection, nil, &books, 50, "")
	if err != nil {
		return domain.Books{}
	}
	return books
}

func (repo *BookRepository) CountBooks(field string, query string) int {
	q := serverDomain.Query{}
	if query != "" {
		if field != "" {
			q[field] = serverDomain.Query{
				"$regex":   fmt.Sprintf("^%v.*", query),
				"$options": "i",
			}
		} else {
			// if not field is specified, we do a text search on pre-defined text index
			q["$text"] = serverDomain.Query{
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
	err := repo.DB.RemoveAll(BooksCollection, serverDomain.Query{"_id": bson.M{"$in": objectIds}})
	return err
}

// DeleteAllBooks Delete all books
func (repo *BookRepository) DeleteAllBooks() error {
	err := repo.DB.DropCollection(BooksCollection)
	return err
}

// GetBook Get book specified by the id
func (repo *BookRepository) GetBookById(id string) (*domain.Book, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var book domain.Book
	err := repo.DB.FindOne(BooksCollection, serverDomain.Query{"_id": bson.ObjectIdHex(id)}, &book)
	return &book, err
}

// UpdateBookById Update book specified by the id
func (repo *BookRepository) UpdateBook(id string, inBook *domain.UpdateBook) (*domain.Book, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	// serialize to a sub-set of allowed domain.Book fields to update
	update := serverDomain.Query{
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
	query := serverDomain.Query{"_id": bson.ObjectIdHex(id)}
	change := serverDomain.Change{
		Update:    serverDomain.Query{"$set": update},
		ReturnNew: true,
	}
	var changedBook domain.Book
	err := repo.DB.Update(BooksCollection, query, change, &changedBook)
	return &changedBook, err
}

// DeleteBook deletes book specified by the id
func (repo *BookRepository) DeleteBook(id string) error {

	if !bson.IsObjectIdHex(id) {
		return errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	err := repo.DB.RemoveOne(BooksCollection, serverDomain.Query{"_id": bson.ObjectIdHex(id)})
	return err
}
