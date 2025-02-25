package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"planilhas/db"
	"planilhas/excel"
	"planilhas/services"
	"strconv"

	"github.com/sqweek/dialog"
)

func getFilePath() string {
	filePath, err := dialog.File().Filter("Planilhas Excel", "xlsx").Title("Selecione o arquivo Excel").Load()
	if err != nil {
		log.Fatalf("Erro ao selecionar o arquivo: %v", err)
	}
	return filePath
}

// Importar Férias: go run main.go ferias loc 2
// Importar Sistemas: go run main.go sistema loc 2

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Uso: go run main.go <comando> <ambiente: loc, hom, dev> <Linha Inicial>")
	}

	comando := os.Args[1]
	ambiente := os.Args[2]

	linhaIni, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("A Linha Inicial deve ser um Número Inteiro! ", err)
		return
	}

	// Conectar ao banco de dados
	conn, err := db.ConnectToDB(ambiente)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer conn.Close(context.Background())

	filePath := getFilePath()
	fmt.Println("Arquivo selecionado:", filePath)

	switch comando {
	case "sistema":
		err = importarSistema(conn, filePath, linhaIni)
		if err != nil {
			log.Fatalf("Erro ao importar dados: %v", err)
		}
	case "ferias":
		err = importarFerias(conn, filePath, linhaIni)
		if err != nil {
			log.Fatalf("Erro ao importar dados: %v", err)
		}
	default:
		log.Fatalf("Comando não reconhecido: %s", comando)
	}
}

// Função para importar Sistemas a partir dos dados do Excel e salvar no banco
func importarSistema(conn db.DBConnection, filePath string, linha int) error {
	// Lê os dados da planilha
	rows, err := excel.LerPlanilha(filePath, "Planilha1")
	if err != nil {
		return fmt.Errorf("erro ao ler planilha: %v", err)
	}

	// Insere os dados no banco
	return services.ProcessarSistemas(conn, rows, linha)
}

// Função para importar os dados do Excel e salvar no banco
func importarFerias(conn db.DBConnection, filePath string, linha int) error {
	// Lê os dados da planilha
	rows, err := excel.LerPlanilha(filePath, "Planilha1")
	if err != nil {
		return fmt.Errorf("erro ao ler planilha: %v", err)
	}

	// Insere os dados no banco
	return services.ProcessarFerias(conn, rows, linha)
}
