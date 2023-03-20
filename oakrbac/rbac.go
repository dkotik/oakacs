/*

Package oakrbac is a simple flexible Role-Based Access Control (RBAC) implementation.

A role is constructed from a set of policies that are sequentially evaluated until one of the policies returns an [Allow] or a [Deny] sentinel value or an error.

# Usage

OakRBAC leans on [context.Context] as the main mechanism for passing access rights through the execution stack.

	// 1. Initialize the [RBAC]:
	var RBAC = oakrbac.Must(oakrbac.New(
		oakrbac.WithNewRole("administrator", oakrbac.AllowEverything)
	))

	// 2. Inject authorization context:
	ctx := RBAC.ContextWithRole("administrator", context.TODO())

	// 3. Authorize an action [Intention]:
	matchedPolicy, err := oakrbac.Authorize(ctx, &oakrbac.Intention{
		Action: oakrbac.ActionCreate,
		ResourcePath: oakrbac.NewResourcePath(
			"myService",
			"user",
			"userUUID",
		)
	})

	// 4. Act on authorization result:
	if err != nil {
		// access denied, log it using [AuthorizationError.Message] method:
		log.Println(err.Message())
		return err
	}

	// when err == nil, [AuthorizationError.Message] method returns "access granted"
	log.Println(err.Message())

# Policies

OakRBAC comes with only two default policies: [AllowEverything] and [DenyEverything]. You will write or generate policies to match your domain logic.

# Predicates

An [Intention] can be created with a set of [Predicate] functions that allow a [Policy] to run code snippets against the resource to examine it during evaluation.

Predicates enable writing incredibly powerful and performant access control policies.

*/
package oakrbac

import (
	"context"
	"errors"
	"fmt"
)

// RBAC is a simple Role Based Access Control system.
type RBAC struct {
	roleRepository       RoleRepository
	contextRoleExtractor ContextRoleExtractor
	listeners            []Listener
}

// GetRole returns a [Role] that matches provided name.
func (r *RBAC) GetRole(name string) (Role, error) {
	return r.roleRepository.GetRole(name)
}

// Authorize matches the named [Role] against an [Intention]. It returns the [Policy] that granted authorization. The second return value is [AuthorizationError] in place of a generic error.
func (r *RBAC) Authorize(ctx context.Context, roleName string, i Intention) error {
	role, err := r.GetRole(roleName)
	if err != nil {
		r.AuthorizationFailed(ctx, i, nil, nil, err)
		return &AuthorizationError{cause: err}
	}
	policy, err := role.Authorize(ctx, i)
	if errors.Is(err, Allow) {
		r.AuthorizationGranted(ctx, []Intention{i}, []Policy{policy}, role)
		return nil
	} else if errors.Is(err, Deny) {
		r.AuthorizationDenied(ctx, []Intention{i}, []Policy{policy}, role)
		return err
	}
	r.AuthorizationFailed(ctx, i, policy, role, err)
	return &AuthorizationError{
		policy: policy,
		cause:  err,
	}
}

func (r *RBAC) AuthorizeEach(ctx context.Context, roleName string, intents ...Intention) (err error) {
	if len(intents) == 0 {
		err = errors.New("cannot authorize an empty list of intents")
		r.AuthorizationFailed(ctx, nil, nil, nil, err)
		return &AuthorizationError{cause: err}
	}
	role, err := r.GetRole(roleName)
	if err != nil {
		r.AuthorizationFailed(ctx, intents[0], nil, nil, err)
		return &AuthorizationError{cause: err}
	}
	policies := make([]Policy, len(intents))
	for i, intent := range intents {
		policy, err := role.Authorize(ctx, intent)
		if errors.Is(err, Allow) {
			policies[i] = policy
			continue
		} else if errors.Is(err, Deny) {
			r.AuthorizationDenied(ctx, intents, []Policy{policy}, role)
			return Deny
		}
		r.AuthorizationFailed(ctx, intent, policy, role, err)
		return &AuthorizationError{
			policy: policy,
			cause:  err,
		}
	}
	r.AuthorizationGranted(ctx, intents, policies, role)
	return nil
}

func (r *RBAC) AuthorizeAny(ctx context.Context, roleName string, intents ...Intention) error {
	role, err := r.GetRole(roleName)
	if err != nil {
		r.AuthorizationFailed(ctx, nil, nil, nil, err)
		return &AuthorizationError{cause: err}
	}
	var policies []Policy
	for _, intent := range intents {
		policy, err := role.Authorize(ctx, intent)
		if errors.Is(err, Allow) {
			r.AuthorizationGranted(ctx, intents, []Policy{policy}, role)
			return nil
		} else if errors.Is(err, Deny) {
			policies = append(policies, policy)
			continue
		}
		r.AuthorizationFailed(ctx, intent, policy, role, err)
		return &AuthorizationError{
			policy: policy,
			cause:  err,
		}
	}
	r.AuthorizationDenied(ctx, intents, policies, role)
	return Deny
}

// New builds an [RBAC] using provided [Option] set.
func New(withOptions ...Option) (rbac *RBAC, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		WithDefaultOptions(),
		func(o *options) (err error) { // validate
			for _, r := range o.roles {
				if err = o.roleRepository.AddRole(r); err != nil {
					return err
				}
			}
			if o.roleRepository.CountRoles() == 0 {
				return errors.New("at least one role is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create OakRBAC: %w", err)
		}
	}
	return &RBAC{
		roleRepository:       o.roleRepository,
		contextRoleExtractor: o.contextRoleExtractor,
		listeners:            o.listeners,
	}, nil
}

// Must panics if an error is associated with [RBAC] constructor. Use together with [New].
func Must(r *RBAC, err error) *RBAC {
	if err != nil {
		panic(err)
	}
	return r
}
