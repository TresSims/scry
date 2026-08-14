package facts

// Facter interface is an interface for gathering facts about a system
type Facter interface {
	// Get a fact about the system
	Get() any
}
