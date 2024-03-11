package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type user struct {
	Name string `json:"name"`
}

var allUsers = []user{
	{Name: "Dan"},
	{Name: "Kate"},
	{Name: "Jack"},
	{Name: "John"},
	{Name: "Alex"},
	{Name: "Bob"},
	{Name: "Tom"},
	{Name: "Sam"},
	{Name: "Alice"},
	{Name: "Jane"},
	{Name: "Mark"},
}

func getIntQueryParam(r *http.Request, name string, defaultValue string) (int, error) {
	param := r.URL.Query().Get(name)
	if param == "" {
		param = defaultValue
	}

	return strconv.Atoi(param)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		page, err := getIntQueryParam(r, "page", "1")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		limit, err := getIntQueryParam(r, "limit", "10")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if limit > len(allUsers) {
			limit = len(allUsers)
		}

		offset := (page - 1) * limit
		users := []user{}
		for i := offset; i < offset+limit; i++ {
			if i >= len(allUsers) {
				break
			}

			users = append(users, allUsers[i])
		}

		var nextPage *int
		if offset+limit < len(allUsers) {
			x := page + 1
			nextPage = &x
		}

		type response struct {
			Data     []user `json:"data"`
			Total    int    `json:"total"`
			NextPage *int   `json:"nextPage"`
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response{
			Data:     users,
			Total:    len(allUsers),
			NextPage: nextPage,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
