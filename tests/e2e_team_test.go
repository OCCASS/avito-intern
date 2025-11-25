package tests

import (
	"slices"
	"testing"

	prDto "github.com/OCCASS/avito-intern/internal/application/pullrequest"
	tDto "github.com/OCCASS/avito-intern/internal/application/team"
	"github.com/OCCASS/avito-intern/internal/domain/pullrequest"
	pullrequestPostgres "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/domain/team"
	teamPostgres "github.com/OCCASS/avito-intern/internal/domain/team/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/domain/user"
	userPostgres "github.com/OCCASS/avito-intern/internal/domain/user/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2ETeam(t *testing.T) {
	if err := CleanDb(db); err != nil {
		t.Fatal("Clean database error.")
	}

	// Repositories
	teamRepository := teamPostgres.NewTeamPostgresRepository(db)
	pullrequestRepository := pullrequestPostgres.NewPullRequestPostgresRepository(db)
	userRepository := userPostgres.NewUserPostgresRepository(db)

	// Services
	pullrequestService := pullrequest.NewPullRequestServices(pullrequestRepository, teamRepository)
	userService := user.NewUserServices(userRepository, teamRepository, pullrequestRepository)
	teamService := team.NewTeamServices(teamRepository, userRepository, pullrequestService, userService)

	// Create
	createTeam := tDto.CreateTeamDto{
		Name: "backend",
		Members: []entity.User{
			{Id: "u-1", Name: "Ivan", IsActive: true},
			{Id: "u-2", Name: "Vanya", IsActive: true},
			{Id: "u-3", Name: "Lena", IsActive: true},
			{Id: "u-4", Name: "Nastya", IsActive: false},
		},
	}
	team, err := teamService.Add(createTeam)
	require.NoError(t, err)
	assert.Equal(t, "backend", team.Name)
	assert.Len(t, team.Members, 4)

	// Create PR
	createPullRequestDto := prDto.CreatePullRequestDto{
		Id:       "pr-1001",
		Name:     "Fix",
		AuthorId: "u-1",
	}
	pr, err := pullrequestService.Create(createPullRequestDto)
	require.NoError(t, err)
	assert.Equal(t, pr.ReviewersIds, []string{"u-2", "u-3"})

	// Create again
	team, err = teamService.Add(createTeam)
	require.Error(t, err)

	// Deactvate all team members
	var membersToDeactivate = []string{"u-2", "u-3"}
	deactivateDto := tDto.DeactivateMembersDto{
		Name:       "backend",
		MembersIds: membersToDeactivate,
	}
	team, err = teamService.DeactivateMembers(deactivateDto)
	require.NoError(t, err)
	var any = false
	for _, member := range team.Members {
		if slices.Contains(membersToDeactivate, member.Id) {
			any = any || member.IsActive
		}
	}
	assert.False(t, any)

	// Check pull request
	pr, err = pullrequestService.Get("pr-1001")
	require.NoError(t, err)
	assert.Empty(t, pr.ReviewersIds)
}
