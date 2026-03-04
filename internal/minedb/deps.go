package minedb

type logger interface {
	Info(...any)
	Warn(...any)
	Error(...any)
	Newlogger()
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
}
