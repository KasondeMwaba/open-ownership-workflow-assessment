# Architecture

## Backend

The Go API separates transport, persistence, authentication, and workflow rules:

- `internal/handlers` owns HTTP routing, middleware, request validation, and responses.
- `internal/services` owns application orchestration such as login, submission operations, and cached dashboard metrics.
- `internal/repositories` owns GORM persistence, transactions, and row-locking queries.
- `internal/models` owns shared response/domain data shapes.
- `internal/workflow` owns status and role transition logic.
- `internal/auth` owns password verification and JWT handling.
- `internal/database` and `internal/cache` own GORM/Postgres and Redis client setup.

The main business invariant is enforced in `Repository.TransitionSubmission`: GORM opens a transaction, locks the submission row with `clause.Locking{Strength: "UPDATE"}`, checks the actor, validates the transition, updates the submission, and inserts the audit event in the same transaction. The service layer invalidates dashboard cache after successful writes.

Models use a small GORM-friendly `BaseModel` for shared `id` and `createdAt` fields. Records with update timestamps embed `UpdatableModel`. The base model avoids soft-delete fields because the current schema does not include `deleted_at`.

Schema migration follows the same broad pattern as the reference auth service: the API runs ordered GORM `AutoMigrate` calls during startup. Demo records are seeded in Go after migration, so the database is no longer initialized from SQL files mounted into the Postgres container.

## Frontend

The React app is designed as an operational review console:

- requester/admin users can create and edit draft or changes-required submissions;
- reviewer/admin users can approve, reject, or request changes on submitted items;
- all users can inspect the audit timeline for submissions they can access.

The browser mirrors the transition map only to hide unavailable buttons. The API remains authoritative.

## Data Model

- `users`: demo identities and roles.
- `submissions`: current state and JSON business payload.
- `audit_events`: append-only history of creation, edits, and decisions.

## Production Considerations

Next hardening steps would include refresh tokens, rate limiting, structured migration tooling, e2e tests, and finer-grained reviewer assignment rules. They are left out to keep the assessment focused and reviewable.
