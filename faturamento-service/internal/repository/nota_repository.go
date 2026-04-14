package repository

import (
	"errors"
	"fmt"

	"faturamento-service/internal/domain"
	"faturamento-service/internal/infra/database"

	"gorm.io/gorm"
)

var (
	ErrNotaNaoEncontrada = errors.New("nota não encontrada")
)

func CriarNota(nota *domain.NotaFiscal) error {
	var numero int

	err := database.DB.Raw("SELECT nextval('seq_numero_nota')").Scan(&numero).Error
	if err != nil {
		return fmt.Errorf("erro ao gerar número da nota: %w", err)
	}

	nota.Numero = numero

	return database.DB.Create(nota).Error
}

func ListarNotas() ([]domain.NotaFiscal, error) {
	var notas []domain.NotaFiscal

	err := database.DB.Preload("Itens").Find(&notas).Error
	return notas, err
}

func BuscarNotaPorID(id uint) (*domain.NotaFiscal, error) {
	var nota domain.NotaFiscal

	err := database.DB.Preload("Itens").First(&nota, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotaNaoEncontrada
	}

	return &nota, err
}

func AtualizarNota(nota *domain.NotaFiscal) error {
	return database.DB.Save(nota).Error
}

func AdicionarItem(item *domain.ItemNota) error {
	return database.DB.Create(item).Error
}