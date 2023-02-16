package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/restapi"
	"cooking.buresovi.net/src/server"
	"github.com/stretchr/testify/assert"
)

func getServerHandlerRecorder(a *app.App) (*restapi.Server, http.Handler, *httptest.ResponseRecorder) {
	server := server.SetupServer(a)
	defer server.Shutdown()

	handler := server.GetHandler()
	recorder := httptest.NewRecorder()

	return server, handler, recorder
}

func TestInsertSingleMeal(t *testing.T) {
	assert := assert.New(t)
	server, handler, recorder := getServerHandlerRecorder(&application)
	defer server.Shutdown()

	d := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)

	reqJson := fmt.Sprintf(`{
		"mealId":5,
		"mealType":"dinner",
		"mealAuthorId":1,
		"mealDate":"%v",
		"mealName":"testing-singleinsert",
		"kcalories":280,
		"consumerIds": [2,3,4]
	}`, d.Format("2006-01-02"))

	var jsonStr = []byte(reqJson)

	req, _ := http.NewRequest("PUT", "/meals", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/meals?date=%v", d.Format("2006-01-02")), nil)
	handler.ServeHTTP(recorder, req)

	assert.Equal(1,
		strings.Count(recorder.Body.String(), "testing-singleinsert"),
		"body of the response needs to contain exactly one element",
	)
	assert.Equal(4,
		strings.Count(recorder.Body.String(), "\"userId\""),
		"response needs to contain exactly three users",
	)
	assert.Equal(1,
		strings.Count(recorder.Body.String(), "pavel.bures@gmail.com"),
		"response must contain the email of the author of the meal",
	)

}

func TestInsertManyMeals(t *testing.T) {
	assert := assert.New(t)
	server, handler, rr := getServerHandlerRecorder(&application)
	defer server.Shutdown()

	d := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 100; i++ {

		reqJson := fmt.Sprintf(`{
			"mealId":5,
			"mealType":"dinner",
			"mealAuthorId":%v,
			"mealDate":"%v",
			"mealName":"testing-gulasovka-%v",
			"kcalories":280,
			"consumerIds": [2,3,4]
		}`, (i%2)+1, d.AddDate(0, 0, i).Format("2006-01-02"), i)

		var jsonStr = []byte(reqJson)

		req, _ := http.NewRequest("PUT", "/meals", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		handler.ServeHTTP(rr, req)

		//TODO: Need to setup some Author into users table.
		if !assert.Equal(http.StatusCreated, rr.Result().StatusCode, "Status code") {
			t.Fatal("status codes should be 201")
		}
	}

	var tests = []struct {
		name          string
		requestString string
		numMeals      int
		numUsers      int
	}{
		{
			"one meal and four users expected",
			"/meals?date=2023-01-01",
			1,
			4,
		},
		{
			"11 meal and four users expected",
			"/meals?date=2023-01-01&daysforward=10",
			11,
			4,
		},
		{
			"10 meals and four users expected",
			"/meals?date=2023-01-01&daysforward=150&limit=10",
			10,
			4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server, handler, rr := getServerHandlerRecorder(&application)
			defer server.Shutdown()

			req, _ := http.NewRequest("GET", tt.requestString, nil)

			handler.ServeHTTP(rr, req)
			assert.Equal(tt.numMeals,
				strings.Count(rr.Body.String(), "testing-gulasovka"),
				fmt.Sprintf("response needs to contain exactly %v element", tt.numMeals),
			)
			assert.Equal(tt.numUsers,
				strings.Count(rr.Body.String(), "\"userId\""),
				fmt.Sprintf("response needs to contain exactly %v users", tt.numUsers),
			)
		})
	}
}

func TestInsetAndFindRestApi(t *testing.T) {
	// assert := assert.New(t)

	server := server.SetupServer(&application)
	defer server.Shutdown()

	// handler := server.GetHandler()

}
