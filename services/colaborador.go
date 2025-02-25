package services

import (
	"context"
	"fmt"
	"planilhas/db"
)

type Colaborador struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	EquipeId int    `json:"equipe_id"`
	Tipo     string `json:"tipo"`
	LiderE   *int   `json:"lider_efetivo"`
	LiderT   *int   `json:"lider_terceirizado"`
	Gestor   *int   `json:"gestor"`
}

type Equipe struct {
	ID                int
	Sigla             string
	LiderEfetivo      int
	LiderTerceirizado int
}

func BuscarColaborador(conn db.DBConnection, tipo string, key string, value string) (*Colaborador, error) {
	var query string

	if tipo == "efetivo" {
		query = fmt.Sprintf(`
            SELECT ce.id, ce.nome, ce.equipe_id, 'efetivo' as tipo, u.colaborador_efetivo_id,
				COALESCE(e.lider_colaborador_efetivo_id, 0) as Lider_E, 
				COALESCE(e.lider_colaborador_terceirizado_id, 0) as Lider_T
            FROM colaborador_efetivo ce
            JOIN equipe e ON ce.equipe_id = e.id
			JOIN colaborador_efetivo_local_trabalhos u on u.local_trabalho_id = e.local_trabalho_id and u.tipo_relacionamento = 'gestor' 
            WHERE lower(ce.%s) = $1
        `, key)
	} else {
		query = fmt.Sprintf(`
            SELECT ct.id, ct.nome, ct.equipe_id, 'terceirizado' AS tipo, u.colaborador_efetivo_id,
				COALESCE(e.lider_colaborador_efetivo_id, 0) as Lider_E, 
				COALESCE(e.lider_colaborador_terceirizado_id, 0) as Lider_T
			FROM colaborador_terceirizado ct
			JOIN equipe e ON ct.equipe_id = e.id
			JOIN colaborador_efetivo_local_trabalhos u on u.local_trabalho_id = e.local_trabalho_id and u.tipo_relacionamento = 'gestor' 
			WHERE lower(ct.%s) = $1;
        `, key)
	}

	//fmt.Printf("Query: %s\n", query)
	//fmt.Printf("Valor passado: %s\n", value)

	var colaborador Colaborador
	err := conn.QueryRow(context.Background(), query, value).Scan(
		&colaborador.ID,
		&colaborador.Nome,
		&colaborador.EquipeId,
		&colaborador.Tipo,
		&colaborador.Gestor,
		&colaborador.LiderE,
		&colaborador.LiderT,
	)
	if err != nil {
		return nil, err
	}

	return &colaborador, nil
}

/*
func BuscarColaborador(conn db.DBConnection, tipo string, key string, value string) (*Colaborador, error) {
	var query string

	if tipo == "efetivo" {
		query = fmt.Sprintf("select id, nome, equipe_id from colaborador_efetivo ce where lower(%s) = $1", key)
	} else {
		query = fmt.Sprintf("select id, nome, equipe_id from colaborador_terceirizado ct where lower(%s) = $1", key)
	}

	fmt.Printf("Query 1: %s\n", query)       // Log da query
	fmt.Printf("Valor passado: %s\n", value) // Log do valor

	var dados SqlColaborador

	err := conn.QueryRow(context.Background(), query, value).Scan(&dados.ID, &dados.Nome, &dados.EquipeId)
	if err != nil {
		return nil, err
	}

	colaborador := Colaborador{
		ID:       dados.ID,
		Nome:     dados.Nome,
		EquipeId: dados.EquipeId,
		Tipo:     &tipo,
	}

	var equipe Equipe
	query = "select id, sigla, lider_colaborador_efetivo_id, lider_colaborador_terceirizado_id from equipe where id = $1"
	err = conn.QueryRow(context.Background(), query, colaborador.EquipeId).Scan(&equipe.ID, &equipe.Sigla, &equipe.LiderEfetivo, &equipe.LiderTerceirizado)
	if err != nil {
		return nil, err
	}

	colaborador.LiderE = equipe.LiderEfetivo
	colaborador.LiderT = equipe.LiderTerceirizado

	return &colaborador, nil
}
*/
