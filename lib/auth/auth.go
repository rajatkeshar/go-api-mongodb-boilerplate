package auth

import (
  "log"
  "time"
  "net/http"
  "encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
  "github.com/go-api-mongodb-boilerplate/models"
)

const (
  mySigningKey = "MySecret"
)

type UsersClaim struct {
  models.User
  jwt.StandardClaims
}

type ResponseBody struct {
    Success   bool    `json:"success"`
    Msg       string  `json:"msg"`
    Data      string  `json:"data"`
}

func GenerateJWT(user models.User) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["authorized"] = true
    claims["_id"] = user.ID.Hex()
    claims["email"] = user.Email
    claims["firstname"] = user.Firstname
    claims["lastname"] = user.Lastname
    claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

    tokenString, err := token.SignedString([]byte(mySigningKey))

    if err != nil {
        log.Println("Something Went Wrong: %s", err.Error())
        return "", err
    }

    return tokenString, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
    var responseBody ResponseBody
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header["Token"] != nil {
            var usersClaim UsersClaim
            token, err := jwt.ParseWithClaims(r.Header["Token"][0], &usersClaim, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    responseBody.Success = false
                    responseBody.Msg = "error in parsing claims"
                    return nil, json.NewEncoder(w).Encode(responseBody)
                }
                return []byte(mySigningKey), nil
            })

            if err != nil {
                responseBody.Success = false
                responseBody.Msg = err.Error()
                json.NewEncoder(w).Encode(responseBody)
            }

            if token.Valid {
                w.Header().Set("_id", usersClaim.ID.Hex())
                w.Header().Set("email", usersClaim.Email)
                w.Header().Set("firstname", usersClaim.Firstname)
                w.Header().Set("lastname", usersClaim.Lastname)
                endpoint(w, r)
            }
        } else {
            responseBody.Success = false
            responseBody.Msg = "authorization error"
            json.NewEncoder(w).Encode(responseBody)
        }
    })
}

func VerifyToken(w http.ResponseWriter, token1 string) error {
    var usersClaim UsersClaim
    token, err := jwt.ParseWithClaims(token1, &usersClaim, func(token *jwt.Token) (interface{}, error) {
        return []byte(mySigningKey), nil
    })

    if token.Valid {
        w.Header().Set("_id", usersClaim.ID.Hex())
        w.Header().Set("email", usersClaim.Email)
        w.Header().Set("firstname", usersClaim.Firstname)
        w.Header().Set("lastname", usersClaim.Lastname)
        return nil
    }
    return err
}
