package log

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
)

func DefaultLevels() Levels {
	return Levels{
		DebugLevel: 10,
		InfoLevel:  20,
		WarnLevel:  40,
		ErrorLevel: 80,
	}
}

func DefaultLevel() Level {
	return WarnLevel
}
