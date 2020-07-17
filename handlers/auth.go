package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jrolfo/restapi/models"
	"golang.org/x/crypto/bcrypt"
)

//Claims struct
// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// @route   POST api/auth
// @desc    Auth user
// @access  Public

//AuthUser function
func AuthUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	jwtKey := []byte(models.Config.JwtKey)

	w.Header().Set("Content-Type", "application/json")
	user := models.User{}
	_ = json.NewDecoder(r.Body).Decode(&user)
	resp := models.AuthResponse{}

	//First read user being passed to the auth function
	userResponse := models.ReadUser(db, user.Username)

	if !userResponse.Success {
		resp.Message = "User not found"
		resp.Success = false
		json.NewEncoder(w).Encode(resp)
		return
	}

	dbPassword := []byte(userResponse.Data[0].Password)
	reqPassword := []byte(user.Password)

	err := bcrypt.CompareHashAndPassword(dbPassword, reqPassword)

	if err != nil {
		resp.Message = "Credentials invalid"
		resp.Success = false
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 10 minutes
	expirationTime := time.Now().Add(models.Config.Expires * time.Minute)
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
	tokenString, errToken := token.SignedString(jwtKey)

	if errToken != nil {
		resp.Message = "Error singing user: " + errToken.Error()
		resp.Success = false
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Token = tokenString
	resp.Success = true
	resp.Message = "User authorized"
	json.NewEncoder(w).Encode(resp)

}
