package dto

type CreateLearningNoteRequest struct {
	Title   string `json:"title" binding:"required,min=3"`
	Content string `json:"content" binding:"required,min=5"`
}
type LearningNoteListResponse struct {
	Success bool               `json:"success"`
	Data    []LearningNoteData `json:"data"`
}
type CreateLearningNoteResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    LearningNoteData `json:"data"`
}
type LearningNoteData struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateLearningNoteRequest struct {
	Title   string `json:"title" binding:"required,min=3"`
	Content string `json:"content" binding:"required,min=5"`
}

type PatchLearningNoteRequest struct {
	Title   *string `json:"title" binding:"omitempty,min=3"`
	Content *string `json:"content" binding:"omitempty,min=5"`
}

type UpdateLearningNoteResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    LearningNoteData `json:"data"`
}

type PatchLearningNoteResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    LearningNoteData `json:"data"`
}
