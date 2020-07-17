package process

// ExitChan is the channel where all internal services should send exit code to shutdown the application
var ExitChan = make(chan int)

// SafeExit shuts down the application safely
func SafeExit(exitCode int) {
	// The ExitChan is expected to be listened by `main` function
	ExitChan <- exitCode
}
