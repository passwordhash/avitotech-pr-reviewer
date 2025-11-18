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

func TestService_SetIsActive(t *testing.T) {
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
					TeamName: "AI",
				}, nil)
			},
			expectedUser: &domain.User{ID: "u1", Username: "Alice", IsActive: true, TeamID: "t1", TeamName: "AI"},
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
			name:     "error - users team not found",
			userID:   "u1",
			isActive: false,
			setupMocks: func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository) {
				u.On("SetIsActive", mock.Anything, "u1", false).Return((*domain.User)(nil), repoErr.ErrTeamNotFound)
			},
			expectedUser:  nil,
			expectedError: svcErr.ErrTeamNotFound,
		},
		{
			name:     "error - unexpected from repo",
			userID:   "u2",
			isActive: true,
			setupMocks: func(u *usermocks.MockUserRepository, tr *usermocks.MockTeamRepository) {
				u.On("SetIsActive", mock.Anything, "u2", true).Return((*domain.User)(nil), errUnexpected)
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
