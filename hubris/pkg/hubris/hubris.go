package hubris

import "context"

// Hubris defines the interface for managing multiple-heuristic decisions machines that scores their candidates based on predefined challenges
type Hubris interface {
	// Add adds a new component to current hubris pipeline
	Add(ctx *context.Context, component interface{})
	// Pick returns the hubris with the id provided in the argument
	Pick(ctx context.Context, id string) []interface{}
	// Fetch returns a list of all the components current defined in the pipeline
	Fetch(ctx context.Context) []interface{}
	// Score returns a map of the hubris id to their score after the pipeline challenges have been executed
	Score(ctx context.Context) (map[string]int, error)
}
