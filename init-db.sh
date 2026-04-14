#!/bin/bash
set -e

echo "📦 Criando bancos de dados da aplicação Korp..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    CREATE DATABASE estoque_db;
    CREATE DATABASE faturamento_db;
    GRANT ALL PRIVILEGES ON DATABASE estoque_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE faturamento_db TO $POSTGRES_USER;
EOSQL

echo "✅ Bancos estoque_db e faturamento_db criados com sucesso."
