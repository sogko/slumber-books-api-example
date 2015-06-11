package hooks

import (
	serverDomain "github.com/sogko/slumber/domain"
)

var HooksMap = serverDomain.ControllerHooksMap{
	PostCreateUserHook: HandlerPostCreateUserHook,
}
