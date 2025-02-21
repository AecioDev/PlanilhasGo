package db

import (
	"context"
	"log"
	"planilhas/config"

	"github.com/jackc/pgx/v4"
)

type DBConnection struct {
	*pgx.Conn
}

// Conecta ao banco de dados
func ConnectToDB(amb string) (DBConnection, error) {

	connStr := config.GetStrConnection(amb)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	return DBConnection{conn}, nil
}
