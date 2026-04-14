package routes

import (
	"github.com/gin-gonic/gin"
	
	"faturamento-service/internal/handler"
)

func SetupRoutes(r *gin.Engine, h *handler.NotaHandler) {
	api := r.Group("/api/v1")
	{
		api.POST("/notas", h.Create)
		api.POST("/notas/:id/itens", h.AddItem)
		api.GET("/notas/:id", h.GetByID)
		api.POST("/notas/:id/imprimir", h.Imprimir)
	}
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "faturamento",
			"status":  "ok",
			"version": "1.0.0",
		})
	})
}
