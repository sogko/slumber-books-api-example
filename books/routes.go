package books

import (
	"github.com/sogko/slumber/domain"
)

const (
	ListBooks  = "ListBooks"
	CreateBook = "CreateBook"
	GetBook    = "GetBook"
	UpdateBook = "UpdateBook"
	DeleteBook = "DeleteBook"
)

func (resource *Resource) generateRoutes() *domain.Routes {
	resource.routes = &domain.Routes{
		domain.Route{
			Name:           ListBooks,
			Method:         "GET",
			Pattern:        "/api/books",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleListBooks_v0,
			},
			ACLHandler: resource.HandlerListBooksACL,
		},
		domain.Route{
			Name:           CreateBook,
			Method:         "POST",
			Pattern:        "/api/books",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleCreateBook_v0,
			},
			ACLHandler: resource.HandlerCreateBookACL,
		},
		domain.Route{
			Name:           GetBook,
			Method:         "GET",
			Pattern:        "/api/books/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetBook_v0,
			},
			ACLHandler: resource.HandlerGetBookACL,
		},
		domain.Route{
			Name:           UpdateBook,
			Method:         "PUT",
			Pattern:        "/api/books/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleUpdateBook_v0,
			},
			ACLHandler: resource.HandlerUpdateBookACL,
		},
		domain.Route{
			Name:           DeleteBook,
			Method:         "DELETE",
			Pattern:        "/api/books/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleDeleteBook_v0,
			},
			ACLHandler: resource.HandlerDeleteBookACL,
		},
	}
	return resource.routes
}
