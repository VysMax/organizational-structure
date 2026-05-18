package usecase

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/VysMax/organizational-structure/models"
)

type OrganizationRepository interface {
	CreateDepartment(department *models.Department) error
	CreateEmployee(employee *models.Employee) error
	GetTree(params *models.RequestTree) (*models.Department, error)
	CheckExistence(parentID int) (bool, error)
}

type Usecase struct {
	repo OrganizationRepository
	log  *slog.Logger
}

func New(repo OrganizationRepository, log *slog.Logger) *Usecase {
	if log == nil {
		log = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}

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

func (uc *Usecase) CreateEmployee(employee *models.Employee) error {
	employee.FullName = strings.TrimSpace(employee.FullName)

	if employee.FullName == "" {
		uc.log.Error("Validation failed", "error", "full name cannot be empty")
		return errors.New("full name cannot be empty")
	}

	if len(employee.FullName) > 200 {
		uc.log.Error("Validation failed", "error", "full name must not contain more than 200 characters")
		return errors.New("full name must not contain more than 200 characters")
	}

	employee.Position = strings.TrimSpace(employee.Position)
	if employee.Position == "" {
		uc.log.Error("Validation failed", "error", "position cannot be empty")
		return errors.New("position cannot be empty")
	}

	if len(employee.Position) > 200 {
		uc.log.Error("Validation failed", "error", "position must not contain more than 200 characters")
		return errors.New("position must not contain more than 200 characters")
	}

	uc.log.Info("Validation successful")

	exists, err := uc.repo.CheckExistence(employee.DepartmentId)
	if err != nil {
		return fmt.Errorf("failed to check department ID existence: %w", err)
	}
	if !exists {
		return errors.New("specified department ID does not exist")
	}

	if err := uc.repo.CreateEmployee(employee); err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) GetTree(params *models.RequestTree) (*models.Department, error) {

	if params.Depth < 1 || params.Depth > 5 {
		uc.log.Error("Validation failed", "error", "allowed depth is from 1 to 5")
		return nil, errors.New("allowed depth is from 1 to 5")
	}

	uc.log.Info("Validation successful")

	tree, err := uc.repo.GetTree(params)
	if err != nil {
		return nil, err
	}

	return tree, nil
}
