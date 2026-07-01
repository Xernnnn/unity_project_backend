package service

import (
	"github.com/RaFYWStud/LearningSessionBackend/config/pkg/errs"
	"github.com/RaFYWStud/LearningSessionBackend/contract"
	"github.com/RaFYWStud/LearningSessionBackend/database"
	"github.com/RaFYWStud/LearningSessionBackend/dto"
	"gorm.io/gorm"
)

type learningNoteService struct {
	learningNoteRepo contract.LearningNoteRepository
}

func ImplLearningNoteService(learningNoteRepo contract.LearningNoteRepository) contract.LearningNoteService {
	return &learningNoteService{learningNoteRepo: learningNoteRepo}
}
func (s *learningNoteService) CreateLearningNote(userID int, req dto.CreateLearningNoteRequest) (*dto.CreateLearningNoteResponse, error) {
	note := &database.LearningNote{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}
	if err := s.learningNoteRepo.Create(note); err != nil {
		return nil, errs.InternalServerError("failed to create learning note")
	}
	return &dto.CreateLearningNoteResponse{
		Success: true,
		Message: "Learning note created successfully",
		Data:    mapLearningNoteToDTO(*note),
	}, nil
}
func (s *learningNoteService) GetMyLearningNotes(userID int) (*dto.LearningNoteListResponse, error) {
	notes, err := s.learningNoteRepo.FindByUserID(userID)
	if err != nil {
		return nil, errs.InternalServerError("failed to fetch learning notes")
	}
	data := make([]dto.LearningNoteData, len(notes))
	for i, note := range notes {
		data[i] = mapLearningNoteToDTO(note)
	}
	return &dto.LearningNoteListResponse{
		Success: true,
		Data:    data,
	}, nil
}
func mapLearningNoteToDTO(note database.LearningNote) dto.LearningNoteData {
	return dto.LearningNoteData{
		ID:        note.ID,
		UserID:    note.UserID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: note.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *learningNoteService) DeleteLearningNote(userID int, noteID int) error {
	err := s.learningNoteRepo.Delete(noteID, userID)
	if err != nil {
		if err.Error() == "record not found" {
			return errs.NotFound("learning note not found or unauthorized to delete")
		}
		return errs.InternalServerError("failed to delete learning note")
	}

	return nil
}

func (s *learningNoteService) UpdateLearningNote(userID int, noteID int, req dto.UpdateLearningNoteRequest) (*dto.UpdateLearningNoteResponse, error) {
	note, err := s.learningNoteRepo.FindByIDAndUserID(noteID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("learning note not found")
		}
		return nil, errs.InternalServerError("failed to fetch learning note")
	}

	note.Title = req.Title
	note.Content = req.Content

	if err := s.learningNoteRepo.Update(note); err != nil {
		return nil, errs.InternalServerError("failed to update learning note")
	}

	return &dto.UpdateLearningNoteResponse{
		Success: true,
		Message: "Learning note updated successfully",
		Data:    mapLearningNoteToDTO(*note),
	}, nil
}

func (s *learningNoteService) PatchLearningNote(userID int, noteID int, req dto.PatchLearningNoteRequest) (*dto.PatchLearningNoteResponse, error) {
	if req.Title == nil && req.Content == nil {
		return nil, errs.BadRequest("at least one field must be provided")
	}

	note, err := s.learningNoteRepo.FindByIDAndUserID(noteID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("learning note not found")
		}
		return nil, errs.InternalServerError("failed to fetch learning note")
	}

	if req.Title != nil {
		note.Title = *req.Title
	}

	if req.Content != nil {
		note.Content = *req.Content
	}

	if err := s.learningNoteRepo.Update(note); err != nil {
		return nil, errs.InternalServerError("failed to patch learning note")
	}

	return &dto.PatchLearningNoteResponse{
		Success: true,
		Message: "Learning note patched successfully",
		Data:    mapLearningNoteToDTO(*note),
	}, nil
}
