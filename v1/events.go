package oakacs

import (
	"context"
	"errors"
	"time"

	"github.com/rs/xid"
)

type (
	// EventType resprents the kind of events that the ACS may issue.
	EventType           uint8
	eventContextKeyType string
)

const (
	eventContextKeyIP         eventContextKeyType = "ip"
	eventContextKeyPermission eventContextKeyType = "permission"

	// EventTypeUnknown indicates an unexpected event, should be treated as fatal error.
	EventTypeUnknown = iota
	// EventTypeSessionExpired marks a session that was checked past deadline.
	EventTypeSessionExpired
	// EventTypeSessionBreached occurs when an attempt to overtake session was detected.
	EventTypeSessionBreached
	// EventTypeAuthenticationSuccess marks a role being succesffully connected to a session.
	EventTypeAuthenticationSuccess
	// EventTypeAuthenticationFailure marks a role being succesffully connected to a session.
	EventTypeAuthenticationFailure
	// EventTypeAuthorizationAllowed marks a successful allowed Permission matched.
	EventTypeAuthorizationAllowed
	// EventTypeAuthorizationDeniedByPermission marks a matched denied Permission.
	EventTypeAuthorizationDeniedByPermission
	// EventTypeAuthorizationDeniedByDefault marks absence of any matched Permissions.
	EventTypeAuthorizationDeniedByDefault
	// EventTypeMaintenance marks clean up and debugging tasks.
	EventTypeMaintenance
	// EventTypeCriticalRepositoryFailure marks a failed critical modification to a backend repository.
	EventTypeCriticalRepositoryFailure
)

func (e EventType) String() string {
	switch e {
	case EventTypeAuthenticationSuccess:
		return "authenticated"
	case EventTypeAuthenticationFailure:
		return "rejected"
	case EventTypeAuthorizationAllowed:
		return "authorized"
	case EventTypeAuthorizationDeniedByPermission:
		return "denied"
	case EventTypeAuthorizationDeniedByDefault:
		return "denied"
	}
	return "<undocumented-event>"
}

type Event struct {
	ctx context.Context // important for contextual unpacking

	When    time.Time
	Type    EventType
	Session xid.ID
	Role    xid.ID
	Error   error
}

func (e *Event) IP() (string, error) {
	if e.ctx != nil { // TODO: is this needed?
		val := e.ctx.Value(eventContextKeyIP)
		switch ip := val.(type) {
		case string:
			return ip, nil
		}
	}
	return "", errors.New("ip address is not associated with context")
}

func (e *Event) Permission() (Permission, error) {
	if e.ctx != nil { // TODO: is this needed?
		val := e.ctx.Value(eventContextKeyPermission)
		switch ip := val.(type) {
		case Permission:
			return ip, nil
		}
	}
	return Permission{}, errors.New("permission is not associated with context")
}

// Broadcast attempts to notify all the subscribers. The dispatch is non-blocking, so if subscriber is busy, the event misses.
func (acs *AccessControlSystem) Broadcast(e Event) {
	e.When = time.Now()
	for _, c := range acs.subscribers {
		select {
		case c <- e:
		default:
		}
	}
}
