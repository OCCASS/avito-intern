package tests

import (
	"testing"

	prDto "github.com/OCCASS/avito-intern/internal/application/pullrequest"
	tDto "github.com/OCCASS/avito-intern/internal/application/team"
	"github.com/OCCASS/avito-intern/internal/domain/pullrequest"
	prPostgres "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/domain/team"
	teamPostgres "github.com/OCCASS/avito-intern/internal/domain/team/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/domain/user"
	userPostgres "github.com/OCCASS/avito-intern/internal/domain/user/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EUser(t *testing.T) {
	CleanDb(db)

	// Repositories
	pullrequestRepository := prPostgres.NewPullRequestPostgresRepository(db)
	teamRepository := teamPostgres.NewTeamPostgresRepository(db)
	userRepository := userPostgres.NewUserPostgresRepository(db)

	// Services
	teamService := team.NewTeamServices(teamRepository, userRepository)
	userService := user.NewUserServices(userRepository, teamRepository, pullrequestRepository)
	pullrequestService := pullrequest.NewPullRequestServices(pullrequestRepository, teamRepository)

	// Create (create users)
	createTeam := tDto.CreateTeamDto{
		Name: "backend",
		Members: []entity.User{
			{Id: "u-1", Name: "Ivan", IsActive: true},
			{Id: "u-2", Name: "Vanya", IsActive: true},
			{Id: "u-3", Name: "Lena", IsActive: true},
			{Id: "u-4", Name: "Nastya", IsActive: false},
		},
	}
	createTeam1 := tDto.CreateTeamDto{
		Name: "frontend",
		Members: []entity.User{
			{Id: "u-5", Name: "Ivan", IsActive: true},
			{Id: "u-6", Name: "Vanya", IsActive: true},
			{Id: "u-7", Name: "Lena", IsActive: true},
			{Id: "u-8", Name: "Nastya", IsActive: false},
		},
	}
	teamService.Add(createTeam)
	teamService.Add(createTeam1)

	// Create pull request
	createDto := prDto.CreatePullRequestDto{
		Id:       "pr-1001",
		Name:     "Fix",
		AuthorId: "u-1",
	}
	createDto1 := prDto.CreatePullRequestDto{
		Id:       "pr-1002",
		Name:     "Fix",
		AuthorId: "u-5",
	}
	pullrequestService.Create(createDto)
	pullrequestService.Create(createDto1)

	// Get pull request where review
	prs, err := userService.GetUserPullRequestsWhereReview("u-2")
	require.NoError(t, err)
	assert.Len(t, prs, 1)
	assert.Equal(t, prs[0].Id, "pr-1001")

	// Get user team
	teamName, err := userService.GetUserTeamName("u-1")
	require.NoError(t, err)
	assert.Equal(t, teamName, "backend")

	teamName1, err := userService.GetUserTeamName("u-5")
	require.NoError(t, err)
	assert.Equal(t, teamName1, "frontend")

	_, err = userService.GetUserTeamName("u-100")
	require.Error(t, err)

	// Set is active
	user, err := userService.SetIsActive("u-4", true)
	require.NoError(t, err)
	assert.True(t, user.IsActive)
}
