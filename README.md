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

---

## Decisões técnicas relevantes

### Numeração de notas via sequence PostgreSQL
Em vez de `MAX(numero) + 1`, que tem race condition sob concorrência, utilizamos `nextval('seq_numero_nota')` — atômico por definição no PostgreSQL. Duas requisições simultâneas nunca recebem o mesmo número.

### SELECT FOR UPDATE na baixa de estoque
O endpoint `PATCH /produtos/:id/saldo` executa dentro de uma transação com `SELECT FOR UPDATE`. Isso bloqueia a linha do produto até o commit, garantindo que duas notas disputando o último item do estoque não causem saldo negativo.

### Impressão em duas fases
A impressão de nota executa em duas fases distintas:
1. **Validação**: consulta o saldo atual de **todos** os itens antes de modificar qualquer coisa
2. **Execução**: só baixa os saldos se todos passaram na validação

Isso evita inconsistência parcial — se validar 3 itens e o 2º falhar, nenhum saldo é modificado.

### Soft delete com partial unique index
Produtos deletados usam soft delete (campo `deleted_at`). O índice único do campo `codigo` é um partial index (`WHERE deleted_at IS NULL`), permitindo reutilizar o mesmo código após deleção.

### Interceptor HTTP centralizado no Angular
Erros HTTP de qualquer serviço passam por um único interceptor funcional que mapeia os status codes para mensagens amigáveis em português e exibe via `MatSnackBar`. Elimina duplicação de tratamento de erro nos componentes.

### BehaviorSubject para estado reativo
Os serviços Angular usam `BehaviorSubject` para manter o estado das listas em memória. Componentes que assinam `produtos$` ou `notas$` são atualizados automaticamente após qualquer operação de escrita, sem re-fetch manual.
