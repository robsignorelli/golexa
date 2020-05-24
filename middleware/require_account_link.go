package middleware

import (
	"context"

	"github.com/robsignorelli/golexa"
	"github.com/robsignorelli/golexa/speech"
	"github.com/sirupsen/logrus"
)

func RequireAccount(options ...RequireAccountOption) golexa.MiddlewareFunc {
	// I really don't expect you to use this text out of the box, but if you want, it's up to you.
	r := requireAccount{
		template: speech.NewTemplate("I'm sorry. You must connect your account using the Alexa app in order to use this feature."),
	}
	for _, opt := range options {
		opt(&r)
	}
	return r.checkAccessToken
}

type RequireAccountOption func(*requireAccount)

func RequireAccountTemplate(t speech.Template) RequireAccountOption {
	return func(r *requireAccount) {
		r.template = t
	}
}

type requireAccount struct {
	template speech.Template
}

func (r requireAccount) checkAccessToken(ctx context.Context, request golexa.Request, next golexa.HandlerFunc) (golexa.Response, error) {
	if request.Context.System.User.AccessToken != "" {
		return next(ctx, request)
	}

	logrus.WithField("label", "golexa").
		WithField("request.id", request.Body.RequestID).
		WithField("user.id", request.Context.System.User.ID).
		WithField("device.id", request.Context.System.Device.ID).
		Info("Missing user access token")

	return golexa.NewResponse(request).
		SpeakTemplate(r.template, request).
		Ok()
}
