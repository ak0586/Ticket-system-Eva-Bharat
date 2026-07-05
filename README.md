# Ticket System

A simple ticket management REST API built with Go, Gin, and SQLite.

## Local Run

```bash
# copy env file and fill in values
cp .env.example .env

go run main.go
```

## Docker Run

```bash
docker build -t ticket-system .
docker run -p 8080:8080 -e JWT_SECRET=your_secret ticket-system
```

## Health Check

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /health | No | Health check |
| POST | /auth/register | No | Register a new user |
| POST | /auth/login | No | Login and get JWT |
| POST | /tickets | Yes | Create a ticket |
| GET | /tickets | Yes | List your tickets |
| GET | /tickets/:id | Yes | Get a specific ticket |
| PATCH | /tickets/:id/status | Yes | Update ticket status |

### Example Requests

**Register**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

**Login**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

**Create Ticket**
```bash
curl -X POST http://localhost:8080/tickets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Login page broken","description":"Users cannot log in on mobile"}'
```

**Update Status**
```bash
curl -X PATCH http://localhost:8080/tickets/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"status":"in_progress"}'
```

## Status Flow

```
open → in_progress → closed
```

Closed tickets cannot be reopened.

## Deployment

Deployed on Render: `https://ticket-system-eva-bharat.onrender.com`

## Assumptions

- A user can only view and update their own tickets
- Accessing another user's ticket returns 404 (not 403) to avoid leaking ticket existence
- SQLite is used for storage — the DB file is persisted in the container
- JWT tokens expire after 24 hours
