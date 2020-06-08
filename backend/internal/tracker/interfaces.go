package tracker

type ITracker interface {
	Poll()
	Home()
	Track()
	Stop()

	IsHomed() bool
	IsHoming() bool
	IsTracking() bool
	IsFinished() bool
}
