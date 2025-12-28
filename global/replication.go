package global

import "github.com/Rx-11/distributed-leaderboard/leaderboard"

type Publisher interface {
	Publish(summary leaderboard.RegionSummary) error
}

type Subscriber interface {
	Receive(summary leaderboard.RegionSummary) error
}

type Replicator interface {
	Tick() error
}
