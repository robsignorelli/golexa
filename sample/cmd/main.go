package main

import (
	"context"

	"github.com/robsignorelli/golexa"
	"github.com/robsignorelli/golexa/sample"
)

func main() {
	skill := golexa.Skill{}
	registerSkillIntents(&skill)
	registerAmazonIntents(&skill)
	golexa.Start(skill)
}

func registerSkillIntents(skill *golexa.Skill) {
	middleware := golexa.Middleware{
		sample.LogRequest,
		sample.ValidateUser,
	}
	todo := sample.NewTodoService()
	skill.RouteIntent(sample.IntentAddTodoItem, middleware.Then(todo.Add))
	skill.RouteIntent(sample.IntentRemoveTodoItem, middleware.Then(todo.Remove))
	skill.RouteIntent(sample.IntentListTodoItems, middleware.Then(todo.List))
}

func registerAmazonIntents(skill *golexa.Skill) {
	skill.RouteIntent(golexa.IntentNameCancel, func(context.Context, golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse().Speak("Canceling.").Ok()
	})
	skill.RouteIntent(golexa.IntentNameStop, func(context.Context, golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse().Ok()
	})
	skill.RouteIntent(golexa.IntentNameHelp, func(context.Context, golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse().Speak("You can ask me to add, remove, or list items.").Ok()
	})
	skill.RouteIntent(golexa.IntentNameFallback, func(context.Context, golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse().Speak("I'm sorry. This skill doesn't know how to do that.").Ok()
	})
	skill.RouteIntent(golexa.IntentNameNavigateHome, func(context.Context, golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse().Ok()
	})
}
