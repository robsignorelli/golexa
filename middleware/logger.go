package middleware

import (
	"context"
	"time"

	"github.com/robsignorelli/golexa"
	"github.com/sirupsen/logrus"
)

// Logger creates a middleware function that logs the start/end of every single utterance you handle
// in your skill. We include all of the "important" things that you probably care about: user/device id,
// elapsed time, request id, etc. You can optionally have the logger spit out the entire JSON of the
// incoming requests as well as what you actually output.
func Logger(options ...LoggerOption) golexa.MiddlewareFunc {
	logger := loggerMiddleware{}
	for _, opt := range options {
		opt(&logger)
	}
	return logger.log
}

// LoggerOption tweaks the settings of your logging middleware. Please use the built-in
// helpers like LogResponseSpeech() and LogRequestJSON().
type LoggerOption func(*loggerMiddleware)

// LogResponseSpeech flips on the setting to include the output speech SSML on every response. To
// better protect your users' privacy, you should probably only turn this on in dev/staging.
func LogResponseSpeech() LoggerOption {
	return func(logger *loggerMiddleware) {
		logger.IncludeSpeech = true
	}
}

// LogRequestJSON flips on the setting to include an attribute that outputs the entire incoming
// request payload on every request. To better protect your users' privacy, you should probably only
// turn this on in dev/staging.
func LogRequestJSON() LoggerOption {
	return func(logger *loggerMiddleware) {
		logger.IncludeRequestJSON = true
	}
}

type loggerMiddleware struct {
	IncludeRequestJSON bool
	IncludeSpeech      bool
}

func (logger loggerMiddleware) log(ctx context.Context, request golexa.Request, next golexa.HandlerFunc) (golexa.Response, error) {
	logger.logStart(request)

	startTime := time.Now()
	response, err := next(ctx, request)
	elapsed := time.Now().Sub(startTime)

	logger.logFinished(request, response, err, elapsed)
	return response, err
}

func (logger loggerMiddleware) logStart(request golexa.Request) {
	entry := logrus.WithField("label", "golexa").
		WithField("request.id", request.Body.RequestID).
		WithField("user.id", request.Context.System.User.ID).
		WithField("device.id", request.Context.System.Device.ID)

	if request.Body.Intent != nil {
		entry = entry.WithField("intent.name", request.Body.Intent.Name)

		for _, slot := range request.Body.Intent.Slots {
			entry = entry.WithField("intent.slot."+slot.Name, slot.Resolve())
		}
	}

	if logger.IncludeRequestJSON {
		entry = entry.WithField("request.json", request)
	}
	entry.Info("Request started")
}

func (logger loggerMiddleware) logFinished(request golexa.Request, response golexa.Response, err error, elapsed time.Duration) {
	entry := logrus.WithField("label", "golexa").
		WithField("request.id", request.Body.RequestID).
		WithField("user.id", request.Context.System.User.ID).
		WithField("device.id", request.Context.System.Device.ID).
		WithField("elapsed", elapsed).
		WithField("elapsed_human", elapsed.String())

	if err != nil {
		entry = entry.WithField("error", err)
	}
	if speak := response.Body.OutputSpeech; logger.IncludeSpeech && speak != nil {
		entry = entry.WithField("response.speech", speak.SSML)
	}

	entry.Info("Request complete")
}
