package sample

import "errors"

// ErrItemNotFound indicates a failure due to the desired item not being in the list. Duh....
var ErrItemNotFound = errors.New("item not found")

// NewTodoRepository creates a new facade for interacting with our fake database. This is just
// an in memory map of "userID->list" which is horrible for a skill since your Lambda's storage
// is ephemeral. The point of this sample is not to show you how to access databases from within
// lambda code - it's to show you how to properly structure your code for a clean skill implementation.
func NewTodoRepository() TodoRepository {
	return TodoRepository{items: map[string][]string{}}
}

// TodoRepository provides our fake database interactions.
type TodoRepository struct {
	items map[string][]string
}

// GetItems fetches the list of items for the given user.
func (r *TodoRepository) GetItems(userID string) []string {
	return r.items[userID]
}

// AddItem adds the specified item to the end of the user's list.
func (r *TodoRepository) AddItem(userID, itemName string) {
	r.items[userID] = append(r.items[userID], itemName)
}

// RemoveItem provides the "hand-waving" for our business logic to remove an item from
// this user's list in the "database".
func (r *TodoRepository) RemoveItem(userID, itemName string) error {
	items := r.items[userID]
	if items == nil {
		return ErrItemNotFound
	}
	for i, userItem := range items {
		if userItem == itemName {
			r.items[userID] = append(items[:i], items[i+1:]...)
			return nil
		}
	}
	return ErrItemNotFound
}
