package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	TenantID  string         `gorm:"index" json:"tenant_id,omitempty"`
	Tenant    *Tenant        `gorm:"foreignKey:TenantID" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FullName  string         `json:"full_name"`
	Avatar    string         `json:"avatar"`
	Role      string         `gorm:"default:user" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}

// Agent represents an AI agent
type Agent struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"index" json:"tenant_id,omitempty"`
	Tenant      *Tenant        `gorm:"foreignKey:TenantID" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"-"`
	Config      string         `gorm:"type:jsonb" json:"config"`
	Status      string         `gorm:"default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (a *Agent) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New().String()
	return nil
}

// Workflow represents a workflow
type Workflow struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"index" json:"tenant_id,omitempty"`
	Tenant      *Tenant        `gorm:"foreignKey:TenantID" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	AgentID     string         `gorm:"not null;index" json:"agent_id"`
	Agent       *Agent         `gorm:"foreignKey:AgentID" json:"-"`
	Definition  string         `gorm:"type:jsonb" json:"definition"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New().String()
	return nil
}

// Tenant represents an organizational tenant.
type Tenant struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Plan        string         `gorm:"default:'free'" json:"plan"`
	Description string         `json:"description"`
	OwnerID     string         `gorm:"not null;index" json:"owner_id"`
	Owner       *User          `gorm:"foreignKey:OwnerID" json:"-"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New().String()
	return nil
}

// Workspace represents a workspace within a tenant.
type Workspace struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"not null;index" json:"tenant_id"`
	Tenant      *Tenant        `gorm:"foreignKey:TenantID" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (w *Workspace) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New().String()
	return nil
}

// Invitation represents an invitation to join a tenant.
type Invitation struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	TenantID   string         `gorm:"not null;index" json:"tenant_id"`
	Tenant     *Tenant        `gorm:"foreignKey:TenantID" json:"-"`
	Email      string         `gorm:"not null;index" json:"email"`
	Token      string         `gorm:"not null;uniqueIndex" json:"token"`
	Role       string         `gorm:"default:'member'" json:"role"`
	Status     string         `gorm:"default:'pending'" json:"status"`
	InvitedBy  string         `gorm:"not null;index" json:"invited_by"`
	AcceptedAt *time.Time     `json:"accepted_at,omitempty"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (i *Invitation) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New().String()
	return nil
}

// Team represents a team within a tenant.
type Team struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"not null;index" json:"tenant_id"`
	Tenant      *Tenant        `gorm:"foreignKey:TenantID" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New().String()
	return nil
}

// TeamMember represents a member of a team within a tenant.
type TeamMember struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	TeamID    string         `gorm:"not null;index" json:"team_id"`
	Team      *Team          `gorm:"foreignKey:TeamID" json:"-"`
	UserID    string         `gorm:"not null;index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"-"`
	Role      string         `gorm:"default:'member'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *TeamMember) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New().String()
	return nil
}

// WorkspaceMember represents a member of a workspace within a tenant.
type WorkspaceMember struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	WorkspaceID string         `gorm:"not null;index" json:"workspace_id"`
	Workspace   *Workspace     `gorm:"foreignKey:WorkspaceID" json:"-"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"-"`
	Role        string         `gorm:"default:'member'" json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *WorkspaceMember) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New().String()
	return nil
}

// APIKey represents an API key credential for external integrations.
type APIKey struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	KeyID     string         `gorm:"uniqueIndex;not null" json:"key_id"`
	UserID    string         `gorm:"not null;index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	KeyHash   string         `gorm:"not null" json:"-"`
	Scopes    string         `json:"scopes,omitempty"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	ExpiresAt *time.Time     `json:"expires_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	k.ID = uuid.New().String()
	k.KeyID = uuid.New().String()
	return nil
}
