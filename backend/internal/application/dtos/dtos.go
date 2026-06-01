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
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
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

type AgentExecutionRequest struct {
	Input interface{} `json:"input" validate:"required"`
}

type AgentExecutionResponse struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
	Subject string `json:"subject"`
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

type CreateTenantRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"omitempty,alphanumdash"`
	Plan        string `json:"plan" validate:"omitempty,oneof=free standard enterprise"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateTenantRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Slug        string `json:"slug" validate:"omitempty,alphanumdash"`
	Plan        string `json:"plan" validate:"omitempty,oneof=free standard enterprise"`
	Description string `json:"description" validate:"omitempty"`
	IsActive    *bool  `json:"is_active" validate:"omitempty"`
}

type TenantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Plan        string `json:"plan"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	IsActive    bool   `json:"is_active"`
}

type TenantMemberResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type CreateWorkspaceRequest struct {
	TenantID    string `json:"tenant_id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateWorkspaceRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	IsActive    *bool  `json:"is_active" validate:"omitempty"`
}

type WorkspaceResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type CreateTeamRequest struct {
	TenantID    string `json:"tenant_id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateTeamRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	IsActive    *bool  `json:"is_active" validate:"omitempty"`
}

type TeamResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type CreateTeamMemberRequest struct {
	TeamID string `json:"team_id" validate:"required,uuid"`
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" validate:"omitempty,oneof=member lead"`
}

type TeamMemberResponse struct {
	ID     string `json:"id"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type CreateWorkspaceMemberRequest struct {
	WorkspaceID string `json:"workspace_id" validate:"required,uuid"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	Role        string `json:"role" validate:"omitempty,oneof=member manager"`
}

type WorkspaceMemberResponse struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	UserID      string `json:"user_id"`
	Role        string `json:"role"`
}

type CreateInvitationRequest struct {
	TenantID  string `json:"tenant_id" validate:"required,uuid"`
	Email     string `json:"email" validate:"required,email"`
	Role      string `json:"role" validate:"omitempty,oneof=member admin"`
	ExpiresIn int    `json:"expires_in" validate:"omitempty,min=3600"`
}

type InvitationResponse struct {
	ID         string  `json:"id"`
	TenantID   string  `json:"tenant_id"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	Status     string  `json:"status"`
	Token      string  `json:"token"`
	InvitedBy  string  `json:"invited_by"`
	ExpiresAt  *string `json:"expires_at,omitempty"`
	AcceptedAt *string `json:"accepted_at,omitempty"`
}

type CreateAPIKeyRequest struct {
	Name      string   `json:"name" validate:"required"`
	Scopes    []string `json:"scopes,omitempty"`
	ExpiresIn int      `json:"expires_in,omitempty"`
}

type APIKeyResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Key       string   `json:"key,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	IsActive  bool     `json:"is_active"`
	ExpiresAt *string  `json:"expires_at,omitempty"`
}

type AcceptInvitationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type AcceptInvitationResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type WorkflowExecutionResponse struct {
	WorkflowID string `json:"workflow_id"`
	AgentID    string `json:"agent_id"`
	Status     string `json:"status"`
	Subject    string `json:"subject"`
}
