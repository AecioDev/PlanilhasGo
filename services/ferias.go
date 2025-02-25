package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"planilhas/db"
	"planilhas/models"
)

// Processa as linhas da planilha e insere no banco
func ProcessarFerias(conn db.DBConnection, rows [][]string, linhaIni int) error {
	for i := linhaIni; i < len(rows); i++ {
		row := rows[i]
		numColunas := 9

		// Garantir que a linha tenha exatamente 'numColunas' elementos
		if len(row) < numColunas {
			for len(row) < numColunas {
				row = append(row, "") // Completa com strings vazias
			}
		} else if len(row) > numColunas {
			row = row[:numColunas] // Trunca caso tenha mais colunas do que o esperado
		}

		if (strings.TrimSpace(row[5]) != "" && strings.TrimSpace(row[6]) != "") || (strings.TrimSpace(row[7]) != "" && strings.TrimSpace(row[8]) != "") { // Verifica se tem pelo menos um período de férias preenchido

			nomeColaborador := strings.ToLower(strings.TrimSpace(row[0]))
			cpfColaborador := strings.TrimSpace(row[1])
			cpfColaborador = strings.ReplaceAll(cpfColaborador, ".", "")
			cpfColaborador = strings.ReplaceAll(cpfColaborador, "-", "")

			colaborador, err := BuscarColaborador(conn, "efetivo", "cpf", cpfColaborador)
			if err != nil {
				if err.Error() == "no rows in result set" {
					// Tenta buscar no colaborador terceirizado
					log.Printf("Colaborador efetivo não encontrado, tentando buscar em terceirizados...")
					colaborador, err = BuscarColaborador(conn, "terceirizado", "cpf", cpfColaborador)
					if err != nil {
						log.Printf("Erro ao buscar colaborador terceirizado: %v", err)
					}
				} else {
					// Outro erro, trata como um problema mais crítico
					log.Printf("Erro ao buscar colaborador efetivo: %v", err)
				}
			}

			if colaborador != nil {

				idStr := fmt.Sprintf("%d", colaborador.ID)

				Empresa := strings.TrimSpace(row[3])
				var EmpresaID int
				switch Empresa {
				case "AZ TECNOLOGIA":
					EmpresaID = 1

				case "MILTEC":
					EmpresaID = 2

				case "GEOI2":
					EmpresaID = 3

				case "DIGIX":
					EmpresaID = 6

				case "GUATÓS":
					EmpresaID = 7

				case "INOVVATI":
					EmpresaID = 8
				}

				var UserEmpresaID int
				switch EmpresaID {
				case 1:
					UserEmpresaID = 2 // 2 - César Augusto Sanches

				case 2:
					UserEmpresaID = 4 // 4 - Guilherme Castilho Matos

				case 3:
					UserEmpresaID = 9 // 9 - Phelipe Gomes de Melo

				case 6:
					UserEmpresaID = 12 // 12 - Digithobrasil Solucoes em Software LTDA

				case 7:
					UserEmpresaID = 0 // 0 - Sem Usuário

				case 8:
					UserEmpresaID = 11 // 11 - Michel Mendes
				}

				//log.Printf("EmpresaID: %s, UserEmpresaID: %d\n", Empresa, UserEmpresaID)

				ferias := models.Ferias{
					Colaborador_id:   idStr,
					Colaborador_nome: nomeColaborador,
					Cargo:            row[2],
					Empresa_id:       row[3],
					Equipe_sigla:     row[4],
					Data_inicio_p1:   row[5],
					Data_fim_p1:      row[6],
					Data_inicio_p2:   row[7],
					Data_fim_p2:      row[8],
					Lider_E_id:       *colaborador.LiderE,
					Lider_T_id:       *colaborador.LiderT,
					Gestor_id:        *colaborador.Gestor,
					UserEmpresa_id:   UserEmpresaID,
				}

				log.Printf("Cadastrando Férias do Colaborador: %s\n", nomeColaborador)

				if err := CadastraFerias(conn, ferias); err != nil {
					log.Printf("Erro ao cadastrar férias do colaborador: %s. Erro: %v", nomeColaborador, err)
				}

			} else {
				log.Printf("Colaborador Nome: %s, não encontrado no BD\n", nomeColaborador)
			}

			//break
		}
	}

	//Verifica se o arquivo 'scriptFerias.sql' existe e executa o script criado
	if _, err := os.Stat("scriptFerias.sql"); err == nil {
		log.Println("Executando scriptFerias.sql")
		err := db.ExecScript(conn, "scriptFerias.sql")
		if err != nil {
			log.Printf("Erro ao executar scriptFerias.sql: %v", err)
		}
	} else {
		log.Println("Arquivo 'scriptFerias.sql' não encontrado")
	}

	return nil
}

