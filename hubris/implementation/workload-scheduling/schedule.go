package workload_scheduling

import "context"

// Decider defined each of the decision components / units of decisions
type Decider interface {
	Score(ctx context.Context) (map[string]int, error)
}

// has to implement the hubris.Hubris interface
type Scheduler struct {
	// weight defines the weight attributed to the scores from each of the deciders
	weight map[Decider]int
}

func NewScheduler(ctx context.Context) *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) scoreLatency(ctx context.Context) *Scheduler {
}
