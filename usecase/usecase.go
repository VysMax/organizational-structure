package usecase

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/VysMax/organizational-structure/models"
)

type OrganizationRepository interface {
	CreateDepartment(department *models.Department) error
	CheckExistence(parentID int) (bool, error)
}

type Usecase struct {
	repo OrganizationRepository
	log  *slog.Logger
}

func New(repo OrganizationRepository, log *slog.Logger) *Usecase {
	return &Usecase{
		repo: repo,
		log:  log,
	}
}

func (uc *Usecase) CreateDepartment(department *models.Department) error {
	name := strings.TrimSpace(department.Name)
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if len(name) > 200 {
		return errors.New("name must not contain more than 200 characters")
	}

	if department.ParentID != nil {
		exists, err := uc.repo.CheckExistence(*department.ParentID)
		if err != nil {
			return fmt.Errorf("failed to check parent ID existence: %w", err)
		}
		if !exists {
			return errors.New("specified parent ID does not exist")
		}
	}

	if err := uc.repo.CreateDepartment(department); err != nil {
		return fmt.Errorf("failed to create department: %w", err)
	}

	return nil
}
