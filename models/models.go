package models

import "time"

type Department struct {
	Id        int          `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Name      string       `json:"name" gorm:"size:200;not null"`
	ParentID  *int         `json:"parent_id,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	Employees []Employee   `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
	Children  []Department `json:"children,omitempty" gorm:"-"`
}

type Employee struct {
	Id           int        `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	DepartmentId int        `json:"departament_id"`
	FullName     string     `json:"full_name" gorm:"size:200;not null"`
	Position     string     `json:"position" gorm:"size:200;not null"`
	HiredAt      *time.Time `json:"hired_at" gorm:"type:date"`
	CreatedAt    time.Time  `json:"created_at"`
}

type DepartmentTree struct {
	Department
	Employees []Employee       `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
	Children  []DepartmentTree `json:"children,omitempty" gorm:"-"`
}

type RequestTree struct {
	Id               int  `json:"id,omitempty"`
	Depth            int  `json:"depth"`
	IncludeEmployees bool `json:"include_employees"`
}

type RequestEmployee struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
	HiredAt  string `json:"hired_at"`
}

type RequestDelete struct {
	Id                     int    `json:"id,omitempty"`
	Mode                   string `json:"mode"`
	ReassignToDepartmentID int    `json:"reassign_to_department_id"`
}
