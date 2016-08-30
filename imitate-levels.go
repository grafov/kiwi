package kiwi

// LevelName allows to change default recVal "level" to any recVal you want.
// Set it to empty string if you want to report level without presetting any name.
var LevelName = "level"

func (l *Logger) Fatal(keyVals ...interface{}) {
	l.Add(keyVals...).Log(LevelName, "fatal")
}

func (l *Logger) Crit(keyVals ...interface{}) {
	l.Add(keyVals...).Log(LevelName, "critical")
}

// Error imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "error". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any recVal you want.
func (l *Logger) Error(keyVals ...interface{}) {
	l.Add(keyVals...).Log(LevelName, "error")
}

// Warn imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "warning". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any recVal you want.
func (l *Logger) Warn(keyVals ...interface{}) {
	l.Add(keyVals...).Log(LevelName, "warning")
}

// Info imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "info". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any value what you want.
func (l *Logger) Info(keyVals ...interface{}) {
	l.Add(keyVals...).Log(LevelName, "info")
}

func (l *Logger) Debug(keyVals ...interface{}) {
	l.Add(keyVals...).Log(LevelName, "debug")
}
