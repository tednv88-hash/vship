-- vShip Initial Schema Migration
-- Cross-border consolidation shipping system
-- PostgreSQL with UUID primary keys, JSONB, GIN indexes, soft-delete
--
-- Table creation order respects foreign key dependencies.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- 1. tenants
-- ============================================================
CREATE TABLE IF NOT EXISTS tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    subdomain   VARCHAR(100) UNIQUE NOT NULL,
    plan        VARCHAR(50) DEFAULT 'free',
    status      VARCHAR(50) DEFAULT 'active',
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_tenants_extra_fields ON tenants USING GIN (extra_fields);

-- ============================================================
-- 2. users
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    phone         VARCHAR(50),
    user_role     VARCHAR(50) DEFAULT 'user',
    status        VARCHAR(50) DEFAULT 'active',
    profile_pic   VARCHAR(500),
    last_login_at TIMESTAMPTZ,
    logged_out_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    trashed_at    TIMESTAMPTZ,
    extra_fields  JSONB DEFAULT '{}'
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_tenant_email ON users (tenant_id, email);
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users (tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_extra_fields ON users USING GIN (extra_fields);

-- ============================================================
-- 3. countries
-- ============================================================
CREATE TABLE IF NOT EXISTS countries (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    name_en      VARCHAR(255),
    code         VARCHAR(2) NOT NULL,
    phone_code   VARCHAR(10),
    status       VARCHAR(50) DEFAULT 'active',
    sort_order   INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_countries_tenant_id ON countries (tenant_id);
CREATE INDEX IF NOT EXISTS idx_countries_extra_fields ON countries USING GIN (extra_fields);

-- ============================================================
-- 4. currencies
-- ============================================================
CREATE TABLE IF NOT EXISTS currencies (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    name          VARCHAR(255) NOT NULL,
    code          VARCHAR(10) NOT NULL,
    symbol        VARCHAR(10),
    exchange_rate DECIMAL(18,6) DEFAULT 1.0,
    is_default    BOOLEAN DEFAULT FALSE,
    status        VARCHAR(50) DEFAULT 'active',
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    trashed_at    TIMESTAMPTZ,
    extra_fields  JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_currencies_tenant_id ON currencies (tenant_id);
CREATE INDEX IF NOT EXISTS idx_currencies_extra_fields ON currencies USING GIN (extra_fields);

-- ============================================================
-- 5. member_levels
-- ============================================================
CREATE TABLE IF NOT EXISTS member_levels (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    name          VARCHAR(255) NOT NULL,
    code          VARCHAR(50),
    discount_rate DECIMAL(5,2) DEFAULT 0,
    min_spend     DECIMAL(12,2) DEFAULT 0,
    sort_order    INT DEFAULT 0,
    status        VARCHAR(50) DEFAULT 'active',
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    trashed_at    TIMESTAMPTZ,
    extra_fields  JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_member_levels_tenant_id ON member_levels (tenant_id);
CREATE INDEX IF NOT EXISTS idx_member_levels_extra_fields ON member_levels USING GIN (extra_fields);

-- ============================================================
-- 6. package_categories
-- ============================================================
CREATE TABLE IF NOT EXISTS package_categories (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    status       VARCHAR(50) DEFAULT 'active',
    sort_order   INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_package_categories_tenant_id ON package_categories (tenant_id);
CREATE INDEX IF NOT EXISTS idx_package_categories_extra_fields ON package_categories USING GIN (extra_fields);

-- ============================================================
-- 7. route_categories
-- ============================================================
CREATE TABLE IF NOT EXISTS route_categories (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    status       VARCHAR(50) DEFAULT 'active',
    sort_order   INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_route_categories_tenant_id ON route_categories (tenant_id);
CREATE INDEX IF NOT EXISTS idx_route_categories_extra_fields ON route_categories USING GIN (extra_fields);

-- ============================================================
-- 8. warehouses
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouses (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID REFERENCES tenants(id),
    name           VARCHAR(255) NOT NULL,
    code           VARCHAR(50),
    country_id     UUID REFERENCES countries(id),
    address        TEXT,
    phone          VARCHAR(50),
    contact_person VARCHAR(255),
    status         VARCHAR(50) DEFAULT 'active',
    is_default     BOOLEAN DEFAULT FALSE,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW(),
    trashed_at     TIMESTAMPTZ,
    extra_fields   JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_warehouses_tenant_id ON warehouses (tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouses_extra_fields ON warehouses USING GIN (extra_fields);

-- ============================================================
-- 9. warehouse_shelves
-- ============================================================
CREATE TABLE IF NOT EXISTS warehouse_shelves (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id),
    warehouse_id  UUID NOT NULL REFERENCES warehouses(id),
    shelf_code    VARCHAR(100) NOT NULL,
    zone          VARCHAR(100),
    capacity      INT DEFAULT 0,
    used_capacity INT DEFAULT 0,
    status        VARCHAR(50) DEFAULT 'active',
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    extra_fields  JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_warehouse_shelves_tenant_id ON warehouse_shelves (tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_shelves_extra_fields ON warehouse_shelves USING GIN (extra_fields);

-- ============================================================
-- 10. shipping_marks
-- ============================================================
CREATE TABLE IF NOT EXISTS shipping_marks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    user_id      UUID NOT NULL REFERENCES users(id),
    code         VARCHAR(100) NOT NULL,
    description  TEXT,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_shipping_marks_tenant_id ON shipping_marks (tenant_id);
CREATE INDEX IF NOT EXISTS idx_shipping_marks_extra_fields ON shipping_marks USING GIN (extra_fields);

-- ============================================================
-- 11. packages
-- ============================================================
CREATE TABLE IF NOT EXISTS packages (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID REFERENCES tenants(id),
    tracking_number   VARCHAR(255) NOT NULL,
    shipping_mark_id  UUID REFERENCES shipping_marks(id),
    user_id           UUID REFERENCES users(id),
    warehouse_id      UUID REFERENCES warehouses(id),
    shelf_id          UUID REFERENCES warehouse_shelves(id),
    category_id       UUID REFERENCES package_categories(id),
    status            VARCHAR(50) DEFAULT 'forecast',
    source            VARCHAR(50),
    weight            DECIMAL(10,3) DEFAULT 0,
    length            DECIMAL(10,2) DEFAULT 0,
    width             DECIMAL(10,2) DEFAULT 0,
    height            DECIMAL(10,2) DEFAULT 0,
    volume_weight     DECIMAL(10,3) DEFAULT 0,
    chargeable_weight DECIMAL(10,3) DEFAULT 0,
    declared_value    DECIMAL(12,2) DEFAULT 0,
    declared_currency VARCHAR(10),
    item_description  TEXT,
    remark            TEXT,
    received_at       TIMESTAMPTZ,
    shelved_at        TIMESTAMPTZ,
    inspected_at      TIMESTAMPTZ,
    packed_at         TIMESTAMPTZ,
    shipped_at        TIMESTAMPTZ,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    trashed_at        TIMESTAMPTZ,
    extra_fields      JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_packages_tenant_id ON packages (tenant_id);
CREATE INDEX IF NOT EXISTS idx_packages_tracking_number ON packages (tracking_number);
CREATE INDEX IF NOT EXISTS idx_packages_extra_fields ON packages USING GIN (extra_fields);

-- ============================================================
-- 12. package_status_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS package_status_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(id),
    package_id  UUID NOT NULL REFERENCES packages(id),
    from_status VARCHAR(50) NOT NULL,
    to_status   VARCHAR(50) NOT NULL,
    remark      TEXT,
    operator_id UUID NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_package_status_logs_tenant_id ON package_status_logs (tenant_id);

-- ============================================================
-- 13. shipping_routes
-- ============================================================
CREATE TABLE IF NOT EXISTS shipping_routes (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID REFERENCES tenants(id),
    name                   VARCHAR(255) NOT NULL,
    origin_country_id      UUID REFERENCES countries(id),
    destination_country_id UUID REFERENCES countries(id),
    origin_warehouse_id    UUID REFERENCES warehouses(id),
    destination_warehouse_id UUID REFERENCES warehouses(id),
    category_id            UUID REFERENCES route_categories(id),
    transport_mode         VARCHAR(50) NOT NULL,
    billing_mode           VARCHAR(50) NOT NULL,
    weight_unit            VARCHAR(10) DEFAULT 'KG',
    volume_weight_ratio    NUMERIC(10,2) DEFAULT 5000,
    rounding_rule          VARCHAR(50),
    multi_box_mode         VARCHAR(50) DEFAULT 'separate',
    estimated_days         INT DEFAULT 0,
    status                 VARCHAR(50) DEFAULT 'active',
    sort_order             INT DEFAULT 0,
    created_at             TIMESTAMPTZ DEFAULT NOW(),
    updated_at             TIMESTAMPTZ DEFAULT NOW(),
    trashed_at             TIMESTAMPTZ,
    extra_fields           JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_shipping_routes_tenant_id ON shipping_routes (tenant_id);
CREATE INDEX IF NOT EXISTS idx_shipping_routes_extra_fields ON shipping_routes USING GIN (extra_fields);

-- ============================================================
-- 14. route_pricing_tiers
-- ============================================================
CREATE TABLE IF NOT EXISTS route_pricing_tiers (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID REFERENCES tenants(id),
    route_id                UUID NOT NULL REFERENCES shipping_routes(id),
    member_level_id         UUID REFERENCES member_levels(id),
    weight_min              NUMERIC(10,2) DEFAULT 0,
    weight_max              NUMERIC(10,2) DEFAULT 0,
    unit_price              NUMERIC(10,2) DEFAULT 0,
    first_weight            NUMERIC(10,2) DEFAULT 0,
    first_weight_price      NUMERIC(10,2) DEFAULT 0,
    additional_weight_price NUMERIC(10,2) DEFAULT 0,
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_route_pricing_tiers_tenant_id ON route_pricing_tiers (tenant_id);

-- ============================================================
-- 15. consolidation_orders
-- ============================================================
CREATE TABLE IF NOT EXISTS consolidation_orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID REFERENCES tenants(id),
    order_number      VARCHAR(100) NOT NULL,
    user_id           UUID REFERENCES users(id),
    shipping_route_id UUID REFERENCES shipping_routes(id),
    warehouse_id      UUID REFERENCES warehouses(id),
    member_level_id   UUID REFERENCES member_levels(id),
    status            VARCHAR(50) DEFAULT 'draft',
    total_weight      NUMERIC(10,3) DEFAULT 0,
    total_volume_weight NUMERIC(10,3) DEFAULT 0,
    chargeable_weight NUMERIC(10,3) DEFAULT 0,
    shipping_fee      NUMERIC(12,2) DEFAULT 0,
    insurance_fee     NUMERIC(12,2) DEFAULT 0,
    service_fee       NUMERIC(12,2) DEFAULT 0,
    consumable_fee    NUMERIC(12,2) DEFAULT 0,
    total_amount      NUMERIC(12,2) DEFAULT 0,
    paid_amount       NUMERIC(12,2) DEFAULT 0,
    currency          VARCHAR(10) DEFAULT 'TWD',
    payment_method    VARCHAR(50),
    payment_status    VARCHAR(50) DEFAULT 'unpaid',
    recipient_name    VARCHAR(255),
    recipient_phone   VARCHAR(50),
    recipient_address TEXT,
    recipient_city    VARCHAR(100),
    recipient_state   VARCHAR(100),
    recipient_zip     VARCHAR(20),
    recipient_country VARCHAR(100),
    remark            TEXT,
    paid_at           TIMESTAMPTZ,
    packed_at         TIMESTAMPTZ,
    shipped_at        TIMESTAMPTZ,
    arrived_at        TIMESTAMPTZ,
    completed_at      TIMESTAMPTZ,
    cancelled_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    trashed_at        TIMESTAMPTZ,
    extra_fields      JSONB DEFAULT '{}'
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_order_number ON consolidation_orders (tenant_id, order_number);
CREATE INDEX IF NOT EXISTS idx_consolidation_orders_tenant_id ON consolidation_orders (tenant_id);
CREATE INDEX IF NOT EXISTS idx_consolidation_orders_order_number ON consolidation_orders (order_number);
CREATE INDEX IF NOT EXISTS idx_consolidation_orders_extra_fields ON consolidation_orders USING GIN (extra_fields);

-- ============================================================
-- 16. order_packages
-- ============================================================
CREATE TABLE IF NOT EXISTS order_packages (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID REFERENCES tenants(id),
    order_id   UUID NOT NULL REFERENCES consolidation_orders(id),
    package_id UUID NOT NULL REFERENCES packages(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_order_packages_tenant_id ON order_packages (tenant_id);

-- ============================================================
-- 17. value_added_services  (moved before order_services which references it)
-- ============================================================
CREATE TABLE IF NOT EXISTS value_added_services (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    price        NUMERIC(10,2) DEFAULT 0,
    price_unit   VARCHAR(50) DEFAULT 'per_item',
    status       VARCHAR(50) DEFAULT 'active',
    sort_order   INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_value_added_services_tenant_id ON value_added_services (tenant_id);
CREATE INDEX IF NOT EXISTS idx_value_added_services_extra_fields ON value_added_services USING GIN (extra_fields);

-- ============================================================
-- 18. order_services
-- ============================================================
CREATE TABLE IF NOT EXISTS order_services (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    order_id     UUID NOT NULL REFERENCES consolidation_orders(id),
    service_id   UUID REFERENCES value_added_services(id),
    service_name VARCHAR(255) NOT NULL,
    quantity     INT DEFAULT 1,
    unit_price   NUMERIC(10,2) DEFAULT 0,
    total_price  NUMERIC(10,2) DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_order_services_tenant_id ON order_services (tenant_id);

-- ============================================================
-- 19. payment_audits
-- ============================================================
CREATE TABLE IF NOT EXISTS payment_audits (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID REFERENCES tenants(id),
    order_id        UUID NOT NULL REFERENCES consolidation_orders(id),
    user_id         UUID REFERENCES users(id),
    amount          NUMERIC(12,2) DEFAULT 0,
    currency        VARCHAR(10) DEFAULT 'TWD',
    payment_method  VARCHAR(50),
    certificate_url VARCHAR(500),
    status          VARCHAR(50) DEFAULT 'pending',
    reviewed_by_id  UUID REFERENCES users(id),
    reviewed_at     TIMESTAMPTZ,
    reject_reason   TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    extra_fields    JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_payment_audits_tenant_id ON payment_audits (tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_audits_extra_fields ON payment_audits USING GIN (extra_fields);

-- ============================================================
-- 20. logistics_companies  (moved before shipping_batches which references it)
-- ============================================================
CREATE TABLE IF NOT EXISTS logistics_companies (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID REFERENCES tenants(id),
    name           VARCHAR(255) NOT NULL,
    code           VARCHAR(100),
    contact_person VARCHAR(255),
    phone          VARCHAR(50),
    email          VARCHAR(255),
    website        VARCHAR(500),
    tracking_url   VARCHAR(500),
    status         VARCHAR(50) DEFAULT 'active',
    sort_order     INT DEFAULT 0,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW(),
    trashed_at     TIMESTAMPTZ,
    extra_fields   JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_logistics_companies_tenant_id ON logistics_companies (tenant_id);
CREATE INDEX IF NOT EXISTS idx_logistics_companies_extra_fields ON logistics_companies USING GIN (extra_fields);

-- ============================================================
-- 21. shipping_batches
-- ============================================================
CREATE TABLE IF NOT EXISTS shipping_batches (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                UUID REFERENCES tenants(id),
    batch_number             VARCHAR(100) NOT NULL,
    name                     VARCHAR(255),
    type                     VARCHAR(50) NOT NULL,
    container_code           VARCHAR(100),
    master_tracking_number   VARCHAR(255),
    origin_warehouse_id      UUID REFERENCES warehouses(id),
    destination_warehouse_id UUID REFERENCES warehouses(id),
    logistics_company_id     UUID REFERENCES logistics_companies(id),
    status                   VARCHAR(50) DEFAULT 'preparing',
    total_weight             NUMERIC(12,3) DEFAULT 0,
    total_volume             NUMERIC(12,3) DEFAULT 0,
    total_orders             INT DEFAULT 0,
    departed_at              TIMESTAMPTZ,
    arrived_at               TIMESTAMPTZ,
    completed_at             TIMESTAMPTZ,
    remark                   TEXT,
    created_at               TIMESTAMPTZ DEFAULT NOW(),
    updated_at               TIMESTAMPTZ DEFAULT NOW(),
    trashed_at               TIMESTAMPTZ,
    extra_fields             JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_shipping_batches_tenant_id ON shipping_batches (tenant_id);
CREATE INDEX IF NOT EXISTS idx_shipping_batches_batch_number ON shipping_batches (batch_number);
CREATE INDEX IF NOT EXISTS idx_shipping_batches_extra_fields ON shipping_batches USING GIN (extra_fields);

-- ============================================================
-- 22. batch_orders
-- ============================================================
CREATE TABLE IF NOT EXISTS batch_orders (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID REFERENCES tenants(id),
    batch_id   UUID NOT NULL REFERENCES shipping_batches(id),
    order_id   UUID NOT NULL REFERENCES consolidation_orders(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_batch_orders_tenant_id ON batch_orders (tenant_id);

-- ============================================================
-- 23. batch_tracking_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS batch_tracking_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(id),
    batch_id    UUID NOT NULL REFERENCES shipping_batches(id),
    status      VARCHAR(50) NOT NULL,
    location    VARCHAR(255),
    description TEXT,
    occurred_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_batch_tracking_logs_tenant_id ON batch_tracking_logs (tenant_id);

-- ============================================================
-- 24. insurance_products
-- ============================================================
CREATE TABLE IF NOT EXISTS insurance_products (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    premium_rate NUMERIC(8,4) DEFAULT 0,
    min_premium  NUMERIC(10,2) DEFAULT 0,
    max_coverage NUMERIC(12,2) DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    sort_order   INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_insurance_products_tenant_id ON insurance_products (tenant_id);
CREATE INDEX IF NOT EXISTS idx_insurance_products_extra_fields ON insurance_products USING GIN (extra_fields);

-- ============================================================
-- 25. consumables
-- ============================================================
CREATE TABLE IF NOT EXISTS consumables (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    price        NUMERIC(10,2) DEFAULT 0,
    unit         VARCHAR(50) DEFAULT '個',
    stock        INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_consumables_tenant_id ON consumables (tenant_id);
CREATE INDEX IF NOT EXISTS idx_consumables_extra_fields ON consumables USING GIN (extra_fields);

-- ============================================================
-- 26. tracking_templates
-- ============================================================
CREATE TABLE IF NOT EXISTS tracking_templates (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    name         VARCHAR(255) NOT NULL,
    content      TEXT NOT NULL,
    sort_order   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_tracking_templates_tenant_id ON tracking_templates (tenant_id);
CREATE INDEX IF NOT EXISTS idx_tracking_templates_extra_fields ON tracking_templates USING GIN (extra_fields);

-- ============================================================
-- 27. banners
-- ============================================================
CREATE TABLE IF NOT EXISTS banners (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    title        VARCHAR(255) NOT NULL,
    image_url    VARCHAR(500) NOT NULL,
    link_url     VARCHAR(500),
    position     VARCHAR(100),
    sort_order   INT DEFAULT 0,
    status       VARCHAR(50) DEFAULT 'active',
    start_at     TIMESTAMPTZ,
    end_at       TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_banners_tenant_id ON banners (tenant_id);
CREATE INDEX IF NOT EXISTS idx_banners_extra_fields ON banners USING GIN (extra_fields);

-- ============================================================
-- 28. articles
-- ============================================================
CREATE TABLE IF NOT EXISTS articles (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID REFERENCES tenants(id),
    title        VARCHAR(255) NOT NULL,
    content      TEXT,
    category     VARCHAR(100),
    cover_image  VARCHAR(500),
    author       VARCHAR(255),
    status       VARCHAR(50) DEFAULT 'draft',
    view_count   INT DEFAULT 0,
    sort_order   INT DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    trashed_at   TIMESTAMPTZ,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_articles_tenant_id ON articles (tenant_id);
CREATE INDEX IF NOT EXISTS idx_articles_extra_fields ON articles USING GIN (extra_fields);

-- ============================================================
-- 29. coupons
-- ============================================================
CREATE TABLE IF NOT EXISTS coupons (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID REFERENCES tenants(id),
    code             VARCHAR(100) NOT NULL,
    name             VARCHAR(255) NOT NULL,
    type             VARCHAR(50) NOT NULL,
    value            NUMERIC(10,2) DEFAULT 0,
    min_order_amount NUMERIC(12,2) DEFAULT 0,
    max_discount     NUMERIC(12,2) DEFAULT 0,
    total_count      INT DEFAULT 0,
    used_count       INT DEFAULT 0,
    member_level_id  UUID REFERENCES member_levels(id),
    status           VARCHAR(50) DEFAULT 'active',
    start_at         TIMESTAMPTZ,
    end_at           TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    trashed_at       TIMESTAMPTZ,
    extra_fields     JSONB DEFAULT '{}'
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_coupon_code ON coupons (tenant_id, code);
CREATE INDEX IF NOT EXISTS idx_coupons_tenant_id ON coupons (tenant_id);
CREATE INDEX IF NOT EXISTS idx_coupons_extra_fields ON coupons USING GIN (extra_fields);

-- ============================================================
-- 30. settings
-- ============================================================
CREATE TABLE IF NOT EXISTS settings (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID REFERENCES tenants(id),
    key        VARCHAR(255) NOT NULL,
    value      TEXT,
    "group"    VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_setting_key ON settings (tenant_id, key);
CREATE INDEX IF NOT EXISTS idx_settings_tenant_id ON settings (tenant_id);
