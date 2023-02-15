package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cooking.buresovi.net/src/server"
	"github.com/stretchr/testify/assert"
)

func TestInsertMeal(t *testing.T) {
	assert := assert.New(t)

	server := server.SetupServer(&application)
	defer server.Shutdown()

	handler := server.GetHandler()
	rr := httptest.NewRecorder()
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

	req, _ := http.NewRequest("GET", fmt.Sprintf("/meals?date=%v", d.Format("2006-01-02")), nil)
	handler.ServeHTTP(rr, req)
	assert.Equal(1,
		strings.Count(rr.Body.String(), "testing-gulasovka"),
		"body of the response needs to contain exactly one element",
	)
	//TODO: The consumers are not written to the consmers_meals table.
	//TODO: Write test for a single simple insert.
	assert.Equal(3,
		strings.Count(rr.Body.String(), "\"userId\""),
		"response needs to contain exactly three users",
	)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/meals?date=%v&daysforward=10", d.Format("2006-01-02")), nil)
	handler.ServeHTTP(rr, req)
	assert.Equal(11,
		strings.Count(rr.Body.String(), "testing-gulasovka"),
		"body of the response needs to contain exactly one element",
	)
	assert.Equal(3,
		strings.Count(rr.Body.String(), "\"userId\""),
		"response needs to contain exactly three users",
	)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/meals?date=%v&daysforward=150&limit=10", d.Format("2006-01-02")), nil)
	handler.ServeHTTP(rr, req)
	assert.Equal(10,
		strings.Count(rr.Body.String(), "testing-gulasovka"),
		"body of the response needs to contain exactly one element",
	)
	assert.Equal(3,
		strings.Count(rr.Body.String(), "\"userId\""),
		"response needs to contain exactly three users",
	)
}

func TestInsetAndFindRestApi(t *testing.T) {
	// assert := assert.New(t)

	server := server.SetupServer(&application)
	defer server.Shutdown()

	// handler := server.GetHandler()

}
