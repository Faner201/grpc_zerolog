package grpc_zerolog

import (
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
)

// LoggableEvent defines the events a log line can be added on.
type LoggableEvent uint

const (
	// StartCall is a loggable event representing start of the gRPC call.
	StartCall LoggableEvent = iota
	// FinishCall is a loggable event representing finish of the gRPC call.
	FinishCall
	// PayloadReceived is a loggable event representing received request (server) or response (client).
	// Log line for this event also includes (potentially big) proto.Message of that payload in
	// "grpc.request.content" (server) or "grpc.response.content" (client) field.
	// NOTE: This can get quite verbose, especially for streaming calls, use with caution (e.g. debug only purposes).
	PayloadReceived
	// PayloadSent is a loggable event representing sent response (server) or request (client).
	// Log line for this event also includes (potentially big) proto.Message of that payload in
	// "grpc.response.content" (server) or "grpc.request.content" (client) field.
	// NOTE: This can get quite verbose, especially for streaming calls, use with caution (e.g. debug only purposes).
	PayloadSent
)

var (
	// DefaultCodeToLevelFunc is the default implementation code to level logic.
	// returns Info on codes.OK and Error in all other cases
	DefaultCodeToLevelFunc CodeToLevel = func(code codes.Code) zerolog.Level {
		if code == codes.OK {
			return zerolog.InfoLevel
		}
		return zerolog.ErrorLevel
	}

	// DefaultDeciderFunc is the default implementation of decider
	// returns true always
	DefaultDeciderFunc Decider = func(fullMethodName string, err error) bool {
		return true
	}

	defaultOptions = &options{
		levelFunc:      DefaultCodeToLevelFunc,
		shouldLog:      DefaultDeciderFunc,
		loggableEvents: []LoggableEvent{StartCall, FinishCall},
	}
)

// CodeToLevel function defines the mapping between gRPC return codes and interceptor log level
type CodeToLevel func(code codes.Code) zerolog.Level

// Decider function defines rules for suppressing any interceptor logs
type Decider func(fullMethodName string, err error) bool

// Option used to configure the interceptors
type Option func(*options)

// WithLevels customizes the function for mapping gRPC return codes and interceptor log level statements
func WithLevels(f CodeToLevel) Option {
	return func(o *options) {
		o.levelFunc = f
	}
}

// WithDecider customizes the function for deciding if the gRPC interceptor logs should log depends on fullMethodName and error from handler
func WithDecider(f Decider) Option {
	return func(o *options) {
		o.shouldLog = f
	}
}

// WithLogOnEvents customizes on what events the gRPC interceptor should log on.
func WithLogOnEvents(events ...LoggableEvent) Option {
	return func(o *options) {
		o.loggableEvents = events
	}
}

type options struct {
	levelFunc      CodeToLevel
	shouldLog      Decider
	loggableEvents []LoggableEvent
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}
