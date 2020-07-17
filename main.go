package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/jrolfo/restapi/handlers"
	"github.com/jrolfo/restapi/models"
)

func main() {

	//Get Configuation
	models.InitConfig()

	// init mux router
	r := mux.NewRouter()

	//Init database
	db := dbConn()

	//TODO All the CORS CRAP

	//route handlers /endpoints

	//Private endopoins are wrapped into the auth function which validates the token in the header and tells if user is //logged or not

	//Books
	r.HandleFunc("/api/books",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.GetBooks(w, r, db) })
		}).Methods("GET")

	r.HandleFunc("/api/books/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.GetBook(w, r, db) })
		}).Methods("GET")

	r.HandleFunc("/api/books",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.CreateBook(w, r, db) })
		}).Methods("POST")

	r.HandleFunc("/api/books/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.UpdateBook(w, r, db) })
		}).Methods("PUT")

	r.HandleFunc("/api/books/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.DeleteBook(w, r, db) })
		}).Methods("DELETE")

	//Users
	r.HandleFunc("/api/users/{username}",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.GetUser(w, r, db) })
		}).Methods("GET")

	r.HandleFunc("/api/users",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.CreateUser(w, r, db) })
		}).Methods("POST")

	r.HandleFunc("/api/users",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.GetUsers(w, r, db) })
		}).Methods("GET")

	r.HandleFunc("/api/users/{username}",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.UpdateUser(w, r, db) })
		}).Methods("PUT")

	r.HandleFunc("/api/users/{username}",
		func(w http.ResponseWriter, r *http.Request) {
			auth(w, r, func() { handlers.DeleteUser(w, r, db) })
		}).Methods("DELETE")

	//Public endpoints are just wrapped in a dummy function so we can pass DB to the handlers, MUX wont let me pass an //extra parm

	//Auth
	r.HandleFunc("/api/auth",
		func(w http.ResponseWriter, r *http.Request) { handlers.AuthUser(w, r, db) }).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", r))
	fmt.Println("Server Listing")
	defer db.Close()

}

//Auth function used as wrapper for private endpoints
func auth(w http.ResponseWriter, r *http.Request, f func()) {
	authorized := models.ValidateToken(r.Header.Get("x-auth-token"))
	if authorized {
		f()
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func dbConn() (db *sql.DB) {
	//Format root:@tcp(127.0.0.1:3306)/test
	dbDriver := "mysql"
	dbUser := models.Config.User
	dbPass := models.Config.Password
	dbServer := models.Config.Server
	dbPort := models.Config.Port
	dbName := models.Config.Database
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbServer+":"+dbPort+")/"+dbName)

	fmt.Println(dbDriver, dbUser+":"+dbPass+"@tcp("+dbServer+":"+dbPort+")/"+dbName)

	if err != nil {
		panic(err.Error())
	}
	return db
}
