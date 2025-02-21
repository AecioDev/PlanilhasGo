package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type DBConnection struct {
	*pgx.Conn
}

// Conecta ao banco de dados
func ConnectToDB(amb string) (DBConnection, error) {

	connStr := "postgres://postgres:12345@localhost:5432/portal_conhecimento"
	switch amb {
	case "hom":
		connStr = "postgres://postgres:MastyJ8675@s0925.ms:5432/portal_conhecimento"
	case "dev":
		connStr = "postgres://postgres:MastyJ8675@s0925.ms:5432/portal_conhecimento_dev"
	}
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	return DBConnection{conn}, nil
}
