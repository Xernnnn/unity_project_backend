package service

import "github.com/RaFYWStud/LearningSessionBackend/contract"

func New(repo *contract.Repository) *contract.Service {
	return &contract.Service{
		Auth:         ImplAuthService(repo.Auth),
		LearningNote: ImplLearningNoteService(repo.LearningNote),
		Todo:         ImplTodoService(repo.Todo),
	}
}
