package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/loads"

	"cooking.buresovi.net/src/gen-server/restapi"
	"cooking.buresovi.net/src/gen-server/restapi/operations"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"cooking.buresovi.net/src/handlers"
	"github.com/stretchr/testify/assert"
)

func TestGetMeals(t *testing.T) {
	assert.Equal(t, "a", "a")

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewCookingAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	api.MealsGetMealsHandler = meals.GetMealsHandlerFunc(handlers.GetMealsHandler)

	server.ConfigureAPI()

	req, err := http.NewRequest("GET", "/meals?date=2023-08-01", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := server.GetHandler()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
