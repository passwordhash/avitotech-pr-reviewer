package team_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var baseURL url.URL

func init() {
	port := os.Getenv("HTTP_PORT")
	baseURL = url.URL{
		Scheme: "http",
		Host:   "localhost:" + port,
		Path:   "/team/",
	}
}

type member struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type team struct {
	TeamName string `json:"team_name"`
	Members  []member `json:"members"`
}

type addTeamResponse struct {
	Team team `json:"team"`
}

func TestAddTeam_OK(t *testing.T) {
	t.Parallel()

	e := httpexpect.Default(t, baseURL.String())

	members := []member{
		{UserID: "1", Username: "alice", IsActive: true},
		{UserID: "2", Username: "bob", IsActive: false},
	}
	teamData := team{
		TeamName: "Dev Team",
		Members:  members,
	}

	var resp addTeamResponse
	e.POST("/add").WithJSON(teamData).Expect().
		Status(201).JSON().Object().ContainsKey("team").Decode(&resp)

	if resp.Team.TeamName != teamData.TeamName {
		t.Errorf("expected team name 'Dev Team', got '%s'", resp.Team.TeamName)
	}
	if len(resp.Team.Members) != len(members) {
		t.Errorf("expected 2 members, got %d", len(resp.Team.Members))
	}
}
