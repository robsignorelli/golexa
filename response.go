package golexa

import (
	"errors"
	"strings"
	"time"

	"github.com/robsignorelli/golexa/speech"
	"github.com/sirupsen/logrus"
)

// NewResponse create a bare-bones response instance that you can continue to expand on
// with additional instructions for Alexa like `Speak()` or `SimpleCard()`.
func NewResponse(request Request) Response {
	r := Response{
		Request: request,
		Version: "1.0",
		Body:    responseBody{},
	}
	return r.EndSession(true)
}

// EndSession indicates whether or not the user is at the end of a dialog w/ Alexa. This
// defaults to 'true' so that all interactions are terminal unless otherwise specified.
func (r Response) EndSession(flag bool) Response {
	r.Body.ShouldEndSession = &flag
	return r
}

// Speak indicates w/ you want the Alexa voice to dictate back to the user. You can provide
// plain text or SSML.
func (r Response) Speak(textOrSSML string) Response {
	r.Body.OutputSpeech = &intentResponse{
		SSML: wrapSSML(textOrSSML),
	}
	return r
}

func (r Response) SpeakTemplate(template speech.Template, value interface{}) Response {
	textOrSSML, err := template.Eval(speech.TemplateContext{
		Language: r.Request.Language(),
		Now:      time.Now(),
		Value:    value,
	})
	if err != nil {
		logrus.Error("unable to speak template: %v", err)
		return r.Speak("I'm sorry. I seem to have trouble with words, today.")
	}
	return r.Speak(textOrSSML)
}

// SimpleCard customizes what the user should see on an Echo device that supports a screen
// or what shows up when they look at their interaction history in the Alexa app.
func (r Response) SimpleCard(title, text string) Response {
	r.Body.Card = &intentResponse{
		Type:    "Simple",
		Title:   title,
		Content: text,
	}
	return r
}

// ElicitSlot keeps the current session open and has the user's echo device go back into capture
// mode. Whatever the user speaks next will be applied to the specified slot and all other slots
// from this request will be sent along to the slot you named. You should use this in conjunction
// with `Speak()` so that Alexa will give a meaningful prompt rather than just displaying a blue ring.
//
// For example, if you're setting up some sort of game, you might respond with the speech "How many players?"
// and `ElicitSlot(req, "StartGameIntent", "num_players")`. Once the user tell you how many players
// then that slot info will be sent to your "StartGameIntent".
func (r Response) ElicitSlot(request Request, intentName, slotName string) Response {
	slots := request.Body.Intent.Slots.Clone()
	slots[slotName] = Slot{Name: slotName, Value: ""}

	r.Body.Directives = append(r.Body.Directives, directive{
		Type:          "Dialog.ElicitSlot",
		SlotToElicit:  slotName,
		UpdatedIntent: &updatedIntent{Name: intentName, Slots: slots, ConfirmationStatus: "NONE"},
	})

	// Need to keep the session open so we can prompt the user and have them respond.
	return r.EndSession(false)
}

// Reprompt should be used in conjunction w/ an `ElicitSlot()` call. If the user doesn't say anything
// when they're asked to fill in one of the slots, this will be a second audio prompt to try to get them
// to say something. If the user actually responded the first time, they won't actually hear this.
func (r Response) Reprompt(textOrSSML string) Response {
	r.Body.Reprompt.OutputSpeech = intentResponse{
		SSML: wrapSSML(textOrSSML),
	}
	return r
}

// Ok simply returns the Response in its current state and a 'nil' error. This is a convenience so
// that you can build your response at the end of your handlers which require a response and an error.
func (r Response) Ok() (Response, error) {
	return r, nil
}

// Fail should be used in only the most dire of unrecoverable circumstances. It will respond
// with no Alexa instructions and an error w/ the given message. You should NOT use this in
// instances where your skill can't give a meaningful response to a question. It should only
// be used for critical, unexpected paths such as receiving a request for an intent your code
// doesn't have registered.
func Fail(message string) (Response, error) {
	return Response{}, errors.New(message)
}

// Response encapsulates all of the various options that your skill can respond with to control
// the audio and visual aspects of the experience. It is a builder that lets you add on additional
// pieces to the experience as you deem necessary. For instance you can invoke `Speak()` to control
// what Alexa will say to the user then call `SimpleCard()` to provide some basic visual feedback.
//
// This struct is an adaptation of the one provided by: https://github.com/arienmalec/alexa-go
//
// Also see https://developer.amazon.com/docs/custom-skills/request-and-response-json-reference.html#response-format
type Response struct {
	Request           Request                `json:"-"`
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              responseBody           `json:"response"`
}

type responseBody struct {
	OutputSpeech     *intentResponse `json:"outputSpeech,omitempty"`
	Card             *intentResponse `json:"card,omitempty"`
	Reprompt         *reprompt       `json:"reprompt,omitempty"`
	Directives       []directive     `json:"directives,omitempty"`
	ShouldEndSession *bool           `json:"shouldEndSession,omitempty"`
}

type intentResponse struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	SSML    string `json:"ssml,omitempty"`
	Content string `json:"content,omitempty"`
}

type directive struct {
	Type          string         `json:"type,omitempty"`
	SlotToElicit  string         `json:"slotToElicit,omitempty"`
	UpdatedIntent *updatedIntent `json:"updatedIntent,omitempty"`
	PlayBehavior  string         `json:"playBehavior,omitempty"`
	AudioItem     struct {
		Stream struct {
			Token                string `json:"token,omitempty"`
			URL                  string `json:"url,omitempty"`
			OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
		} `json:"stream,omitempty"`
	} `json:"audioItem,omitempty"`
}

type updatedIntent struct {
	Name               string `json:"name,omitempty"`
	ConfirmationStatus string `json:"confirmationStatus,omitempty"`
	Slots              Slots  `json:"slots,omitempty"`
}

type reprompt struct {
	OutputSpeech intentResponse `json:"outputSpeech,omitempty"`
}

// wrapSSML ensures that the text you want Alexa to speak is SSML. It allows you to
// utilize the same attribute in the response whether you are simply giving plain text
// or you built your own SSML markup.
func wrapSSML(textOrSSML string) string {
	if strings.HasPrefix(textOrSSML, "<speak") {
		return textOrSSML
	}
	return "<speak>" + textOrSSML + "</speak>"
}
