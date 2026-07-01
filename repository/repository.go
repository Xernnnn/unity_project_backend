package repository

import (
	"github.com/RaFYWStud/LearningSessionBackend/contract"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *contract.Repository {
	return &contract.Repository{
		Auth:         ImplAuthRepository(db),
		LearningNote: ImplLearningNoteRepository(db),
		Todo:         ImplTodoRepository(db),
	}
}
