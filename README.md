## Gotchas
In the documentation, the `configure_cooking.go` congifures the API. 

In the `func configureAPI(api *operations.CookingAPI) http.Handler` there is an condition:
`if api.TodosAddOneHandler == nil {` on the existence of the handler. 

Remove the check, otherwiser it won't get configured.