package sample

import (
	"context"
	"fmt"
	"strings"

	"github.com/robsignorelli/golexa"
)

const SlotItemName = "item_name"
const IntentAddTodoItem = "AddTodoItem"
const IntentRemoveTodoItem = "RemoveTodoItem"
const IntentListTodoItems = "ListTodoItems"

func NewTodoService() TodoService {
	return TodoService{items: map[string][]string{}}
}

type TodoService struct {
	// items maps a user id to their personal list of items. Please do not do this in a "real" skill
	// since the resources in your Lambda are ephemeral. Ideally you should be interacting w/ some API
	// or database service to handle long-term storage of user/skill data. This sample is meant to
	// highlight how you can set up the Alexa-specific bits of your skill without worrying too much about
	// how to structure your business logic code.
	items map[string][]string
}

// Add appends the item that the user uttered to their personal to-do list. It responds to an
// utterance such as "Add laundry to my to-do list" where "laundry" is the value for the {item_name}
// slot. Additionally, this supports an interaction such as "Update my list" where there is no item
// slot data. In that case, we'll have Alexa ask the user to speak the item name and try this
// intent/action again.
func (service *TodoService) Add(ctx context.Context, request golexa.Request) (golexa.Response, error) {
	// Don't "fail hard" - just have Alexa say something to the user indicating why we couldn't do this.
	userID := request.Session.User.ID
	if userID == "" {
		return golexa.NewResponse().Speak("I'm sorry. I don't know who you are.").Ok()
	}

	// Not a failure. Have Alexa prompt the user for what the items should be. Once the user
	// responds, this intent should be re-invoked, but this time with the name filled in.
	itemName := request.Body.Intent.Slots.Resolve(SlotItemName)
	if itemName == "" {
		return golexa.NewResponse().
			Speak("What would you like to add to your list?").
			ElicitSlot(request, IntentAddTodoItem, SlotItemName).
			Ok()
	}

	service.items[userID] = append(service.items[userID], itemName)
	return golexa.NewResponse().
		Speak(fmt.Sprintf(`Okay. I have added "%s" to your list.`, itemName)).
		Ok()
}

func (service *TodoService) Remove(ctx context.Context, request golexa.Request) (golexa.Response, error) {
	// Don't "fail hard" - just have Alexa say something to the user indicating why we couldn't do this.
	userID := request.Session.User.ID
	if userID == "" {
		return golexa.NewResponse().Speak("I'm sorry. I don't know who you are.").Ok()
	}

	// Not a failure. Have Alexa prompt the user for what the items should be. Once the user
	// responds, this intent should be re-invoked, but this time with the name filled in.
	itemName := request.Body.Intent.Slots.Resolve(SlotItemName)
	if itemName == "" {
		return golexa.NewResponse().
			Speak("What would you like to remove from your list?").
			ElicitSlot(request, IntentAddTodoItem, SlotItemName).
			Ok()
	}

	service.removeItem(userID, itemName)
	return golexa.NewResponse().
		Speak(fmt.Sprintf(`Okay. I have added "%s" to your list.`, itemName)).
		Ok()
}

func (service *TodoService) List(ctx context.Context, request golexa.Request) (golexa.Response, error) {
	// Don't "fail hard" - just have Alexa say something to the user indicating why we couldn't do this.
	userID := request.Session.User.ID
	if userID == "" {
		return golexa.NewResponse().Speak("I'm sorry. I don't know who you are.").Ok()
	}

	items := service.items[userID]
	if len(items) == 0 {
		return golexa.NewResponse().Speak("Hmm. Your list is empty.").Ok()
	}

	text := strings.Builder{}
	text.WriteString(fmt.Sprintf("I found %d items in your list.\n", len(items)))
	for _, item := range items {
		text.WriteString(item)
		text.WriteString(".\n") // insert punctuation so Alexa doesn't make it one run-on sentence.
	}
	return golexa.NewResponse().Speak(text.String()).Ok()
}

func (service *TodoService) removeItem(userID, removeMe string) {
	items := service.items[userID]
	if items == nil {
		return
	}
	for i, userItem := range items {
		if userItem == removeMe {
			service.items[userID] = append(items[:i], items[i+1:]...)
			return
		}
	}
}
