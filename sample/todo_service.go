package sample

import (
	"context"
	"strings"

	"github.com/robsignorelli/golexa"
	"github.com/robsignorelli/golexa/speech"
	"golang.org/x/text/language"
)

// SlotItemName is the name of the slot where users specify items to add/remove.
const SlotItemName = "item_name"

// IntentAddTodoItem is the name of the intent where we add new items to the user's list
const IntentAddTodoItem = "AddTodoItem"

// IntentRemoveTodoItem is the name of the intent where we remove items from the user's list
const IntentRemoveTodoItem = "RemoveTodoItem"

// IntentListTodoItems is the name of the intent where we have Alexa rattle off all of a user's items
const IntentListTodoItems = "ListTodoItems"

// NewTodoService creates a controller/service that handles all of the intents related to managing
// your items list.
func NewTodoService(repository TodoRepository) TodoService {
	service := TodoService{repository: repository}

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

	// When you hit the "RemoveTodoItem" intent but we couldn't find an item w/ that name.
	service.templateRemoveNotFound = speech.NewTemplate(
		`I'm sorry. I couldn't find "{{.Value}}" in the list.`,
		speech.WithTranslation(language.Spanish, `Lo siento. No pude encontrar "{{.Value}}" en la lista.`))

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

// TodoService wrangles all of of the dependencies for our list management business logic as well as our
// handlers and response templates for the various interactions we support.
type TodoService struct {
	repository TodoRepository

	templateAddElicit      speech.Template
	templateAddSuccess     speech.Template
	templateRemoveElicit   speech.Template
	templateRemoveSuccess  speech.Template
	templateRemoveNotFound speech.Template
	templateListEmpty      speech.Template
	templateListSuccess    speech.Template
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
	service.repository.AddItem(request.Session.User.ID, itemName)

	// Have Alexa speak some sort of confirmation.
	return golexa.NewResponse(request).
		SpeakTemplate(service.templateAddSuccess, itemName).
		Ok()
}

// Remove obviously removes an item from the user's list who made the utterance. Just like Add(),
// this will have Alexa ask the user to specify which item they want to remove.
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

	// Do your "business logic" to handle the user's request.
	if err := service.repository.RemoveItem(request.Session.User.ID, itemName); err == ErrItemNotFound {
		return golexa.NewResponse(request).
			SpeakTemplate(service.templateRemoveNotFound, itemName).
			Ok()
	}

	// Have Alexa speak some sort of confirmation.
	return golexa.NewResponse(request).
		SpeakTemplate(service.templateRemoveSuccess, itemName).
		Ok()
}

// List simply has Alexa rattle off ALL of the items on your list. This is just a sample skill, so
// this would be a terrible experience if the the list were any longer than 3 or 4 items.
func (service *TodoService) List(_ context.Context, request golexa.Request) (golexa.Response, error) {
	items := service.repository.GetItems(request.Session.User.ID)
	if len(items) == 0 {
		return golexa.NewResponse(request).
			SpeakTemplate(service.templateListEmpty, nil).
			Ok()
	}

	return golexa.NewResponse(request).
		SpeakTemplate(service.templateListSuccess, items).
		Ok()
}
