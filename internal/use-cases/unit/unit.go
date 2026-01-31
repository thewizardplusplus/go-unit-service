package unitUseCase

import (
	"context"
	"errors"
	"fmt"

	unitEntity "go-unit-service/internal/entities/unit"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

var (
	// ErrBadParams indicates invalid input parameters.
	ErrBadParams = errors.New("bad params")
	// ErrUserIDMismatch indicates the update user does not match the unit owner.
	ErrUserIDMismatch = errors.New("unit user id does not match")
	// ErrUserIDMissing indicates the unit has no user id set.
	ErrUserIDMissing = errors.New("unit user id is missing")
)

// Repository defines persistence operations for units.
type Repository interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]unitEntity.Unit, error)
	GetAll(ctx context.Context, userID uuid.UUID, substring mo.Option[string]) ([]unitEntity.Unit, error)
	Create(ctx context.Context, unit unitEntity.Unit) error
	Update(ctx context.Context, unit unitEntity.Unit) error
}

// UseCase provides business logic for units.
type UseCase struct {
	repo Repository
}

// NewUseCase builds a UseCase.
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// GetByIDs returns units by ids.
func (u *UseCase) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]unitEntity.Unit, error) {
	if len(ids) == 0 {
		return nil, errors.Join(ErrBadParams, fmt.Errorf("ids must not be empty"))
	}

	units, err := u.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get units by ids: %w", err)
	}

	return units, nil
}

// GetAll returns all units, optionally filtered by name substring.
func (u *UseCase) GetAll(ctx context.Context, userID uuid.UUID, substring mo.Option[string]) ([]unitEntity.Unit, error) {
	if substring.IsPresent() && substring.MustGet() == "" {
		return nil, errors.Join(ErrBadParams, fmt.Errorf("substring must not be empty when set"))
	}

	units, err := u.repo.GetAll(ctx, userID, substring)
	if err != nil {
		return nil, fmt.Errorf("get all units: %w", err)
	}

	return units, nil
}

// Create creates a new unit.
func (u *UseCase) Create(ctx context.Context, userID uuid.UUID, name string) (unitEntity.Unit, error) {
	unit := unitEntity.New(mo.Some(userID), name)
	if err := unit.Validate(); err != nil {
		return unitEntity.Unit{}, errors.Join(ErrBadParams, fmt.Errorf("validate unit create: %w", err))
	}

	if err := u.repo.Create(ctx, unit); err != nil {
		return unitEntity.Unit{}, fmt.Errorf("create unit: %w", err)
	}

	return unit, nil
}

// Update updates a unit by id and user id.
func (u *UseCase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name string) (unitEntity.Unit, error) {
	units, err := u.repo.GetByIDs(ctx, []uuid.UUID{id})
	if err != nil {
		return unitEntity.Unit{}, fmt.Errorf("get unit by id for update: %w", err)
	}

	if len(units) != 1 {
		return unitEntity.Unit{}, fmt.Errorf("get unit by id for update: expected 1 unit, got %d", len(units))
	}

	unit := units[0]
	if !unit.UserID.IsPresent() {
		return unitEntity.Unit{}, ErrUserIDMissing
	}

	if unit.UserID.MustGet() != userID {
		return unitEntity.Unit{}, ErrUserIDMismatch
	}

	unit.Touch()
	unit.Name = name

	if err := unit.Validate(); err != nil {
		return unitEntity.Unit{}, errors.Join(ErrBadParams, fmt.Errorf("validate unit update: %w", err))
	}

	if err := u.repo.Update(ctx, unit); err != nil {
		return unitEntity.Unit{}, fmt.Errorf("update unit: %w", err)
	}

	return unit, nil
}

// Delete marks the unit as deleted and saves it.
func (u *UseCase) Delete(ctx context.Context, id uuid.UUID) (unitEntity.Unit, error) {
	units, err := u.repo.GetByIDs(ctx, []uuid.UUID{id})
	if err != nil {
		return unitEntity.Unit{}, fmt.Errorf("get unit by id for delete: %w", err)
	}

	if len(units) != 1 {
		return unitEntity.Unit{}, fmt.Errorf("get unit by id for delete: expected 1 unit, got %d", len(units))
	}

	unit := units[0]
	unit.MarkDeleted()

	if err := u.repo.Update(ctx, unit); err != nil {
		return unitEntity.Unit{}, fmt.Errorf("update unit for delete: %w", err)
	}

	return unit, nil
}
