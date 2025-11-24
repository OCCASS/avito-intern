package tests

import (
	"testing"

	tDto "github.com/OCCASS/avito-intern/internal/application/team"
	"github.com/OCCASS/avito-intern/internal/domain/team"
	teamPostgres "github.com/OCCASS/avito-intern/internal/domain/team/repository/postgres"
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
	userRepository := userPostgres.NewUserPostgresRepository(db)

	// Services
	teamService := team.NewTeamServices(teamRepository, userRepository)

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

	// Create again
	team, err = teamService.Add(createTeam)
	require.Error(t, err)
}
