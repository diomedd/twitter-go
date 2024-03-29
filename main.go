package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/diomedd/twitter-go/awsgo"
	"github.com/diomedd/twitter-go/bd"
	"github.com/diomedd/twitter-go/handlers"
	"github.com/diomedd/twitter-go/models"
	"github.com/diomedd/twitter-go/secretmanager"
)

func main() {

	//lambda.Start(EjecutoLambda)
	//bd.BaseConectada()

	/*	clientOptions := options.Client().ApplyURI("mongodb+srv://ddiomede:Pinguin.2021@twitter.urcagqz.mongodb.net/")
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Comprobar la conexión
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Conexión exitosa a MongoDB")

		// Cerrar la conexión al finalizar
		defer func() {
			if err := client.Disconnect(context.Background()); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Conexión cerrada correctamente")
		}()*/
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	var res *events.APIGatewayProxyResponse
	awsgo.InicializoAWS()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno. deben incluir 'SecretName', 'BucketName",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["twittergo"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.DataBase)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// Chequeo Conexión a la BD o Conecto la BD

	bd.ConectarBD(awsgo.Ctx)

	respAPI := handlers.Manejadores(awsgo.Ctx, request)

	fmt.Println("Sali de Manejadores")
	if respAPI.CustomResp == nil {
		headersResp := map[string]string{
			"Content-Type": "application/json",
		}
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       string(respAPI.Message),
			Headers:    headersResp,
		}
		return res, nil
	} else {
		return respAPI.CustomResp, nil
	}
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}
	return traeParametro
}
