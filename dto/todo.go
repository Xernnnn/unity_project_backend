package dto

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,min=3"`
	Description string `json:"description"`
}

type TodoData struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
	CreatedAt   string `json:"created_at"`
}

type CreateTodoResponse struct {
	Message string   `json:"message"`
	Data    TodoData `json:"data"`
}

type TodoListResponse struct {
	Message string     `json:"message"`
	Data    []TodoData `json:"data"`
}
