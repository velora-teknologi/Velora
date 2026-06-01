# API Standards

Base URL

/api/v1

Response Format

{
  "success": true,
  "message": "success",
  "data": {}
}

Error Format

{
  "success": false,
  "message": "error",
  "errors": []
}

Authentication

Authorization: Bearer TOKEN

---

## Auth / User

POST /api/v1/auth/login
- Request: `{ "email": "user@example.com", "password": "secret" }`
- Response: `{ "success": true, "data": { "token": "<jwt>", "user": { ... } } }`

POST /api/v1/users
- Request: `{ "name": "Admin User", "email": "admin@example.com", "password": "securePassword" }`
- Response: create a new user and return user details

GET /api/v1/users/:id
- Protected: requires Authorization header
- Response: user detail data

GET /api/v1/users
- Protected: requires Authorization header
- Response: paginated user list

PUT /api/v1/users/:id
- Protected: requires Authorization header
- Request: `{ "name": "Updated Name", "email": "updated@example.com" }`

DELETE /api/v1/users/:id
- Protected: requires Authorization header

---

## Agents

GET /api/v1/agents
- Protected
- Response: list of agents belonging to the authenticated user

POST /api/v1/agents
- Protected
- Request: `{ "name": "My Agent", "description": "Agent description", "config": { ... } }`
- Response: created agent details

GET /api/v1/agents/:id
- Protected
- Response: agent detail data

GET /api/v1/agents/:id/workflows
- Protected
- Response: list of workflows for the selected agent
- Query params: `skip`, `limit`

PUT /api/v1/agents/:id
- Protected
- Request: `{ "name": "Updated Agent", "description": "Updated description", "config": { ... } }`

DELETE /api/v1/agents/:id
- Protected
- Response: empty body with success status

---

## Workflows

GET /api/v1/workflows
- Protected
- Response: list of workflows belonging to the authenticated user

POST /api/v1/workflows
- Protected
- Request: `{ "name": "My Workflow", "description": "Workflow description", "agent_id": "<agent-id>", "definition": { ... } }`
- Response: created workflow details

GET /api/v1/workflows/:id
- Protected
- Response: workflow detail data

PUT /api/v1/workflows/:id
- Protected
- Request: `{ "name": "Updated Workflow", "description": "Updated description", "agent_id": "<agent-id>", "definition": { ... } }`

DELETE /api/v1/workflows/:id
- Protected

---

## Common Headers

- `Content-Type: application/json`
- `Authorization: Bearer <token>`

## Pagination

Optional query parameters for list endpoints:
- `skip`: offset for pagination
- `limit`: max number of items returned

---

## Tenants

POST /api/v1/tenants
- Protected: requires Authorization header
- Request: `{ "name": "Tenant Name", "description": "Optional description" }`
- Response: created tenant details

GET /api/v1/tenants
- Protected
- Response: paginated tenant list for the authenticated owner

GET /api/v1/tenants/:id
- Protected
- Response: tenant details

PUT /api/v1/tenants/:id
- Protected
- Request: `{ "name": "Updated Name", "description": "Updated description" }`

DELETE /api/v1/tenants/:id
- Protected

GET /api/v1/tenants/:id/members
- Protected: requires Authorization header
- Description: list members (users) that belong to a tenant. Only the tenant owner can list members.
- Query params: `skip`, `limit`
- Response example:

```
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": "user-1",
      "email": "user1@example.com",
      "full_name": "User One",
      "role": "member",
      "is_active": true
    }
  ]
}
```

POST /api/v1/tenants/:id/members
- Protected: requires Authorization header
- Description: add a user to the tenant (by owner/admin). Request body accepts an existing `user_id` and `role`.
- Request example:

```
POST /api/v1/tenants/tenant-1/members
Content-Type: application/json
Authorization: Bearer <token>

{
  "user_id": "user-2",
  "role": "member"
}
```

- Response example:

```
{
  "success": true,
  "message": "Member added",
  "data": {
    "id": "member-123",
    "user_id": "user-2",
    "role": "member",
    "created_at": "2026-06-01T12:00:00Z"
  }
}
```

DELETE /api/v1/tenants/:id/members/:member_id
- Protected: requires Authorization header
- Description: remove a member from the tenant (owner/admin only)
- Response example:

```
{
  "success": true,
  "message": "Member removed",
  "data": null
}
```

---

## Teams Membership

POST /api/v1/teams/:id/members
- Protected: requires Authorization header
- Description: add an existing user to a team within the same tenant.
- Request example:

```
POST /api/v1/teams/team-1/members
Content-Type: application/json
Authorization: Bearer <token>

{
  "user_id": "user-3",
  "role": "member"
}
```

- Response example (created):

```
{
  "success": true,
  "message": "Team member added",
  "data": {
    "id": "tm-456",
    "team_id": "team-1",
    "user_id": "user-3",
    "role": "member"
  }
}
```

GET /api/v1/teams/:id/members
- Protected: requires Authorization header
- Query params: `skip`, `limit`
- Response example: similar to tenant members list

DELETE /api/v1/teams/:id/members/:member_id
- Protected: requires Authorization header
- Response: success message

---

## Workspaces Membership

POST /api/v1/workspaces/:id/members
- Protected: requires Authorization header
- Request example:

```
POST /api/v1/workspaces/workspace-1/members
Content-Type: application/json
Authorization: Bearer <token>

{
  "user_id": "user-4",
  "role": "editor"
}
```

- Response example:

```
{
  "success": true,
  "message": "Workspace member added",
  "data": {
    "id": "wm-789",
    "workspace_id": "workspace-1",
    "user_id": "user-4",
    "role": "editor"
  }
}
```

GET /api/v1/workspaces/:id/members
- Protected: requires Authorization header
- Query params: `skip`, `limit`
- Response example: similar to tenant members list

DELETE /api/v1/workspaces/:id/members/:member_id
- Protected: requires Authorization header
- Response: success message
