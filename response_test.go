package golexa_test

import (
	"testing"

	"github.com/robsignorelli/golexa"
	"github.com/robsignorelli/golexa/speech"
	"github.com/stretchr/testify/suite"
	"golang.org/x/text/language"
)

func TestResponseSuite(t *testing.T) {
	suite.Run(t, new(ResponseSuite))
}

type ResponseSuite struct {
	suite.Suite
}

func (suite ResponseSuite) TestNewResponse() {
	req := golexa.Request{Version: "X"}
	res := golexa.NewResponse(req)
	suite.Equal("X", res.Request.Version,
		"Original request should be maintained")
	suite.Equal("1.0", res.Version,
		"Response should default to version 1.0")
	suite.Nil(res.Body.OutputSpeech,
		"Response should have no output speech")
	suite.Nil(res.Body.Card,
		"Response should have no card")

	suite.Require().NotNil(res.Body.ShouldEndSession,
		"ShouldEndSession should be true by default")
	suite.True(*res.Body.ShouldEndSession,
		"ShouldEndSession should be true by default")
}

func (suite ResponseSuite) TestShouldEndSession() {
	res := golexa.NewResponse(golexa.Request{})

	res = res.EndSession(true)
	suite.True(*res.Body.ShouldEndSession,
		"ShouldEndSession should be true when flipped on")

	res = res.EndSession(false)
	suite.False(*res.Body.ShouldEndSession,
		"ShouldEndSession should be false when flipped off")

	res = res.EndSession(true)
	suite.True(*res.Body.ShouldEndSession,
		"ShouldEndSession should be true when back flipped on")

	res.EndSession(false)
	suite.True(*res.Body.ShouldEndSession,
		"EndSession() should not mutate the response")
}

func (suite ResponseSuite) TestSpeak() {
	res := golexa.NewResponse(golexa.Request{})

	res = res.Speak("Hello")
	suite.Equal("<speak>Hello</speak>", res.Body.OutputSpeech.SSML,
		"Should automatically wrap plain text in <speak> SSML tag.")
	suite.Equal("", res.Body.OutputSpeech.Text,
		"Should use the SSML speech attribute, not the Text attribute.")
	suite.Nil(res.Body.Card,
		"Should not create any Card information")

	res = res.Speak("<speak>The cow goes woof</speak>")
	suite.Equal("<speak>The cow goes woof</speak>", res.Body.OutputSpeech.SSML,
		"Should leave SSML speech alone.")

	res.Speak("Goodbye")
	suite.Equal("<speak>The cow goes woof</speak>", res.Body.OutputSpeech.SSML,
		"Speak does not mutate the original response")
}

func (suite ResponseSuite) TestSpeakTemplate() {
	run := func(locale string) golexa.Response {
		t := speech.NewTemplate("Hello World {{.Value}}",
			speech.WithTranslation(language.Spanish, "Hola Mundo {{.Value}}"))

		req := golexa.Request{}
		req.Body.Locale = locale
		return golexa.NewResponse(req).SpeakTemplate(t, "Foo")
	}

	res := run("")
	suite.Equal("<speak>Hello World Foo</speak>", res.Body.OutputSpeech.SSML,
		"Should evaluate the default translation when the request's local is blank")

	res = run("en-US")
	suite.Equal("<speak>Hello World Foo</speak>", res.Body.OutputSpeech.SSML,
		"Should evaluate the English translation when request's locale is en-US")

	res = run("es-MX")
	suite.Equal("<speak>Hola Mundo Foo</speak>", res.Body.OutputSpeech.SSML,
		"Should evaluate the Spanish translation when the request's local is es-MX")

	res = run("es-ES")
	suite.Equal("<speak>Hola Mundo Foo</speak>", res.Body.OutputSpeech.SSML,
		"Should evaluate the Spanish translation when the request's local is es-ES")

	res = run("it-IT")
	suite.Equal("<speak>Hello World Foo</speak>", res.Body.OutputSpeech.SSML,
		"Should evaluate the default translation when the request's local is it-IT")
}

func (suite ResponseSuite) TestSimpleCard() {
	run := func(title, text string) golexa.Response {
		return golexa.NewResponse(golexa.Request{}).SimpleCard(title, text)
	}

	res := run("", "")
	suite.Equal("Simple", res.Body.Card.Type,
		"Should set the card to 'Simple'")
	suite.Equal("", res.Body.Card.Title,
		"Should allow you to leave the title blank")
	suite.Equal("", res.Body.Card.Content,
		"Should allow you to leave the content blank")

	res = run("Hello", "World, Bro")
	suite.Equal("Simple", res.Body.Card.Type,
		"Should set the card to 'Simple'")
	suite.Equal("Hello", res.Body.Card.Title,
		"Should set the card title to the first argument")
	suite.Equal("World, Bro", res.Body.Card.Content,
		"Should set the card content to the second argument")

	res.SimpleCard("FOO", "BAR")
	suite.Equal("Hello", res.Body.Card.Title,
		"Should not mutate the original Response")
	suite.Equal("World, Bro", res.Body.Card.Content,
		"Should not mutate the original Response")
}

