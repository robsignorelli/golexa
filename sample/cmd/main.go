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
	todoService := sample.NewTodoService()
	skill.RouteIntent(sample.IntentAddTodoItem, todoService.Add)
	skill.RouteIntent(sample.IntentRemoveTodoItem, todoService.Remove)
	skill.RouteIntent(sample.IntentListTodoItems, todoService.List)
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
