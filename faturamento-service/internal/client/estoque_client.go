// Este arquivo é o client HTTP para o estoque-service. Ele encapsula as chamadas ao serviço de estoque, lidando com erros e deserialização.

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

// aqui ativo o httpClient com timeout explícito que sem isso, uma requisição para um serviço

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

//aqui Mantice separado do model local para não criar acoplamento entre serviços.
type ProdutoEstoque struct {
	ID        uint   `json:"id"`
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

func estoqueURL() string {
	return os.Getenv("ESTOQUE_BASE_URL")
}

// BuscarProduto ele consulta o estoque-service pelo ID do produto.
func BuscarProduto(produtoID uint) (*ProdutoEstoque, error) {
	url := fmt.Sprintf("%s/api/v1/produtos/%d", estoqueURL(), produtoID)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, errors.New("serviço de estoque indisponível")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	
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

// ja o BaixarSaldo é para baixar o saldo do produto no estoque, ele faz uma requisição PATCH para o endpoint específico de saldo.

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

// essa funçao VerificarSaude checa se o estoque-service está respondendo.
func VerificarSaude() bool {
	resp, err := httpClient.Get(estoqueURL() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}