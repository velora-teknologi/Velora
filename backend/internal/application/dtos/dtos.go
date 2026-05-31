package dtos

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	FullName string `json:"full_name" validate:"omitempty"`
	Avatar   string `json:"avatar" validate:"omitempty,url"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type CreateAgentRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description" validate:"omitempty"`
	Config      interface{} `json:"config" validate:"omitempty"`
}

type UpdateAgentRequest struct {
	Name        string      `json:"name" validate:"omitempty"`
	Description string      `json:"description" validate:"omitempty"`
	Config      interface{} `json:"config" validate:"omitempty"`
}

type AgentResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	Config      interface{} `json:"config"`
}

type CreateWorkflowRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description" validate:"omitempty"`
	AgentID     string      `json:"agent_id" validate:"required"`
	Definition  interface{} `json:"definition" validate:"required"`
}

type UpdateWorkflowRequest struct {
	Name        string      `json:"name" validate:"omitempty"`
	Description string      `json:"description" validate:"omitempty"`
	Definition  interface{} `json:"definition" validate:"omitempty"`
	IsActive    bool        `json:"is_active" validate:"omitempty"`
}

type WorkflowResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Definition  interface{} `json:"definition"`
	IsActive    bool        `json:"is_active"`
}
