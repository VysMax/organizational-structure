package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/VysMax/organizational-structure/models"
	"github.com/VysMax/organizational-structure/usecase"
	"github.com/VysMax/organizational-structure/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ErrNoName = errors.New("missing name")
)

func TestDepartment_Create(t *testing.T) {
	parentID := 1

	tests := []struct {
		name      string
		parentID  *int
		setupMock func(*mocks.MockOrganizationRepository)
		wantErr   bool
	}{
		{
			name:     "HR",
			parentID: nil,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CreateDepartment(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "subHR",
			parentID: &parentID,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CheckExistence(parentID).Return(true, nil)
				repo.EXPECT().CreateDepartment(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "",
			parentID: nil,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
			},
			wantErr: true,
		},
		{
			name:     "subHR",
			parentID: &parentID,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CheckExistence(parentID).Return(false, errors.New("no such parent"))
			},
			wantErr: true,
		},
	}

	department := &models.Department{
		Name:     "",
		ParentID: nil,
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)

		mockedRepository := mocks.NewMockOrganizationRepository(ctrl)
		tt.setupMock(mockedRepository)

		service := usecase.New(mockedRepository, nil)

		department.Name = tt.name
		department.ParentID = tt.parentID

		err := service.CreateDepartment(department)

		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}

func TestEmployee_Create(t *testing.T) {
	tests := []struct {
		fullName  string
		position  string
		hiredAt   time.Time
		setupMock func(*mocks.MockOrganizationRepository)
		wantErr   bool
	}{
		{
			fullName: "Andrey Ivanov",
			position: "HR manager",
			hiredAt:  time.Now(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CheckExistence(gomock.Any()).Return(true, nil)
				repo.EXPECT().CreateEmployee(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			fullName: "Andrey Ivanov",
			position: "HR manager",
			hiredAt:  time.Now(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CheckExistence(gomock.Any()).Return(true, nil)
				repo.EXPECT().CreateEmployee(gomock.Any()).Return(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			fullName: "Andrey Ivanov",
			position: "HR manager",
			hiredAt:  time.Now(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CheckExistence(gomock.Any()).Return(false, errors.New("no such department"))
			},
			wantErr: true,
		},
		{
			fullName:  "",
			position:  "HR manager",
			hiredAt:   time.Now(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {},
			wantErr:   true,
		},
		{
			fullName:  "Andrey Ivanov",
			position:  "",
			hiredAt:   time.Now(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {},
			wantErr:   true,
		},
	}

	employee := &models.Employee{
		FullName: "",
		Position: "",
		HiredAt:  nil,
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)

		mockedRepository := mocks.NewMockOrganizationRepository(ctrl)
		tt.setupMock(mockedRepository)

		service := usecase.New(mockedRepository, nil)

		employee.FullName = tt.fullName
		employee.Position = tt.position
		employee.HiredAt = &tt.hiredAt

		err := service.CreateEmployee(employee)

		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestGetTree(t *testing.T) {

	response := models.Department{}

	tests := []struct {
		id                int
		depth             int
		include_employees bool
		setupMock         func(*mocks.MockOrganizationRepository)
		wantErr           bool
	}{
		{
			id:                1,
			depth:             2,
			include_employees: false,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().GetTree(gomock.Any()).Return(&response, nil)
			},
			wantErr: false,
		},
		{
			id:                1,
			depth:             6,
			include_employees: false,
			setupMock:         func(repo *mocks.MockOrganizationRepository) {},
			wantErr:           true,
		},
	}

	var req models.RequestTree

	for _, tt := range tests {
		ctrl := gomock.NewController(t)

		mockedRepository := mocks.NewMockOrganizationRepository(ctrl)
		tt.setupMock(mockedRepository)

		service := usecase.New(mockedRepository, nil)

		req.Id = tt.id
		req.Depth = tt.depth
		req.IncludeEmployees = tt.include_employees

		tree, err := service.GetTree(&req)

		if tt.wantErr {
			assert.Error(t, err)
			assert.Nil(t, tree)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		}
	}
}
