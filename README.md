# Open Ownership Full-Stack Assessment

Assignment B implementation: a two-sided submission and approval workflow with strict backend state transitions, role-aware UI, and a complete audit trail.

## Stack

- Backend: Go, Echo, GORM, JWT auth
- Frontend: React, TypeScript, Tailwind CSS, Vite
- Data: PostgreSQL
- Cache: Redis for dashboard metrics
- Runtime: Docker Compose
- Migrations: GORM `AutoMigrate` on API startup

## Live Demo

The hosted frontend is available at:

```text
https://open-ownership-workflow-assessment.vercel.app/login
```

Use the demo credentials below to test the requester, reviewer, and admin workflows.

## Quick Start

```bash
docker compose up --build
```

Open the app at `http://localhost:5173`.

The API is exposed on `http://localhost:18080` to avoid common local conflicts with port `8080`.

Docker Compose includes inline development defaults so reviewers can run the project immediately. Example environment files are included for production-style configuration:

- `.env.example`
- `backend/.env.example`
- `frontend/.env.example`

Real secrets should live in `.env` or deployment secret management, not in committed files.

For hosted deployment notes, see [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md).

Demo login details:

For the Requester use below login credentials:
Email:requester@example.com
Password:password 

For the Reviewer use below login credentials:
Email:reviewer@example.com
Password:password 

For the Admin use below login credentials:
Email:admin@example.com
Password:password

These demo accounts are seeded from environment variables. The values above are local defaults for reviewer convenience and can be overridden with `SEED_REQUESTER_*`, `SEED_REVIEWER_*`, and `SEED_ADMIN_*` variables.

## What to Review

This submission implements **Assignment B: Submission & Approval Workflow**.

Please test all three seeded account types. The experience is intentionally different for each role so reviewers can confirm permissions, workflow enforcement, and audit visibility.

### Reviewer Testing Script

#### 1. Requester account

Sign in with:

```text
requester@example.com
password
```

Recommended checks:

1. Confirm the requester can create a new submission.
2. Confirm the requester can edit their own `draft` or `changes_required` submissions.
3. Submit a draft for review.
4. Confirm the requester cannot approve or reject submissions.
5. Confirm dashboard totals only reflect submissions visible to that requester.
6. Open **Audit Trail / Submission Audit** and confirm only visible submission audit records are shown.

Expected result: requester workflow is scoped to their own submissions.

#### 2. Reviewer account

Sign in with:

```text
reviewer@example.com
password
```

Recommended checks:

1. Open **Review Queue**.
2. Open a submitted record.
3. Move it to `Changes required`, `Approved`, or `Rejected`.
4. Confirm the reviewer cannot create users, roles, or permissions.
5. Confirm reviewer actions create submission audit records.

Expected result: reviewer can process submitted records but cannot perform administration.

#### 3. Admin account

Sign in with:

```text
admin@example.com
password
```

Recommended checks:

1. Open **Administration / Role Management**.
2. Create a custom role with selected permissions.
3. Open **Administration / User Management**.
4. Create a user with the custom role.
5. Disable and enable a user account.
6. Open all Audit Trail sections:
   - **Submission Audit**
   - **Activity Audit**
   - **Session Audit**
   - **System Audit**
7. Try one failed login, then return as admin and confirm it appears in **Session Audit**.
8. Move around the app as another account, then return as admin and confirm the movement appears in **Activity Audit**.

Expected result: admin can manage accounts, roles, permissions, and see complete audit evidence.

### End-to-End Scenario

For the clearest review, run this full flow:

1. Admin creates a custom requester-like role with `submissions:create` and `dashboard:view`.
2. Admin creates a user assigned to that role.
3. New custom user signs in and creates a submission.
4. Reviewer signs in and requests changes.
5. Custom user edits and resubmits.
6. Reviewer approves or rejects.
7. Admin signs in and reviews Submission, Activity, Session, and System audit records.

## Architecture

Backend structure:

```text
backend/
  cmd/api/
  internal/
    auth/
    config/
    models/
    workflow/
    repositories/
    services/
    handlers/
    routes/
    dto/
    cache/
    database/
```

The backend separates responsibilities intentionally:

