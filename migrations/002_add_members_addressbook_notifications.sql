-- Migration 002: Add members, address_books, and notifications tables
-- Run: psql -U postgres -d vship -f migrations/002_add_members_addressbook_notifications.sql

-- Members table
CREATE TABLE IF NOT EXISTS members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(200) NOT NULL,
    email VARCHAR(200),
    phone VARCHAR(50),
    member_level_id UUID REFERENCES member_levels(id),
    balance DECIMAL(12,2) DEFAULT 0,
    avatar_url TEXT,
    status VARCHAR(50) DEFAULT 'active',
    remark TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trashed_at TIMESTAMP WITH TIME ZONE,
    extra_fields JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_members_tenant_id ON members(tenant_id);
CREATE INDEX IF NOT EXISTS idx_members_member_level_id ON members(member_level_id);
CREATE INDEX IF NOT EXISTS idx_members_status ON members(status);
CREATE INDEX IF NOT EXISTS idx_members_email ON members(email);

-- Address Books table
CREATE TABLE IF NOT EXISTS address_books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    member_id UUID NOT NULL REFERENCES members(id),
    recipient_name VARCHAR(200) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    country_id UUID REFERENCES countries(id),
    province VARCHAR(200),
    city VARCHAR(200),
    district VARCHAR(200),
    address TEXT NOT NULL,
    postal_code VARCHAR(50),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trashed_at TIMESTAMP WITH TIME ZONE,
    extra_fields JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_address_books_tenant_id ON address_books(tenant_id);
CREATE INDEX IF NOT EXISTS idx_address_books_member_id ON address_books(member_id);
CREATE INDEX IF NOT EXISTS idx_address_books_country_id ON address_books(country_id);

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    title VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'system',
    content TEXT NOT NULL,
    target_type VARCHAR(50) DEFAULT 'all',
    target_id UUID,
    status VARCHAR(50) DEFAULT 'draft',
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trashed_at TIMESTAMP WITH TIME ZONE,
    extra_fields JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id ON notifications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
