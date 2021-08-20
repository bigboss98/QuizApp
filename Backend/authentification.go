/*
 * Authentification File used to manage all stuff related to Authentification and JWT Token
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const secretkey = "GiorgioneMagoDelGuanciale"

type Authentication struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Name        string `json:"name"`
	TokenString string `json:"token"`
}

func GenerateJWT(name string, role string) (string, error) {
	/*
	 * Generate JWT token given name and role(admin or user)
	 *
	 * Params:
	 * -name(string): Name of User whom we generate JWT token
	 * -role(string): Role of User (admin or user)
	 *
	 * Return the Token as string with also an error object to indicate whether some error occurs
	 * during JWT generation.
	 */
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["name"] = name
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func GeneratehashPassword(password string) (string, error) {
	/*
	 * Generate an Hash Password given visible password
	 *
	 * Params:
	 * -password(string): password of User
	 *
	 * Return Hashed password as string with also an error object
	 * which indicates eventually error happened in Hash generation.
	 */
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func encodeSignIn(token Token, indent string, prefix string) string {
	json_user, _ := json.MarshalIndent(token, prefix, indent)
	return string(json_user)
}

func IsAuthorized(string_token string) bool {
	if string_token == "" {
		log.Printf("No Token Found")
		return false
	}
	var mySigningKey = []byte(secretkey)

	token, err := jwt.Parse(string_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		log.Printf("Your Token has been expired")
		return false
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
		/*
			if claims["role"] == "admin" {

				request.Header.Set("Role", "admin")
				handler.ServeHTTP(response, request)
				return

			} else if claims["role"] == "user" {

				request.Header.Set("Role", "user")
				handler.ServeHTTP(response, request)
				return
			}
		*/
	}
	log.Printf("Not Authorized")
	return false
}
