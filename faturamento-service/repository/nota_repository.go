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

func CriarNota(nota *domain.Nota) error {
	if err := database.DB.Create(nota).Error; err != nil {
		return fmt.Errorf("erro ao criar nota: %w", err)
	}
	return nil
}

func ListarNotas() ([]domain.Nota, error) {
	var notas []domain.Nota

	if err := database.DB.Preload("Itens").Find(&notas).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar notas: %w", err)
	}

	return notas, nil
}

func BuscarNotaPorID(id uint) (*domain.Nota, error) {
	var nota domain.Nota

	err := database.DB.Preload("Itens").First(&nota, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotaNaoEncontrada
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar nota: %w", err)
	}

	return &nota, nil
}

func AtualizarNota(nota *domain.Nota) error {
	if err := database.DB.Save(nota).Error; err != nil {
		return fmt.Errorf("erro ao atualizar nota: %w", err)
	}
	return nil
}