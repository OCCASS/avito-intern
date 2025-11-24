package tests

import "github.com/OCCASS/avito-intern/internal/database"

func CleanDb(db *database.Database) {
	db.Conn.Exec(`TRUNCATE "user", team, pullrequest, team_member, pullrequest_reviewer RESTART IDENTITY CASCADE`)
}
