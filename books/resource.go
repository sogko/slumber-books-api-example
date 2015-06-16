package books

import (
	. "github.com/sogko/slumber-books-api-example/books/domain"
	"github.com/sogko/slumber/domain"
	"net/http"
	"github.com/sogko/slumber/middlewares/mongodb"
	"github.com/sogko/slumber/middlewares/renderer"
)

type Options struct {}

func NewResource(ctx domain.IContext, options *Options) *Resource {
	resource := &Resource{ctx, options, nil, nil, nil}
	resource.generateRoutes()
	return resource
}

// Resource implements domain.IResource and domain.IMiddleware
type Resource struct {
	ctx      domain.IContext
	options  *Options
	routes   *domain.Routes
	databaseCtx domain.IDatabase
	rendererCtx domain.IRenderer
}

func (resource *Resource) Context() domain.IContext {
	return resource.ctx
}

func (resource *Resource) Routes() *domain.Routes {
	return resource.routes
}

func (resource *Resource) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	resource.rendererCtx.Render(w, req, status, v)
}

func (resource *Resource) BookRepository(req *http.Request) *BookRepository {
	return NewBookRepository(resource.databaseCtx)
}

// book resource middleware that get resource-specific context from upstream middlewares
func (resource *Resource) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	resource.databaseCtx = mongodb.GetMongoDbCtx(resource.Context(), req)
	resource.rendererCtx = renderer.GetRendererCtx(resource.Context(), req)
	next(w, req)
}