- `handlers`: HTTP parsing, status codes, routing, response shape.
- `routes`: endpoint registration split by API area.
- `dto`: request/response payloads used at the HTTP boundary.
- `services`: business rules, validation, permissions, workflow orchestration.
- `repositories`: GORM/database access and scoped queries.
- `workflow`: status definitions and transition rules.
- `models`: GORM models and JSON response shapes.

Frontend structure:

- `pages`: route-level workflow screens.
- `components`: shared UI primitives such as headers, pagination, status badges, confirmation dialogs.
- `api`: typed API client.
- `utils`: permission and workflow display helpers.

## Roles and Permissions

The seeded roles are:

| Role | Permissions |
| --- | --- |
| admin | all permissions |
| reviewer | `submissions:review`, `dashboard:view` |
| requester | `submissions:create`, `dashboard:view` |

Admins can create additional roles and permissions. Custom users are scoped by permissions:

- A user with `submissions:create` but without `submissions:review` sees only their own submissions and own dashboard totals.
- A user with `submissions:review` can see the review queue.
- Admins see global administration and audit data.

## Workflow Rules

The backend is the source of truth for all status changes.

| Current status | Allowed next statuses | Roles |
| --- | --- | --- |
| draft | submitted, withdrawn | requester, admin |
| changes_required | submitted, withdrawn | requester, admin |
| submitted | changes_required, approved, rejected | reviewer, admin |
| approved | none | none |
| rejected | none | none |
| withdrawn | none | none |

Every creation, edit, and status transition creates an audit event. Transition updates lock the submission row inside a transaction to avoid conflicting reviewer decisions.

The schema is created by the Go API on startup using GORM migrations, and demo data is seeded by Go code after migration.

## Audit Trails

The app records three kinds of audit data:

- **Submission audit**: creation, edits, status transitions.
- **System audit**: user, role, and permission administration.
- **Session audit**: successful login, failed login, disabled-account attempts, and logout.

Session audit captures timestamp, email, matched user where available, IP address, browser, full user agent, result, and reason.

Non-admin users only see audit records for submissions they are allowed to view. Admins see submission, system, and session audit sections.

## Pagination and Filtering

The submission queue uses backend pagination through:

```http
GET /api/submissions?page=1&pageSize=8&status=submitted
```

The response includes:

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "pageSize": 8
}
```

Audit, user, role, and permission screens include UI pagination. Audit Trail also includes search, event type filtering, and session result filtering.

## Useful Commands

```bash
cd backend
go test ./...
```

```bash
cd frontend
npm install
npm run lint
npm run build
```

## API Overview

- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`
- `GET /api/dashboard`
- `GET /api/submissions?page=1&pageSize=8`
- `POST /api/submissions`
- `GET /api/submissions/{id}`
- `PUT /api/submissions/{id}`
- `POST /api/submissions/{id}/transition`
- `GET /api/submissions/{id}/audit`

## Notes for Reviewers

This project intentionally keeps workflow validation in the Go domain layer and repeats only allowed action discovery in the React UI for usability. If the UI is bypassed, invalid transitions still fail at the API.

Redis is used as a small but visible production-style optimization: dashboard metrics are cached for 30 seconds and invalidated after writes.

## Use of AI Tools

I used AI tools during development, primarily ChatGPT/Codex, as a pair-programming assistant.

AI assistance was used for:

- Scaffolding and iterating on backend and frontend structure.
- Reviewing implementation choices such as Echo routing, GORM models, DTOs, services, repositories, and audit trail design.
- Debugging issues during local development and deployment setup.
- Improving UI layout, workflow screens, documentation, and reviewer guidance.
- Suggesting tests and verification steps.

I personally reviewed, edited, and verified the submitted code. In particular, I checked the workflow rules, role and permission behavior, audit trail visibility, seeded account behavior, deployment configuration, and local build/test results. I understand the submitted implementation and can explain the backend services, frontend workflows, audit logging, permissions, and deployment setup.

Key implementation notes:

- GORM `AutoMigrate` is used to keep the assessment easy to run locally and in review environments. For a long-running production system, I would move to versioned migrations.
- Logout is recorded in the session audit trail. For higher-security production deployments, I would add refresh tokens or server-side token revocation.
- The highest-growth workflow list, submissions, uses API pagination. Smaller administration datasets use UI pagination in this assessment, with a clear path to backend pagination if those datasets grow.
