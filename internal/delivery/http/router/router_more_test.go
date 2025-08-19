package router

import (
	"testing"
	"write_base/internal/delivery/http/controller"

	"github.com/gin-gonic/gin"
)

func TestRegisterOtherRouters_NoPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Construct controllers with nil usecases; we won't hit handlers
	comment := controller.NewCommentController(nil)
	reaction := controller.NewReactionController(nil)
	follow := controller.NewFollowController(nil)
	report := controller.NewReportController(nil)
	ai := controller.NewAIController(nil)

	RegisterCommentRoutes(r, comment)
	RegisterReactionRoutes(r, reaction)
	RegisterFollowRoutes(r, follow)
	RegisterReportRoutes(r, report)
	RegisterAIRoutes(r, ai)
}
