package hooks

import (
	"errors"
	"fmt"
	serverDomain "github.com/sogko/slumber/domain"
	"log"
	"net/http"
)

func HandlerPostCreateUserHook(w http.ResponseWriter, req *http.Request, ctx serverDomain.IContext, payload interface{}) error {
	if payload == nil {
		return errors.New("Missing payload")
	}

	// cast payload to PostCreateUserHookPayload struct
	p := payload.(*serverDomain.PostCreateUserHookPayload)
	user := p.User
	if user == nil {
		return errors.New("Missing user payload")
	}

	// print confirmation link to console
	// ideally, you could write some logic to send this confirmation link to user's email
	log.Println("Confirmation link:", fmt.Sprintf("/api/users/%v/confirm?code=%v", user.ID.Hex(), user.ConfirmationCode))

	return nil
}
