package models

import (
	"time"
)

type Ferias struct {
	ID             uint
	Colaborador_id string
	Lider_E_id     int
	Lider_T_id     int
	Cargo          string
	Empresa_id     string
	Equipe_sigla   string
	Data_inicio_p1 string
	Data_fim_p1    string
	Data_inicio_p2 string
	Data_fim_p2    string
}

type SolicitacaoFerias struct {
	ID             uint
	Status         string
	Created_at     time.Time
	Updated_at     time.Time
	TerceirizadoID uint
}

func (SolicitacaoFerias) TableName() string {
	return "solicitacao_ferias"
}

type PeriodoFeriasColaboradores struct {
	ID                        uint
	DataInicio                time.Time
	DataFim                   time.Time
	ColaboradorEfetivoID      int
	ColaboradorTerceirizadoID int
}

func (PeriodoFeriasColaboradores) TableName() string {
	return "periodo_ferias_colaboradores"
}
