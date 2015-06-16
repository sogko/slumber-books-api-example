package books

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/sogko/slumber-books-api-example/books/domain"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const CurrentBookKey = "CurrentBookKey"

//---- Book Request API v0 ----

type ListBooksResponse_v0 struct {
	Books   Books  `json:"books"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type CreateBookRequest_v0 struct {
	Book NewBook `json:"book"`
}

type CreateBookResponse_v0 struct {
	Book    Book   `json:"book,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type GetBookResponse_v0 struct {
	Book    Book   `json:"book,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type UpdateBookRequest_v0 struct {
	Book ChangeBook `json:"book"`
}

type UpdateBookResponse_v0 struct {
	Book    Book   `json:"book,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type DeleteBookResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type ErrorResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// getBookContext Private helper function to get book object from request context
func (resource *Resource) getBookContext(req *http.Request) (*Book, error) {
	ctx := resource.Context()

	params := mux.Vars(req)
	id := params["id"]

	// try to check if current book context exists first (cache)
	obj := ctx.Get(req, CurrentBookKey)
	if obj != nil {
		return obj.(*Book), nil
	}

	repo := resource.BookRepository(req)
	book, err := repo.GetBookById(id)
	if err == nil {
		ctx.Set(req, CurrentBookKey, book)
	}
	return book, err
}

func (resource *Resource) DecodeRequestBody(w http.ResponseWriter, req *http.Request, target interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(target)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Request body parse error: %v", err.Error()))
		return err
	}
	return nil
}

func (resource *Resource) RenderError(w http.ResponseWriter, req *http.Request, status int, message string) {
	resource.Render(w, req, status, ErrorResponse_v0{
		Message: message,
		Success: false,
	})
}

func (resource *Resource) HandleListBooks_v0(w http.ResponseWriter, req *http.Request) {

	repo := resource.BookRepository(req)
	users := repo.GetBooks()

	resource.Render(w, req, http.StatusOK, ListBooksResponse_v0{
		Books:   users,
		Message: "Book list retrieved",
		Success: true,
	})
}

func (resource *Resource) HandleCreateBook_v0(w http.ResponseWriter, req *http.Request) {
	ctx := resource.Context()
	repo := resource.BookRepository(req)

	// user should never be nil (ACL takes care of that)
	user := ctx.GetCurrentUserCtx(req)

	var body CreateBookRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	// you can add more business logic here to parse values
	// for example: allow only certain values for Format
	var newBook = Book{
		Creator:     bson.ObjectIdHex(user.GetID()),
		Author:      body.Book.Author,
		Name:        body.Book.Name,
		Description: body.Book.Description,
		ISBN:        body.Book.ISBN,
		Format:      body.Book.Format,
	}

	err = repo.CreateBook(&newBook)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "Failed to save book object")
		return
	}

	resource.Render(w, req, http.StatusCreated, CreateBookResponse_v0{
		Book:    newBook,
		Message: "Book created",
		Success: true,
	})
}

func (resource *Resource) HandleGetBook_v0(w http.ResponseWriter, req *http.Request) {

	book, err := resource.getBookContext(req)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "Book not found")
		return
	}

	resource.Render(w, req, http.StatusOK, GetBookResponse_v0{
		Book:    *book,
		Message: "Book retrieved",
		Success: true,
	})

}

// HandleUpdateBook_v0 updates book object
func (resource *Resource) HandleUpdateBook_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	var body UpdateBookRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	repo := resource.BookRepository(req)
	book, err := repo.UpdateBook(id, &body.Book)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	resource.Render(w, req, http.StatusOK, UpdateBookResponse_v0{
		Book:    *book,
		Message: "User updated",
		Success: true,
	})
}

// HandleDeleteBook_v0 deletes book object
func (resource *Resource) HandleDeleteBook_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	repo := resource.BookRepository(req)

	err := repo.DeleteBook(id)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	resource.Render(w, req, http.StatusOK, DeleteBookResponse_v0{
		Message: "Book deleted",
		Success: true,
	})
}
