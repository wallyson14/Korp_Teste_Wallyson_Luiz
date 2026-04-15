// aqui sao os arquivos, funçoes e rotas do http, onde tem a integraçao com o estoque, e onde tem as validaçoes

package http

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "faturamento",
			"version": os.Getenv("APP_VERSION"),
		})
	})

	api := r.Group("/api/v1")
	{
		notas := api.Group("/notas")
		{
			notas.GET("", ListarNotas)
			notas.GET("/:id", BuscarNota)
			notas.POST("", CriarNota)
			
			notas.POST("/:id/itens", AdicionarItemNota)
		}
	}
}