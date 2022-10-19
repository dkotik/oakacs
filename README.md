# Oak Access Control System

## Features

1. Safety: static types, immutability, and proper defaults.
2. Minimalism: tracks the least amount of information possible without compromising safety.
3. Consistency: assumes single fully synchronized source of truth. Revocations are instant.
4. Flexibility: simple, independent, and configurable models that support multiple back-ends.

## Packages

- [oakrbac](pkg.go.dev/github.com/dkotik/oakacs/oakrbac): role-based access control

## Functions

1. Humanity Recognition
2. Throttling
   - Prevent password-reuse?
3. Timing modulation
4. Registration
5. Authentication
   - Password policy
   - Revocation
   - Kill switch
   - Recovery
6. Authorization
7. Observability
   - Logging

## Model

- Identity: provides authentication.
- Group: enumerates roles which are available for identities.
- Role: provides authorization by granular permissions.
  - Permission
    - Service
    - Domain
    - Resource
    - Action
- Session: the result of pairing identity to a role.
- Token: one-time utility codes.

## Security

The library is created in ways that anticipate misconfiguration by aiming at simplicity.

1. All roles and policies deny access by default.
2. Policies and predicates must return explicit sentinel value `Allow`.
3. Comes with a code generation tool that helps build tight access control policies and **test cases**.

## Logging

Logging can be approached in several different ways:

1. By writing a Policy wrapper. Use the function `WithLogger` for an example.
2. Inside the policies themselves.
3. At a higher level with request logs.

## Interesting Access Control Projects

### Authorization

- [goRBAC](https://github.com/mikespook/gorbac)
- [authzed SpiceDB](https://github.com/authzed)
- [Keto](https://github.com/ory/keto)

### Authentication

- [Authelia](https://github.com/authelia/authelia)
