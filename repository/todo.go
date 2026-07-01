package repository

import (
	"github.com/RaFYWStud/LearningSessionBackend/contract"
	"github.com/RaFYWStud/LearningSessionBackend/database"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func ImplTodoRepository(db *gorm.DB) contract.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *database.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) FindAll() ([]database.Todo, error) {
	var todos []database.Todo
	err := r.db.Order("created_at DESC").Find(&todos).Error
	return todos, err
}

func (r *todoRepository) Delete(id int) error {
	result := r.db.Where("id = ?", id).Delete(&database.Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
