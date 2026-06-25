# Deployment Guide

This project can be deployed from one GitHub repository with separate hosting targets:

- Backend on Railway
- Frontend on Vercel
- PostgreSQL and Redis on Railway

Hosted frontend:

```text
https://open-ownership-workflow-assessment.vercel.app/login
```

## Repository Layout

Use one repository and point each platform to the correct root directory.

```text
backend/   Railway service root
frontend/  Vercel project root
```

Commit these files:

```text
docker-compose.yml
.env.example
backend/.env.example
frontend/.env.example
```

Do not commit real `.env` files.

## Railway Backend

1. Create a new Railway project.
2. Add a PostgreSQL database.
3. Add a Redis service.
4. Add a backend service from the GitHub repository.
5. Set the service root directory to:

```text
backend
```

6. Add the backend environment variables:

```env
PORT=8080
DATABASE_URL=<Railway PostgreSQL connection URL>
REDIS_URL=<Railway Redis connection URL>
JWT_SECRET=<long random production secret>
CORS_ORIGIN=<Vercel frontend URL>

SEED_REQUESTER_NAME=Amina Requester
SEED_REQUESTER_EMAIL=requester@example.com
SEED_REQUESTER_PASSWORD=<demo password or reviewer password>
SEED_REVIEWER_NAME=Noah Reviewer
SEED_REVIEWER_EMAIL=reviewer@example.com
SEED_REVIEWER_PASSWORD=<demo password or reviewer password>
SEED_ADMIN_NAME=Sam Admin
SEED_ADMIN_EMAIL=admin@example.com
SEED_ADMIN_PASSWORD=<demo password or reviewer password>
```

7. Deploy the service.
8. Confirm the health endpoint responds:

```text
https://your-railway-api-url/healthz
```

## Vercel Frontend

1. Create a new Vercel project from the same GitHub repository.
2. Set the root directory to:

```text
frontend
```

3. Add the frontend environment variable:

```env
VITE_API_URL=<Railway backend URL>
```

Example:

```env
VITE_API_URL=https://open-ownership-api.up.railway.app
```

4. Deploy the frontend.
5. Copy the final Vercel URL and set it as `CORS_ORIGIN` on Railway.
6. Redeploy the Railway backend after changing `CORS_ORIGIN`.

## Local Docker Run

For local review:

```bash
docker compose up --build
```

Frontend:

```text
http://localhost:5173
```

Backend:

```text
http://localhost:18080
```

## Common Issues

### Railway says a directory does not exist

Check the Railway service root directory. It should be:

```text
backend
```

not `backed` or the repository root.

### Vercel shows a blank page

Check:

- Vercel root directory is `frontend`.
- `VITE_API_URL` points to the Railway backend URL.
- The frontend was redeployed after changing env variables.
- SPA routing is enabled through `frontend/vercel.json`.

### Login fails after deployment

Check:

- `DATABASE_URL` points to the Railway PostgreSQL database.
- `JWT_SECRET` is set.
- Seed credentials in Railway match the credentials being tested.
- The backend deployed successfully after env changes.

### Browser blocks API requests

Check:

- Railway `CORS_ORIGIN` exactly matches the Vercel frontend URL.
- Include the protocol, for example `https://...vercel.app`.
- Redeploy Railway after changing `CORS_ORIGIN`.

## Reviewer Accounts

The default reviewer-friendly accounts are:

```text
requester@example.com / password
reviewer@example.com / password
admin@example.com / password
```

These can be overridden with the `SEED_*` environment variables listed above.
