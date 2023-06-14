package bd

import (
	"context"

	"github.com/diomedd/twitter-go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertoRegistro(u models.Usuario) (string, bool, error) {

	ctx := context.TODO()

	db := MongCN.Database(DatabaseName)
	col := db.Collection("usuarios")

	u.Password, _ = EncriptarPassword(u.Password)

	//funcion de mongo, insertone solo inserta 1 registro
	result, err := col.InsertOne(ctx, u)
	if err != nil {
		return "", false, err
	}

	//objid captura el id y lo inserta
	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil

}
