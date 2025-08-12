package router

import (
	"write_base/internal/delivery/http/controller"

	"github.com/gin-gonic/gin"
)

func RegisterTagRouter(r *gin.Engine, tagHandler *controller.TagHandler) {
	tags := r.Group("/tags")
	{
		tags.POST("/new", tagHandler.CreateTag)
		tags.GET("", tagHandler.ListTags)
		tags.PATCH("/:id/approve", tagHandler.ApproveTag)
		tags.PATCH("/:id/reject", tagHandler.RejectTag)
		tags.DELETE("/:id", tagHandler.DeleteTag)
	}
}