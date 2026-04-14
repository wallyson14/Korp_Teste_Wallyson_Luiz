package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type EstoqueClient interface {
	VerificarProduto(produtoID int) (map[string]interface{}, error)
	DecreaseStock(produtoID int, quantidade int) error
}

type estoqueClient struct {
	baseURL string
	client  *http.Client
}

func NewEstoqueClient() EstoqueClient {
	baseURL := os.Getenv("ESTOQUE_BASE_URL")
	if baseURL == "" {
		baseURL = "http://korp_estoque:8081"
	}
	
	return &estoqueClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *estoqueClient) VerificarProduto(produtoID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/produtos/%d", c.baseURL, produtoID)
	
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar estoque: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("produto %d não encontrado", produtoID)
	}
	
	var produto map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}
	
	return produto, nil
}

func (c *estoqueClient) DecreaseStock(produtoID int, quantidade int) error {
	url := fmt.Sprintf("%s/api/v1/estoque/decrease", c.baseURL)
	
	body := map[string]interface{}{
		"produto_id": produtoID,
		"quantidade": quantidade,
	}
	
	jsonBody, _ := json.Marshal(body)
	
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("erro ao chamar serviço de estoque: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("estoque retornou erro: status %d", resp.StatusCode)
	}
	
	return nil
}
