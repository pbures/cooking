package server

import (
	"log"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/restapi"
	"cooking.buresovi.net/src/gen-server/restapi/operations"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"cooking.buresovi.net/src/handlers"
	"github.com/go-openapi/loads"
)

func SetupServer(theApp *app.App) *restapi.Server {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewCookingAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.Port = 60987

	api.MealsGetMealsHandler = meals.GetMealsHandlerFunc(handlers.NewGetMealHandler(*theApp))
	api.MealsInsertOneHandler = meals.InsertOneHandlerFunc(handlers.NewInsertOneHandler(*theApp))

	server.ConfigureAPI()
	return server
}
