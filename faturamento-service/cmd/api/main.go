package main

import (
	"log"
	"os"

	"faturamento-service/internal/infra/database"
	httpTransport "faturamento-service/internal/transport/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// carrega variáveis do .env (se existir)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env não encontrado, usando variáveis do sistema")
	}

	// modo release
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// conexão com banco
	database.ConnectDB()

	r := gin.New()

	// segurança
	r.SetTrustedProxies(nil)

	// middlewares
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// rotas
	httpTransport.Setup(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("🚀 faturamento-service iniciado na porta %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Falha ao iniciar servidor: %v", err)
	}
}