package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jrolfo/restapi/models"
)

//GetBooks function
func GetBooks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	resp := models.ReadBooks(db)
	json.NewEncoder(w).Encode(resp)
}

//GetBook Function
func GetBook(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parmID, _ := strconv.Atoi(params["id"])
	resp := models.ReadBook(db, parmID)
	json.NewEncoder(w).Encode(resp)
}

//CreateBook function
func CreateBook(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	book := models.Book{}
	_ = json.NewDecoder(r.Body).Decode(&book)
	resp := models.CreateBook(db, book)
	json.NewEncoder(w).Encode(resp)

}

//UpdateBook function
func UpdateBook(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parmID, _ := strconv.Atoi(params["id"])
	book := models.Book{}
	_ = json.NewDecoder(r.Body).Decode(&book)
	ret := models.UpdateBook(db, parmID, book)
	json.NewEncoder(w).Encode(ret)
}

//DeleteBook Function
func DeleteBook(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parmID, _ := strconv.Atoi(params["id"])

	ret := models.DeleteBook(db, parmID)
	json.NewEncoder(w).Encode(ret)
}
