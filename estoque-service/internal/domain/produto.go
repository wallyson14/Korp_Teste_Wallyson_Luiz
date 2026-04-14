package domain

import (
	"time"

	"gorm.io/gorm"
)

// Produto representa um item do catálogo de estoque.
//
// Nota sobre o índice único de 'codigo':
// Utilizamos um partial unique index (uniqueIndex com where) para que o
// soft delete funcione corretamente. Sem isso, um produto deletado (deleted_at != NULL)
// continua ocupando o índice, impedindo o cadastro do mesmo código novamente.
// A tag GORM não suporta partial index diretamente — o índice correto é criado
// via AutoMigrate + trigger ou via migration manual no init-db.sh.
// A constraint no nível do GORM (uniqueIndex) serve como salvaguarda adicional.
type Produto struct {
	ID        uint           `json:"id"         gorm:"primarykey"`
	Codigo    string         `json:"codigo"     gorm:"not null;size:50;index:idx_produto_codigo_ativo,unique,where:deleted_at IS NULL"`
	Descricao string         `json:"descricao"  gorm:"not null;size:200"`
	Saldo     int            `json:"saldo"      gorm:"not null;default:0;check:saldo >= 0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"          gorm:"index"`
}