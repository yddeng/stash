package serverType

const (
	Dir        = uint32(1)
	Login      = uint32(2)
	Gate       = uint32(3)
	Game       = uint32(4)
	World      = uint32(5)
	Map        = uint32(6)
	Team       = uint32(7)
	WebService = uint32(8)
	Rank       = uint32(14)
)

var (
	serverTypeName = map[uint32]string{
		Dir:        "Dir",
		Login:      "Login",
		Gate:       "Gate",
		Game:       "Game",
		World:      "World",
		Map:        "Map",
		Team:       "Team",
		WebService: "WebService",
		Rank:       "Rank",
	}
)

func Type2Name(t uint32) string {
	return serverTypeName[t]
}
