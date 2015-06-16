package hooks

import (
	"errors"
	"fmt"
	"github.com/sogko/slumber-users"
	"log"
	"net/http"
)

func HandlerPostCreateUserHook(resource *users.Resource, w http.ResponseWriter, req *http.Request, payload *users.PostCreateUserHookPayload) error {

	if payload == nil || payload.User == nil {
		return errors.New("Missing user payload")
	}

	user := payload.User.(*users.User)
	// print confirmation link to console
	// ideally, you could write some logic to send this confirmation link to user's email
	log.Println("Confirmation link:", fmt.Sprintf("/api/users/%v/confirm?code=%v", user.ID.Hex(), user.ConfirmationCode))

	return nil
}
