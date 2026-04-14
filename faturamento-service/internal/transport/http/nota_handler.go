package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"faturamento-service/internal/domain"
	"faturamento-service/internal/repository"

	"github.com/gin-gonic/gin"
)

// =======================
// NOTAS
// =======================

func CriarNota(c *gin.Context) {
	nota := &domain.NotaFiscal{
		Status: domain.StatusAberta,
	}

	if err := repository.CriarNota(nota); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar nota"})
		return
	}

	c.JSON(http.StatusCreated, nota)
}

func ListarNotas(c *gin.Context) {
	notas, err := repository.ListarNotas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao listar notas"})
		return
	}

	c.JSON(http.StatusOK, notas)
}

func BuscarNota(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	nota, err := repository.BuscarNotaPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nota não encontrada"})
		return
	}

	c.JSON(http.StatusOK, nota)
}

// =======================
// ITENS DA NOTA (INTEGRADO COM ESTOQUE)
// =======================

func AdicionarItemNota(c *gin.Context) {
	notaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da nota inválido"})
		return
	}

	var input struct {
		ProdutoID  uint `json:"produto_id"`
		Quantidade int  `json:"quantidade"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}

	if input.Quantidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantidade deve ser maior que zero"})
		return
	}

	// =========================
	// 1. BUSCAR PRODUTO NO ESTOQUE
	// =========================
	estoqueURL := os.Getenv("ESTOQUE_URL")
	if estoqueURL == "" {
		estoqueURL = "http://localhost:8081"
	}

	resp, err := http.Get(estoqueURL + "/api/v1/produtos/" + strconv.Itoa(int(input.ProdutoID)))
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "produto não encontrado no estoque"})
		return
	}
	defer resp.Body.Close()

	var produto struct {
		ID          uint   `json:"id"`
		Codigo      string `json:"codigo"`
		Descricao   string `json:"descricao"`
		Saldo       int    `json:"saldo"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao ler produto"})
		return
	}

	// =========================
	// 2. VALIDAR ESTOQUE
	// =========================
	if produto.Saldo < input.Quantidade {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estoque insuficiente"})
		return
	}

	// =========================
	// 3. BAIXAR ESTOQUE
	// =========================
	body := map[string]int{
		"quantidade": input.Quantidade,
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"PATCH",
		estoqueURL+"/api/v1/produtos/"+strconv.Itoa(int(input.ProdutoID))+"/saldo",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao baixar estoque"})
		return
	}

	// =========================
	// 4. SALVAR ITEM NA NOTA
	// =========================
	item := domain.ItemNota{
		NotaFiscalID: uint(notaID),
		ProdutoID:    input.ProdutoID,
		CodigoProduto: produto.Codigo,
		Descricao:     produto.Descricao,
		Quantidade:    input.Quantidade,
	}

	if err := repository.AdicionarItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao adicionar item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}