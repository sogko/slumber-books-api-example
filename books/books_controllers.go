package books

import (
	"github.com/gorilla/mux"
	"github.com/sogko/slumber-books-api-example/domain"
	"github.com/sogko/slumber-books-api-example/repositories"
	"github.com/sogko/slumber/controllers"
	serverDomain "github.com/sogko/slumber/domain"
	"net/http"
)

const CurrentBookKey = "CurrentBookKey"

//---- Book Request API v0 ----

type ListBooksResponse_v0 struct {
	Books   domain.Books `json:"books"`
	Message string       `json:"message,omitempty"`
	Success bool         `json:"success"`
}

type CreateBookRequest_v0 struct {
	Book domain.NewBook `json:"book"`
}

type CreateBookResponse_v0 struct {
	Book    domain.Book `json:"book,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

type GetBookResponse_v0 struct {
	Book    domain.Book `json:"book,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

type UpdateBookRequest_v0 struct {
	Book domain.UpdateBook `json:"book"`
}

type UpdateBookResponse_v0 struct {
	Book    domain.Book `json:"book,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

type DeleteBookResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// getBookContext Private helper function to get book object from request context
func getBookContext(req *http.Request, ctx serverDomain.IContext) (*domain.Book, error) {

	params := mux.Vars(req)
	id := params["id"]

	// try to check if current book context exists first (cache)
	obj := ctx.Get(req, CurrentBookKey)
	if obj != nil {
		return obj.(*domain.Book), nil
	}

	db := ctx.GetDbCtx(req)
	repo := repositories.BookRepository{db}

	book, err := repo.GetBookById(id)
	if err == nil {
		ctx.Set(req, CurrentBookKey, book)
	}
	return book, err
}

func HandleListBooks_v0(w http.ResponseWriter, req *http.Request, ctx serverDomain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)

	repo := repositories.BookRepository{db}
	users := repo.GetBooks()

	r.JSON(w, http.StatusOK, ListBooksResponse_v0{
		Books:   users,
		Message: "Book list retrieved",
		Success: true,
	})
}

func HandleCreateBook_v0(w http.ResponseWriter, req *http.Request, ctx serverDomain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	repo := repositories.BookRepository{db}
	user := ctx.GetCurrentUserCtx(req)

	var body CreateBookRequest_v0
	err := controllers.DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	// you can add more business logic here to parse values
	// for example: allow only certain values for Format
	var newBook = domain.Book{
		Creator:     user.ID,
		Author:      body.Book.Author,
		Name:        body.Book.Name,
		Description: body.Book.Description,
		ISBN:        body.Book.ISBN,
		Format:      body.Book.Format,
	}

	err = repo.CreateBook(&newBook)
	if err != nil {
		controllers.RenderErrorResponseHelper(w, req, r, "Failed to save book object")
		return
	}

	r.JSON(w, http.StatusCreated, CreateBookResponse_v0{
		Book:    newBook,
		Message: "Book created",
		Success: true,
	})
}

func HandleGetBook_v0(w http.ResponseWriter, req *http.Request, ctx serverDomain.IContext) {
	r := ctx.GetRendererCtx(req)

	book, err := getBookContext(req, ctx)
	if err != nil {
		controllers.RenderErrorResponseHelper(w, req, r, "Book not found")
		return
	}

	r.JSON(w, http.StatusOK, GetBookResponse_v0{
		Book:    *book,
		Message: "Book retrieved",
		Success: true,
	})

}

// HandleUpdateBook_v0 updates book object
func HandleUpdateBook_v0(w http.ResponseWriter, req *http.Request, ctx serverDomain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	var body UpdateBookRequest_v0
	err := controllers.DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	repo := repositories.BookRepository{db}
	book, err := repo.UpdateBook(id, &body.Book)
	if err != nil {
		controllers.RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, UpdateBookResponse_v0{
		Book:    *book,
		Message: "User updated",
		Success: true,
	})
}

// HandleDeleteBook_v0 deletes book object
func HandleDeleteBook_v0(w http.ResponseWriter, req *http.Request, ctx serverDomain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	repo := repositories.BookRepository{db}

	err := repo.DeleteBook(id)
	if err != nil {
		controllers.RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, DeleteBookResponse_v0{
		Message: "Book deleted",
		Success: true,
	})
}
