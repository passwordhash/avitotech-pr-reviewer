package pullrequest

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"avitotech-pr-reviewer/internal/domain"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
	"avitotech-pr-reviewer/internal/service/pullrequest/mocks"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:funlen
func TestService_CreatePullRequest(t *testing.T) {
	tests := []struct {
		name          string
		prID          string
		prName        string
		authorID      string
		setupMock     func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository)
		expectedPR    *domain.PullRequest
		expectedError error
	}{
		{
			name:     "success - pull request created with reviewers",
			prID:     "pr-123",
			prName:   "Add new feature",
			authorID: "user-1",
			setupMock: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				um.On("GetByID", mock.Anything, "user-1").
					Return(&domain.User{ID: "user-1", TeamID: "team-1"}, nil)

				tm.On("GetActiveMembersByTeamID", mock.Anything, "team-1").
					Return([]domain.Member{
						{ID: "user-2", Username: "Reviewer1", IsActive: true},
						{ID: "user-3", Username: "Reviewer2", IsActive: true},
					}, nil)

				expectedPR := &domain.PullRequest{
					ID:        "pr-123",
					Name:      "Add new feature",
					AuthorID:  "user-1",
					Status:    domain.PRStatusOpen,
					Reviewers: []string{"user-2", "user-3"},
				}
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.PullRequest")).
					Return(expectedPR, nil)
			},
			expectedPR: &domain.PullRequest{
				ID:        "pr-123",
				Name:      "Add new feature",
				AuthorID:  "user-1",
				Status:    domain.PRStatusOpen,
				Reviewers: []string{"user-2", "user-3"},
			},
			expectedError: nil,
		},
		{
			name:     "error - pull request already exists",
			prID:     "pr-123",
			prName:   "Add new feature",
			authorID: "user-1",
			setupMock: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				um.On("GetByID", mock.Anything, "user-1").
					Return(&domain.User{ID: "user-1", TeamID: "team-1"}, nil)

				tm.On("GetActiveMembersByTeamID", mock.Anything, "team-1").
					Return([]domain.Member{
						{ID: "user-2", Username: "Reviewer1", IsActive: true},
					}, nil)

				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.PullRequest")).
					Return(nil, repoErr.ErrPRExists)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrPRExists,
		},
		{
			name:     "error - author not found",
			prID:     "pr-123",
			prName:   "Add new feature",
			authorID: "user-unknown",
			setupMock: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				um.On("GetByID", mock.Anything, "user-unknown").
					Return(nil, repoErr.ErrUserNotFound)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrRepo := mocks.NewMockPrRepository(t)
			mockUserRepo := mocks.NewMockUserRepository(t)
			mockTeamRepo := mocks.NewMockTeamRepository(t)

			tt.setupMock(mockPrRepo, mockUserRepo, mockTeamRepo)

			lgr := slog.New(slog.DiscardHandler)

			svc := &Service{
				lgr:          lgr,
				prRepo:       mockPrRepo,
				userRepo:     mockUserRepo,
				teamRepo:     mockTeamRepo,
				maxReviewers: 2,
			}

			ctx := context.Background()
			pr, err := svc.CreatePullRequest(ctx, tt.prID, tt.prName, tt.authorID)

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, pr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPR, pr)
			}

			mockPrRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockTeamRepo.AssertExpectations(t)
		})
	}
}
