package books

import (
	serverDomain "github.com/sogko/slumber/domain"
	"net/http"
)

var ACL = serverDomain.ACLMap{
	ListBooks: func(user *serverDomain.User, req *http.Request, ctx serverDomain.IContext) (bool, string) {
		// allow anonymous access
		return true, ""
	},
	CreateBook: func(user *serverDomain.User, req *http.Request, ctx serverDomain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != serverDomain.StatusActive {
			// must be an active user
			return false, ""
		}
		return true, ""
	},
	GetBook: func(user *serverDomain.User, req *http.Request, ctx serverDomain.IContext) (bool, string) {
		// allow anonymous access
		return true, ""
	},
	UpdateBook: func(user *serverDomain.User, req *http.Request, ctx serverDomain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != serverDomain.StatusActive {
			// must be an active user
			return false, "Inactive/pending user: please confirm your user account"
		}
		book, err := getBookContext(req, ctx)
		if err != nil {
			// book must be valid
			return false, ""
		}
		if book.Creator.Hex() != user.ID.Hex() {
			// only user who created the book entity is allowed to update a book
			return false, "Only creator of this book entry can make modifications"
		}
		return true, ""
	},
	DeleteBook: func(user *serverDomain.User, req *http.Request, ctx serverDomain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != serverDomain.StatusActive {
			// must be an active user
			return false, "Inactive/pending user: please confirm your user account"
		}
		book, err := getBookContext(req, ctx)
		if err != nil {
			// book must be valid
			return false, ""
		}
		if book.Creator.Hex() != user.ID.Hex() {
			// only user who created the book entity is allowed to delete a book
			return false, "Only creator of this book entry can delete it"
		}
		return true, ""
	},
}
