package controller

import (
	"context"
	"net/http"

	"write_base/internal/delivery/http/controller/dto"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type AIController struct {
   Usecase domain.IAIUsecase
}

func NewAIController(usecase domain.IAIUsecase) *AIController {
   return &AIController{Usecase: usecase}
}

// POST /ai/suggest
func (c *AIController) Suggest(ctx *gin.Context) {
   var reqDTO dto.SuggestionRequestDTO
   if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
      ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
      return
   }

   domainReq := &domain.SuggestionRequest{Prompt: reqDTO.Prompt}
   resp, err := c.Usecase.GetSuggestions(context.Background(), domainReq)
   if err != nil {
      ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
   }

   ctx.JSON(http.StatusOK, dto.SuggestionResponseDTO{
      Suggestions:  resp.Suggestions,
      Improvements: resp.Improvements,
   })
}
// POST /ai/generate-content
func (c *AIController) GenerateContent(ctx *gin.Context) {
   var reqDTO dto.GenerateContentRequestDTO
   if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
       ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
       return
   }

   domainReq := &domain.GenerateContentRequest{Prompt: reqDTO.Prompt}
   resp, err := c.Usecase.GenerateContent(context.Background(), domainReq)
   if err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   ctx.JSON(http.StatusOK, dto.GenerateContentResponseDTO{
       Content: resp.Content,
   })
}