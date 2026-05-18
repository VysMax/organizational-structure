package repository

import (
	"errors"
	"log/slog"

	"github.com/VysMax/organizational-structure/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
			return nil, err
		}
		tree.Children = append(tree.Children, *subtree)
	}

	return &tree, nil
}

func (r *Repo) UpdateParent(department *models.Department) error {

	var result *gorm.DB

	switch department.Name {
	case "":
		result = r.db.Model(&models.Department{}).
			Where("id = ?", department.Id).
			Update("parent_id", department.ParentID).Scan(&department)

	default:
		result = r.db.Model(&models.Department{}).
			Where("id = ?", department.Id).
			Clauses(clause.Returning{}).
			Update("name", department.Name).
			Update("parent_id", department.ParentID).Scan(&department)
	}

	r.log.Debug("model", "department", department)

	if err := result.Error; err != nil {
		r.log.Error("Failed to update parent", "error", err)
		return err
	}

	return nil
}

func (r *Repo) DeleteDepartment(params *models.RequestDelete) error {

	var err error

	switch params.Mode {
	case "cascade":
		err = r.db.Delete(&models.Department{}, params.Id).Error
	case "reassign":
		tx := r.db.Begin()

		if tx.Error != nil {
			r.log.Error("Failed to begin transaction", "error", err)
			return tx.Error
		}

		if err = tx.Model(&models.Employee{}).
			Where("department_id = ?", params.Id).
			Update("department_id", params.ReassignToDepartmentID).Error; err != nil {
			tx.Rollback()
			r.log.Error("Failed to change employees' department", "error", err)
			return err
		}

		var eliminatedDept models.Department
		if err := tx.First(&eliminatedDept, params.Id).Error; err != nil {
			tx.Rollback()
			r.log.Error("Failed to get eliminated department parent ID", "error", err)
			return err
		}

		r.log.Debug("Got", "parentID", eliminatedDept.ParentID)

		if err = tx.Model(&models.Department{}).
			Where("parent_id = ?", params.Id).
			Update("parent_id", eliminatedDept.ParentID).Error; err != nil {
			tx.Rollback()
			r.log.Error("Failed to change employees' department", "error", err)
			return err
		}

		if err = tx.Delete(&models.Department{}, params.Id).Error; err != nil {
			tx.Rollback()
			r.log.Error("Failed to delete department", "error", err)
			return err
		}

		if err = tx.Commit().Error; err != nil {
			r.log.Error("Failed to commit transaction", "error", err)
			return err
		}
	}

	return nil
}
