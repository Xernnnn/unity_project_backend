package controller

import (
	"github.com/RaFYWStud/LearningSessionBackend/config/middleware"
	"github.com/RaFYWStud/LearningSessionBackend/config/pkg/errs"
	"github.com/RaFYWStud/LearningSessionBackend/contract"
	"github.com/RaFYWStud/LearningSessionBackend/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LearningNoteController struct {
	service contract.LearningNoteService
}

func (lc *LearningNoteController) GetPrefix() string {
	return "/learning-notes"
}
func (lc *LearningNoteController) InitService(service *contract.Service) {
	lc.service = service.LearningNote
}
func (lc *LearningNoteController) InitRoute(app *gin.RouterGroup) {
	auth := app.Group("")
	auth.Use(middleware.Auth())
	{
		auth.GET("", lc.getMyLearningNotes)
		auth.POST("", lc.createLearningNote)
		auth.DELETE("/:id", lc.deleteLearningNote)
		auth.PUT("/:id", lc.updateLearningNote)
		auth.PATCH("/:id", lc.patchLearningNote)
	}
}
func (lc *LearningNoteController) getMyLearningNotes(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		HandlerError(ctx, errs.Unauthorized("user not authenticated"))
		return
	}
	response, err := lc.service.GetMyLearningNotes(userID.(int))
	if err != nil {
		HandlerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
func (lc *LearningNoteController) createLearningNote(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		HandlerError(ctx, errs.Unauthorized("user not authenticated"))
		return
	}
	var payload dto.CreateLearningNoteRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		HandlerError(ctx, errs.BadRequest("invalid request payload"))
		return
	}
	response, err := lc.service.CreateLearningNote(userID.(int), payload)
	if err != nil {
		HandlerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (lc *LearningNoteController) deleteLearningNote(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		HandlerError(ctx, errs.Unauthorized("user not authenticated"))
		return
	}

	noteIDParam := ctx.Param("id")
	noteID, err := strconv.Atoi(noteIDParam)
	if err != nil {
		HandlerError(ctx, errs.BadRequest("invalid note id format"))
		return
	}

	err = lc.service.DeleteLearningNote(userID.(int), noteID)
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Learning note deleted successfully",
	})
}

func parseIDParam(ctx *gin.Context) (int, error) {
	return strconv.Atoi(ctx.Param("id"))
}

func (lc *LearningNoteController) updateLearningNote(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		HandlerError(ctx, errs.Unauthorized("user not authenticated"))
		return
	}

	noteID, err := parseIDParam(ctx)
	if err != nil {
		HandlerError(ctx, errs.BadRequest("invalid learning note id"))
		return
	}

	var payload dto.UpdateLearningNoteRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		HandlerError(ctx, errs.BadRequest("invalid request payload"))
		return
	}

	response, err := lc.service.UpdateLearningNote(userID.(int), noteID, payload)
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (lc *LearningNoteController) patchLearningNote(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		HandlerError(ctx, errs.Unauthorized("user not authenticated"))
		return
	}

	noteID, err := parseIDParam(ctx)
	if err != nil {
		HandlerError(ctx, errs.BadRequest("invalid learning note id"))
		return
	}

	var payload dto.PatchLearningNoteRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		HandlerError(ctx, errs.BadRequest("invalid request payload"))
		return
	}

	response, err := lc.service.PatchLearningNote(userID.(int), noteID, payload)
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
