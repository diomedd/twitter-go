package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diomedd/twitter-go/bd"
	"github.com/diomedd/twitter-go/models"
)

func Registro(ctx context.Context) models.RestApi {

	var t models.Usuario
	var r models.RestApi
	r.Status = 400

	fmt.Println("entre a registro")

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)

	if err != nil {

		r.Message = err.Error()
		fmt.Println((r.Message))
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "debe especificar el email"
		fmt.Println(r.Message)
		return r
	}

	if len(t.Password) < 6 {
		r.Message = "debe especificar una contraseña de al menos 6 caracteres"
		fmt.Println(r.Message)
		return r
	}

	_, encontrado, _ := bd.ChequeoYaExisteUsuario(t.Email)

	if encontrado {
		r.Message = "ya existe un usuario registrado con ese email"
		fmt.Println(r.Message)
		return r
	}

	_, status, err := bd.InsertoRegistro(t)

	if err != nil {

		r.Message = "Ocurrio un error al intentar realizar el registro del usuario " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if !status {
		r.Message = "no se logro insertar el registro del usuario " + err.Error()
		fmt.Println(r.Message)
		return r
	}
	r.Status = 200
	r.Message = "registro ok"
	fmt.Println(r.Message)
	return r
}
