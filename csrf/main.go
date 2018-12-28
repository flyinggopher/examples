// Copyright (c) 2018. Flying Gopher Authors
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"database/sql"

	"github.com/flyingjamnik/csrf"
	"github.com/flyingjamnik/errors"

	_ "github.com/mattn/go-sqlite3"
)

type data struct {
	Token string
}

var (
	Base                 *sql.DB
	CSRFTokenDoesntExist = errors.New("CSRF token doesn't exist.")
)

func main() {
	log.Println("Connecting with database.")
	storage := csrf.NewStorage(":memory:")
	err := storage.CreateTables()

	log.Println("Starting website.")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tmpl, err = template.ParseFiles("index.html")
			if err != nil {
				http.Error(w, errors.InternalServerErr.Error(), 500)
				return
			}

			c := csrf.RegisterCSRF()

			err = storage.SaveCSRF(c)

			if err != nil {
				http.Error(w, errors.InternalServerErr.Error(), 500)
				return
			}

			tokenToHidden := data{c.Token}

			err = tmpl.Execute(w, tokenToHidden)

			if err != nil {
				http.Error(w, errors.InternalServerErr.Error(), 500)
				return
			}

		case http.MethodPost:
			tokenFromHidden := r.PostFormValue("token-hidden")

		default:
			http.Error(w, errors.UnauthorizedRequestErr.Error(), 401)
		}
	})

	http.ListenAndServe(":8080", mux)
}
