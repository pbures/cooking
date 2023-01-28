// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"cooking.buresovi.net/src/server/models"
	"cooking.buresovi.net/src/server/restapi/operations"
	"cooking.buresovi.net/src/server/restapi/operations/todos"
)

//go:generate swagger generate server --target ../../server --name Cooking --spec ../../swagger/swagger.yml --principal interface{}

var exampleFlags = struct {
	Example1 string `long:"example1" description:"Sample for showing how to configure cmd-line flags"`
	Example2 string `long:"example2" description:"Further info at https://github.com/jessevdk/go-flags"`
}{}

func configureFlags(api *operations.CookingAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Example Flags",
			LongDescription:  "",
			Options:          &exampleFlags,
		},
	}
}

var items = make(map[int64]*models.Item)
var lastID int64

var itemsLock = &sync.Mutex{}

func newItemID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

func addItem(item *models.Item) error {
	if item == nil {
		return errors.New(500, "item must be present")
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	newID := newItemID()
	item.ID = newID
	items[newID] = item

	return nil
}

func updateItem(id int64, item *models.Item) error {
	if item == nil {
		return errors.New(500, "item must be present")
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	_, exists := items[id]
	if !exists {
		return errors.NotFound("not found: item %d", id)
	}

	item.ID = id
	items[id] = item
	return nil
}

func deleteItem(id int64) error {
	itemsLock.Lock()
	defer itemsLock.Unlock()

	_, exists := items[id]
	if !exists {
		return errors.NotFound("not found: item %d", id)
	}

	delete(items, id)
	return nil
}

func allItems(since int64, limit int32) (result []*models.Item) {
	result = make([]*models.Item, 0)
	for id, item := range items {
		if len(result) >= int(limit) {
			return
		}
		if since == 0 || id > since {
			result = append(result, item)
		}
	}
	return
}

func configureAPI(api *operations.CookingAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	//api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.TodosAddOneHandler == nil {
		api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
			return middleware.NotImplemented("operation todos.AddOne has not yet been implemented")
		})
	}
	if api.TodosDestroyOneHandler == nil {
		api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
			return middleware.NotImplemented("operation todos.DestroyOne has not yet been implemented")
		})
	}

	api.TodosFindTodosHandler = todos.FindTodosHandlerFunc(func(params todos.FindTodosParams) middleware.Responder {
		return todos.NewFindTodosOK()
		//return middleware.NotImplemented("oOOperation todos.FindTodos has not yet been implemented")
	})

	if api.TodosUpdateOneHandler == nil {
		api.TodosUpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams) middleware.Responder {
			return middleware.NotImplemented("operation todos.UpdateOne has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
