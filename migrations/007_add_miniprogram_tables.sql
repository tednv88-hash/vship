-- Migration 007: Add miniprogram feature tables
-- ============================================================
-- Tables: cart_items, favorites, browsing_histories, feedbacks,
--         signin_logs, content_pages
-- ============================================================

-- ============================================================
-- 1. Cart Items
-- ============================================================
CREATE TABLE IF NOT EXISTS cart_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID,
    user_id      UUID NOT NULL,
    goods_id     UUID NOT NULL,
    sku_id       UUID,
    quantity     INT DEFAULT 1,
    selected     BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_cart_items_tenant_id ON cart_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_goods_id ON cart_items(goods_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_extra_fields ON cart_items USING GIN (extra_fields);

-- ============================================================
-- 2. Favorites
-- ============================================================
CREATE TABLE IF NOT EXISTS favorites (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID,
    user_id      UUID NOT NULL,
    goods_id     UUID NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_tenant_user_goods ON favorites(tenant_id, user_id, goods_id);
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_extra_fields ON favorites USING GIN (extra_fields);

-- ============================================================
-- 3. Browsing Histories
-- ============================================================
CREATE TABLE IF NOT EXISTS browsing_histories (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID,
    user_id      UUID NOT NULL,
    goods_id     UUID NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_browsing_histories_tenant_id ON browsing_histories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_browsing_histories_user_id ON browsing_histories(user_id);
CREATE INDEX IF NOT EXISTS idx_browsing_histories_extra_fields ON browsing_histories USING GIN (extra_fields);

-- ============================================================
-- 4. Feedbacks
-- ============================================================
CREATE TABLE IF NOT EXISTS feedbacks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID,
    user_id      UUID NOT NULL,
    type         VARCHAR(50),
    content      TEXT,
    images       JSONB,
    contact_info VARCHAR(255),
    status       VARCHAR(50) DEFAULT 'pending',
    reply        TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_feedbacks_tenant_id ON feedbacks(tenant_id);
CREATE INDEX IF NOT EXISTS idx_feedbacks_user_id ON feedbacks(user_id);
CREATE INDEX IF NOT EXISTS idx_feedbacks_status ON feedbacks(status);
CREATE INDEX IF NOT EXISTS idx_feedbacks_extra_fields ON feedbacks USING GIN (extra_fields);

-- ============================================================
-- 5. Signin Logs (daily check-in)
-- ============================================================
CREATE TABLE IF NOT EXISTS signin_logs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID,
    user_id          UUID NOT NULL,
    signin_date      VARCHAR(10) NOT NULL,
    points_earned    INT DEFAULT 0,
    consecutive_days INT DEFAULT 0,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    extra_fields     JSONB DEFAULT '{}'::jsonb
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_signin_logs_tenant_user_date ON signin_logs(tenant_id, user_id, signin_date);
CREATE INDEX IF NOT EXISTS idx_signin_logs_user_id ON signin_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_signin_logs_extra_fields ON signin_logs USING GIN (extra_fields);

-- ============================================================
-- 6. Content Pages (about, privacy, terms, etc.)
-- ============================================================
CREATE TABLE IF NOT EXISTS content_pages (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID,
    slug         VARCHAR(100) NOT NULL,
    title        VARCHAR(255),
    content      TEXT,
    status       VARCHAR(50) DEFAULT 'active',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'::jsonb
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_content_pages_tenant_slug ON content_pages(tenant_id, slug);
CREATE INDEX IF NOT EXISTS idx_content_pages_extra_fields ON content_pages USING GIN (extra_fields);

-- ============================================================
-- 7. Seed default content pages for main tenant
-- ============================================================
INSERT INTO content_pages (tenant_id, slug, title, content, status)
SELECT '432ba247-0944-40e8-a94f-e58603859020'::uuid, 'about', 'About Us', '<h2>About vShip</h2><p>vShip is a modern shipping and logistics platform designed to make international shipping simple, transparent, and affordable.</p>', 'active'
WHERE NOT EXISTS (SELECT 1 FROM content_pages WHERE tenant_id = '432ba247-0944-40e8-a94f-e58603859020'::uuid AND slug = 'about');

INSERT INTO content_pages (tenant_id, slug, title, content, status)
SELECT '432ba247-0944-40e8-a94f-e58603859020'::uuid, 'privacy', 'Privacy Policy', '<h2>Privacy Policy</h2><p>We are committed to protecting your privacy. This policy describes how we collect, use, and safeguard your personal information.</p>', 'active'
WHERE NOT EXISTS (SELECT 1 FROM content_pages WHERE tenant_id = '432ba247-0944-40e8-a94f-e58603859020'::uuid AND slug = 'privacy');

INSERT INTO content_pages (tenant_id, slug, title, content, status)
SELECT '432ba247-0944-40e8-a94f-e58603859020'::uuid, 'terms', 'Terms of Service', '<h2>Terms of Service</h2><p>By using vShip, you agree to the following terms and conditions governing your use of our platform and services.</p>', 'active'
WHERE NOT EXISTS (SELECT 1 FROM content_pages WHERE tenant_id = '432ba247-0944-40e8-a94f-e58603859020'::uuid AND slug = 'terms');
