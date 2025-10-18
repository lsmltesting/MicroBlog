package logger

type Level string
type Source string

const (
	LevelInfo     Level = "INFO"
	LevelError    Level = "ERROR"
	LevelDebug    Level = "DEBUG"
	LevelWarning  Level = "WARN"
	LevelBusiness Level = "BUSINESS"
)

const (
	SourceMain    Source = "MAIN"
	SourceHandler Source = "HANDLER"
	SourceService Source = "SERVICE"
	SourceRepo    Source = "REPO"
	SourceDB      Source = "DB"
)
