# Database Schema

## users

id UUID PK

tenant_id UUID

email VARCHAR

password_hash TEXT

created_at TIMESTAMP

updated_at TIMESTAMP

deleted_at TIMESTAMP

---

## tenants

id UUID PK

name VARCHAR

slug VARCHAR

plan VARCHAR

created_at TIMESTAMP

updated_at TIMESTAMP

---

## workflows

id UUID PK

tenant_id UUID

name VARCHAR

status VARCHAR

created_at TIMESTAMP

updated_at TIMESTAMP

---

## agents

id UUID PK

tenant_id UUID

name VARCHAR

provider VARCHAR

model VARCHAR

created_at TIMESTAMP

updated_at TIMESTAMP

---

## teams

id UUID PK

tenant_id UUID FK -> tenants.id

name VARCHAR

description TEXT

created_at TIMESTAMP

updated_at TIMESTAMP

---

## workspaces

id UUID PK

tenant_id UUID FK -> tenants.id

name VARCHAR

description TEXT

created_at TIMESTAMP

updated_at TIMESTAMP

---

## team_members

id UUID PK

team_id UUID FK -> teams.id

user_id UUID FK -> users.id

role VARCHAR

created_at TIMESTAMP

updated_at TIMESTAMP

---

## workspace_members

id UUID PK

workspace_id UUID FK -> workspaces.id

user_id UUID FK -> users.id

role VARCHAR

created_at TIMESTAMP

updated_at TIMESTAMP
