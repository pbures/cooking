package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/restapi"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
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

	result := meals.GetMealsOKBody{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Nil(err)

	assert.Equal(1, len(result.Meals))
	assert.Equal(4, len(result.Users))

	assert.Equal(result.Meals[0].MealName, "testing-singleinsert")
	assert.Equal(1, int(result.Meals[0].MealAuthorID))
	assert.Equal(int64(280), result.Meals[0].Kcalories)

	for _, u := range result.Users {
		if u.UserID == result.Meals[0].MealAuthorID {
			assert.Equal("pavel.bures@gmail.com", string(u.Email))
		}
	}
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
			"date and daysforward",
			"/meals?date=2023-01-01&daysforward=1",
			11,
			4,
		},
		{
			"date, daysforward and limit",
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

			result := meals.GetMealsOKBody{}
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			assert.Nil(err)

			assert.Equal(tt.numMeals, len(result.Meals))
			assert.Equal(tt.numUsers, len(result.Users))
		})
	}
}

func TestInsetAndUpdate(t *testing.T) {
	assert := assert.New(t)

	server, handler, recorder := getServerHandlerRecorder(&application)
	defer server.Shutdown()

	d := time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC)

	getMealJson := func(authorId int, mealType string, mealName string, consumerIds []int) string {
		cis := "["
		comma := ""
		for _, ci := range consumerIds {
			cis = cis + comma + fmt.Sprint(ci)
			comma = ","
		}
		cis = cis + "]"

		return fmt.Sprintf(`{
			"mealType":"%v",
			"mealAuthorId":%v,
			"mealDate":"%v",
			"mealName":"%v",
			"kcalories":280,
			"consumerIds": %v
		}`, mealType, authorId, d.Format("2006-01-02"), mealName, cis)
	}
	/* Create the meal of type dinner */
	req, _ := http.NewRequest("PUT",
		"/meals",
		bytes.NewBuffer([]byte(getMealJson(1, "dinner", "testing-insert-and-update-1", []int{1, 2}))),
	)
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, req)

	/* Update to change the name and consumers */
	req, _ = http.NewRequest("PUT",
		"/meals",
		bytes.NewBuffer([]byte(getMealJson(2, "dinner", "testing-insert-and-update-2", []int{2, 3, 4}))),
	)
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, req)

	/* Insert the same date, but different type - this creates a new meal */
	req, _ = http.NewRequest("PUT",
		"/meals",
		bytes.NewBuffer([]byte(getMealJson(2, "breakfast", "testing-insert-and-update-3", []int{1}))),
	)
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/meals?date=%v", d.Format("2006-01-02")), nil)
	handler.ServeHTTP(recorder, req)

	assert.Equal(
		2,
		strings.Count(recorder.Body.String(), "mealId"),
	)
	mids, err := parseMealIds(recorder.Body.String(), "mealId")
	assert.Equal(1, len(mids), "exactly one mealId should be returned")
	assert.Nil(err)

	cids, err := parseMealIds(recorder.Body.String(), "consumerIds")
	assert.Nil(err)
	assert.Equal(1, len(cids), "exactly one consuerrIds array should be returned")
	assert.Equal("[2,3,4],", cids[0])
}

func parseMealIds(body string, keystr string) ([]string, error) {
	r, _ := regexp.Compile(fmt.Sprintf("\"%v\":([0-9\\[\\],]+)", keystr))

	sm := r.FindStringSubmatch(body)
	if len(sm) < 2 {
		return nil, errors.New("could not parse Id from the meals json")
	}

	return sm[1:], nil
}
