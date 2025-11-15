package team

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
	"avitotech-pr-reviewer/internal/service/team/mocks"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
)

var ErrUnexpected = errors.New("unexpected error")

func TestService_TeamWithMembers(t *testing.T) { //nolint:funlen
	tests := []struct {
		name          string
		teamName      string
		setupMock     func(m1 *mocks.MockTeamRepository, m2 *mocks.MockUserRepository)
		expectedTeam  *domain.Team
		expectedError error
	}{
		{
			name:     "success - team found with members",
			teamName: "devops",
			setupMock: func(m1 *mocks.MockTeamRepository, m2 *mocks.MockUserRepository) {
				m1.On("GetByName", mock.Anything, "devops").Return(&domain.Team{
					ID:   "team-001",
					Name: "devops",
				}, nil)

				m2.On("ListByTeamID", mock.Anything, "team-001").
					Return([]domain.Member{
						{ID: "u1", Username: "Charlie", IsActive: true},
						{ID: "u2", Username: "Dana", IsActive: true},
					}, nil)
			},
			expectedTeam: &domain.Team{
				ID:   "team-001",
				Name: "devops",
				Members: []domain.Member{
					{ID: "u1", Username: "Charlie", IsActive: true},
					{ID: "u2", Username: "Dana", IsActive: true},
				},
			},
			expectedError: nil,
		},
		{
			name:     "error - team not found",
			teamName: "nonexistent",
			setupMock: func(m1 *mocks.MockTeamRepository, m2 *mocks.MockUserRepository) {
				m1.On("GetByName", mock.Anything, "nonexistent").
					Return(nil, repoErr.ErrTeamNotFound)
			},
			expectedTeam:  nil,
			expectedError: svcErr.ErrTeamNotFound,
		},
		{
			name:     "success - team found with no members",
			teamName: "qa",
			setupMock: func(m1 *mocks.MockTeamRepository, m2 *mocks.MockUserRepository) {
				m1.On("GetByName", mock.Anything, "qa").Return(&domain.Team{
					ID:   "team-002",
					Name: "qa",
				}, nil)

				m2.On("ListByTeamID", mock.Anything, "team-002").
					Return([]domain.Member{}, nil)
			},
			expectedTeam: &domain.Team{
				ID:      "team-002",
				Name:    "qa",
				Members: []domain.Member{},
			},
			expectedError: nil,
		},
		{
			name:     "error - unexpected repository error",
			teamName: "infra",
			setupMock: func(m1 *mocks.MockTeamRepository, m2 *mocks.MockUserRepository) {
				m1.On("GetByName", mock.Anything, "infra").
					Return(nil, ErrUnexpected)
			},
			expectedTeam:  nil,
			expectedError: ErrUnexpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTeamRepo := new(mocks.MockTeamRepository)
			mockUserRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockTeamRepo, mockUserRepo)

			service := &Service{
				lgr:       slog.New(slog.DiscardHandler),
				teamsRepo: mockTeamRepo,
				usersRepo: mockUserRepo,
			}

			result, err := service.TeamWithMembers(context.Background(), tt.teamName)

			if tt.expectedError != nil {
				require.Error(t, err)
				if errors.Is(tt.expectedError, ErrUnexpected) {
					require.Contains(t, err.Error(), tt.expectedError.Error())
				} else {
					require.ErrorIs(t, err, tt.expectedError)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.expectedTeam, result)
			}

			mockTeamRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestService_CreateTeam(t *testing.T) { //nolint:funlen
	tests := []struct {
		name          string
		teamName      string
		members       []domain.Member
		setupMock     func(m *mocks.MockTeamRepository)
		expectedTeam  *domain.Team
		expectedError error
	}{
		{
			name:     "success - team created with members",
			teamName: "backend",
			members: []domain.Member{
				{ID: "u1", Username: "Alice", IsActive: true},
				{ID: "u2", Username: "Bob", IsActive: true},
			},
			setupMock: func(m *mocks.MockTeamRepository) {
				expectedTeam := &domain.Team{
					ID:   "team-123",
					Name: "backend",
					Members: []domain.Member{
						{ID: "u1", Username: "Alice", IsActive: true},
						{ID: "u2", Username: "Bob", IsActive: true},
					},
				}
				m.On("CreateWithMembers", mock.Anything, "backend", mock.Anything).
					Return(expectedTeam, nil)
			},
			expectedTeam: &domain.Team{
				ID:   "team-123",
				Name: "backend",
				Members: []domain.Member{
					{ID: "u1", Username: "Alice", IsActive: true},
					{ID: "u2", Username: "Bob", IsActive: true},
				},
			},
			expectedError: nil,
		},
		{
			name:     "error - team already exists",
			teamName: "existing-team",
			members: []domain.Member{
				{ID: "u1", Username: "Alice", IsActive: true},
			},
			setupMock: func(m *mocks.MockTeamRepository) {
				m.On("CreateWithMembers", mock.Anything, "existing-team", mock.Anything).
					Return(nil, repoErr.ErrTeamExists)
			},
			expectedTeam:  nil,
			expectedError: svcErr.ErrTeamExists,
		},
		{
			name:     "success - team created with empty members",
			teamName: "empty-team",
			members:  []domain.Member{},
			setupMock: func(m *mocks.MockTeamRepository) {
				expectedTeam := &domain.Team{
					ID:      "team-456",
					Name:    "empty-team",
					Members: []domain.Member{},
				}
				m.On("CreateWithMembers", mock.Anything, "empty-team", mock.Anything).
					Return(expectedTeam, nil)
			},
			expectedTeam: &domain.Team{
				ID:      "team-456",
				Name:    "empty-team",
				Members: []domain.Member{},
			},
			expectedError: nil,
		},
		{
			name:     "success - team with inactive members",
			teamName: "mixed-team",
			members: []domain.Member{
				{ID: "u1", Username: "Alice", IsActive: true},
				{ID: "u2", Username: "Bob", IsActive: false},
			},
			setupMock: func(m *mocks.MockTeamRepository) {
				expectedTeam := &domain.Team{
					ID:   "team-789",
					Name: "mixed-team",
					Members: []domain.Member{
						{ID: "u1", Username: "Alice", IsActive: true},
						{ID: "u2", Username: "Bob", IsActive: false},
					},
				}
				m.On("CreateWithMembers", mock.Anything, "mixed-team", mock.Anything).
					Return(expectedTeam, nil)
			},
			expectedTeam: &domain.Team{
				ID:   "team-789",
				Name: "mixed-team",
				Members: []domain.Member{
					{ID: "u1", Username: "Alice", IsActive: true},
					{ID: "u2", Username: "Bob", IsActive: false},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockTeamRepository(t)
			tt.setupMock(mockRepo)

			service := &Service{
				lgr:       slog.New(slog.DiscardHandler),
				teamsRepo: mockRepo,
			}

			result, err := service.CreateTeam(context.Background(), tt.teamName, tt.members)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.expectedTeam, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
