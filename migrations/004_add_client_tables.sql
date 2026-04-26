-- Migration 004: Add client module tables (web_menus, web_links, wechat_menus, languages, page_categories)

-- Web navigation menus (hierarchical)
CREATE TABLE IF NOT EXISTS web_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    parent_id UUID REFERENCES web_menus(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    type INT DEFAULT 10, -- 10=单页, 20=列表, 30=关于我们, 40=仓库地址
    link_id VARCHAR(255) DEFAULT '',
    sort_order INT DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trashed_at TIMESTAMP WITH TIME ZONE,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_web_menus_tenant ON web_menus(tenant_id);
CREATE INDEX IF NOT EXISTS idx_web_menus_parent ON web_menus(parent_id);

-- Friendly links (友情链接)
CREATE TABLE IF NOT EXISTS web_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(500) DEFAULT '',
    image_url VARCHAR(500) DEFAULT '',
    link_type VARCHAR(50) DEFAULT 'image', -- image, text
    sort_order INT DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trashed_at TIMESTAMP WITH TIME ZONE,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_web_links_tenant ON web_links(tenant_id);

-- WeChat official account custom menus
CREATE TABLE IF NOT EXISTS wechat_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    parent_id UUID REFERENCES wechat_menus(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) DEFAULT 'click', -- click, view, miniprogram, scancode_push, media_id
    key VARCHAR(255) DEFAULT '',
    url VARCHAR(500) DEFAULT '',
    appid VARCHAR(255) DEFAULT '',
    pagepath VARCHAR(255) DEFAULT '',
    backup_url VARCHAR(500) DEFAULT '',
    media_id VARCHAR(255) DEFAULT '',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trashed_at TIMESTAMP WITH TIME ZONE,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_wechat_menus_tenant ON wechat_menus(tenant_id);

-- Supported languages
CREATE TABLE IF NOT EXISTS languages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    name VARCHAR(100) NOT NULL,
    enname VARCHAR(100) DEFAULT '',
    langto VARCHAR(50) DEFAULT '',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trashed_at TIMESTAMP WITH TIME ZONE,
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_languages_tenant ON languages(tenant_id);

-- Page category template settings
CREATE TABLE IF NOT EXISTS page_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    category_style VARCHAR(50) DEFAULT '20', -- 10=一级分类(大图), 11=一级分类(小图), 20=二级分类
    share_title VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    extra_fields JSONB DEFAULT '{}'
);
CREATE INDEX IF NOT EXISTS idx_page_categories_tenant ON page_categories(tenant_id);
