package proxy

import (
	ec "github.com/labstack/echo/v4"
)

// Echos is an interface that provides a proxy of the methods of echo.
type Echos interface {
	NewEcho() (Echo, Logger)
}

// echosProxy is a proxy struct that implements the Echos interface.
type echosProxy struct{}

// NewEchos returns a new instance of the Echos interface.
func NewEchos() Echos {
	return &echosProxy{}
}

// NewEcho returns a new instance of the echo.Echo and echo.Logger.
func (*echosProxy) NewEcho() (Echo, Logger) {
	echo := ec.New()
	return &ehco{echo}, &logger{echo.Logger}
}

// Echo is an interface that provides a proxy of the methods of echo.Echo.
type Echo interface {
	Get(path string, h ec.HandlerFunc, m ...ec.MiddlewareFunc)
	Group(prefix string, m ...ec.MiddlewareFunc) Group
	Start(address string) error
	Use(middleware ...ec.MiddlewareFunc)
}

// ehco is a proxy struct that implements the Echo interface.
type ehco struct {
	*ec.Echo
}

// Get adds a GET route to the echo server.
func (e *ehco) Get(path string, h ec.HandlerFunc, m ...ec.MiddlewareFunc) {
	e.GET(path, h, m...)
}

// Group returns a new instance of the Group interface.
func (e *ehco) Group(prefix string, m ...ec.MiddlewareFunc) Group {
	return &group{e.Echo.Group(prefix, m...)}
}

// Start starts the echo server.
func (e *ehco) Start(address string) error {
	return e.Echo.Start(address)
}

// Use adds middleware to the echo server.
func (e *ehco) Use(middleware ...ec.MiddlewareFunc) {
	e.Echo.Use(middleware...)
}

// Logger is an interface that provides a proxy of the methods of echo.Logger.
type Logger interface {
	Fatal(err error)
}

// logger is a proxy struct that implements the Logger interface.
type logger struct {
	ec.Logger
}

// Fatal logs the error and exits the application.
func (l *logger) Fatal(err error) {
	l.Logger.Fatal(err)
}

// Group is an interface that provides a proxy of the methods of echo.Group.
type Group interface {
	GET(path string, h ec.HandlerFunc, m ...ec.MiddlewareFunc)
}

// group is a proxy struct that implements the Group interface.
type group struct {
	*ec.Group
}

// GET adds a GET route to the echo server.
func (g *group) GET(path string, h ec.HandlerFunc, m ...ec.MiddlewareFunc) {
	g.Group.GET(path, h, m...)
}
