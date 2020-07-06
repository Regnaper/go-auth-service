package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var signingKey = []byte(os.Getenv("SECRET_KEY")) // signing key for JWT
var domain = "rocky-escarpment-09841.herokuapp.com" // server domain name for cookie
var port = os.Getenv("PORT")
var databaseName = "tokensdb"

type Token struct {
	GUID         string
	AccessToken  string
	RefreshToken string
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/secret/{guid}", RefreshHandler).Methods("POST")
	router.HandleFunc("/secret/{guid}", DeleteHandler).Methods("DELETE")
	router.HandleFunc("/secret/{guid}", CreateHandler).Methods("GET")
	router.HandleFunc("/secrets/{guid}", DeleteAllHandler).Methods("DELETE")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Application server")
	})

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	if port == "" {
		port = "9000" // Default port if not specified
	}
	http.ListenAndServe(":"+port, nil)
}
