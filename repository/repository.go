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
			return false, nil
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
