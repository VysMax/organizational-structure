package repository

import (
	"errors"
	"log/slog"

	"github.com/VysMax/organizational-structure/models"
	"gorm.io/gorm"
)

type Repo struct {
	db  *gorm.DB
	log *slog.Logger
}

func New(db *gorm.DB, log *slog.Logger) *Repo {
	return &Repo{
		db:  db,
		log: log,
	}
}

func (r *Repo) CheckExistence(parentID int) (bool, error) {
	err := r.db.Model(&models.Department{}).Where("id = ?", parentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Info("ParentID not found")
			return false, err
		}
		r.log.Error("failed to check parentID existence:", "error", err)
		return false, err
	}

	r.log.Info("parentID exists", "id", parentID)

	return true, nil
}

func (r *Repo) CreateDepartment(department *models.Department) error {
	if err := r.db.Create(department).Error; err != nil {
		r.log.Error("failed to create department", "error", err)
		return err
	}

	r.log.Info("Department created", "id", department.Id, "name", department.Name)

	return nil
}

func (r *Repo) CreateEmployee(employee *models.Employee) error {
	if err := r.db.Create(employee).Error; err != nil {
		r.log.Error("failed to create employee", "error", err)
		return err
	}

	r.log.Info("Employee created", "id", employee.Id, "name", employee.DepartmentId)

	return nil
}

func (r *Repo) GetTree(params *models.RequestTree) (*models.Department, error) {

	var tree models.Department

	query := r.db
	if params.IncludeEmployees {
		query = query.Preload("Employees", func(db *gorm.DB) *gorm.DB {
			return db.Order("full_name ASC")
		})
	}

	err := query.First(&tree, params.Id).Error
	if err != nil {
		r.log.Error("Failed to get head department", "error", err)
		return nil, err
	}

	r.log.Debug("current params", "depth", params.Depth, "head_id", params.Id)

	if params.Depth == 0 {
		tree.Children = []models.Department{}
		return &tree, nil
	}

	var children []models.Department
	err = r.db.Where("parent_id = ?", params.Id).Find(&children).Error
	if err != nil {
		r.log.Error("Failed to get children", "error", err)
		return nil, err
	}

	r.log.Debug("Got", "head_id", params.Id, "children", children)

	for i, child := range children {
		if i == 0 {
			params.Depth--
		}
		params.Id = child.Id

		subtree, err := r.GetTree(params)
		if err != nil {
			r.log.Error("Failed to get subtree", "depth", params.Depth, "error", err)
		}
		tree.Children = append(tree.Children, *subtree)
	}

	return &tree, nil
}
