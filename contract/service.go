package contract

import "github.com/RaFYWStud/LearningSessionBackend/dto"

type Service struct {
	Auth         AuthService
	LearningNote LearningNoteService
	Todo         TodoService
}

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(userID int) (*dto.ProfileResponse, error)
}

type LearningNoteService interface {
	CreateLearningNote(userID int, req dto.CreateLearningNoteRequest) (*dto.CreateLearningNoteResponse, error)
	GetMyLearningNotes(userID int) (*dto.LearningNoteListResponse, error)
	UpdateLearningNote(userID int, noteID int, req dto.UpdateLearningNoteRequest) (*dto.UpdateLearningNoteResponse, error)
	PatchLearningNote(userID int, noteID int, req dto.PatchLearningNoteRequest) (*dto.PatchLearningNoteResponse, error)
	DeleteLearningNote(userID int, noteID int) error
}

type TodoService interface {
	CreateTodo(req dto.CreateTodoRequest) (*dto.CreateTodoResponse, error)
	GetAllTodos() (*dto.TodoListResponse, error)
	DeleteTodo(id int) error
}
