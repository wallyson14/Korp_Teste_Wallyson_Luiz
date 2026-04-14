package http

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func Setup(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "estoque",
			"version": os.Getenv("APP_VERSION"),
		})
	})

	api := r.Group("/api/v1")
	{

		p := api.Group("/produtos")
		{
			p.GET("", ListarProdutos)
			p.GET("/:id", BuscarProduto)
			p.POST("", CriarProduto)
			p.PUT("/:id", AtualizarProduto)
			p.DELETE("/:id", DeletarProduto)
			p.PATCH("/:id/saldo", BaixarSaldo)
		}
	}
}