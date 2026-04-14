package models

import "time"

type NotaFiscal struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Numero    int       `gorm:"unique" json:"numero"`
	Status    string    `gorm:"default:'Aberta'" json:"status"`
	Itens     []NotaItem `gorm:"foreignKey:NotaFiscalID" json:"itens"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotaItem struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	NotaFiscalID  int       `gorm:"index" json:"nota_fiscal_id"`
	ProdutoID     int       `json:"produto_id"`
	CodigoProduto string    `json:"codigo_produto"`
	Descricao     string    `json:"descricao"`
	Quantidade    int       `json:"quantidade"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
