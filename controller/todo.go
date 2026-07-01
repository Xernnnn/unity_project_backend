package controller

import (
	"net/http"
	"strconv"

	"github.com/RaFYWStud/LearningSessionBackend/config/middleware"
	"github.com/RaFYWStud/LearningSessionBackend/config/pkg/errs"
	"github.com/RaFYWStud/LearningSessionBackend/contract"
	"github.com/RaFYWStud/LearningSessionBackend/dto"
	"github.com/gin-gonic/gin"
)

type TodoController struct {
	service contract.TodoService
}

func (tc *TodoController) GetPrefix() string {
	return "/api/todos"
}

func (tc *TodoController) InitService(service *contract.Service) {
	tc.service = service.Todo
}

func (tc *TodoController) InitRoute(app *gin.RouterGroup) {
	auth := app.Group("")
	auth.Use(middleware.Auth()) // INI SATPAMNYA!
	{
		auth.POST("", tc.createTodo)
		auth.GET("", tc.getAllTodos)
		auth.DELETE("/:id", tc.deleteTodo)
	}
}

func (tc *TodoController) createTodo(ctx *gin.Context) {
	var payload dto.CreateTodoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		// Simplifikasi pesan error validasi
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Validasi gagal",
			"errors": gin.H{
				"title": "Title wajib diisi dan minimal 3 karakter",
			},
		})
		return
	}

	response, err := tc.service.CreateTodo(payload)
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (tc *TodoController) getAllTodos(ctx *gin.Context) {
	response, err := tc.service.GetAllTodos()
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (tc *TodoController) deleteTodo(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		HandlerError(ctx, errs.BadRequest("Invalid ID format"))
		return
	}

	err = tc.service.DeleteTodo(id)
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Todo berhasil dihapus",
	})
}
