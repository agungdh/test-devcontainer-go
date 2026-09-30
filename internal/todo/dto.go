package todo

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,max=200"`
	Description string `json:"description" validate:"max=1000"`
	IsCompleted *bool  `json:"is_completed"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=200"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	IsCompleted *bool   `json:"is_completed"`
}

type ListFilter struct {
	Q           string
	IsCompleted *bool
	Limit       int
	Offset      int
}

func (r UpdateTodoRequest) Empty() bool {
	return r.Title == nil && r.Description == nil && r.IsCompleted == nil
}
