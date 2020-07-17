package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jrolfo/restapi/models"
)

// @route   POST api/users
// @desc    Register new user
// @access  Public

//CreateUser function
func CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	user := models.User{}
	_ = json.NewDecoder(r.Body).Decode(&user)
	resp := models.CreateUser(db, user)
	json.NewEncoder(w).Encode(resp)
}

// @route   get api/users
// @desc    get one user
// @access  private

//GetUser Function
func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parmID, _ := params["username"]
	resp := models.ReadUser(db, parmID)
	json.NewEncoder(w).Encode(resp)
}

// @route   PUT api/users
// @desc    Update user
// @access  private

//UpdateUser function
func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parmID, _ := params["username"]
	user := models.User{}
	_ = json.NewDecoder(r.Body).Decode(&user)
	ret := models.UpdateUser(db, parmID, user)
	json.NewEncoder(w).Encode(ret)
}

// @route   get api/users
// @desc    get all users
// @access  private

//GetUsers function
func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	resp := models.ReadUsers(db)
	json.NewEncoder(w).Encode(resp)
}

// @route   delete api/users
// @desc    delete one user
// @access  private

//DeleteUser Function
func DeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parmID, _ := params["username"]

	ret := models.DeleteUser(db, parmID)
	json.NewEncoder(w).Encode(ret)
}
