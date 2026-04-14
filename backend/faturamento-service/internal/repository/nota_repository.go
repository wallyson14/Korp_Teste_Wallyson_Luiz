package repository

import (
	"faturamento-service/internal/models"
	
	"gorm.io/gorm"
)

type NotaRepository interface {
	Create(nota *models.NotaFiscal) error
	FindByIDWithItems(id int) (*models.NotaFiscal, error)
	AddItem(item *models.NotaItem) error
	Update(nota *models.NotaFiscal) error
}

type notaRepository struct {
	db *gorm.DB
}

func NewNotaRepository(db *gorm.DB) NotaRepository {
	return &notaRepository{db: db}
}

func (r *notaRepository) Create(nota *models.NotaFiscal) error {
	return r.db.Create(nota).Error
}

func (r *notaRepository) FindByIDWithItems(id int) (*models.NotaFiscal, error) {
	var nota models.NotaFiscal
	err := r.db.Preload("Itens").First(&nota, id).Error
	return &nota, err
}

func (r *notaRepository) AddItem(item *models.NotaItem) error {
	return r.db.Create(item).Error
}

func (r *notaRepository) Update(nota *models.NotaFiscal) error {
	return r.db.Save(nota).Error
}
