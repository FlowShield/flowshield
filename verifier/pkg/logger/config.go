package logger

func LogLevel() string {
	return "debug"
}

func LogFile() string {
	return LogPath() + LogName()
}

func SendLogToFile() bool {
	return true
}

func LogPath() string {
	return "logs/"
}

func LogName() string {
	return "verifier.log"
}
