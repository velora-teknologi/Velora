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
