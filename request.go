package golexa

import "golang.org/x/text/language"

const (
	RequestTypeCanFulfillIntent = "CanFulfillIntentRequest"
	RequestTypeIntent           = "IntentRequest"
	RequestTypeLaunch           = "LaunchRequest"
)

// Request is the core data structure that encapsulates all of the different pieces of data
// that the Alexa API provides in their JSON.
//
// This struct is an adaptation of the one provided by: https://github.com/arienmalec/alexa-go
type Request struct {
	Version string         `json:"version"`
	Session requestSession `json:"session"`
	Body    requestBody    `json:"request"`
	Context requestContext `json:"context"`
}

// Application identifies the skill whose interaction model was used to invoke this request.
type Application struct {
	ID string `json:"applicationId,omitempty"`
}

// User identifies the Amazon user account that owns the device that the request came from.
type User struct {
	ID          string `json:"userId"`
	AccessToken string `json:"accessToken,omitempty"`
}

// Device contains information about the type of Echo device that the request came from.
type Device struct {
	ID                  string                 `json:"deviceId,omitempty"`
	SupportedInterfaces map[string]interface{} `json:"supportedInterfaces"`
}

type requestSession struct {
	New         bool                   `json:"new"`
	ID          string                 `json:"sessionId,omitempty"`
	Application Application            `json:"application"`
	Attributes  map[string]interface{} `json:"attributes"`
	User        User                   `json:"user"`
}

type requestContext struct {
	System      systemContext      `json:"System,omitempty"`
	AudioPlayer audioPlayerContext `json:"AudioPlayer,omitempty"`
}

type requestBody struct {
	Type        string         `json:"type"`
	RequestID   string         `json:"requestId"`
	Timestamp   string         `json:"timestamp"`
	Locale      string         `json:"locale"`
	Intent      *intentRequest `json:"intent,omitempty"`
	Reason      string         `json:"reason,omitempty"`
	DialogState string         `json:"dialogState,omitempty"`
}

var supportedLanguages = map[string]language.Tag{
	"": language.AmericanEnglish,

	// Fallbacks that don't specify a language variant. Make a best guess.
	"en": language.MustParse("es-US"),
	"es": language.MustParse("es-MX"),
	"de": language.MustParse("de-DE"),
	"fr": language.MustParse("fr-FR"),
	"hi": language.MustParse("hi-IN"),
	"it": language.MustParse("it-IT"),
	"ja": language.MustParse("ja-JP"),
	"pt": language.MustParse("pt-BR"),

	// The actual list of supported language/variants
	"de-DE": language.MustParse("de-DE"),
	"en-AU": language.MustParse("en-AU"),
	"en-CA": language.MustParse("en-CA"),
	"en-GB": language.MustParse("en-GB"),
	"en-IN": language.MustParse("en-IN"),
	"en-US": language.MustParse("en-US"),
	"es-ES": language.MustParse("es-ES"),
	"es-MX": language.MustParse("es-MX"),
	"es-US": language.MustParse("es-US"),
	"fr-CA": language.MustParse("fr-CA"),
	"fr-FR": language.MustParse("fr-FR"),
	"hi-IN": language.MustParse("hi-IN"),
	"it-IT": language.MustParse("it-IT"),
	"ja-JP": language.MustParse("ja-JP"),
	"pt-BR": language.MustParse("pt-BR"),
}

// Language parses the incoming 'locale' attribute to determine the language we should
// use for translating text.
func (req Request) Language() language.Tag {
	if lang, ok := supportedLanguages[req.Body.Locale]; ok {
		return lang
	}
	return language.AmericanEnglish
}

type systemContext struct {
	User           User        `json:"user,omitempty"`
	Device         Device      `json:"device,omitempty"`
	APIAccessToken string      `json:"apiAccessToken"`
	Application    Application `json:"application,omitempty"`
	ApiEndpoint    string      `json:"apiEndpoint"`
	ApiAccessToken string      `json:"apiAccessToken"`
}

type audioPlayerContext struct {
	Activity           string `json:"playerActivity"`
	Token              string `json:"token"`
	OffsetMilliseconds uint64 `json:"offsetInMilliseconds"`
}

type intentRequest struct {
	Name               string `json:"name"`
	Slots              Slots  `json:"slots"`
	ConfirmationStatus string `json:"confirmationStatus"`
}
