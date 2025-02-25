package models

import (
	"time"
)

type Ferias struct {
	ID               uint   `json:"id"`
	Colaborador_id   string `json:"colaborador_id"`
	Colaborador_nome string `json:"colaborador_nome"`
	Lider_E_id       int    `json:"lider_e_id"`
	Lider_T_id       int    `json:"lider_t_id"`
	Gestor_id        int    `json:"gestor_id"`
	Cargo            string `json:"cargo"`
	Empresa_id       string `json:"empresa_id"`
	UserEmpresa_id   int    `json:"user_empresa_id"`
	Equipe_sigla     string `json:"equipe_sigla"`
	Data_inicio_p1   string `json:"data_inicio_p1"`
	Data_fim_p1      string `json:"data_fim_p1"`
	Data_inicio_p2   string `json:"data_inicio_p2"`
	Data_fim_p2      string `json:"data_fim_p2"`
}

type SolicitacaoFerias struct {
	ID             uint      `json:"id"`
	Status         string    `json:"status"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
	TerceirizadoID uint      `json:"terceirizado_id"`
}

func (SolicitacaoFerias) TableName() string {
	return "solicitacao_ferias"
}

type PeriodoFeriasColaboradores struct {
	ID                        uint      `json:"id"`
	DataInicio                time.Time `json:"data_inicio"`
	DataFim                   time.Time `json:"data_fim"`
	ColaboradorEfetivoID      int       `json:"colaborador_efetivo_id"`
	ColaboradorTerceirizadoID int       `json:"colaborador_terceirizado_id"`
}

func (PeriodoFeriasColaboradores) TableName() string {
	return "periodo_ferias_colaboradores"
}
