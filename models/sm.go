package models

type Secret struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"passwor"`
	JWTSign  string `json:"jwtsign"`
	DataBase string `json:"database"`
}
