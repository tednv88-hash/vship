-- Migration 003: Add remaining tables for complete cross-border consolidation shipping management system
-- Run: psql -U postgres -d vship -f migrations/003_add_remaining_tables.sql
--
-- This migration adds 48 new tables and extends several existing tables.
-- All statements are idempotent (IF NOT EXISTS / IF NOT EXISTS).

-- ============================================================
-- 1. goods_categories
-- ============================================================
CREATE TABLE IF NOT EXISTS goods_categories (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    parent_id    UUID REFERENCES goods_categories(id),
    image_url    VARCHAR(500),
    sort_order   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_goods_categories_tenant_id ON goods_categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_goods_categories_parent_id ON goods_categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_goods_categories_status ON goods_categories(status);
CREATE INDEX IF NOT EXISTS idx_goods_categories_extra_fields ON goods_categories USING GIN (extra_fields);

-- ============================================================
-- 2. goods
-- ============================================================
CREATE TABLE IF NOT EXISTS goods (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID REFERENCES tenants(id),
    category_id    UUID REFERENCES goods_categories(id),
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    image_url      VARCHAR(500),
    images         JSONB DEFAULT '[]'::jsonb,
    price          NUMERIC(12,2) NOT NULL DEFAULT 0,
    original_price NUMERIC(12,2) DEFAULT 0,
    stock          INT DEFAULT 0,
    sales_count    INT DEFAULT 0,
    unit           VARCHAR(50),
    weight         NUMERIC(10,3) DEFAULT 0,
    sort_order     INT DEFAULT 0,
    status         VARCHAR(50) DEFAULT 'active',
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW(),
    trashed_at     TIMESTAMPTZ,
    extra_fields   JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_goods_tenant_id ON goods(tenant_id);
CREATE INDEX IF NOT EXISTS idx_goods_category_id ON goods(category_id);
CREATE INDEX IF NOT EXISTS idx_goods_status ON goods(status);
CREATE INDEX IF NOT EXISTS idx_goods_extra_fields ON goods USING GIN (extra_fields);

-- ============================================================
-- 3. goods_reviews
-- ============================================================
CREATE TABLE IF NOT EXISTS goods_reviews (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    goods_id     UUID REFERENCES goods(id),
    user_id      UUID REFERENCES users(id),
    order_id     UUID,
    content      TEXT,
    rating       INT DEFAULT 5,
    images       JSONB DEFAULT '[]'::jsonb,
    status       VARCHAR(50) DEFAULT 'pending',
    reply        TEXT,
    replied_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_goods_reviews_tenant_id ON goods_reviews(tenant_id);
CREATE INDEX IF NOT EXISTS idx_goods_reviews_goods_id ON goods_reviews(goods_id);
CREATE INDEX IF NOT EXISTS idx_goods_reviews_user_id ON goods_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_goods_reviews_status ON goods_reviews(status);
CREATE INDEX IF NOT EXISTS idx_goods_reviews_extra_fields ON goods_reviews USING GIN (extra_fields);

-- ============================================================
-- 4. shop_orders
-- ============================================================
CREATE TABLE IF NOT EXISTS shop_orders (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID REFERENCES tenants(id),
    order_no         VARCHAR(100) NOT NULL,
    user_id          UUID REFERENCES users(id),
    status           VARCHAR(50) DEFAULT 'pending_payment',
    total_amount     NUMERIC(12,2) DEFAULT 0,
    discount_amount  NUMERIC(12,2) DEFAULT 0,
    shipping_fee     NUMERIC(12,2) DEFAULT 0,
    pay_amount       NUMERIC(12,2) DEFAULT 0,
    pay_method       VARCHAR(50),
    pay_time         TIMESTAMPTZ,
    delivery_time    TIMESTAMPTZ,
    receive_time     TIMESTAMPTZ,
    remark           TEXT,
    receiver_name    VARCHAR(255),
    receiver_phone   VARCHAR(50),
    receiver_address TEXT,
    express_company  VARCHAR(255),
    express_no       VARCHAR(255),
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    trashed_at       TIMESTAMPTZ,
    extra_fields     JSONB DEFAULT '{}'::jsonb
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_shop_orders_tenant_order_no ON shop_orders(tenant_id, order_no);
CREATE INDEX IF NOT EXISTS idx_shop_orders_tenant_id ON shop_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_shop_orders_user_id ON shop_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_shop_orders_status ON shop_orders(status);
CREATE INDEX IF NOT EXISTS idx_shop_orders_extra_fields ON shop_orders USING GIN (extra_fields);

-- ============================================================
-- 5. shop_order_items
-- ============================================================
CREATE TABLE IF NOT EXISTS shop_order_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    order_id     UUID REFERENCES shop_orders(id),
    goods_id     UUID REFERENCES goods(id),
    goods_name   VARCHAR(255),
    goods_image  VARCHAR(500),
    price        NUMERIC(12,2) DEFAULT 0,
    quantity     INT DEFAULT 1,
    total_amount NUMERIC(12,2) DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_shop_order_items_tenant_id ON shop_order_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_shop_order_items_order_id ON shop_order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_shop_order_items_goods_id ON shop_order_items(goods_id);
CREATE INDEX IF NOT EXISTS idx_shop_order_items_extra_fields ON shop_order_items USING GIN (extra_fields);

-- ============================================================
-- 6. order_refunds
-- ============================================================
CREATE TABLE IF NOT EXISTS order_refunds (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    order_id      UUID REFERENCES shop_orders(id),
    order_item_id UUID,
    user_id       UUID REFERENCES users(id),
    refund_no     VARCHAR(100),
    type          VARCHAR(50) DEFAULT 'refund_only',
    reason        VARCHAR(500),
    description   TEXT,
    images        JSONB DEFAULT '[]'::jsonb,
    amount        NUMERIC(12,2) DEFAULT 0,
    status        VARCHAR(50) DEFAULT 'pending',
    audit_remark  TEXT,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    trashed_at    TIMESTAMPTZ,
    extra_fields  JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_order_refunds_tenant_id ON order_refunds(tenant_id);
CREATE INDEX IF NOT EXISTS idx_order_refunds_order_id ON order_refunds(order_id);
CREATE INDEX IF NOT EXISTS idx_order_refunds_user_id ON order_refunds(user_id);
CREATE INDEX IF NOT EXISTS idx_order_refunds_status ON order_refunds(status);
CREATE INDEX IF NOT EXISTS idx_order_refunds_extra_fields ON order_refunds USING GIN (extra_fields);

-- ============================================================
-- 7. article_categories
-- ============================================================
CREATE TABLE IF NOT EXISTS article_categories (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    sort_order   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_article_categories_tenant_id ON article_categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_article_categories_status ON article_categories(status);
CREATE INDEX IF NOT EXISTS idx_article_categories_extra_fields ON article_categories USING GIN (extra_fields);

-- ============================================================
-- 8. file_groups (created before files so FK can reference it)
-- ============================================================
CREATE TABLE IF NOT EXISTS file_groups (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    sort_order   INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_file_groups_tenant_id ON file_groups(tenant_id);
CREATE INDEX IF NOT EXISTS idx_file_groups_extra_fields ON file_groups USING GIN (extra_fields);

-- ============================================================
-- 9. files
-- ============================================================
CREATE TABLE IF NOT EXISTS files (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    group_id     UUID REFERENCES file_groups(id),
    file_name    VARCHAR(500) NOT NULL,
    file_path    VARCHAR(1000),
    file_url     VARCHAR(1000),
    file_type    VARCHAR(100),
    file_size    BIGINT DEFAULT 0,
    storage_type VARCHAR(50),
    uploader_id  UUID REFERENCES users(id),
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_files_tenant_id ON files(tenant_id);
CREATE INDEX IF NOT EXISTS idx_files_group_id ON files(group_id);
CREATE INDEX IF NOT EXISTS idx_files_uploader_id ON files(uploader_id);
CREATE INDEX IF NOT EXISTS idx_files_extra_fields ON files USING GIN (extra_fields);

-- ============================================================
-- 10. points_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS points_logs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    user_id       UUID REFERENCES users(id),
    value         INT NOT NULL DEFAULT 0,
    balance       INT NOT NULL DEFAULT 0,
    type          VARCHAR(50) NOT NULL DEFAULT 'earn',
    description   TEXT,
    related_id    UUID,
    related_type  VARCHAR(100),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    extra_fields  JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_points_logs_tenant_id ON points_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_points_logs_user_id ON points_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_points_logs_extra_fields ON points_logs USING GIN (extra_fields);

-- ============================================================
-- 11. recharge_orders
-- ============================================================
CREATE TABLE IF NOT EXISTS recharge_orders (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    user_id      UUID REFERENCES users(id),
    order_no     VARCHAR(100) NOT NULL,
    amount       NUMERIC(12,2) NOT NULL DEFAULT 0,
    gift_amount  NUMERIC(12,2) DEFAULT 0,
    pay_method   VARCHAR(50),
    pay_time     TIMESTAMPTZ,
    status       VARCHAR(50) DEFAULT 'pending',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_tenant_id ON recharge_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_user_id ON recharge_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_status ON recharge_orders(status);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_extra_fields ON recharge_orders USING GIN (extra_fields);

-- ============================================================
-- 12. recharge_plans
-- ============================================================
CREATE TABLE IF NOT EXISTS recharge_plans (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    amount       NUMERIC(12,2) NOT NULL DEFAULT 0,
    gift_amount  NUMERIC(12,2) DEFAULT 0,
    gift_points  INT DEFAULT 0,
    sort_order   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_recharge_plans_tenant_id ON recharge_plans(tenant_id);
CREATE INDEX IF NOT EXISTS idx_recharge_plans_status ON recharge_plans(status);
CREATE INDEX IF NOT EXISTS idx_recharge_plans_extra_fields ON recharge_plans USING GIN (extra_fields);

-- ============================================================
-- 13. user_balance_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS user_balance_logs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    user_id       UUID REFERENCES users(id),
    amount        NUMERIC(12,2) NOT NULL DEFAULT 0,
    balance       NUMERIC(12,2) NOT NULL DEFAULT 0,
    type          VARCHAR(50) NOT NULL DEFAULT 'recharge',
    description   TEXT,
    related_id    UUID,
    related_type  VARCHAR(100),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    extra_fields  JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_user_balance_logs_tenant_id ON user_balance_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_balance_logs_user_id ON user_balance_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_balance_logs_extra_fields ON user_balance_logs USING GIN (extra_fields);

-- ============================================================
-- 14. user_discounts
-- ============================================================
CREATE TABLE IF NOT EXISTS user_discounts (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    user_id       UUID REFERENCES users(id),
    route_id      UUID REFERENCES shipping_routes(id),
    discount_rate NUMERIC(5,2) DEFAULT 0,
    remark        TEXT,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    trashed_at    TIMESTAMPTZ,
    extra_fields  JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_user_discounts_tenant_id ON user_discounts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_discounts_user_id ON user_discounts(user_id);
CREATE INDEX IF NOT EXISTS idx_user_discounts_route_id ON user_discounts(route_id);
CREATE INDEX IF NOT EXISTS idx_user_discounts_extra_fields ON user_discounts USING GIN (extra_fields);

-- ============================================================
-- 15. user_marks
-- ============================================================
CREATE TABLE IF NOT EXISTS user_marks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    user_id      UUID REFERENCES users(id),
    mark_name    VARCHAR(255) NOT NULL,
    remark       TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_user_marks_tenant_id ON user_marks(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_marks_user_id ON user_marks(user_id);
CREATE INDEX IF NOT EXISTS idx_user_marks_extra_fields ON user_marks USING GIN (extra_fields);

-- ============================================================
-- 16. coupon_receive_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS coupon_receive_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    coupon_id    UUID REFERENCES coupons(id),
    user_id      UUID REFERENCES users(id),
    status       VARCHAR(50) DEFAULT 'unused',
    used_at      TIMESTAMPTZ,
    order_id     UUID,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_coupon_receive_logs_tenant_id ON coupon_receive_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_coupon_receive_logs_coupon_id ON coupon_receive_logs(coupon_id);
CREATE INDEX IF NOT EXISTS idx_coupon_receive_logs_user_id ON coupon_receive_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_coupon_receive_logs_status ON coupon_receive_logs(status);
CREATE INDEX IF NOT EXISTS idx_coupon_receive_logs_extra_fields ON coupon_receive_logs USING GIN (extra_fields);

-- ============================================================
-- 17. blind_box_activities
-- ============================================================
CREATE TABLE IF NOT EXISTS blind_box_activities (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    image_url    VARCHAR(500),
    cost_points  INT DEFAULT 0,
    prizes       JSONB DEFAULT '[]'::jsonb,
    start_at     TIMESTAMPTZ,
    end_at       TIMESTAMPTZ,
    status       VARCHAR(50) DEFAULT 'active',
    total_draws  INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_blind_box_activities_tenant_id ON blind_box_activities(tenant_id);
CREATE INDEX IF NOT EXISTS idx_blind_box_activities_status ON blind_box_activities(status);
CREATE INDEX IF NOT EXISTS idx_blind_box_activities_extra_fields ON blind_box_activities USING GIN (extra_fields);

-- ============================================================
-- 18. blind_box_draws
-- ============================================================
CREATE TABLE IF NOT EXISTS blind_box_draws (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    activity_id  UUID REFERENCES blind_box_activities(id),
    user_id      UUID REFERENCES users(id),
    prize_info   JSONB DEFAULT '{}'::jsonb,
    cost_points  INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_blind_box_draws_tenant_id ON blind_box_draws(tenant_id);
CREATE INDEX IF NOT EXISTS idx_blind_box_draws_activity_id ON blind_box_draws(activity_id);
CREATE INDEX IF NOT EXISTS idx_blind_box_draws_user_id ON blind_box_draws(user_id);
CREATE INDEX IF NOT EXISTS idx_blind_box_draws_extra_fields ON blind_box_draws USING GIN (extra_fields);

-- ============================================================
-- 19. dealer_applications
-- ============================================================
CREATE TABLE IF NOT EXISTS dealer_applications (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    user_id      UUID REFERENCES users(id),
    real_name    VARCHAR(255),
    phone        VARCHAR(50),
    reason       TEXT,
    status       VARCHAR(50) DEFAULT 'pending',
    audit_remark TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_dealer_applications_tenant_id ON dealer_applications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dealer_applications_user_id ON dealer_applications(user_id);
CREATE INDEX IF NOT EXISTS idx_dealer_applications_status ON dealer_applications(status);
CREATE INDEX IF NOT EXISTS idx_dealer_applications_extra_fields ON dealer_applications USING GIN (extra_fields);

-- ============================================================
-- 20. dealer_levels (created before dealers so FK can reference it)
-- ============================================================
CREATE TABLE IF NOT EXISTS dealer_levels (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID REFERENCES tenants(id),
    name              VARCHAR(255) NOT NULL,
    level_no          INT NOT NULL DEFAULT 1,
    commission_rate_1 NUMERIC(5,2) DEFAULT 0,
    commission_rate_2 NUMERIC(5,2) DEFAULT 0,
    commission_rate_3 NUMERIC(5,2) DEFAULT 0,
    upgrade_condition JSONB DEFAULT '{}'::jsonb,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    trashed_at        TIMESTAMPTZ,
    extra_fields      JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_dealer_levels_tenant_id ON dealer_levels(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dealer_levels_extra_fields ON dealer_levels USING GIN (extra_fields);

-- ============================================================
-- 21. dealers
-- ============================================================
CREATE TABLE IF NOT EXISTS dealers (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID REFERENCES tenants(id),
    user_id              UUID REFERENCES users(id),
    level_id             UUID REFERENCES dealer_levels(id),
    parent_id            UUID REFERENCES dealers(id),
    commission_rate      NUMERIC(5,2) DEFAULT 0,
    total_commission     NUMERIC(12,2) DEFAULT 0,
    withdrawn_commission NUMERIC(12,2) DEFAULT 0,
    status               VARCHAR(50) DEFAULT 'active',
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW(),
    trashed_at           TIMESTAMPTZ,
    extra_fields         JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_dealers_tenant_id ON dealers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dealers_user_id ON dealers(user_id);
CREATE INDEX IF NOT EXISTS idx_dealers_level_id ON dealers(level_id);
CREATE INDEX IF NOT EXISTS idx_dealers_parent_id ON dealers(parent_id);
CREATE INDEX IF NOT EXISTS idx_dealers_status ON dealers(status);
CREATE INDEX IF NOT EXISTS idx_dealers_extra_fields ON dealers USING GIN (extra_fields);

-- ============================================================
-- 22. dealer_orders
-- ============================================================
CREATE TABLE IF NOT EXISTS dealer_orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID REFERENCES tenants(id),
    dealer_id         UUID REFERENCES dealers(id),
    order_id          UUID,
    order_type        VARCHAR(50),
    commission_amount NUMERIC(12,2) DEFAULT 0,
    status            VARCHAR(50) DEFAULT 'pending',
    settled_at        TIMESTAMPTZ,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    extra_fields      JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_dealer_orders_tenant_id ON dealer_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dealer_orders_dealer_id ON dealer_orders(dealer_id);
CREATE INDEX IF NOT EXISTS idx_dealer_orders_status ON dealer_orders(status);
CREATE INDEX IF NOT EXISTS idx_dealer_orders_extra_fields ON dealer_orders USING GIN (extra_fields);

-- ============================================================
-- 23. dealer_withdrawals
-- ============================================================
CREATE TABLE IF NOT EXISTS dealer_withdrawals (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    dealer_id    UUID REFERENCES dealers(id),
    amount       NUMERIC(12,2) NOT NULL DEFAULT 0,
    method       VARCHAR(50) DEFAULT 'bank',
    account_info JSONB DEFAULT '{}'::jsonb,
    status       VARCHAR(50) DEFAULT 'pending',
    audit_remark TEXT,
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_dealer_withdrawals_tenant_id ON dealer_withdrawals(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dealer_withdrawals_dealer_id ON dealer_withdrawals(dealer_id);
CREATE INDEX IF NOT EXISTS idx_dealer_withdrawals_status ON dealer_withdrawals(status);
CREATE INDEX IF NOT EXISTS idx_dealer_withdrawals_extra_fields ON dealer_withdrawals USING GIN (extra_fields);

-- ============================================================
-- 24. admin_roles
-- ============================================================
CREATE TABLE IF NOT EXISTS admin_roles (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    permissions  JSONB DEFAULT '[]'::jsonb,
    sort_order   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_admin_roles_tenant_id ON admin_roles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_admin_roles_status ON admin_roles(status);
CREATE INDEX IF NOT EXISTS idx_admin_roles_extra_fields ON admin_roles USING GIN (extra_fields);

-- ============================================================
-- 25. warehouse_addresses
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_addresses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    warehouse_id UUID REFERENCES warehouses(id),
    label        VARCHAR(255),
    recipient    VARCHAR(255),
    phone        VARCHAR(50),
    country_id   UUID REFERENCES countries(id),
    province     VARCHAR(200),
    city         VARCHAR(200),
    district     VARCHAR(200),
    address      TEXT,
    postal_code  VARCHAR(50),
    is_default   BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_warehouse_addresses_tenant_id ON warehouse_addresses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_addresses_warehouse_id ON warehouse_addresses(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_addresses_extra_fields ON warehouse_addresses USING GIN (extra_fields);

-- ============================================================
-- 26. warehouse_applications
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_applications (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID REFERENCES tenants(id),
    applicant_name VARCHAR(255) NOT NULL,
    phone          VARCHAR(50),
    warehouse_name VARCHAR(255),
    address        TEXT,
    status         VARCHAR(50) DEFAULT 'pending',
    audit_remark   TEXT,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW(),
    trashed_at     TIMESTAMPTZ,
    extra_fields   JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_warehouse_applications_tenant_id ON warehouse_applications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_applications_status ON warehouse_applications(status);
CREATE INDEX IF NOT EXISTS idx_warehouse_applications_extra_fields ON warehouse_applications USING GIN (extra_fields);

-- ============================================================
-- 27. warehouse_clerks
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_clerks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    warehouse_id UUID REFERENCES warehouses(id),
    user_id      UUID REFERENCES users(id),
    name         VARCHAR(255),
    phone        VARCHAR(50),
    role         VARCHAR(50) DEFAULT 'clerk',
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerks_tenant_id ON warehouse_clerks(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerks_warehouse_id ON warehouse_clerks(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerks_user_id ON warehouse_clerks(user_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerks_status ON warehouse_clerks(status);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerks_extra_fields ON warehouse_clerks USING GIN (extra_fields);

-- ============================================================
-- 28. warehouse_clerk_reviews
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_clerk_reviews (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    clerk_id     UUID REFERENCES warehouse_clerks(id),
    user_id      UUID REFERENCES users(id),
    order_id     UUID,
    rating       INT DEFAULT 5,
    content      TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerk_reviews_tenant_id ON warehouse_clerk_reviews(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerk_reviews_clerk_id ON warehouse_clerk_reviews(clerk_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerk_reviews_user_id ON warehouse_clerk_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_clerk_reviews_extra_fields ON warehouse_clerk_reviews USING GIN (extra_fields);

-- ============================================================
-- 29. warehouse_capital_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_capital_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    warehouse_id UUID REFERENCES warehouses(id),
    amount       NUMERIC(12,2) NOT NULL DEFAULT 0,
    balance      NUMERIC(12,2) NOT NULL DEFAULT 0,
    type         VARCHAR(50) NOT NULL DEFAULT 'income',
    description  TEXT,
    related_id   UUID,
    related_type VARCHAR(100),
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_warehouse_capital_logs_tenant_id ON warehouse_capital_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_capital_logs_warehouse_id ON warehouse_capital_logs(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_capital_logs_extra_fields ON warehouse_capital_logs USING GIN (extra_fields);

-- ============================================================
-- 30. warehouse_bonuses
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_bonuses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    warehouse_id UUID REFERENCES warehouses(id),
    amount       NUMERIC(12,2) NOT NULL DEFAULT 0,
    type         VARCHAR(50) DEFAULT 'service',
    description  TEXT,
    month        VARCHAR(20),
    status       VARCHAR(50) DEFAULT 'pending',
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_warehouse_bonuses_tenant_id ON warehouse_bonuses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_bonuses_warehouse_id ON warehouse_bonuses(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_bonuses_status ON warehouse_bonuses(status);
CREATE INDEX IF NOT EXISTS idx_warehouse_bonuses_extra_fields ON warehouse_bonuses USING GIN (extra_fields);

-- ============================================================
-- 31. warehouse_withdrawals
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_withdrawals (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    warehouse_id UUID REFERENCES warehouses(id),
    amount       NUMERIC(12,2) NOT NULL DEFAULT 0,
    method       VARCHAR(50),
    account_info JSONB DEFAULT '{}'::jsonb,
    status       VARCHAR(50) DEFAULT 'pending',
    audit_remark TEXT,
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_warehouse_withdrawals_tenant_id ON warehouse_withdrawals(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_withdrawals_warehouse_id ON warehouse_withdrawals(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_withdrawals_status ON warehouse_withdrawals(status);
CREATE INDEX IF NOT EXISTS idx_warehouse_withdrawals_extra_fields ON warehouse_withdrawals USING GIN (extra_fields);

-- ============================================================
-- 32. batch_templates
-- ============================================================
CREATE TABLE IF NOT EXISTS batch_templates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID REFERENCES tenants(id),
    name             VARCHAR(255) NOT NULL,
    route_id         UUID REFERENCES shipping_routes(id),
    description      TEXT,
    default_settings JSONB DEFAULT '{}'::jsonb,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    trashed_at       TIMESTAMPTZ,
    extra_fields     JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_batch_templates_tenant_id ON batch_templates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_batch_templates_route_id ON batch_templates(route_id);
CREATE INDEX IF NOT EXISTS idx_batch_templates_extra_fields ON batch_templates USING GIN (extra_fields);

-- ============================================================
-- 33. order_reviews
-- ============================================================
CREATE TABLE IF NOT EXISTS order_reviews (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    order_id     UUID REFERENCES consolidation_orders(id),
    user_id      UUID REFERENCES users(id),
    rating       INT DEFAULT 5,
    content      TEXT,
    images       JSONB DEFAULT '[]'::jsonb,
    reply        TEXT,
    replied_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_order_reviews_tenant_id ON order_reviews(tenant_id);
CREATE INDEX IF NOT EXISTS idx_order_reviews_order_id ON order_reviews(order_id);
CREATE INDEX IF NOT EXISTS idx_order_reviews_user_id ON order_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_order_reviews_extra_fields ON order_reviews USING GIN (extra_fields);

-- ============================================================
-- 34. sms_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS sms_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    phone        VARCHAR(50) NOT NULL,
    content      TEXT,
    type         VARCHAR(50),
    status       VARCHAR(50) DEFAULT 'pending',
    sent_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_sms_logs_tenant_id ON sms_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sms_logs_status ON sms_logs(status);
CREATE INDEX IF NOT EXISTS idx_sms_logs_extra_fields ON sms_logs USING GIN (extra_fields);

-- ============================================================
-- 35. email_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS email_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    email        VARCHAR(255) NOT NULL,
    subject      VARCHAR(500),
    content      TEXT,
    type         VARCHAR(50),
    status       VARCHAR(50) DEFAULT 'pending',
    sent_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_email_logs_tenant_id ON email_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_email_logs_status ON email_logs(status);
CREATE INDEX IF NOT EXISTS idx_email_logs_extra_fields ON email_logs USING GIN (extra_fields);

-- ============================================================
-- 36. barcode_settings
-- ============================================================
CREATE TABLE IF NOT EXISTS barcode_settings (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    barcode_type VARCHAR(50) NOT NULL,
    prefix       VARCHAR(50),
    length       INT DEFAULT 12,
    current_seq  BIGINT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_barcode_settings_tenant_id ON barcode_settings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_barcode_settings_extra_fields ON barcode_settings USING GIN (extra_fields);

-- ============================================================
-- 37. bank_accounts
-- ============================================================
CREATE TABLE IF NOT EXISTS bank_accounts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    bank_name    VARCHAR(255) NOT NULL,
    account_name VARCHAR(255),
    account_no   VARCHAR(100),
    branch       VARCHAR(255),
    swift_code   VARCHAR(50),
    currency     VARCHAR(10),
    is_default   BOOLEAN DEFAULT FALSE,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_tenant_id ON bank_accounts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_status ON bank_accounts(status);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_extra_fields ON bank_accounts USING GIN (extra_fields);

-- ============================================================
-- 38. remittance_certificates
-- ============================================================
CREATE TABLE IF NOT EXISTS remittance_certificates (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID REFERENCES tenants(id),
    user_id           UUID REFERENCES users(id),
    order_id          UUID,
    amount            NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency          VARCHAR(10),
    bank_account_id   UUID REFERENCES bank_accounts(id),
    certificate_image VARCHAR(1000),
    remark            TEXT,
    status            VARCHAR(50) DEFAULT 'pending',
    audit_remark      TEXT,
    auditor_id        UUID REFERENCES users(id),
    audited_at        TIMESTAMPTZ,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    trashed_at        TIMESTAMPTZ,
    extra_fields      JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_remittance_certificates_tenant_id ON remittance_certificates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_remittance_certificates_user_id ON remittance_certificates(user_id);
CREATE INDEX IF NOT EXISTS idx_remittance_certificates_status ON remittance_certificates(status);
CREATE INDEX IF NOT EXISTS idx_remittance_certificates_extra_fields ON remittance_certificates USING GIN (extra_fields);

-- ============================================================
-- 39. payment_flows
-- ============================================================
CREATE TABLE IF NOT EXISTS payment_flows (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    order_type   VARCHAR(50),
    order_id     UUID,
    order_no     VARCHAR(100),
    amount       NUMERIC(12,2) NOT NULL DEFAULT 0,
    pay_method   VARCHAR(50),
    pay_no       VARCHAR(255),
    status       VARCHAR(50) DEFAULT 'pending',
    pay_time     TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_payment_flows_tenant_id ON payment_flows(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_flows_order_id ON payment_flows(order_id);
CREATE INDEX IF NOT EXISTS idx_payment_flows_status ON payment_flows(status);
CREATE INDEX IF NOT EXISTS idx_payment_flows_extra_fields ON payment_flows USING GIN (extra_fields);

-- ============================================================
-- 40. channels
-- ============================================================
CREATE TABLE IF NOT EXISTS channels (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    code         VARCHAR(100),
    type         VARCHAR(50),
    config       JSONB DEFAULT '{}'::jsonb,
    status       VARCHAR(50) DEFAULT 'active',
    sort_order   INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_channels_tenant_id ON channels(tenant_id);
CREATE INDEX IF NOT EXISTS idx_channels_status ON channels(status);
CREATE INDEX IF NOT EXISTS idx_channels_extra_fields ON channels USING GIN (extra_fields);

-- ============================================================
-- 41. navigation_items
-- ============================================================
CREATE TABLE IF NOT EXISTS navigation_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    url          VARCHAR(1000),
    icon         VARCHAR(255),
    parent_id    UUID REFERENCES navigation_items(id),
    position     VARCHAR(50) DEFAULT 'header',
    sort_order   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_navigation_items_tenant_id ON navigation_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_navigation_items_parent_id ON navigation_items(parent_id);
CREATE INDEX IF NOT EXISTS idx_navigation_items_status ON navigation_items(status);
CREATE INDEX IF NOT EXISTS idx_navigation_items_extra_fields ON navigation_items USING GIN (extra_fields);

-- ============================================================
-- 42. help_articles
-- ============================================================
CREATE TABLE IF NOT EXISTS help_articles (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    title        VARCHAR(500) NOT NULL,
    content      TEXT,
    category     VARCHAR(255),
    sort_order   INT DEFAULT 0,
    is_hot       BOOLEAN DEFAULT FALSE,
    view_count   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_help_articles_tenant_id ON help_articles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_help_articles_status ON help_articles(status);
CREATE INDEX IF NOT EXISTS idx_help_articles_extra_fields ON help_articles USING GIN (extra_fields);

-- ============================================================
-- 43. subscribe_messages
-- ============================================================
CREATE TABLE IF NOT EXISTS subscribe_messages (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    template_id  VARCHAR(255),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    content      JSONB DEFAULT '{}'::jsonb,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_subscribe_messages_tenant_id ON subscribe_messages(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscribe_messages_status ON subscribe_messages(status);
CREATE INDEX IF NOT EXISTS idx_subscribe_messages_extra_fields ON subscribe_messages USING GIN (extra_fields);

-- ============================================================
-- 44. sharing_orders
-- ============================================================
CREATE TABLE IF NOT EXISTS sharing_orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID REFERENCES tenants(id),
    user_id           UUID REFERENCES users(id),
    sharer_id         UUID REFERENCES users(id),
    order_id          UUID,
    order_type        VARCHAR(50),
    commission_amount NUMERIC(12,2) DEFAULT 0,
    status            VARCHAR(50) DEFAULT 'pending',
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    extra_fields      JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_sharing_orders_tenant_id ON sharing_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sharing_orders_user_id ON sharing_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_sharing_orders_sharer_id ON sharing_orders(sharer_id);
CREATE INDEX IF NOT EXISTS idx_sharing_orders_status ON sharing_orders(status);
CREATE INDEX IF NOT EXISTS idx_sharing_orders_extra_fields ON sharing_orders USING GIN (extra_fields);

-- ============================================================
-- 45. sharing_verifications
-- ============================================================
CREATE TABLE IF NOT EXISTS sharing_verifications (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    order_id     UUID,
    user_id      UUID REFERENCES users(id),
    verifier_id  UUID REFERENCES users(id),
    verified_at  TIMESTAMPTZ,
    status       VARCHAR(50) DEFAULT 'pending',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_sharing_verifications_tenant_id ON sharing_verifications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sharing_verifications_user_id ON sharing_verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_sharing_verifications_status ON sharing_verifications(status);
CREATE INDEX IF NOT EXISTS idx_sharing_verifications_extra_fields ON sharing_verifications USING GIN (extra_fields);

-- ============================================================
-- 46. app_settings
-- ============================================================
CREATE TABLE IF NOT EXISTS app_settings (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    setting_type VARCHAR(50) NOT NULL,
    config       JSONB DEFAULT '{}'::jsonb,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_app_settings_tenant_id ON app_settings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_app_settings_setting_type ON app_settings(setting_type);
CREATE INDEX IF NOT EXISTS idx_app_settings_extra_fields ON app_settings USING GIN (extra_fields);

-- ============================================================
-- 47. page_designs
-- ============================================================
CREATE TABLE IF NOT EXISTS page_designs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    type         VARCHAR(50) DEFAULT 'home',
    page_data    JSONB DEFAULT '{}'::jsonb,
    status       VARCHAR(50) DEFAULT 'active',
    is_default   BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_page_designs_tenant_id ON page_designs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_page_designs_status ON page_designs(status);
CREATE INDEX IF NOT EXISTS idx_page_designs_extra_fields ON page_designs USING GIN (extra_fields);

-- ============================================================
-- 48. user_birthdays
-- ============================================================
CREATE TABLE IF NOT EXISTS user_birthdays (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID REFERENCES tenants(id),
    user_id        UUID REFERENCES users(id),
    birthday       DATE,
    lunar_birthday VARCHAR(20),
    send_coupon    BOOLEAN DEFAULT FALSE,
    send_points    BOOLEAN DEFAULT FALSE,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW(),
    extra_fields   JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_user_birthdays_tenant_id ON user_birthdays(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_birthdays_user_id ON user_birthdays(user_id);
CREATE INDEX IF NOT EXISTS idx_user_birthdays_extra_fields ON user_birthdays USING GIN (extra_fields);


-- ============================================================
-- ALTER EXISTING TABLES: Add new columns
-- ============================================================

-- users: add balance, points, grade_id, birthday, referrer_id, is_dealer
ALTER TABLE users ADD COLUMN IF NOT EXISTS balance NUMERIC(12,2) DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS points INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS grade_id UUID;
ALTER TABLE users ADD COLUMN IF NOT EXISTS birthday DATE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referrer_id UUID;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_dealer BOOLEAN DEFAULT FALSE;

-- articles: add category_id
ALTER TABLE articles ADD COLUMN IF NOT EXISTS category_id UUID;
-- Add FK constraint only if it does not already exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_articles_category_id'
          AND table_name = 'articles'
    ) THEN
        ALTER TABLE articles
            ADD CONSTRAINT fk_articles_category_id
            FOREIGN KEY (category_id) REFERENCES article_categories(id);
    END IF;
END
$$;

-- shipping_batches: add template_id
ALTER TABLE shipping_batches ADD COLUMN IF NOT EXISTS template_id UUID;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_shipping_batches_template_id'
          AND table_name = 'shipping_batches'
    ) THEN
        ALTER TABLE shipping_batches
            ADD CONSTRAINT fk_shipping_batches_template_id
            FOREIGN KEY (template_id) REFERENCES batch_templates(id);
    END IF;
END
$$;

-- packages: add is_returned, is_problem, problem_reason, is_appointment, appointment_at
ALTER TABLE packages ADD COLUMN IF NOT EXISTS is_returned BOOLEAN DEFAULT FALSE;
ALTER TABLE packages ADD COLUMN IF NOT EXISTS is_problem BOOLEAN DEFAULT FALSE;
ALTER TABLE packages ADD COLUMN IF NOT EXISTS problem_reason TEXT;
ALTER TABLE packages ADD COLUMN IF NOT EXISTS is_appointment BOOLEAN DEFAULT FALSE;
ALTER TABLE packages ADD COLUMN IF NOT EXISTS appointment_at TIMESTAMPTZ;

-- consolidation_orders: add is_overdue, overdue_at, arrears_amount, offline_pay_status, inpack_type
ALTER TABLE consolidation_orders ADD COLUMN IF NOT EXISTS is_overdue BOOLEAN DEFAULT FALSE;
ALTER TABLE consolidation_orders ADD COLUMN IF NOT EXISTS overdue_at TIMESTAMPTZ;
ALTER TABLE consolidation_orders ADD COLUMN IF NOT EXISTS arrears_amount NUMERIC(12,2) DEFAULT 0;
ALTER TABLE consolidation_orders ADD COLUMN IF NOT EXISTS offline_pay_status VARCHAR(50);
ALTER TABLE consolidation_orders ADD COLUMN IF NOT EXISTS inpack_type INT DEFAULT 1;