// Cadastra a Féria do Colaborador no Banco de Dados
func CadastraFerias(conn db.DBConnection, ferias models.Ferias) error {

	LiderE := fmt.Sprintf("%d", ferias.Lider_E_id)
	if ferias.Lider_E_id <= 0 {
		LiderE = "null"
	}

	LiderT := fmt.Sprintf("%d", ferias.Lider_T_id)
	if ferias.Lider_T_id <= 0 {
		LiderT = "null"
	}

	/*
		WITH nova_solicitacao AS (
			INSERT INTO public.solicitacao_ferias(
					created_at, updated_at, data_avaliacao_lider, status, terceirizado_id, lider_colaborador_efetivo_id, lider_colaborador_terceirizado_id)
				VALUES(now(), now(), now(), 'aprovado_por_lider', 3, 10, null) RETURNING id)
				INSERT INTO periodo_ferias (data_inicio, data_fim, solicitacao_ferias_id) VALUES('2025-10-13', '2025-10-23', (SELECT id FROM nova_solicitacao));

		se tiver mais de um período adiciona:
				query = query + ', ('2025-10-13', '2025-10-23', (SELECT id FROM nova_solicitacao))';

		INSERT INTO public.solicitacao_ferias
			(id, created_at, updated_at, deleted_at, status, data_avaliacao_lider, data_avaliacao_gestor, data_avaliacao_empresa, descricao_avaliacao_lider, descricao_avaliacao_gestor, descricao_avaliacao_empresa, terceirizado_id, lider_colaborador_efetivo_id, lider_colaborador_terceirizado_id, efetivo_gestor_id, usuario_externo_empresa_terceirizada_id)
			VALUES(nextval('solicitacao_ferias_id_seq'::regclass), '', '', '', '', '', '', '', '', '', '', 0, 0, 0, 0, 0);

	*/

	query := fmt.Sprintf("/*%s*/ ", ferias.Colaborador_nome)

	query += fmt.Sprintf("WITH nova_solicitacao AS ("+
		"INSERT INTO public.solicitacao_ferias("+
		"created_at, updated_at, status, terceirizado_id, "+
		"data_avaliacao_lider, lider_colaborador_efetivo_id, lider_colaborador_terceirizado_id, descricao_avaliacao_lider, "+
		"data_avaliacao_gestor, efetivo_gestor_id, descricao_avaliacao_gestor, "+
		"data_avaliacao_empresa, usuario_externo_empresa_terceirizada_id, descricao_avaliacao_empresa) "+
		"VALUES(now(), now(), 'aprovado_por_empresa', %s, "+
		"now(), %s, %s, 'Ferias aprovadas via importação de planilha excel', "+
		"now(), %d, 'Ferias aprovadas via importação de planilha excel', "+
		"now(), %d, 'Ferias aprovadas via importação de planilha excel') RETURNING id) ", ferias.Colaborador_id, LiderE, LiderT, ferias.Gestor_id, ferias.UserEmpresa_id)

	if ferias.Data_inicio_p1 != "" && ferias.Data_fim_p1 != "" {
		query = query + fmt.Sprintf(`INSERT INTO periodo_ferias (data_inicio, data_fim, solicitacao_ferias_id) VALUES('%s', '%s', (SELECT id FROM nova_solicitacao))`, ferias.Data_inicio_p1, ferias.Data_fim_p1)
	}

	if ferias.Data_inicio_p2 != "" && ferias.Data_fim_p2 != "" {
		query = query + fmt.Sprintf(`, ('%s', '%s', (SELECT id FROM nova_solicitacao))`, ferias.Data_inicio_p2, ferias.Data_fim_p2)
	}

	query = query + ";\n"

	//fmt.Printf("Query Férias: %s\n", query)

	// Grava a query em um arquivo
	f, err := os.OpenFile("scriptFerias.sql", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Erro ao abrir o arquivo: %v", err)
		return err
	}

	if _, err := f.WriteString(query); err != nil {
		log.Printf("Erro ao escrever no arquivo: %v", err)
		return err
	}

	defer f.Close()

	return nil
}

/*
func formataData(dataStr string) string {
	// Definir o layout correto para o formato de entrada (DD-MM-YY)
	layoutEntrada := "02-01-06"

	// Converte a string para o formato de data
	data, err := time.Parse(layoutEntrada, dataStr)
	if err != nil {
		fmt.Println("Erro ao converter data:", err)
		return ""
	}

	// Formatar para o formato desejado (YYYY-MM-DD)
	layoutSaida := "2006-01-02"
	dataFormatada := data.Format(layoutSaida)
	return dataFormatada
}
*/
