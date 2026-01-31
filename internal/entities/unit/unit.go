package unit

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/samber/mo"
)

// Unit represents a business entity.
type Unit struct {
	ID        uuid.UUID            `validate:"required"`
	Version   int                  `validate:"required,gt=0"`
	CreatedAt time.Time            `validate:"required"`
	UpdatedAt time.Time            `validate:"required,gtecsfield=CreatedAt"`
	DeletedAt mo.Option[time.Time] `validate:"omitempty"`
	UserID    mo.Option[uuid.UUID] `validate:"omitempty"`
	Name      string               `validate:"required"`
}

// New creates a new Unit with required defaults.
func New(userID mo.Option[uuid.UUID], name string) Unit {
	now := time.Now()

	return Unit{
		ID:        uuid.New(),
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		Name:      name,
	}
}

// Validate validates the Unit invariants.
func (u Unit) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// Touch increments the version and updates UpdatedAt.
func (u *Unit) Touch() {
	if u == nil {
		return
	}

	u.Version++
	u.UpdatedAt = time.Now()
}

// MarkDeleted marks the unit as deleted.
func (u *Unit) MarkDeleted() {
	if u == nil {
		return
	}

	now := time.Now()
	u.UpdatedAt = now
	u.DeletedAt = mo.Some(now)
}
