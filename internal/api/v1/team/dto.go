package team

import "avitotech-pr-reviewer/internal/domain"

type userReq struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (u *userReq) ToDomain() domain.Member {
	return domain.Member{
		ID:       u.UserID,
		Username: u.Username,
		IsActive: u.IsActive,
	}
}

type addReq struct {
	TeamName string    `json:"team_name" binding:"required"`
	Members  []userReq `json:"members"`
}

func (a *addReq) ToDomainMembers() []domain.Member {
	domainMembers := make([]domain.Member, len(a.Members))
	for i, member := range a.Members {
		domainMembers[i] = member.ToDomain()
	}

	return domainMembers
}

type userResp struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// func fromDomainUser(u domain.User) userResp {
//	return userResp{
//		UserID:   u.ID,
//		Username: u.Username,
//		IsActive: u.IsActive,
//	}
// }

func fromDomainMember(m domain.Member) userResp {
	return userResp{
		UserID:   m.ID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}

type teamDTO struct {
	TeamName string     `json:"team_name"`
	Members  []userResp `json:"members"`
}

func fromDomainTeam(t *domain.Team) teamDTO {
	members := make([]userResp, len(t.Members))
	for i, member := range t.Members {
		members[i] = fromDomainMember(member)
	}

	return teamDTO{
		TeamName: t.Name,
		Members:  members,
	}
}
