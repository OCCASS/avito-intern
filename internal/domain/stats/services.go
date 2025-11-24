package stats

import (
	"github.com/OCCASS/avito-intern/internal/application/stats"
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"
)

type StatsServices struct {
	userRepository userRepository.UserRepository
}

func NewStatsServices(
	ur userRepository.UserRepository,
) StatsServices {
	return StatsServices{
		userRepository: ur,
	}
}

func (s StatsServices) UsersStats() (stats.UsersStatsDto, error) {
	return s.userRepository.GetStats()
}
