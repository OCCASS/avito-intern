package postgres

import (
	"database/sql"

	"github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/domain/team/repository"
	"github.com/OCCASS/avito-intern/internal/entity"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

type TeamPostgresRepository struct {
	db *database.Database
}

func NewTeamPostgresRepository(db *database.Database) *TeamPostgresRepository {
	return &TeamPostgresRepository{db}
}

func (r TeamPostgresRepository) Create(team entity.Team) (entity.Team, error) {
	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return entity.Team{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	queryTeam := `INSERT INTO team(name) VALUES ($1)`
	if _, err := tx.Exec(queryTeam, team.Name); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case pgerrcode.UniqueViolation:
				return entity.Team{}, repository.ErrTeamAlreadyExists
			}
		}
		return entity.Team{}, err
	}

	membersIds := make([]string, 0, len(team.Members))
	for i := 0; i < len(team.Members); i++ {
		membersIds = append(membersIds, team.Members[i].Id)
	}

	queryTeamMembers := `INSERT INTO team_member(team_name, member_id) SELECT $1, UNNEST($2::TEXT[])`
	if _, err := tx.Exec(queryTeamMembers, team.Name, pq.Array(membersIds)); err != nil {
		return entity.Team{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.Team{}, err
	}

	return team, nil
}

func (r TeamPostgresRepository) Get(name string) (entity.Team, error) {
	var newTeam entity.Team

	queryTeam := `SELECT name FROM team WHERE name=$1`
	if err := r.db.Conn.Get(&newTeam, queryTeam, name); err != nil {
		if err == sql.ErrNoRows {
			return entity.Team{}, repository.ErrTeamNotFound
		}
		return entity.Team{}, err
	}

	var members []entity.User
	queryMembers := `SELECT u.id, u.name, u.is_active
	FROM team_member tm
	JOIN "user" u ON u.id = tm.member_id
	WHERE tm.team_name=$1`
	if err := r.db.Conn.Select(&members, queryMembers, name); err != nil {
		return entity.Team{}, err
	}

	newTeam.Members = members

	return newTeam, nil
}

func (r TeamPostgresRepository) GetByUser(userId string) (entity.Team, error) {
	var teamName string
	if err := r.db.Conn.Get(&teamName, `SELECT team_name FROM team_member WHERE member_id=$1`, userId); err != nil {
		if err == sql.ErrNoRows {
			return entity.Team{}, repository.ErrTeamNotFound
		}
		return entity.Team{}, err
	}
	return r.Get(teamName)
}
