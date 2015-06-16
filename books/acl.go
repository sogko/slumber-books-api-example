package books

import (
	"github.com/sogko/slumber-users"
	"github.com/sogko/slumber/domain"
	"net/http"
)

func (resource *Resource) HandlerListBooksACL(req *http.Request, user domain.IUser) (bool, string) {
	// allow anonymous access
	return true, ""
}

func (resource *Resource) HandlerCreateBookACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		// enforce authenticated access
		return false, ""
	}
	u := user.(*users.User)
	if u.Status != users.StatusActive {
		// must be an active user
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandlerGetBookACL(req *http.Request, user domain.IUser) (bool, string) {
	// allow anonymous access
	return true, ""
}

func (resource *Resource) HandlerUpdateBookACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		// enforce authenticated access
		return false, ""
	}
	u := user.(*users.User)
	if u.Status != users.StatusActive {
		// must be an active user
		return false, "Inactive/pending user: please confirm your user account"
	}
	book, err := resource.getBookContext(req)
	if err != nil {
		// book must be valid
		return false, "Invalid book"
	}
	if book.Creator.Hex() != u.ID.Hex() {
		// only user who created the book entity is allowed to update a book
		return false, "Only creator of this book entry can make modifications"
	}
	return true, ""
}

func (resource *Resource) HandlerDeleteBookACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		// enforce authenticated access
		return false, ""
	}
	u := user.(*users.User)
	if u.Status != users.StatusActive {
		// must be an active user
		return false, "Inactive/pending user: please confirm your user account"
	}
	book, err := resource.getBookContext(req)
	if err != nil {
		// book must be valid
		return false, ""
	}
	if book.Creator.Hex() != u.ID.Hex() {
		// only user who created the book entity is allowed to delete a book
		return false, "Only creator of this book entry can delete it"
	}
	return true, ""
}
