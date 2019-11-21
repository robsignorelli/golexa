package golexa

// The standard intents that your skill should implement to handle standard Alexa interactions.
const (
	IntentNameCancel       = "AMAZON.CancelIntent"
	IntentNameFallback     = "AMAZON.FallbackIntent"
	IntentNameHelp         = "AMAZON.HelpIntent"
	IntentNameNavigateHome = "AMAZON.NavigateHomeIntent"
	IntentNameStop         = "AMAZON.StopIntent"
)

type intentRoute struct {
	handlerFunc HandlerFunc
	name        string
}
