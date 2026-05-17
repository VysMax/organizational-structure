package usecase_test

import (
	"errors"
	"testing"

	"github.com/VysMax/organizational-structure/models"
	"github.com/VysMax/organizational-structure/usecase"
	"github.com/VysMax/organizational-structure/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ErrNoName = errors.New("missing name")
)

func TestOrganization_Create(t *testing.T) {
	parentID := 1

	tests := []struct {
		scenario  string
		name      string
		parentID  *int
		setupMock func(*mocks.MockOrganizationRepository)
		wantErr   bool
	}{
		{
			scenario: "success",
			name:     "HR",
			parentID: nil,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CreateDepartment(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			scenario: "success with parent",
			name:     "subHR",
			parentID: &parentID,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CheckExistence(parentID).Return(true, nil)
				repo.EXPECT().CreateDepartment(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			scenario: "no name",
			name:     "",
			parentID: nil,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
			},
			wantErr: true,
		},
		{
			scenario: "no such parent",
			name:     "subHR",
			parentID: &parentID,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().CheckExistence(parentID).Return(false, nil)
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
		defer ctrl.Finish()

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
