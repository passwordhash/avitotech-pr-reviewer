package pullrequest

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"avitotech-pr-reviewer/internal/domain"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
	"avitotech-pr-reviewer/internal/service/pullrequest/mocks"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestService_SetMerged(t *testing.T) {
	createdTime := time.Now().Add(-1 * time.Hour)
	mergedTime := time.Now()

	tests := []struct {
		name          string
		prID          string
		setupMock     func(m *mocks.MockPrRepository)
		expectedPR    *domain.PullRequest
		expectedError error
	}{
		{
			name: "success - pull request set merged",
			prID: "pr-100",
			setupMock: func(m *mocks.MockPrRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Add cart",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u100", "u101"},
						CreatedAt: createdTime,
						MergedAt:  nil,
					}, nil)

				m.On("SetMerged", mock.Anything, "pr-100", mock.Anything).
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Add cart",
						AuthorID:  "u123",
						Status:    domain.PRStatusMerged,
						Reviewers: []string{"u100", "u101"},
						CreatedAt: createdTime,
						MergedAt:  &mergedTime,
					}, nil)
			},
			expectedPR: &domain.PullRequest{
				ID:        "pr-100",
				Name:      "Add cart",
				AuthorID:  "u123",
				Status:    domain.PRStatusMerged,
				Reviewers: []string{"u100", "u101"},
				CreatedAt: createdTime,
				MergedAt:  &mergedTime,
			},
			expectedError: nil,
		},
		{
			name: "success - pr already merged, do nothing",
			prID: "pr-100",
			setupMock: func(m *mocks.MockPrRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Add cart",
						AuthorID:  "u123",
						Status:    domain.PRStatusMerged,
						Reviewers: []string{"u100", "u101"},
						CreatedAt: createdTime,
						MergedAt:  &mergedTime,
					}, nil)
			},
			expectedPR: &domain.PullRequest{
				ID:        "pr-100",
				Name:      "Add cart",
				AuthorID:  "u123",
				Status:    domain.PRStatusMerged,
				Reviewers: []string{"u100", "u101"},
				CreatedAt: createdTime,
				MergedAt:  &mergedTime,
			},
			expectedError: nil,
		},
		{
			name: "success - but no reviewrs",
			prID: "pr-100",
			setupMock: func(m *mocks.MockPrRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Add cart",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: nil,
						CreatedAt: createdTime,
						MergedAt:  nil,
					}, nil)

				m.On("SetMerged", mock.Anything, "pr-100", mock.Anything).
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Add cart",
						AuthorID:  "u123",
						Status:    domain.PRStatusMerged,
						Reviewers: nil,
						CreatedAt: createdTime,
						MergedAt:  &mergedTime,
					}, nil)
			},
			expectedPR: &domain.PullRequest{
				ID:        "pr-100",
				Name:      "Add cart",
				AuthorID:  "u123",
				Status:    domain.PRStatusMerged,
				Reviewers: nil,
				CreatedAt: createdTime,
				MergedAt:  &mergedTime,
			},
			expectedError: nil,
		},
		{
			name: "error - pr not founded",
			prID: "pr-000",
			setupMock: func(m *mocks.MockPrRepository) {
				m.On("GetByID", mock.Anything, "pr-000").
					Return(nil, repoErr.ErrPRNotFound)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrPRNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrRepo := mocks.NewMockPrRepository(t)

			tt.setupMock(mockPrRepo)

			lgr := slog.New(slog.DiscardHandler)

			svc := &Service{
				lgr:    lgr,
				prRepo: mockPrRepo,
			}

			ctx := context.Background()
			pr, err := svc.SetMerged(ctx, tt.prID)

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, pr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPR, pr)
			}

			mockPrRepo.AssertExpectations(t)
		})
	}
}

