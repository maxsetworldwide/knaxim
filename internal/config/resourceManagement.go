package config

var resources chan struct{}

// GetResourceTracker returns a buffered channel to act as a count of the active connections to subsystems and block when too many connections are open
func GetResourceTracker() chan struct{} {
	if V.ActiveFileProcessing == 0 {
		resource := make(chan struct{}, 1)
		resource <- struct{}{}
		return resource
	}
	return resources
}
