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

// UserID traverses the request structure to extract the id of the Amazon/Alexa user making the call.
func (r Request) UserID() string {
	return r.Context.System.User.ID
}

// UserAccessToken traverses the request structure to extract the account-linked access token for the caller.
func (r Request) UserAccessToken() string {
	return r.Context.System.User.AccessToken
}

// DeviceID traverses the request structure to extract the id of the device making the call.
func (r Request) DeviceID() string {
	return r.Context.System.Device.ID
}

// SkillID returns the id of the skill that was invoked to handle this request.
func (r Request) SkillID() string {
	return r.Context.System.Application.ID
}

// SessionID traverses the request structure to extract the id of the Amazon/Alexa session making the call.
func (r Request) SessionID() string {
	return r.Session.ID
}

// Language parses the incoming 'locale' attribute to determine the language we should
// use for translating text.
func (r Request) Language() language.Tag {
	if lang, ok := supportedLanguages[r.Body.Locale]; ok {
		return lang
	}
	return language.AmericanEnglish
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

// NewIntentRequest creates a minimal request instance you can use to write
// unit tests for your intent requests.
func NewIntentRequest(intentName string, slots Slots) Request {
	return Request{
		Version: "1.0",
		Body: requestBody{
			Intent: &intentRequest{
				Name:  intentName,
				Slots: slots,
			},
		},
	}
}