func TestService_ReassignReviewer(t *testing.T) {
	tests := []struct {
		name               string
		prID               string
		oldReviewID        string
		mockSetup          func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository)
		expectedPR         *domain.PullRequest
		expectedReplacedBy []string
		expectedError      error
	}{
		{
			name:        "success - команда из трех, только один был ревьювером",
			prID:        "pr-100",
			oldReviewID: "u100",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u100"},
					}, nil)

				um.On("GetByID", mock.Anything, "u123").
					Return(&domain.User{ID: "u123", TeamID: "team-1"}, nil)

				tm.On("GetActiveMembersByTeamID", mock.Anything, "team-1").
					Return([]domain.Member{
						{ID: "u100", Username: "Reviewer1", IsActive: true},
						{ID: "u101", Username: "Reviewer2", IsActive: true},
						{ID: "u123", Username: "Author", IsActive: true},
					}, nil)

				m.On("UpdateReviewer", mock.Anything, "pr-100", "u100", "u101").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u101"},
					}, nil)
			},
			expectedPR: &domain.PullRequest{
				ID:        "pr-100",
				Name:      "Improve UI",
				AuthorID:  "u123",
				Status:    domain.PRStatusOpen,
				Reviewers: []string{"u101"},
			},
			expectedReplacedBy: []string{"u101"},
			expectedError:      nil,
		},
		{
			name:        "success - команда из четырех, 2 бывших ревьювера и автор, 1 кандидат",
			prID:        "pr-100",
			oldReviewID: "u100",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u100", "u101"},
					}, nil)

				um.On("GetByID", mock.Anything, "u123").
					Return(&domain.User{ID: "u123", TeamID: "team-1"}, nil)

				tm.On("GetActiveMembersByTeamID", mock.Anything, "team-1").
					Return([]domain.Member{
						{ID: "u100", Username: "Reviewer1", IsActive: true},
						{ID: "u101", Username: "Reviewer2", IsActive: true},
						{ID: "u102", Username: "Reviewer3", IsActive: true},
						{ID: "u123", Username: "Author", IsActive: true},
					}, nil)

				m.On("UpdateReviewer", mock.Anything, "pr-100", "u100", "u102").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u101", "u102"},
					}, nil)
			},
			expectedPR: &domain.PullRequest{
				ID:        "pr-100",
				Name:      "Improve UI",
				AuthorID:  "u123",
				Status:    domain.PRStatusOpen,
				Reviewers: []string{"u101", "u102"},
			},
			expectedReplacedBy: []string{"u102"},
			expectedError:      nil,
		},
		{
			name:        "success - команда из пятерых, 2 бывших ревьювера и автор, 2 кандидата",
			prID:        "pr-100",
			oldReviewID: "u100",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u100", "u101"},
					}, nil)

				um.On("GetByID", mock.Anything, "u123").
					Return(&domain.User{ID: "u123", TeamID: "team-1"}, nil)

				tm.On("GetActiveMembersByTeamID", mock.Anything, "team-1").
					Return([]domain.Member{
						{ID: "u100", Username: "Reviewer1", IsActive: true},
						{ID: "u101", Username: "Reviewer2", IsActive: true},
						{ID: "u102", Username: "Reviewer3", IsActive: true},
						{ID: "u103", Username: "Reviewer4", IsActive: true},
						{ID: "u123", Username: "Author", IsActive: true},
					}, nil)

				m.On("UpdateReviewer", mock.Anything, "pr-100", "u100", mock.Anything).
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u101", "u102"}, // или "u101", newID — не критично для теста
					}, nil)
			},
			expectedPR: &domain.PullRequest{
				ID:        "pr-100",
				Name:      "Improve UI",
				AuthorID:  "u123",
				Status:    domain.PRStatusOpen,
				Reviewers: []string{"u101", "u102"},
			},
			expectedReplacedBy: []string{"u102", "u103"},
			expectedError:      nil,
		},
		{
			name:        "error - нет кандидатов (в команде только два активных участника, один из которых - автор)",
			prID:        "pr-100",
			oldReviewID: "u100",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u100"},
					}, nil)

				um.On("GetByID", mock.Anything, "u123").
					Return(&domain.User{ID: "u123", TeamID: "team-1"}, nil)

				tm.On("GetActiveMembersByTeamID", mock.Anything, "team-1").
					Return([]domain.Member{
						{ID: "u100", Username: "Reviewer1", IsActive: true},
						{ID: "u123", Username: "Author", IsActive: true},
					}, nil)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrPRNoCandidates,
		},
		{
			name:        "error - нет кандидатов (в команде три активных участника, один из которых - автор, остальные уже ревьюверы)",
			prID:        "pr-100",
			oldReviewID: "u100",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u100", "u101"},
					}, nil)

				um.On("GetByID", mock.Anything, "u123").
					Return(&domain.User{ID: "u123", TeamID: "team-1"}, nil)

				tm.On("GetActiveMembersByTeamID", mock.Anything, "team-1").
					Return([]domain.Member{
						{ID: "u100", Username: "Reviewer1", IsActive: true},
						{ID: "u101", Username: "Reviewer2", IsActive: true},
						{ID: "u123", Username: "Author", IsActive: true},
					}, nil)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrPRNoCandidates,
		},
		{
			name:        "error - pr не найден",
			prID:        "pr-999",
			oldReviewID: "u100",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-999").
					Return(nil, repoErr.ErrPRNotFound)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrPRNotFound,
		},
		{
			name:        "error - старвый ревьювер не найден",
			prID:        "pr-100",
			oldReviewID: "u999",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusOpen,
						Reviewers: []string{"u100"},
					}, nil)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrUserNotFound,
		},
		{
			name:        "error - pr уже был смерджин",
			prID:        "pr-100",
			oldReviewID: "u100",
			mockSetup: func(m *mocks.MockPrRepository, um *mocks.MockUserRepository, tm *mocks.MockTeamRepository) {
				m.On("GetByID", mock.Anything, "pr-100").
					Return(&domain.PullRequest{
						ID:        "pr-100",
						Name:      "Improve UI",
						AuthorID:  "u123",
						Status:    domain.PRStatusMerged,
						Reviewers: []string{"u100"},
					}, nil)
			},
			expectedPR:    nil,
			expectedError: svcErr.ErrPRAlreadyMerged,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrRepo := mocks.NewMockPrRepository(t)
			mockUserRepo := mocks.NewMockUserRepository(t)
			mockTeamRepo := mocks.NewMockTeamRepository(t)

			tt.mockSetup(mockPrRepo, mockUserRepo, mockTeamRepo)

			lgr := slog.New(slog.DiscardHandler)

			svc := &Service{
				lgr:      lgr,
				prRepo:   mockPrRepo,
				userRepo: mockUserRepo,
				teamRepo: mockTeamRepo,
			}

			ctx := context.Background()
			pr, newAddedReviewer, err := svc.ReassignReviewer(ctx, tt.prID, tt.oldReviewID)
			_ = newAddedReviewer

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, pr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPR, pr)
			}

			mockPrRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockTeamRepo.AssertExpectations(t)
		})
	}
}
