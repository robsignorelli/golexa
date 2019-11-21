package sample

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	// Not a failure. Have Alexa prompt the user for what the items should be. Once the user
	// responds, this intent should be re-invoked, but this time with the name filled in.
	itemName := request.Body.Intent.Slots.Resolve(SlotItemName)
	if itemName == "" {
		return golexa.NewResponse().
			Speak("What would you like to add to your list?").
			ElicitSlot(request, IntentAddTodoItem, SlotItemName).
			Ok()
	}

	userID := request.Session.User.ID
	service.items[userID] = append(service.items[userID], itemName)
	return golexa.NewResponse().
		Speak(fmt.Sprintf(`Okay. I have added "%s" to your list.`, itemName)).
		Ok()
}

func (service *TodoService) Remove(ctx context.Context, request golexa.Request) (golexa.Response, error) {
	// Not a failure. Have Alexa prompt the user for what the items should be. Once the user
	// responds, this intent should be re-invoked, but this time with the name filled in.
	itemName := request.Body.Intent.Slots.Resolve(SlotItemName)
	if itemName == "" {
		return golexa.NewResponse().
			Speak("What would you like to remove from your list?").
			ElicitSlot(request, IntentAddTodoItem, SlotItemName).
			Ok()
	}

	service.removeItem(request.Session.User.ID, itemName)
	return golexa.NewResponse().
		Speak(fmt.Sprintf(`Okay. I have added "%s" to your list.`, itemName)).
		Ok()
}

func (service *TodoService) List(ctx context.Context, request golexa.Request) (golexa.Response, error) {
	items := service.items[request.Session.User.ID]
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

// ValidateUser ensures that the user/device has gone through "Account Linking" so we've linked
// the Amazon user to a user in our system. When they're linked, Amazon will provide an OAuth2 access
// token with the request, so we can validate the existence of that field.
func ValidateUser(ctx context.Context, request golexa.Request, next golexa.HandlerFunc) (golexa.Response, error) {
	if request.Session.User.AccessToken == "" {
		return golexa.NewResponse().Speak("I'm sorry. Please link your account through the Alexa app.").Ok()
	}
	return next(ctx, request)
}

// LogRequest prints some JSON logging that includes the current time and the incoming request JSON. It also
// writes a second log line once the request is complete, outputting how long the request took.
func LogRequest(ctx context.Context, request golexa.Request, next golexa.HandlerFunc) (golexa.Response, error) {
	fmt.Println(fmt.Sprintf(`{"logger": "golexa", "timestamp":"%v", "requestId": %v", "payload": %+v}`,
		time.Now().Format(time.RFC3339),
		request.Body.RequestID,
		request))

	// The call to next() doesn't need to be the last line, so we can do more work after the "real"
	// request handling work has been done.
	startTime := time.Now()
	response, err := next(ctx, request)
	elapsed := time.Now().Sub(startTime)

	fmt.Println(fmt.Sprintf(`{"logger": "golexa", "timestamp":"%v", "requestId": %v", "elapsed": "%v"}`,
		time.Now().Format(time.RFC3339),
		request.Body.RequestID,
		elapsed))

	return response, err
}
