package handler

import (
	"fmt"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"faturamento-service/internal/client"
	"faturamento-service/internal/models"
	"faturamento-service/internal/repository"
)

type NotaHandler struct {
	repo    repository.NotaRepository
	estoque client.EstoqueClient
}

func NewNotaHandler(repo repository.NotaRepository, estoque client.EstoqueClient) *NotaHandler {
	return &NotaHandler{
		repo:    repo,
		estoque: estoque,
	}
}

func (h *NotaHandler) Create(c *gin.Context) {
	var req struct {
		Numero int `json:"numero"`
	}
	
	nota := &models.NotaFiscal{
		Status: "Aberta",
	}
	
	if err := c.ShouldBindJSON(&req); err == nil && req.Numero > 0 {
		nota.Numero = req.Numero
	}
	
	if err := h.repo.Create(nota); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, nota)
}

func (h *NotaHandler) AddItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	var req struct {
		ProdutoID  int `json:"produto_id"`
		Quantidade int `json:"quantidade"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	
	// Verificar se o produto existe no estoque
	produto, err := h.estoque.VerificarProduto(req.ProdutoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("produto não encontrado no estoque: %v", err)})
		return
	}
	
	// Adicionar item à nota
	item := &models.NotaItem{
		NotaFiscalID:  id,
		ProdutoID:     req.ProdutoID,
		Quantidade:    req.Quantidade,
		CodigoProduto: produto["codigo"].(string),
		Descricao:     produto["descricao"].(string),
	}
	
	if err := h.repo.AddItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, item)
}

func (h *NotaHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	nota, err := h.repo.FindByIDWithItems(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota não encontrada"})
		return
	}
	
	c.JSON(http.StatusOK, nota)
}

func (h *NotaHandler) Imprimir(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	nota, err := h.repo.FindByIDWithItems(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota não encontrada"})
		return
	}
	
	if nota.Status != "Aberta" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nota já está fechada"})
		return
	}
	
	// Atualizar estoque para cada item
	for _, item := range nota.Itens {
		if err := h.estoque.DecreaseStock(item.ProdutoID, item.Quantidade); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Erro ao atualizar estoque: %v", err),
			})
			return
		}
	}
	
	// Fechar nota
	nota.Status = "Fechada"
	if err := h.repo.Update(nota); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Nota fiscal fechada com sucesso"})
}
