package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//User struct
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Password string `json:"password"`
}

//UserResponse struct
type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []User `json:"data"`
	Token   string `json:"token"`
}

var expires = Config.Expires

//Claims struct
// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//CreateUser function
func CreateUser(db *sql.DB, user User) UserResponse {

	res := UserResponse{}
	b := []byte(user.Password)

	hash, hassErr := hashAndSalt(b)

	if hassErr != nil {

		return handleUserError(res, "Error generating hash & salt: "+hassErr.Error())
	}

	insForm, prepareErr := db.Prepare("INSERT INTO users(username, name, lastname, password) VALUES(?,?,?,?)")
	if prepareErr != nil {
		return handleUserError(res, "Error inserting: "+prepareErr.Error())

	}
	_, insertErr := insForm.Exec(user.Username, user.Name, user.Lastname, hash)

	if insertErr != nil {
		return handleUserError(res, "Error creating user: "+insertErr.Error())

	}

	// Declare the expiration time of the token
	// here, we have kept it as 10 minutes
	expirationTime := time.Now().Add(expires * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	jwtKey := []byte(Config.JwtKey)
	tokenString, errToken := token.SignedString(jwtKey)

	if errToken != nil {
		return handleUserError(res, "Error singing user: "+errToken.Error())
	}

	res.Token = tokenString
	res.Data = append(res.Data, user)
	res.Success = true
	res.Message = "User created"

	return res
}

//ReadUsers function
func ReadUsers(db *sql.DB) UserResponse {

	data := []User{}
	user := User{}
	res := UserResponse{}

	results, err := db.Query("SELECT * FROM users")

	if err != nil {
		return handleUserError(res, err.Error())

	}

	for results.Next() {
		var userName, lastName, name, password string

		err = results.Scan(&userName, &name, &lastName, &password)
		if err != nil {
			return handleUserError(res, err.Error())
		}
		user.Username = userName
		user.Name = name
		user.Lastname = lastName
		user.Password = password

		data = append(data, user)
	}
	res.Success = true
	res.Message = ""
	res.Data = data

	return res
}

//ReadUser function
func ReadUser(db *sql.DB, parmID string) UserResponse {

	data := []User{}
	user := User{}
	res := UserResponse{}

	results, err := db.Query("SELECT * FROM users WHERE username=?", parmID)

	if err != nil {
		return handleUserError(res, err.Error())
	}

	for results.Next() {
		var userName, lastName, name, password string

		err = results.Scan(&userName, &name, &lastName, &password)
		if err != nil {
			return handleUserError(res, err.Error())
		}
		user.Username = userName
		user.Name = name
		user.Lastname = lastName
		user.Password = password

		data = append(data, user)
	}
	//If no data then that means that the user was not found
	if len(data) == 0 {
		res.Data = data
		return handleUserError(res, "User not found")
	}

	res.Success = true
	res.Message = ""
	res.Data = data

	return res
}

//UpdateUser func
func UpdateUser(db *sql.DB, paramID string, user User) UserResponse {
	data := []User{}
	res := UserResponse{}

	q := ReadUser(db, paramID)

	if len(q.Data) == 0 {
		return handleUserError(res, "User not found")
	}
	if user.Password != "" {
		b := []byte(user.Password)
		hash, hassErr := hashAndSalt(b)

		if hassErr != nil {
			return handleUserError(res, "Error on generating hash & salt: "+hassErr.Error())
		}

		insForm, err := db.Prepare("UPDATE users SET name=?, lastname=?, password=? WHERE username=?")

		if err != nil {
			return handleUserError(res, "Error updating: "+err.Error())
		}

		insForm.Exec(user.Name, user.Lastname, hash, paramID)

	} else {
		insForm, err := db.Prepare("UPDATE users SET name=?, lastname=? WHERE username=?")
		if err != nil {
			return handleUserError(res, "Error updating: "+err.Error())
		}

		insForm.Exec(user.Name, user.Lastname, paramID)
	}

	user.Username = paramID
	data = append(data, user)
	res.Success = true
	res.Message = "Successfull update"
	res.Data = data

	return res
}

//DeleteUser func
func DeleteUser(db *sql.DB, paramID string) UserResponse {
	data := []User{}

	res := UserResponse{}
	res.Data = data

	q := ReadUser(db, paramID)

	if len(q.Data) == 0 {
		return handleUserError(res, "User not found")
	}

	insForm, err := db.Prepare("DELETE FROM users WHERE username=?")

	if err != nil {
		return handleUserError(res, "Error deleting: "+err.Error())

	}
	insForm.Exec(paramID)

	res.Message = "Document deleted"
	res.Success = true

	return res
}

//ValidateToken func
func ValidateToken(tknStr string) bool {
	jwtKey := []byte(Config.JwtKey)

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println(err.Error())
			return false
		}
		fmt.Println(err.Error())
		return false
	}
	if !tkn.Valid {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// hash and salt passwords
func hashAndSalt(pwd []byte) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), err
}

func handleUserError(res UserResponse, err string) UserResponse {
	res.Success = false
	res.Message = err

	return res
}
