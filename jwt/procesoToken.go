package jwt

import (
	"errors"
	"strings"

	"github.com/diomedd/twitter-go/models"
	jwt "github.com/golang-jwt/jwt/v5"
)

/*Email valor de Email usado en todos los EndPoints */
var Email string

/*IDUsuario es el ID devuelto del modelo, que se usará en todos los EndPoints */
var IDUsuario string

func ProcesoToken(tk string, JWTSign string) (*models.Claim, bool, string, error) {

	miClave := []byte(JWTSign)
	var claims models.Claim

	splitToken := strings.Split(tk, "Bearer")

	if len(splitToken) != 2 {

		return &claims, false, string(""), errors.New("formato de token invalido")
	}

	tk = strings.TrimSpace(splitToken[1])

	tkn, err := jwt.ParseWithClaims(tk, &claims, func(token *jwt.Token) (interface{}, error) {

		return miClave, nil

	})

	if err == nil {

		//rutina para chequear contra la bd
	}
	if !tkn.Valid {

		return &claims, false, string(""), errors.New("token invalido")
	}

	return &claims, false, string(""), err

}
