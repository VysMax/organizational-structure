package models

import "time"

type Department struct {
	Id        int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:200;not null"`
	ParentID  *int      `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Employee struct {
	Id           int
	DepartmentId int
	FullName     string
	Position     string
	HiredAt      time.Time
	CreatedAt    time.Time
}
