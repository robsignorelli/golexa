package golexa

import (
	"context"
)

// HandlerFunc defines a core operation of your skill. It takes the request with all incoming
// user interaction information and returns a Response w/ instructions for how Alexa should
// react to the user's utterance/request.
type HandlerFunc func(ctx context.Context, request Request) (Response, error)

// The standard intents that your skill should implement to handle standard Alexa interactions.
const (
	IntentNameCancel       = "AMAZON.CancelIntent"
	IntentNameFallback     = "AMAZON.FallbackIntent"
	IntentNameHelp         = "AMAZON.HelpIntent"
	IntentNameNavigateHome = "AMAZON.NavigateHomeIntent"
	IntentNameStop         = "AMAZON.StopIntent"
)

// Skill is the root data structure for your program. It wrangles all of the handlers for the
// different types of requests your skill is expected to encounter.
type Skill struct {
	Name       string
	intents    map[string]intentRoute
	canFulfill HandlerFunc
	launch     HandlerFunc
}

// RouteIntent indicates that any "IntentRequest" with the specified intent name should be handled
// by the given function.
func (skill *Skill) RouteIntent(intentName string, handlerFunc HandlerFunc) {
	if skill.intents == nil {
		skill.intents = map[string]intentRoute{}
	}
	skill.intents[intentName] = intentRoute{
		name:        intentName,
		handlerFunc: handlerFunc,
	}
}

// CanFulfillIntent allows you to support the "pre-flight" CanFulfillIntentRequest if you want to
// have your skill undergo name free certification.
//
// See: https://developer.amazon.com/docs/custom-skills/understand-name-free-interaction-for-custom-skills.html
func (skill *Skill) CanFulfillIntent(handlerFunc HandlerFunc) {
	skill.canFulfill = handlerFunc
}

// Launch registers the handler for when the user utters "Alexa, open XXX" to launch your skill.
func (skill *Skill) Launch(handlerFunc HandlerFunc) {
	skill.launch = handlerFunc
}

// Handle routes the incoming Alexa request to the correct, registered handler.
func (skill Skill) Handle(ctx context.Context, request Request) (Response, error) {
	switch request.Body.Type {
	case RequestTypeIntent:
		return skill.handleIntent(ctx, request)
	case RequestTypeCanFulfillIntent:
		return skill.handleCanFulfillIntent(ctx, request)
	case RequestTypeLaunch:
		return skill.handleLaunch(ctx, request)
	default:
		return Fail("golexa: unsupported request type: " + request.Body.Type)
	}
}

func (skill Skill) handleIntent(ctx context.Context, request Request) (Response, error) {
	if request.Body.Intent == nil {
		return Fail("golexa: body is missing intent data for IntentRequest")
	}

	name := request.Body.Intent.Name
	intentRoute, ok := skill.intents[name]
	if !ok {
		return Fail("golexa: no handler registered for intent: " + name)
	}
	return intentRoute.handlerFunc(ctx, request)
}

func (skill Skill) handleCanFulfillIntent(ctx context.Context, request Request) (Response, error) {
	if skill.canFulfill == nil {
		return Fail("golexa: no handler registered for CanFulfillIntentRequest")
	}
	return skill.canFulfill(ctx, request)
}

func (skill Skill) handleLaunch(ctx context.Context, request Request) (Response, error) {
	if skill.launch == nil {
		return Fail("golexa: no handler registered for LaunchRequest")
	}
	return skill.launch(ctx, request)
}

type intentRoute struct {
	handlerFunc HandlerFunc
	name        string
}
