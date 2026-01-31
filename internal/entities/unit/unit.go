package unitEntity

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/samber/mo"
)

// Unit represents a business entity.
type Unit struct {
	ID        uuid.UUID `validate:"required"`
	Version   int       `validate:"required,gt=0"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required,gtecsfield=CreatedAt"`
	DeletedAt mo.Option[time.Time]
	UserID    mo.Option[uuid.UUID]
	Name      string `validate:"required"`
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
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("validate unit: %w", err)
	}

	if u.DeletedAt.IsPresent() && u.DeletedAt.MustGet().IsZero() {
		return fmt.Errorf("validate unit deleted at: value is zero")
	}

	if u.UserID.IsPresent() && u.UserID.MustGet() == uuid.Nil {
		return fmt.Errorf("validate unit user id: value is zero")
	}

	return nil
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
