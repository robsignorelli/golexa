package main

import (
	"context"

	"github.com/robsignorelli/golexa"
	"github.com/robsignorelli/golexa/middleware"
	"github.com/robsignorelli/golexa/sample"
)

func main() {
	skill := golexa.Skill{}
	registerSkillIntents(&skill)
	registerAmazonIntents(&skill)
	golexa.Start(skill)
}

func registerSkillIntents(skill *golexa.Skill) {
	mw := golexa.Middleware{
		middleware.Logger(
			middleware.LogRequestJSON(),
			middleware.LogResponseSpeech()),
		sample.ValidateUser,
	}
	todo := sample.NewTodoService()
	skill.RouteIntent(sample.IntentAddTodoItem, mw.Then(todo.Add))
	skill.RouteIntent(sample.IntentRemoveTodoItem, mw.Then(todo.Remove))
	skill.RouteIntent(sample.IntentListTodoItems, mw.Then(todo.List))
}

func registerAmazonIntents(skill *golexa.Skill) {
	skill.RouteIntent(golexa.IntentNameCancel, func(_ context.Context, req golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse(req).Speak("Canceling.").Ok()
	})
	skill.RouteIntent(golexa.IntentNameStop, func(_ context.Context, req golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse(req).Ok()
	})
	skill.RouteIntent(golexa.IntentNameHelp, func(_ context.Context, req golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse(req).Speak("You can ask me to add, remove, or list items.").Ok()
	})
	skill.RouteIntent(golexa.IntentNameFallback, func(_ context.Context, req golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse(req).Speak("I'm sorry. This skill doesn't know how to do that.").Ok()
	})
	skill.RouteIntent(golexa.IntentNameNavigateHome, func(_ context.Context, req golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse(req).Ok()
	})
}
