# Korp ERP — Sistema de Emissão de Notas Fiscais

> Projeto técnico desenvolvido para o processo seletivo da **Korp ERP**.

---

## Arquitetura

O sistema é estruturado em **arquitetura de microsserviços**, composta por quatro containers independentes orquestrados via Docker Compose:

```
┌─────────────────────────────────────────────────────────┐ 
│                      Docker Network                     │
│                                                         │
│  ┌──────────────┐    ┌──────────────────┐               │
│  │   Angular    │    │ faturamento-svc  │               │
│  │  (nginx:80)  │───▶│   Gin + GORM     │               │
│  └──────────────┘    │   porta 8082     │               │
│         │            └────────┬─────────┘               │
│         │                     │ HTTP interno             │
│         │            ┌────────▼─────────┐               │
│         └───────────▶│  estoque-svc     │               │
│                      │  Gin + GORM      │               │
│                      │  porta 8081      │               │
│                      └────────┬─────────┘               │
│                               │                         │
│                      ┌────────▼─────────┐               │
│                      │   PostgreSQL 15  │               │
│                      │  estoque_db      │               │
│                      │  faturamento_db  │               │
│                      └──────────────────┘               │
└─────────────────────────────────────────────────────────┘
```

---

## Stack Tecnológica

| Camada | Tecnologia | Justificativa |
|---|---|---|
| Frontend | Angular 17 + Angular Material | Requisito do desafio |
| Backend | Go 1.21 + Gin | Simplicidade, performance e tipagem forte |
| ORM | GORM | Produtividade com PostgreSQL em Go |
| Banco | PostgreSQL 15 | Robusto, suporte nativo a sequences e transações |
| Infra | Docker + Docker Compose | Ambiente reproduzível e isolado |
| Proxy | nginx | Serve o Angular e faz proxy reverso para os serviços |

---

## Requisitos implementados

### Obrigatórios
- [x] Cadastro de Produtos (código, descrição, saldo)
- [x] Cadastro de Notas Fiscais com numeração sequencial automática
- [x] Inclusão de múltiplos produtos em uma nota
- [x] Impressão de notas com botão intuitivo e spinner de processamento
- [x] Bloqueio de impressão para notas já fechadas
- [x] Atualização automática de saldo após impressão
- [x] Arquitetura de microsserviços (estoque + faturamento)
- [x] Tratamento de falhas com feedback ao usuário
- [x] Persistência real em banco de dados PostgreSQL

### Opcionais
- [x] **Tratamento de Concorrência** — `SELECT FOR UPDATE` garante que duas requisições simultâneas nunca deduzam o mesmo saldo
- [x] **Idempotência parcial** — sequence nativa do PostgreSQL garante numeração única mesmo sob concorrência

---

## Como executar

### Pré-requisitos
- Docker Desktop instalado e rodando
- Git

### Passo a passo

```bash
git clone https://github.com/wallyson14/Korp_Teste_Wallyson_Luiz.git
cd Korp_Teste_Wallyson_Luiz

chmod +x init-db.sh

cd frontend
npm install --legacy-peer-deps
cd ..

docker-compose up --build
```

> **Usuários Windows/WSL2:** se o `localhost` não responder nas portas 8081/8082 após o build, execute no PowerShell como Administrador: `netsh interface portproxy reset` e reinicie o terminal WSL.

Após o build (primeira vez ~3 minutos), acesse:

| Serviço | URL |
|---|---|
| Frontend | http://localhost:4200 |
| Estoque API | http://localhost:8081 |
| Faturamento API | http://localhost:8082 |

---

## Endpoints da API

### estoque-service (`localhost:8081`)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/produtos` | Lista todos os produtos |
| `POST` | `/api/v1/produtos` | Cria um produto |
| `PUT` | `/api/v1/produtos/:id` | Atualiza descrição e saldo |
| `DELETE` | `/api/v1/produtos/:id` | Remove produto (soft delete) |
| `PATCH` | `/api/v1/produtos/:id/saldo` | Deduz saldo (uso interno) |
| `GET` | `/health` | Health check |

### faturamento-service (`localhost:8082`)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/notas` | Lista todas as notas |
| `POST` | `/api/v1/notas` | Cria nova nota (número automático) |
| `DELETE` | `/api/v1/notas/:id` | Remove nota aberta |
| `POST` | `/api/v1/notas/:id/itens` | Adiciona produto à nota |
| `DELETE` | `/api/v1/notas/:id/itens/:itemId` | Remove item da nota |
| `POST` | `/api/v1/notas/:id/imprimir` | **Imprime, fecha nota e baixa estoque** |
| `GET` | `/health` | Health check + status do estoque |
