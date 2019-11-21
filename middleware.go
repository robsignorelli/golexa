package golexa

import "context"

// MiddlewareFunc defines a standardized unit of work that you want to execute after an Alexa
// request comes in, but before your actual intent handler fires. By simply not invoking the 'next'
// function, you can short circuit the execution without running your intent handler at all (think
// restricting access to intents based on whether they're account linked or not).
type MiddlewareFunc func(ctx context.Context, request Request, next HandlerFunc) (Response, error)

// Middleware represents a chain of units of work to execute before your intent/request handler.
type Middleware []MiddlewareFunc

// Then creates a wrapper function that forces a request to run through a your gauntlet of middleware
// functions before finally executing the intent/request handler you're registering.
func (m Middleware) Then(handler HandlerFunc) HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		mw := m[i]
		if mw == nil {
			continue
		}
		next := handler
		handler = func(ctx context.Context, request Request) (Response, error) {
			return mw(ctx, request, next)
		}
	}
	return handler
}
