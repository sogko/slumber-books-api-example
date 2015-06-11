package books

import (
	serverDomain "github.com/sogko/slumber/domain"
)

const (
	ListBooks  = "ListBooks"
	CreateBook = "CreateBook"
	GetBook    = "GetBook"
	UpdateBook = "UpdateBook"
	DeleteBook = "DeleteBook"
)

var Routes = serverDomain.Routes{
	serverDomain.Route{
		Name:           ListBooks,
		Method:         "GET",
		Pattern:        "/api/books",
		DefaultVersion: "0.0",
		RouteHandlers: serverDomain.RouteHandlers{
			"0.0": HandleListBooks_v0,
		},
	},
	serverDomain.Route{
		Name:           CreateBook,
		Method:         "POST",
		Pattern:        "/api/books",
		DefaultVersion: "0.0",
		RouteHandlers: serverDomain.RouteHandlers{
			"0.0": HandleCreateBook_v0,
		},
	},
	serverDomain.Route{
		Name:           GetBook,
		Method:         "GET",
		Pattern:        "/api/books/{id}",
		DefaultVersion: "0.0",
		RouteHandlers: serverDomain.RouteHandlers{
			"0.0": HandleGetBook_v0,
		},
	},
	serverDomain.Route{
		Name:           UpdateBook,
		Method:         "PUT",
		Pattern:        "/api/books/{id}",
		DefaultVersion: "0.0",
		RouteHandlers: serverDomain.RouteHandlers{
			"0.0": HandleUpdateBook_v0,
		},
	},
	serverDomain.Route{
		Name:           DeleteBook,
		Method:         "DELETE",
		Pattern:        "/api/books/{id}",
		DefaultVersion: "0.0",
		RouteHandlers: serverDomain.RouteHandlers{
			"0.0": HandleDeleteBook_v0,
		},
	},
}
