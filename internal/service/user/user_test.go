package user

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"avitotech-pr-reviewer/internal/domain"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
	usermocks "avitotech-pr-reviewer/internal/service/user/mocks"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
)

var errUnexpected = errors.New("unexpected error")

func TestService_SetIsActive(t *testing.T) { //nolint:funlen
	tests := []struct {
		name          string
		userID        string
		isActive      bool
		setupMocks    func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository)
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:     "success - status updated and team resolved",
			userID:   "u1",
			isActive: true,
			setupMocks: func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository) {
				u.On("SetIsActive", mock.Anything, "u1", true).Return(&domain.User{
					ID:       "u1",
					Username: "Alice",
					IsActive: true,
					TeamID:   "t1",
				}, nil)
				tr.On("GetByID", mock.Anything, "t1").Return(&domain.Team{
					ID:   "t1",
					Name: "devops",
				}, nil)
			},
			expectedUser: &domain.User{ID: "u1", Username: "Alice", IsActive: true, TeamID: "t1", TeamName: "devops"},
		},
		{
			name:     "error - user not found",
			userID:   "nope",
			isActive: false,
			setupMocks: func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository) {
				u.On("SetIsActive", mock.Anything, "nope", false).Return((*domain.User)(nil), repoErr.ErrUserNotFound)
			},
			expectedUser:  nil,
			expectedError: svcErr.ErrUserNotFound,
		},
		{
			name:     "error - unexpected from user repo",
			userID:   "u2",
			isActive: true,
			setupMocks: func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository) {
				u.On("SetIsActive", mock.Anything, "u2", true).Return((*domain.User)(nil), errUnexpected)
			},
			expectedUser:  nil,
			expectedError: errUnexpected,
		},
		{
			name:     "error - team not found",
			userID:   "u3",
			isActive: true,
			setupMocks: func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository) {
				u.On("SetIsActive", mock.Anything, "u3", true).Return(&domain.User{
					ID:       "u3",
					Username: "Bob",
					IsActive: true,
					TeamID:   "t404",
				}, nil)
				tr.On("GetByID", mock.Anything, "t404").Return((*domain.Team)(nil), repoErr.ErrTeamNotFound)
			},
			expectedUser:  nil,
			expectedError: svcErr.ErrTeamNotFound,
		},
		{
			name:     "error - unexpected from team repo",
			userID:   "u4",
			isActive: true,
			setupMocks: func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository) {
				u.On("SetIsActive", mock.Anything, "u4", true).Return(&domain.User{
					ID:       "u4",
					Username: "Eve",
					IsActive: true,
					TeamID:   "t1",
				}, nil)
				tr.On("GetByID", mock.Anything, "t1").Return((*domain.Team)(nil), errUnexpected)
			},
			expectedUser:  nil,
			expectedError: errUnexpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := usermocks.NewMockUserRepository(t)
			tr := usermocks.NewMockTeamRepository(t)
			tt.setupMocks(ur, tr)

			svc := &Service{
				lgr:        slog.New(slog.DiscardHandler),
				userRepo:   ur,
				teamRepo:   tr,
				adminToken: "secret",
			}

			got, err := svc.SetIsActive(context.Background(), tt.userID, tt.isActive)

			if tt.expectedError != nil {
				require.Error(t, err)
				if errors.Is(tt.expectedError, errUnexpected) {
					require.Contains(t, err.Error(), tt.expectedError.Error())
				} else {
					require.ErrorIs(t, err, tt.expectedError)
				}
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, tt.expectedUser, got)
			}
		})
	}
}

func TestService_VerifyAdminAccess(t *testing.T) {
	svc := &Service{
		lgr:        slog.New(slog.DiscardHandler),
		adminToken: "tok123",
	}

	ok, err := svc.VerifyAdminAccess(context.Background(), "tok123")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = svc.VerifyAdminAccess(context.Background(), "wrong")
	require.NoError(t, err)
	assert.False(t, ok)
}
