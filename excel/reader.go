package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// Lê a planilha e retorna os dados
func LerPlanilha(filePath, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir a planilha: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter as linhas da aba: %v", err)
	}
	return rows, nil
}