func (suite ResponseSuite) TestElicitSlot() {
	run := func(targetIntentName, slotName string, slots golexa.Slots) golexa.Response {
		req := golexa.NewIntentRequest("Foo", slots)
		return golexa.NewResponse(req).
			EndSession(true).
			ElicitSlot(targetIntentName, slotName)
	}

	res := run("", "name", golexa.Slots{})
	suite.True(*res.Body.ShouldEndSession,
		"Should leave 'ShouldEndSession' alone if we gave a blank intent name")
	suite.Len(res.Body.Directives, 0,
		"Should have 0 directives on the response with a blank intent name")

	res = run("Foo", "", golexa.Slots{})
	suite.True(*res.Body.ShouldEndSession,
		"Should leave 'ShouldEndSession' alone if we gave a blank slot name")
	suite.Len(res.Body.Directives, 0,
		"Should have 0 directives on the response with a blank slot name")

	res = run("Foo", "name", golexa.Slots{})
	suite.Require().Len(res.Body.Directives, 1,
		"Should have 0 directives on the response with a blank slot name")
	name, ok := res.Body.Directives[0].UpdatedIntent.Slots["name"]
	suite.Equal("", name.Value,
		"Directive's slots should include a blank 'name' entry even if request's slots were 0-length")
	suite.True(ok,
		"Directive's slots should include a blank 'name' entry even if request's slots were 0-length")

	// Standard, ok case.
	res = run("Foo", "name", golexa.Slots{
		"name": golexa.Slot{Value: "Bob"},
		"age":  golexa.Slot{Value: "99"},
	})
	suite.False(*res.Body.ShouldEndSession,
		"Should keep the session open even if it was previously set to end it")
	suite.Require().Len(res.Body.Directives, 1,
		"Should have 1 directive on the response")
	suite.Equal("Dialog.ElicitSlot", res.Body.Directives[0].Type,
		"Directive should ben a 'Dialog.ElicitSlot' type")
	suite.Equal("name", res.Body.Directives[0].SlotToElicit,
		"Directive should elicit the 'name' slot")
	suite.Equal("Foo", res.Body.Directives[0].UpdatedIntent.Name,
		"Directive should redirect back to the Foo intent")
	suite.Equal("NONE", res.Body.Directives[0].UpdatedIntent.ConfirmationStatus,
		"Directive have confirmation status of NONE")
	suite.Equal("", res.Body.Directives[0].UpdatedIntent.Slots["name"].Value,
		"Directive should blank out the elicited slot")
	suite.Equal("99", res.Body.Directives[0].UpdatedIntent.Slots["age"].Value,
		"Directive should preserve other slots")
}

func (suite ResponseSuite) TestReprompt() {
	run := func(textOrSSML string) golexa.Response {
		return golexa.NewResponse(golexa.Request{}).Reprompt(textOrSSML)
	}

	res := run("")
	suite.Require().NotNil(res.Body.Reprompt,
		"Should have a non-nil 'Reprompt' attribute")
	suite.Equal("<speak></speak>", res.Body.Reprompt.OutputSpeech.SSML,
		"Should provide empty <speak> SSML if reprompt text is empty")

	res = run("Hello World")
	suite.Require().NotNil(res.Body.Reprompt,
		"Should have a non-nil 'Reprompt' attribute")
	suite.Equal("<speak>Hello World</speak>", res.Body.Reprompt.OutputSpeech.SSML,
		"Should wrap your plain text w/ SSML")

	res = run("<speak>Moo.</speak>")
	suite.Require().NotNil(res.Body.Reprompt,
		"Should have a non-nil 'Reprompt' attribute")
	suite.Equal("<speak>Moo.</speak>", res.Body.Reprompt.OutputSpeech.SSML,
		"Should preserve SSML as-is")
}

func (suite ResponseSuite) TestOk() {
	res, err := golexa.NewResponse(golexa.Request{}).Ok()
	suite.Require().Nil(err,
		"Ok should never return a not-nil error")
	suite.Equal("1.0", res.Version,
		"Ok should return the same response even when it has no meaningful content")

	res, err = golexa.NewResponse(golexa.Request{}).Speak("Moo").Ok()
	suite.Require().Nil(err,
		"Ok should never return a not-nil error")
	suite.Equal("<speak>Moo</speak>", res.Body.OutputSpeech.SSML,
		"Ok should return the same response you've been constructing")

	res, err = golexa.NewResponse(golexa.Request{}).
		Speak("Moo").
		SimpleCard("Cow Sounds", "I said Moo.").
		Ok()
	suite.Require().Nil(err,
		"Ok should never return a not-nil error")
	suite.Equal("<speak>Moo</speak>", res.Body.OutputSpeech.SSML,
		"Ok should return the same response you've been constructing")
	suite.Equal("Cow Sounds", res.Body.Card.Title,
		"Ok should return the same response you've been constructing")
	suite.Equal("I said Moo.", res.Body.Card.Content,
		"Ok should return the same response you've been constructing")
}
