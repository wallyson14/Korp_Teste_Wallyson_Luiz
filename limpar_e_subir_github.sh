#!/bin/bash

echo "=== LIMPANDO SISTEMA ==="
cd "/mnt/c/Users/WALLISON/Downloads/Korp_Teste_Wallyson_Luiz"
docker compose down -v
docker system prune -f
rm -rf frontend/node_modules

echo -e "\n=== VERIFICANDO GIT ==="
git remote -v

echo -e "\n=== ADICIONANDO ARQUIVOS ==="
git add .

echo -e "\n=== COMMIT ==="
git commit -m "feat: Sistema de Notas Fiscais Korp ERP"

echo -e "\n=== ENVIANDO PARA GITHUB ==="
git push -u origin main

echo -e "\n=== SUBINDO SISTEMA ==="
docker compose up -d

echo -e "\n=== VERIFICANDO ==="
sleep 10
docker ps

echo -e "\n✅ Processo concluído!"
