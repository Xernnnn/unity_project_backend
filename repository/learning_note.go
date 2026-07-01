package repository

import (
	"github.com/RaFYWStud/LearningSessionBackend/contract"
	"github.com/RaFYWStud/LearningSessionBackend/database"
	"gorm.io/gorm"
)

type learningNoteRepository struct {
	db *gorm.DB
}

func ImplLearningNoteRepository(db *gorm.DB) contract.LearningNoteRepository {
	return &learningNoteRepository{db: db}
}
func (r *learningNoteRepository) Create(note *database.LearningNote) error {
	return r.db.Create(note).Error
}
func (r *learningNoteRepository) FindByUserID(userID int) ([]database.LearningNote, error) {
	var notes []database.LearningNote
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notes).Error
	return notes, err
}

func (r *learningNoteRepository) Delete(id int, userID int) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&database.LearningNote{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *learningNoteRepository) FindByIDAndUserID(id int, userID int) (*database.LearningNote, error) {
	var note database.LearningNote

	err := r.db.Where("id = ? AND user_id = ?", id, userID).
		First(&note).Error

	return &note, err
}

func (r *learningNoteRepository) Update(note *database.LearningNote) error {
	return r.db.Save(note).Error
}
