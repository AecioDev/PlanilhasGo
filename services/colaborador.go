package services

import (
	"context"
	"fmt"
	"planilhas/db"
)

type Colaborador struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Tipo     string `json:"tipo"`
	EquipeId int    `json:"equipe_id"`
	LiderE   int    `json:"lider_efetivo"`
	LiderT   int    `json:"lider_terceirizado"`
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
		query = fmt.Sprintf("select id, nome, equipe_id from colaborador_efetivo ce where lower(%s) = $1", key)
	} else {
		query = fmt.Sprintf("select id, nome, equipe_id from colaborador_terceirizado ct where lower(%s) = $1", key)
	}

	var colaborador Colaborador

	err := conn.QueryRow(context.Background(), query, value).Scan(&colaborador.ID, &colaborador.Nome, &colaborador.EquipeId)
	if err != nil {
		return nil, err
	}

	colaborador.Tipo = tipo

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
