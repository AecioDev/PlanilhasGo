package models

type Atualizacoes struct {
	Titulo           string `json:"titulo" form:"titulo" binding:"required,min=3,max=200"`
	Descricao        string `json:"descricao" form:"descricao" binding:"required,max=300"`
	Versao_Sistema   string `json:"versao_sistema" binding:"required"`
	Nome_Responsavel string `json:"nome_responsavel" binding:"required"`
	Data_Atualizacao string `json:"data_atualizacao" binding:"required"`
}

type EquipeRelacionada struct {
	EquipeId     uint `json:"equipeId" binding:"gte=1"`
	IsAtuador    bool `json:"atuador" form:"atuador"`
	IsUtilizador bool `json:"utilizador" form:"utilizador"`
}

type EndpointServiceSistema struct {
	Metodo      string `json:"metodo" binding:"required,oneof=GET POST PUT DELETE PATCH"`
	Descricao   string `json:"descricao" binding:"required"`
	UrlEndpoint string `json:"url_endpoint" binding:"required"`
}

type ServiceSistema struct {
	Nome               string `json:"nome" binding:"required,min=3,max=200"`
	Descricao          string `json:"descricao" form:"descricao" binding:"required,max=300"`
	UrlBase            string `json:"url_base" binding:"required,url"`
	Status             bool   `json:"status" binding:"required"`
	MetodoAutenticacao string `json:"metodo_auth" binding:"omitempty"`
	InfoAutenticacao   string `json:"info_auth" binding:"omitempty"`
	RateLimit          string `json:"rate_limit" binding:"omitempty"`
	DocumentacaoApi    string `json:"documentacao_api" binding:"omitempty,url"`

	EndpointsRelacionados []EndpointServiceSistema `json:"endpoints" binding:"dive"`
}

type Sistema struct {
	// Dados Principais
	Sid            string `json:"sid" gorm:"column:sid"`
	Nome           string `json:"nome" gorm:"not null"`
	Descricao      string `json:"descricao" gorm:"not null"`
	Local          string `json:"local" gorm:"not null"`
	Servidor       string `json:"servidor"`
	LinkPrincipal  string `json:"link_principal"`
	LinkDnsSistema string `json:"link_dns_sistema"`

	// Entidades Relacionadas ao Sistema | // relacionamento N-N
	Atualizacoes            []Atualizacoes      `json:"atualizacoes" binding:"omitempty,dive"`
	EquipesRelacionadas     []EquipeRelacionada `json:"equipes_relacionadas" binding:"omitempty,dive"`
	ServicosRelacionados    []ServiceSistema    `json:"servicos_relacionados" binding:"omitempty,dive"`
	SistemasRelacionadosIDs []uint              `json:"sistemasRelacionadosIds" binding:"omitempty"`
	EfetivosIDs             []uint              `json:"efetivosIds" binding:"omitempty"`
	TerceirizadosIDs        []uint              `json:"terceirizadosIds" binding:"omitempty"`
}
