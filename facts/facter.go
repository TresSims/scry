package facts

// Facter interface is an interface for gathering facts about a system
type Facter func() (any, error)

type Fact struct {
	Facter Facter

	Cache any
}
