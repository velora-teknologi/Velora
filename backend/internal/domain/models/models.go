package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
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
