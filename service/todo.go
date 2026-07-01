package service

import (
	"github.com/RaFYWStud/LearningSessionBackend/config/pkg/errs"
	"github.com/RaFYWStud/LearningSessionBackend/contract"
	"github.com/RaFYWStud/LearningSessionBackend/database"
	"github.com/RaFYWStud/LearningSessionBackend/dto"
	"gorm.io/gorm"
)

type todoService struct {
	todoRepo contract.TodoRepository
}

func ImplTodoService(todoRepo contract.TodoRepository) contract.TodoService {
	return &todoService{todoRepo: todoRepo}
}

func (s *todoService) CreateTodo(req dto.CreateTodoRequest) (*dto.CreateTodoResponse, error) {
	todo := &database.Todo{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := s.todoRepo.Create(todo); err != nil {
		return nil, errs.InternalServerError("failed to create todo")
	}

	return &dto.CreateTodoResponse{
		Message: "Todo berhasil dibuat",
		Data:    mapTodoToDTO(*todo),
	}, nil
}

func (s *todoService) GetAllTodos() (*dto.TodoListResponse, error) {
	todos, err := s.todoRepo.FindAll()
	if err != nil {
		return nil, errs.InternalServerError("failed to fetch todos")
	}

	data := make([]dto.TodoData, len(todos))
	for i, t := range todos {
		data[i] = mapTodoToDTO(t)
	}

	message := "Daftar todo berhasil diambil"
	if len(todos) == 0 {
		message = "Belum ada todo"
	}

	return &dto.TodoListResponse{
		Message: message,
		Data:    data,
	}, nil
}

func (s *todoService) DeleteTodo(id int) error {
	err := s.todoRepo.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errs.NotFound("Todo tidak ditemukan")
		}
		return errs.InternalServerError("failed to delete todo")
	}
	return nil
}

func mapTodoToDTO(todo database.Todo) dto.TodoData {
	return dto.TodoData{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		IsDone:      todo.IsDone,
		CreatedAt:   todo.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
