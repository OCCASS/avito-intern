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

func TestE2EPullRequest(t *testing.T) {
	if err := CleanDb(db); err != nil {
		t.Fatal("Clean database error.")
	}

	// Repositories
	pullrequestRepository := prPostgres.NewPullRequestPostgresRepository(db)
	teamRepository := teamPostgres.NewTeamPostgresRepository(db)
	userRepository := userPostgres.NewUserPostgresRepository(db)

	// Services
	pullrequestService := pullrequest.NewPullRequestServices(pullrequestRepository, teamRepository)
	userService := user.NewUserServices(userRepository, teamRepository, pullrequestRepository)
	teamService := team.NewTeamServices(teamRepository, userRepository, pullrequestService, userService)

	// Create team
	createTeam := tDto.CreateTeamDto{
		Name: "backend",
		Members: []entity.User{
			{Id: "u-1", Name: "Ivan", IsActive: true},
			{Id: "u-2", Name: "Vanya", IsActive: true},
			{Id: "u-3", Name: "Lena", IsActive: true},
			{Id: "u-4", Name: "Nastya", IsActive: false},
		},
	}
	_, err := teamService.Add(createTeam)
	require.NoError(t, err)

	// Create pull request
	createDto := prDto.CreatePullRequestDto{
		Id:       "pr-1001",
		Name:     "Fix",
		AuthorId: "u-1",
	}
	pr, err := pullrequestService.Create(createDto)
	require.NoError(t, err)
	assert.Equal(t, []string{"u-2", "u-3"}, pr.ReviewersIds)
	assert.Equal(t, entity.StatusOpen, pr.Status)

	// Reassign pull request reviewer
	reassignDto := prDto.ReassignPullRequestDto{
		PullRequestId: "pr-1001",
		OldReviewerId: "u-2",
	}
	pr, newReviewerId, err := pullrequestService.Reassign(reassignDto)
	require.Error(t, err)
	assert.Nil(t, newReviewerId)

	// Merge
	mergeDto := prDto.MergePullRequestDto{
		Id: "pr-1001",
	}
	pr, err = pullrequestService.Merge(mergeDto)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusMerged, pr.Status)

	// Merge again
	pr1, err := pullrequestService.Merge(mergeDto)
	require.NoError(t, err)
	assert.Equal(t, pr, pr1)

	// Create another team
	createTeam = tDto.CreateTeamDto{
		Name: "frontend",
		Members: []entity.User{
			{Id: "u-5", Name: "Ivan", IsActive: true},
			{Id: "u-6", Name: "Vanya", IsActive: true},
			{Id: "u-7", Name: "Lena", IsActive: true},
			{Id: "u-8", Name: "Nastya", IsActive: true},
			{Id: "u-9", Name: "Nastya", IsActive: true},
			{Id: "u-10", Name: "Nastya", IsActive: true},
		},
	}
	_, err = teamService.Add(createTeam)
	require.NoError(t, err)

	// Create pull request
	createDto = prDto.CreatePullRequestDto{
		Id:       "pr-1002",
		Name:     "Fix_2",
		AuthorId: "u-5",
	}
	pr, err = pullrequestService.Create(createDto)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusOpen, pr.Status)

	// Reassign
	reassignDto = prDto.ReassignPullRequestDto{
		PullRequestId: "pr-1002",
		OldReviewerId: pr.ReviewersIds[0],
	}
	newPr, newReviewerId, err := pullrequestService.Reassign(reassignDto)
	require.NoError(t, err)
	assert.NotContains(t, newPr.ReviewersIds, pr.ReviewersIds[0])
	assert.NotContains(t, newPr.ReviewersIds, newPr.AuthorId)
	assert.Contains(t, newPr.ReviewersIds, *newReviewerId)
}
