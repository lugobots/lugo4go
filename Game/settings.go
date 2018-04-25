package Game

type TeamName string

const HomeTeam TeamName = "home"
const AwayTeam TeamName = "away"

type StartUpSettings struct {
	IncomeTopic string
	OutcomeTopic string
	TeamName TeamName
	Port int
	TeamTitle string
}

type PlayerNumber string

const (
	POSITION_A PlayerNumber = "1"
	POSITION_B PlayerNumber = "2"
	POSITION_C PlayerNumber = "3"
	POSITION_D PlayerNumber = "4"
	POSITION_E PlayerNumber = "5"
)
