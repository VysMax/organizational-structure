package models

import "time"

type Department struct {
	Id        int
	Name      string
	ParentID  int
	CreatedAt time.Time
}

type Employee struct {
	Id           int
	DepartmentId int
	FullName     string
	Position     string
	HiredAt      time.Time
	CreatedAt    time.Time
}
