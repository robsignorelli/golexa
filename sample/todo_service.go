package sample

import (
	"context"
	"strings"

	"github.com/robsignorelli/golexa"
	"github.com/robsignorelli/golexa/speech"
	"golang.org/x/text/language"
)

const SlotItemName = "item_name"
const IntentAddTodoItem = "AddTodoItem"
const IntentRemoveTodoItem = "RemoveTodoItem"
const IntentListTodoItems = "ListTodoItems"

func NewTodoService() TodoService {
	service := TodoService{
		items: map[string][]string{},
	}

	// When you hit the "AddTodoItem" intent but didn't specify an item name.
	service.templateAddElicit = speech.NewTemplate(
		`What would you like to add to your list?`,
		speech.WithTranslation(language.Spanish, `¿Qué te gustaría agregar a tu lista?`))

	// The success confirmation for the "AddTodoItem" intent
	service.templateAddSuccess = speech.NewTemplate(
		`Okay. I have added "{{.Value}}" to your list.`,
		speech.WithFunc("normalize", strings.ToUpper),
		speech.WithTranslation(language.Spanish, `Bueno. He agregado "{{.Value | normalize}}" a su lista.`))

	// When you hit the "RemoveTodoItem" intent but didn't specify an item name.
	service.templateRemoveElicit = speech.NewTemplate(
		`What would you like to remove from your list?`,
		speech.WithTranslation(language.Spanish, `¿Qué te gustaría eliminar de tu lista?`))

	// The success confirmation for the "RemoveTodoItem" intent
	service.templateRemoveSuccess = speech.NewTemplate(
		`Okay. I have removed "{{.Value}}" from your list.`,
		speech.WithFunc("normalize", strings.ToUpper),
		speech.WithTranslation(language.Spanish, `Bueno. He eliminado "{{.Value | normalize}}" de su lista.`))

	// When you hit the "ListTodoItems" intent, but don't have any items in the list.
	service.templateListEmpty = speech.NewTemplate(
		`Hmm. Your list is empty.`,
		speech.WithTranslation(language.Spanish, `Tu lista esta vacia.`))

	// Confirmation speech for when you hit the "ListTodoItems".
	service.templateListSuccess = speech.NewTemplate(
		`I found {{len .Value}} items in your list. {{range .Value}}{{. | normalize}}. {{end}}`,
		speech.WithFunc("normalize", strings.ToUpper),
		speech.WithTranslation(language.Spanish,
			`Encontré {{len .Value}} artículos en tu lista. {{range .Value}}{{. | normalize}}. {{end}}`))

	return service
}

type TodoService struct {
	// items maps a user id to their personal list of items. Please do not do this in a "real" skill
	// since the resources in your Lambda are ephemeral. Ideally you should be interacting w/ some API
	// or database service to handle long-term storage of user/skill data. This sample is meant to
	// highlight how you can set up the Alexa-specific bits of your skill without worrying too much about
	// how to structure your business logic code.
	items map[string][]string

	templateAddElicit     speech.Template
	templateAddSuccess    speech.Template
	templateRemoveElicit  speech.Template
	templateRemoveSuccess speech.Template
	templateListEmpty     speech.Template
	templateListSuccess   speech.Template
}

// Add appends the item that the user uttered to their personal to-do list. It responds to an
// utterance such as "Add laundry to my to-do list" where "laundry" is the value for the {item_name}
// slot. Additionally, this supports an interaction such as "Update my list" where there is no item
// slot data. In that case, we'll have Alexa ask the user to speak the item name and try this
// intent/action again.
func (service *TodoService) Add(_ context.Context, request golexa.Request) (golexa.Response, error) {
	// Not a failure. Have Alexa prompt the user for what the items should be. Once the user
	// responds, this intent should be re-invoked, but this time with the name filled in.
	itemName := request.Body.Intent.Slots.Resolve(SlotItemName)
	if itemName == "" {
		return golexa.NewResponse(request).
			SpeakTemplate(service.templateAddElicit, nil).
			ElicitSlot(request, IntentAddTodoItem, SlotItemName).
			Ok()
	}

	// Do your "business logic" to handle the user's request.
	userID := request.Session.User.ID
	service.items[userID] = append(service.items[userID], itemName)

	// Have Alexa speak some sort of confirmation.
	return golexa.NewResponse(request).
		SpeakTemplate(service.templateAddSuccess, itemName).
		Ok()
}

func (service *TodoService) Remove(_ context.Context, request golexa.Request) (golexa.Response, error) {
	// Not a failure. Have Alexa prompt the user for what the items should be. Once the user
	// responds, this intent should be re-invoked, but this time with the name filled in.
	itemName := request.Body.Intent.Slots.Resolve(SlotItemName)
	if itemName == "" {
		return golexa.NewResponse(request).
			SpeakTemplate(service.templateRemoveElicit, nil).
			ElicitSlot(request, IntentAddTodoItem, SlotItemName).
			Ok()
	}

	service.removeItem(request.Session.User.ID, itemName)
	return golexa.NewResponse(request).
		SpeakTemplate(service.templateRemoveSuccess, itemName).
		Ok()
}

func (service *TodoService) List(_ context.Context, request golexa.Request) (golexa.Response, error) {
	items := service.items[request.Session.User.ID]
	if len(items) == 0 {
		return golexa.NewResponse(request).
			SpeakTemplate(service.templateListEmpty, nil).
			Ok()
	}

	return golexa.NewResponse(request).
		SpeakTemplate(service.templateListSuccess, items).
		Ok()
}

// removeItem provides the "hand-waving" for our business logic to remove an item from
// this user's list in the "database"
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
	if request.Context.System.User.AccessToken == "" {
		return golexa.NewResponse(request).Speak("I'm sorry. Please link your account through the Alexa app.").Ok()
	}
	return next(ctx, request)
}
