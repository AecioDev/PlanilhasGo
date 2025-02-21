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
func ProcessarSistemas(conn db.DBConnection, rows [][]string, linhaIni int) error {
	for i := linhaIni; i < len(rows); i++ {
		row := rows[i]

		/*
			A - SID
			B - Servidor
			C - Local
			D - Nome da Aplicação
			E - Descrição
			F - DNS
			G - Consumer
			H - Email
		*/

		s := models.Sistema{}

		s.Sid = row[0]
		s.Servidor = row[1]
		s.Local = row[2]
		s.Nome = row[3]
		s.Descricao = row[4]
		s.LinkDnsSistema = row[5]

		// Exemplo dos dados q estão na planilha
		// cargadfeservice cargadfenaoencontradoservice cargadfenaoprocessadoservice cargatriagemservice cons-cte-referenciado cons-mdfe-referenciado cons-nfe-referenciado dfe-gerencial efronteiras
		//consumer := strings.Split(row[6], " ")
		//fmt.Printf("Linha %d: Consumer: %v\n", i, consumer)

		// Busca um sistema pela Descrição para Adicionar como Sistema Relacionado

		// Dados do Lider na Planilha pra obter a Equipe Atuante
		emailLider := strings.ToLower(strings.TrimSpace(row[7]))

		fmt.Printf("Linha %d: E-mail Lider: %s\n", i, emailLider)

		colaborador, err := BuscarColaborador(conn, "efetivo", "email", emailLider)
		if err != nil {
			if err.Error() == "no rows in result set" {
				// Tenta buscar no colaborador terceirizado
				log.Printf("Colaborador efetivo não encontrado, tentando buscar em terceirizados...")
				colaborador, err = BuscarColaborador(conn, "terceirizado", "email", emailLider)
				if err != nil {
					log.Printf("Erro ao buscar colaborador terceirizado: %v", err)
				}
			} else {
				// Outro erro, trata como um problema mais crítico
				log.Printf("Erro ao buscar colaborador efetivo: %v", err)
			}
		}

		if colaborador != nil {
			equipeRelacionada := models.EquipeRelacionada{
				IsAtuador:    true,
				IsUtilizador: false,
				EquipeId:     uint(colaborador.EquipeId),
			}
			s.EquipesRelacionadas = append(s.EquipesRelacionadas, equipeRelacionada)
		}

		fmt.Printf("Cadastrando Sistema SID: %s\n", s.Sid)

		// Cadastra o Sistema no Banco de Dados
		CadastraSistema(conn, s)
	}
	return nil
}

func CadastraSistema(conn db.DBConnection, sistema models.Sistema) error {

	// Cadastra o Sistema no Banco de Dados
	query := "INSERT INTO sistema (created_at, updated_at, deleted_at, sid, nome, descricao, local, servidor, link_principal, link_dns_sistema) VALUES (now(), now(), NULL, $1, $2, $3, $4, $5, $6, $7) RETURNING id"
	_, err := conn.Exec(context.Background(), query, sistema.Sid, sistema.Nome, sistema.Descricao, sistema.Local, sistema.Servidor, sistema.LinkPrincipal, sistema.LinkDnsSistema)
	if err != nil {
		return err
	}

	// Busca o ID do Sistema Cadastrado
	var id uint
	err = conn.QueryRow(context.Background(), "SELECT id FROM sistema WHERE sid = $1", sistema.Sid).Scan(&id)
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
