package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// httpClient com timeout explícito — sem isso, uma requisição para um serviço
// lento pode segurar uma goroutine indefinidamente, esgotando o pool de workers.
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// ProdutoEstoque é o DTO de resposta do estoque-service.
// Mantemos separado do model local para não criar acoplamento entre serviços.
type ProdutoEstoque struct {
	ID        uint   `json:"id"`
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

func estoqueURL() string {
	return os.Getenv("ESTOQUE_BASE_URL")
}

// BuscarProduto consulta o estoque-service pelo ID do produto.
// BUG-03 corrigido: rota usa /api/v1/produtos, alinhada com o estoque-service.
func BuscarProduto(produtoID uint) (*ProdutoEstoque, error) {
	url := fmt.Sprintf("%s/api/v1/produtos/%d", estoqueURL(), produtoID)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, errors.New("serviço de estoque indisponível")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok, prossegue
	case http.StatusNotFound:
		return nil, errors.New("produto não encontrado no serviço de estoque")
	default:
		return nil, fmt.Errorf("resposta inesperada do estoque-service: HTTP %d", resp.StatusCode)
	}

	var produto ProdutoEstoque
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		return nil, errors.New("erro ao deserializar resposta do estoque-service")
	}
	return &produto, nil
}

// BaixarSaldo deduz quantidade do saldo de um produto via estoque-service.
// BUG-03 corrigido: rota /api/v1/produtos/:id/saldo.
func BaixarSaldo(produtoID uint, quantidade int) error {
	url := fmt.Sprintf("%s/api/v1/produtos/%d/saldo", estoqueURL(), produtoID)

	payload, _ := json.Marshal(map[string]int{"quantidade": quantidade})

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao montar request para estoque-service: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.New("serviço de estoque indisponível ao tentar baixar saldo")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnprocessableEntity:
		return errors.New("saldo insuficiente no estoque")
	case http.StatusNotFound:
		return fmt.Errorf("produto id=%d não encontrado no estoque", produtoID)
	default:
		return fmt.Errorf("erro inesperado ao baixar saldo: HTTP %d", resp.StatusCode)
	}
}

// VerificarSaude checa se o estoque-service está respondendo.
func VerificarSaude() bool {
	resp, err := httpClient.Get(estoqueURL() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}