package db

import (
	"bufio"
	"context"
	"log"
	"os"
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

// Criar a função ExecScript para executar o script a partir de um arquivo
func ExecScript(conn DBConnection, script string) error {
	file, err := os.Open(script)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback(context.Background())

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			_, err := tx.Exec(context.Background(), scanner.Text())
			if err != nil {
				return err
			}
		}

		err = tx.Commit(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
