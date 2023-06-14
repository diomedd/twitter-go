package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/diomedd/twitter-go/awsgo"
	"github.com/diomedd/twitter-go/bd"
	"github.com/diomedd/twitter-go/handlers"
	"github.com/diomedd/twitter-go/models"
	"github.com/diomedd/twitter-go/secretmanager"
	//"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {

	lambda.Start(EjecutoLambda)
	//awsgo.InicializoAWS()

}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "error en las variables de entorno, deben incliuir 'SecretName', 'BacketName', 'UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "application-json",
			},
		}
		return res, nil
	}

	SecretModels, err := secretmanager.GetSecrets(os.Getenv("SecretName"))

	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "error en la lectura de secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application-json",
			},
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["twitter-dd"], os.Getenv("UrlPrefix"), "", -1)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModels.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModels.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModels.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModels.DataBase)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModels.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	//chequeo con con la bd
	err = bd.ConectarBD(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "error conectando en la base de datos " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "aplication-json",
			},
		}
		return res, nil
	}
	restAPI := handlers.Manejadores(awsgo.Ctx, request)
	if restAPI.CustomResp == nil {
		if err != nil {
			res = &events.APIGatewayProxyResponse{
				StatusCode: restAPI.Status,
				Body:       restAPI.Message,
				Headers: map[string]string{
					"Content-Type": "aplication-json",
				},
			}
			return res, nil
		} else {
			return restAPI.CustomResp, nil
		}

	}
	return res, nil
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

	return true
}
