package contract

import (
	"github.com/RaFYWStud/LearningSessionBackend/database"
)

type Repository struct {
	Auth         AuthRepository
	LearningNote LearningNoteRepository
	Todo         TodoRepository 
}

type AuthRepository interface {
	CreateUser(user *database.User) error
	FindByEmail(email string) (*database.User, error)
	FindByID(id int) (*database.User, error)
}

type LearningNoteRepository interface {
	Create(note *database.LearningNote) error
	FindByUserID(userID int) ([]database.LearningNote, error)
	FindByIDAndUserID(id int, userID int) (*database.LearningNote, error)
	Update(note *database.LearningNote) error
	Delete(id int, userID int) error
}

type TodoRepository interface {
	Create(todo *database.Todo) error
	FindAll() ([]database.Todo, error)
	Delete(id int) error
}
