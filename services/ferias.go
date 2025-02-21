package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"planilhas/db"
	"planilhas/models"
)

// Processa as linhas da planilha e insere no banco
func ProcessarFerias(conn db.DBConnection, rows [][]string, linhaIni int) error {
	for i := linhaIni; i < len(rows); i++ {
		row := rows[i]
		numColunas := 8

		// Garantir que a linha tenha exatamente 'numColunas' elementos
		if len(row) < numColunas {
			for len(row) < numColunas {
				row = append(row, "") // Completa com strings vazias
			}
		} else if len(row) > numColunas {
			row = row[:numColunas] // Trunca caso tenha mais colunas do que o esperado
		}

		fmt.Printf("Linha %d: %v\n", i, row)

		nomeColaborador := strings.ToLower(strings.TrimSpace(row[0]))

		//fmt.Printf("Linha %d: Nome Colaborador: %s\n", i, nomeColaborador)

		colaborador, err := BuscarColaborador(conn, "efetivo", "nome", nomeColaborador)
		if err != nil {
			if err.Error() == "no rows in result set" {
				// Tenta buscar no colaborador terceirizado
				log.Printf("Colaborador efetivo não encontrado, tentando buscar em terceirizados...")
				colaborador, err = BuscarColaborador(conn, "terceirizado", "nome", nomeColaborador)
				if err != nil {
					log.Printf("Erro ao buscar colaborador terceirizado: %v", err)
				}
			} else {
				// Outro erro, trata como um problema mais crítico
				log.Printf("Erro ao buscar colaborador efetivo: %v", err)
			}
		}

		if colaborador != nil {

			ferias := models.Ferias{
				Colaborador_id: row[0],
				Cargo:          row[1],
				Empresa_id:     row[2],
				Equipe_sigla:   row[3],
				Data_inicio_p1: row[4],
				Data_fim_p1:    row[5],
				Data_inicio_p2: row[6],
				Data_fim_p2:    row[7],
				Lider_E_id:     colaborador.LiderE,
				Lider_T_id:     colaborador.LiderT,
			}

			log.Printf("Cadastrando Férias do Colaborador: %s\n", nomeColaborador)

			if err := CadastraFerias(conn, ferias); err != nil {
				log.Printf("Erro ao cadastrar férias do colaborador: %s. Erro: %v", nomeColaborador, err)
			}

			break

		} else {
			log.Printf("Colaborador Nome: %s, não encontrado no BD\n", nomeColaborador)
		}

	}
	return nil
}

func CadastraFerias(conn db.DBConnection, ferias models.Ferias) error {

	// Cadastra a Féria do Colaborador no Banco de Dados
	/*

		WITH nova_solicitacao AS (
			INSERT INTO public.solicitacao_ferias(id, created_at, status, terceirizado_id)
				VALUES(nextval('solicitacao_ferias_id_seq'::regclass), NOW(), 'aberto', 16) RETURNING id)

			INSERT INTO periodo_ferias (id, data_inicio, data_fim, solicitacao_ferias_id)
				VALUES(nextval('periodo_ferias_id_seq'::regclass), '2024-11-18', '2024-12-03', (SELECT id FROM nova_solicitacao));

	*/
	query := `INSERT INTO public.solicitacao_ferias(created_at, updated_at, status, 
				data_avaliacao_lider, terceirizado_id, lider_colaborador_efetivo_id, lider_colaborador_terceirizado_id,) 
				VALUES (now(), now(), 'aprovado_por_lider', now(), $1, $2, $3)`
	_, err := conn.Exec(context.Background(), query, ferias.Colaborador_id, ferias.Lider_E_id, ferias.Lider_T_id)
	if err != nil {
		return err
	}

	// Busca o ID do Registro de Férias Cadastrado
	var id uint
	err = conn.QueryRow(context.Background(), "SELECT id FROM solicitacao_ferias WHERE sid = $1", sistema.Sid).Scan(&id)
	if err != nil {
		return err
	}

	// Cadastra as Equipes Relacionadas
	for _, equipe := range sistema.EquipesRelacionadas {
		_, err = conn.Exec(context.Background(), "INSERT INTO equipe_sistemas (created_at, updated_at, sistema_id, equipe_id, is_atuador, is_utilizador) VALUES (now(), now(), $1, $2, $3, $4)", id, equipe.EquipeId, equipe.IsAtuador, equipe.IsUtilizador)
		if err != nil {
			return err
		}
	}

	return nil

}
