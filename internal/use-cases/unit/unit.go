package unit

import (
	"context"
	"errors"
	"fmt"
	"strings"

	unitEntity "go-unit-service/internal/entities/unit"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

// Repository defines persistence operations for units.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (unitEntity.Unit, error)
	GetAll(ctx context.Context) ([]unitEntity.Unit, error)
	Create(ctx context.Context, unit unitEntity.Unit) error
	Update(ctx context.Context, unit unitEntity.Unit) error
}

// Service provides business logic for units.
type Service struct {
	repo Repository
}

// NewService builds a Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ErrUserIDMismatch indicates the update user does not match the unit owner.
var ErrUserIDMismatch = errors.New("unit user id does not match")

// ErrUserIDMissing indicates the unit has no user id set.
var ErrUserIDMissing = errors.New("unit user id is missing")

// GetByID returns a unit by id.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (unitEntity.Unit, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAll returns all units without pagination.
func (s *Service) GetAll(ctx context.Context) ([]unitEntity.Unit, error) {
	return s.repo.GetAll(ctx)
}

// GetAllFiltered returns all units filtered by name substring.
func (s *Service) GetAllFiltered(ctx context.Context, substring string) ([]unitEntity.Unit, error) {
	units, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if substring == "" {
		return units, nil
	}

	filtered := make([]unitEntity.Unit, 0, len(units))
	for _, unit := range units {
		if strings.Contains(unit.Name, substring) {
			filtered = append(filtered, unit)
		}
	}

	return filtered, nil
}

// Create creates a new unit.
func (s *Service) Create(ctx context.Context, userID mo.Option[uuid.UUID], name string) (unitEntity.Unit, error) {
	unit := unitEntity.New(userID, name)
	if err := unit.Validate(); err != nil {
		return unitEntity.Unit{}, err
	}

	if err := s.repo.Create(ctx, unit); err != nil {
		return unitEntity.Unit{}, err
	}

	return unit, nil
}

// Update updates a unit by id and user id.
func (s *Service) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name string) (unitEntity.Unit, error) {
	unit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return unitEntity.Unit{}, err
	}

	if !unit.UserID.IsPresent() {
		return unitEntity.Unit{}, ErrUserIDMissing
	}

	if unit.UserID.MustGet() != userID {
		return unitEntity.Unit{}, ErrUserIDMismatch
	}

	unit.Touch()
	unit.Name = name

	if err := unit.Validate(); err != nil {
		return unitEntity.Unit{}, fmt.Errorf("validate unit update: %w", err)
	}

	if err := s.repo.Update(ctx, unit); err != nil {
		return unitEntity.Unit{}, err
	}

	return unit, nil
}

// Delete marks the unit as deleted and saves it.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) (unitEntity.Unit, error) {
	unit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return unitEntity.Unit{}, err
	}

	unit.MarkDeleted()

	if err := s.repo.Update(ctx, unit); err != nil {
		return unitEntity.Unit{}, err
	}

	return unit, nil
}
