package game

type Option func(*Game)

func WithIDGenerator(gener IDGenerator) Option {
	return func(g *Game) {
		g.IDGenerator = gener
	}
}
