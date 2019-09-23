package auth

import (
  "fmt"
  "time"
  "net/http"
  "encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
  "github.com/GORest-API-MongoDB/models"
)

const (
  mySigningKey = "MySecret"
)

type UsersClaim struct {
  models.User
  jwt.StandardClaims
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
        fmt.Println("Something Went Wrong: %s", err.Error())
        return "", err
    }

    return tokenString, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header["Token"] != nil {
            var usersClaim UsersClaim
            token, err := jwt.ParseWithClaims(r.Header["Token"][0], &usersClaim, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("There was an error")
                }
                return []byte(mySigningKey), nil
            })

            if err != nil {
                //fmt.Fprintf(w, err.Error())
                json.NewEncoder(w).Encode(err.Error())
            }

            if token.Valid {
                w.Header().Set("_id", usersClaim.ID.Hex())
                w.Header().Set("email", usersClaim.Email)
                w.Header().Set("firstname", usersClaim.Firstname)
                w.Header().Set("lastname", usersClaim.Lastname)
                endpoint(w, r)
            }
        } else {
            fmt.Fprintf(w, "Not Authorized")
        }
    })
}
