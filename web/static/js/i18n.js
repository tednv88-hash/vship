// vShip i18n - Internationalization support (Traditional Chinese / Simplified Chinese / English)
// Based on vWork i18n architecture, adapted for vShip admin panel

const SafeStorage = {
    get: function(key) {
        try { return localStorage.getItem(key); } catch (e) { return null; }
    },
    set: function(key, value) {
        try { localStorage.setItem(key, value); } catch (e) { /* ignore */ }
    }
};

const I18n = {
    currentLang: 'zh',
    translations: {},
    _ready: false,
    _readyPromise: null,
    _readyResolve: null,
    _initStarted: false,
    _unlockTimer: null,
    _pendingPatches: [],

    _ensureReadyPromise: function() {
        if (!this._readyPromise) {
            this._readyPromise = new Promise((resolve) => { this._readyResolve = resolve; });
        }
    },

    // Initialize i18n
    init: function() {
        if (this._initStarted) return;
        this._initStarted = true;

        this._ensureReadyPromise();
        this._ready = false;

        // Safety unlock timer - prevent blank page if i18n fails to load
        try {
            if (this._unlockTimer) clearTimeout(this._unlockTimer);
            this._unlockTimer = setTimeout(() => {
                try {
                    if (document.body && document.body.classList.contains('i18n-loading')) {
                        console.warn('i18n: unlock timeout reached, removing i18n-loading');
                        document.body.classList.remove('i18n-loading');
                    }
                } catch (e) { /* ignore */ }
            }, 5000);
        } catch (e) { /* ignore */ }

        // Determine language: localStorage > browser language > default (zh)
        const savedLang = SafeStorage.get('vship_lang');
        const browserLang = navigator.language || navigator.userLanguage;
        const isSupportedLang = (l) => l === 'zh' || l === 'zh-CN' || l === 'en';

        if (savedLang && isSupportedLang(savedLang)) {
            this.currentLang = savedLang;
        } else {
            if (browserLang.startsWith('zh-CN') || browserLang.startsWith('zh-Hans')) {
                this.currentLang = 'zh-CN';
            } else if (browserLang.startsWith('zh')) {
                this.currentLang = 'zh';
            } else {
                this.currentLang = 'en';
            }
            SafeStorage.set('vship_lang', this.currentLang);
        }

        this.loadLanguage(this.currentLang);
        this.bindLangSwitcher();
    },

    // Wait for translations to be loaded
    whenReady: function(timeoutMs = 3000) {
        if (this._ready) return Promise.resolve();
        this._ensureReadyPromise();
        const p = this._readyPromise;
        if (!timeoutMs || timeoutMs <= 0) return p;
        return Promise.race([
            p,
            new Promise((resolve) => setTimeout(resolve, timeoutMs))
        ]);
    },

    // Bind language switcher buttons ([data-lang])
    bindLangSwitcher: function() {
        document.querySelectorAll('[data-lang]').forEach(el => {
            el.addEventListener('click', (e) => {
                e.preventDefault();
                const lang = el.getAttribute('data-lang');
                if (lang) {
                    this.switchLanguage(lang);
                }
            });
        });
    },

    // Load language file
    loadLanguage: async function(lang) {
        try {
            const langFile = lang === 'zh-CN' ? 'zh-CN' : (lang === 'zh' ? 'zh' : lang);
            const url = `/static/locales/${langFile}.json`;

            const fetchWithTimeout = async (u, timeoutMs = 4000) => {
                if (typeof AbortController === 'undefined') {
                    return await fetch(u);
                }
                const controller = new AbortController();
                const timer = setTimeout(() => {
                    try { controller.abort(); } catch (e) { /* ignore */ }
                }, timeoutMs);
                try {
                    return await fetch(u, { signal: controller.signal });
                } finally {
                    clearTimeout(timer);
                }
            };

            const response = await fetchWithTimeout(url, 4000);
            if (!response.ok) {
                throw new Error(`Failed to load ${langFile}.json: ${response.status}`);
            }

            const translations = await response.json();
            this.translations = translations;
            this.currentLang = lang;
            SafeStorage.set('vship_lang', lang);
            this._ready = true;
            if (typeof this._readyResolve === 'function') {
                try { this._readyResolve(); } catch (e) { /* ignore */ }
            }

            // Dispatch language change event
            if (typeof window !== 'undefined' && window.dispatchEvent) {
                window.dispatchEvent(new CustomEvent('languageChanged', { detail: { lang: lang } }));
            }

            this._flushPendingPatches();

            requestAnimationFrame(() => {
                try {
                    this.updatePage();
                } catch (e) {
                    console.error('i18n updatePage error:', e);
                    if (document.body) {
                        document.body.classList.remove('i18n-loading');
                    }
                }
            });
        } catch (error) {
            console.error('Failed to load language file:', error);
            if (lang !== 'zh') {
                // Fallback to Traditional Chinese
                await this.loadLanguage('zh');
            } else {
                if (document.body) {
                    document.body.classList.remove('i18n-loading');
                }
                this.translations = {};
                this._ready = true;
                if (typeof this._readyResolve === 'function') {
                    try { this._readyResolve(); } catch (e) { /* ignore */ }
                }
                requestAnimationFrame(() => {
                    try { this.updatePage(); } catch (e) { /* ignore */ }
                });
            }
        } finally {
            try {
                if (this._unlockTimer) {
                    clearTimeout(this._unlockTimer);
                    this._unlockTimer = null;
                }
                if (document.body && document.body.classList.contains('i18n-loading')) {
                    document.body.classList.remove('i18n-loading');
                }
            } catch (e) { /* ignore */ }
        }
    },

    // Switch language (reloads page to ensure full update)
    switchLanguage: function(lang) {
        this.currentLang = lang;
        SafeStorage.set('vship_lang', lang);
        window.location.reload();
    },

    // Flush pending patches after translations are loaded
    _flushPendingPatches: function() {
        if (!this._pendingPatches || this._pendingPatches.length === 0) return;
        const patches = this._pendingPatches;
        this._pendingPatches = [];
        try {
            for (const { el, key, attr } of patches) {
                try {
                    if (!el || !el.nodeType) continue;
                    const translated = this.t(key);
                    if (!translated || translated === key) continue;
                    if (attr === 'placeholder') {
                        el.placeholder = translated;
                    } else if (attr === 'title') {
                        el.setAttribute('title', translated);
                    } else {
                        const hasElementChildren = Array.from(el.childNodes || []).some(n => n.nodeType === Node.ELEMENT_NODE);
                        if (hasElementChildren) {
                            for (const node of el.childNodes) {
                                if (node.nodeType === Node.TEXT_NODE && node.textContent.trim()) {
                                    node.textContent = translated;
                                    break;
                                }
                            }
                        } else {
                            el.textContent = translated;
                        }
                    }
                } catch (e) { /* ignore single patch failure */ }
            }
        } catch (e) {
            console.warn('i18n _flushPendingPatches error:', e);
        }
    },

    // Get translation by dot-notation key
    t: function(key, el, attr) {
        const keys = key.split('.');
        let value = this.translations;

        for (const k of keys) {
            if (value && typeof value === 'object') {
                value = value[k];
            } else {
                if (!this._ready && el && el.nodeType) {
                    this._pendingPatches.push({ el, key, attr: attr || null });
                }
                return key;
            }
        }

        if (!this._ready && el && el.nodeType && (value === undefined || value === null)) {
            this._pendingPatches.push({ el, key, attr: attr || null });
        }
        return value || key;
    },

    // Apply translations to a specific DOM subtree
    applyTranslations: function(rootEl) {
        if (!rootEl || !rootEl.querySelectorAll) return;

        const toText = (translation, fallbackText) => {
            if (typeof translation === 'string') return translation;
            if (typeof translation === 'number' || typeof translation === 'boolean') return String(translation);
            return fallbackText || '';
        };

        try {
            rootEl.querySelectorAll('[data-i18n]').forEach(el => {
                try {
                    const key = el.getAttribute('data-i18n');
                    const translation = this.t(key);
                    if (translation && translation !== key) {
                        el.textContent = toText(translation, el.textContent || '');
                    }
                } catch (e) { /* ignore */ }
            });
            rootEl.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
                try {
                    const key = el.getAttribute('data-i18n-placeholder');
                    const translation = this.t(key);
                    if (translation && translation !== key) {
                        el.placeholder = toText(translation, el.placeholder || '');
                    }
                } catch (e) { /* ignore */ }
            });
            rootEl.querySelectorAll('[data-i18n-title]').forEach(el => {
                try {
                    const key = el.getAttribute('data-i18n-title');
                    const translation = this.t(key);
                    if (translation && translation !== key) {
                        el.setAttribute('title', toText(translation, ''));
                    }
                } catch (e) { /* ignore */ }
            });
        } catch (e) {
            console.warn('i18n applyTranslations error:', e);
        }
    },

    // Update the entire page with translations
    updatePage: function() {
        const toText = (translation, fallbackText = '') => {
            if (typeof translation === 'string') return translation;
            if (typeof translation === 'number' || typeof translation === 'boolean') return String(translation);
            return fallbackText || '';
        };

        const safeUpdate = (fn) => {
            try { fn(); } catch (e) { console.warn('i18n update element error:', e); }
        };

        // Update data-i18n-placeholder elements
        try {
            document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
                safeUpdate(() => {
                    const key = element.getAttribute('data-i18n-placeholder');
                    const translation = this.t(key);
                    if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
                        const fallback = element.placeholder || '';
                        element.placeholder = (translation === key) ? fallback : toText(translation, fallback);
                    }
                });
            });
        } catch (e) { /* ignore */ }

        // Update data-i18n elements
        try {
            document.querySelectorAll('[data-i18n]').forEach(element => {
                safeUpdate(() => {
                    const key = element.getAttribute('data-i18n');
                    const translation = this.t(key);
                    const fallback = element.textContent || '';
                    const translationStr = (translation === key) ? fallback : toText(translation, fallback);

                    if (element.tagName === 'SPAN') {
                        element.textContent = translationStr;
                        return;
                    }

                    if (element.tagName === 'BUTTON' || element.tagName === 'A') {
                        const icon = element.querySelector('i');
                        const svg = element.querySelector('svg');
                        if (icon) {
                            // Preserve icon, update text
                            const iconHTML = icon.outerHTML;
                            element.innerHTML = iconHTML + ' ' + translationStr;
                        } else if (svg) {
                            const svgHTML = svg.outerHTML;
                            element.innerHTML = svgHTML + ' ' + translationStr;
                        } else {
                            element.textContent = translationStr;
                        }
                        return;
                    }

                    element.textContent = translationStr;
                });
            });
        } catch (e) { /* ignore */ }

        // Update page title
        try {
            const pageTitle = document.querySelector('title[data-i18n-title]');
            if (pageTitle) {
                const key = pageTitle.getAttribute('data-i18n-title');
                const translated = this.t(key);
                if (translated && translated !== key) {
                    document.title = translated + ' - vShip';
                }
            }
        } catch (e) { /* ignore */ }

        // Update data-i18n-title elements (tooltip titles)
        try {
            document.querySelectorAll('[data-i18n-title]').forEach(element => {
                if (element.tagName === 'TITLE') return; // Skip <title> tag
                safeUpdate(() => {
                    const key = element.getAttribute('data-i18n-title');
                    const translation = this.t(key);
                    if (translation && translation !== key) {
                        element.setAttribute('title', translation);
                    }
                });
            });
        } catch (e) { /* ignore */ }

        // Update language switcher active state
        try {
            document.querySelectorAll('[data-lang]').forEach(btn => {
                safeUpdate(() => {
                    if (btn.getAttribute('data-lang') === this.currentLang) {
                        btn.classList.add('active');
                    } else {
                        btn.classList.remove('active');
                    }
                });
            });
        } catch (e) { /* ignore */ }

        // Remove loading class
        try {
            if (document.body) {
                document.body.classList.remove('i18n-loading');
            }
        } catch (e) { /* ignore */ }
    }
};

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', function() {
    I18n.init();
});
