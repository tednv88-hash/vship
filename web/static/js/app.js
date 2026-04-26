/**
 * vShip CMS - Common JavaScript Utilities
 */

const API_BASE = '/api/v1';

/**
 * Get auth token from localStorage (matching vWork pattern)
 */
function getToken() {
    return localStorage.getItem('auth_token') || '';
}

/**
 * API request wrapper with auth token
 * @param {string} method - HTTP method
 * @param {string} url - API endpoint (relative to API_BASE)
 * @param {object|null} data - Request body
 * @returns {Promise<any>} Response data
 */
async function apiRequest(method, url, data = null) {
    const fullUrl = url.startsWith('http') ? url : API_BASE + url;
    const options = {
        method: method.toUpperCase(),
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'same-origin',
    };
    // Add Authorization header from localStorage (matching vWork pattern)
    const token = localStorage.getItem('auth_token');
    if (token) {
        options.headers['Authorization'] = 'Bearer ' + token;
    }
    if (data && (method !== 'GET' && method !== 'HEAD')) {
        options.body = JSON.stringify(data);
    }
    const response = await fetch(fullUrl, options);
    if (response.status === 401) {
        // Smart 401 handling: check error message before redirecting (matching vWork)
        try {
            const errData = await response.json();
            const errMsg = (errData.message || errData.error || '').toLowerCase();
            if (errMsg.includes('session') || errMsg.includes('expired') || errMsg.includes('invalid') || errMsg.includes('unauthorized') || errMsg === '') {
                localStorage.removeItem('auth_token');
                localStorage.removeItem('user');
                localStorage.removeItem('tenant_id');
                window.location.href = '/login';
                return null;
            }
        } catch (e) {
            localStorage.removeItem('auth_token');
            localStorage.removeItem('user');
            localStorage.removeItem('tenant_id');
            window.location.href = '/login';
            return null;
        }
    }
    const result = await response.json();
    if (!response.ok) {
        const errMsg = result.message || result.error || 'Request failed';
        showToast(errMsg, 'error');
        throw new Error(errMsg);
    }
    return result;
}

/**
 * Show Bootstrap toast notification
 * @param {string} message - Toast message
 * @param {string} type - Toast type: success, error, warning, info
 */
function showToast(message, type = 'success') {
    const container = document.getElementById('toastContainer');
    if (!container) return;

    const icons = {
        success: 'fa-check-circle',
        error: 'fa-exclamation-circle',
        warning: 'fa-exclamation-triangle',
        info: 'fa-info-circle',
    };

    const id = 'toast_' + Date.now();
    const html = `
        <div id="${id}" class="toast toast-${type}" role="alert" aria-live="assertive" aria-atomic="true" data-bs-delay="4000">
            <div class="toast-body d-flex align-items-center gap-2">
                <i class="fas ${icons[type] || icons.info}"></i>
                <span>${message}</span>
                <button type="button" class="btn-close btn-close-white ms-auto" data-bs-dismiss="toast" aria-label="Close" style="font-size:0.65rem;"></button>
            </div>
        </div>
    `;
    container.insertAdjacentHTML('beforeend', html);

    const toastEl = document.getElementById(id);
    const toast = new bootstrap.Toast(toastEl);
    toast.show();
    toastEl.addEventListener('hidden.bs.toast', () => toastEl.remove());
}

/**
 * Format ISO date string to locale display
 * @param {string} dateStr - ISO date string
 * @returns {string} Formatted date
 */
function formatDate(dateStr) {
    if (!dateStr) return '--';
    try {
        const d = new Date(dateStr);
        if (isNaN(d.getTime())) return dateStr;
        return d.getFullYear() + '-' +
            String(d.getMonth() + 1).padStart(2, '0') + '-' +
            String(d.getDate()).padStart(2, '0') + ' ' +
            String(d.getHours()).padStart(2, '0') + ':' +
            String(d.getMinutes()).padStart(2, '0');
    } catch {
        return dateStr;
    }
}

/**
 * Format money amount
 * @param {number} amount - Money amount
 * @param {string} currency - Currency symbol
 * @returns {string} Formatted money
 */
function formatMoney(amount, currency = '¥') {
    if (amount === null || amount === undefined) return currency + '0.00';
    return currency + Number(amount).toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

/**
 * Render pagination controls
 * @param {number} total - Total items
 * @param {number} page - Current page
 * @param {number} limit - Items per page
 * @param {string} containerId - Container element ID
 * @param {function} callback - Page change callback
 */
function renderPagination(total, page, limit, containerId, callback) {
    const container = document.getElementById(containerId);
    if (!container) return;

    const totalPages = Math.ceil(total / limit) || 1;
    const start = (page - 1) * limit + 1;
    const end = Math.min(page * limit, total);

    let html = '<div class="pagination-wrapper">';
    const paginationText = (typeof I18n !== 'undefined' && I18n.t)
        ? I18n.t('common.pagination').replace('{start}', total > 0 ? start : 0).replace('{end}', end).replace('{total}', total)
        : `顯示 ${total > 0 ? start : 0}-${end} / 共 ${total} 筆`;
    html += `<span class="pagination-info">${paginationText}</span>`;
    html += '<nav><ul class="pagination mb-0">';

    // Previous
    html += `<li class="page-item${page <= 1 ? ' disabled' : ''}">
        <a class="page-link" href="#" data-page="${page - 1}">&laquo;</a>
    </li>`;

    // Page numbers
    const maxVisible = 5;
    let startPage = Math.max(1, page - Math.floor(maxVisible / 2));
    let endPage = Math.min(totalPages, startPage + maxVisible - 1);
    if (endPage - startPage < maxVisible - 1) {
        startPage = Math.max(1, endPage - maxVisible + 1);
    }

    if (startPage > 1) {
        html += `<li class="page-item"><a class="page-link" href="#" data-page="1">1</a></li>`;
        if (startPage > 2) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
    }

    for (let i = startPage; i <= endPage; i++) {
        html += `<li class="page-item${i === page ? ' active' : ''}">
            <a class="page-link" href="#" data-page="${i}">${i}</a>
        </li>`;
    }

    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
        html += `<li class="page-item"><a class="page-link" href="#" data-page="${totalPages}">${totalPages}</a></li>`;
    }

    // Next
    html += `<li class="page-item${page >= totalPages ? ' disabled' : ''}">
        <a class="page-link" href="#" data-page="${page + 1}">&raquo;</a>
    </li>`;

    html += '</ul></nav></div>';
    container.innerHTML = html;

    // Bind click events
    container.querySelectorAll('.page-link[data-page]').forEach(link => {
        link.addEventListener('click', function (e) {
            e.preventDefault();
            const p = parseInt(this.dataset.page);
            if (p >= 1 && p <= totalPages && p !== page) {
                callback(p);
            }
        });
    });
}

/**
 * Debounce utility
 * @param {function} fn - Function to debounce
 * @param {number} delay - Delay in ms
 * @returns {function} Debounced function
 */
function debounce(fn, delay = 300) {
    let timer;
    return function (...args) {
        clearTimeout(timer);
        timer = setTimeout(() => fn.apply(this, args), delay);
    };
}

/**
 * Confirmation dialog
 * @param {string} message - Confirmation message
 * @returns {boolean} User confirmed
 */
function confirmDelete(message) {
    const defaultMsg = (typeof I18n !== 'undefined' && I18n.t) ? I18n.t('common.confirm_delete') : '確定要刪除嗎？此操作無法撤銷。';
    return confirm(message || defaultMsg);
}

/**
 * Load table data from API and render
 * @param {string} url - API endpoint
 * @param {string} tableBodyId - Table body element ID
 * @param {function} renderRow - Function to render a single row (item) => HTML
 * @param {string} paginationId - Pagination container ID
 * @param {number} page - Current page
 * @param {number} limit - Items per page
 * @param {string} extraParams - Additional query params
 */
async function loadTableData(url, tableBodyId, renderRow, paginationId, page = 1, limit = 20, extraParams = '') {
    const tbody = document.getElementById(tableBodyId);
    if (!tbody) return;

    const _t = (key, fallback) => (typeof I18n !== 'undefined' && I18n.t) ? I18n.t(key) : fallback;
    tbody.innerHTML = '<tr><td colspan="20" class="text-center py-4"><div class="spinner-border spinner-border-sm text-primary" role="status"></div> ' + _t('common.loading', '載入中...') + '</td></tr>';

    try {
        const separator = url.includes('?') ? '&' : '?';
        const fullUrl = `${url}${separator}page=${page}&limit=${limit}${extraParams ? '&' + extraParams : ''}`;
        const data = await apiRequest('GET', fullUrl);

        const items = data.data || data.items || data || [];
        const total = data.total || data.count || items.length || 0;

        if (items.length === 0) {
            tbody.innerHTML = '<tr><td colspan="20" class="text-center py-4 text-muted">' + _t('common.no_data', '暫無數據') + '</td></tr>';
        } else {
            tbody.innerHTML = items.map(renderRow).join('');
        }

        if (paginationId) {
            renderPagination(total, page, limit, paginationId, (newPage) => {
                loadTableData(url, tableBodyId, renderRow, paginationId, newPage, limit, extraParams);
            });
        }

        return { items, total };
    } catch (err) {
        tbody.innerHTML = '<tr><td colspan="20" class="text-center py-4 text-danger">' + _t('common.load_failed_retry', '載入失敗，請重試') + '</td></tr>';
        console.error('Failed to load table data:', err);
        return null;
    }
}

/**
 * Get status badge HTML
 * @param {string} status - Status string
 * @returns {string} Badge HTML
 */
function getStatusBadge(status) {
    const _t = (key, fallback) => (typeof I18n !== 'undefined' && I18n.t) ? I18n.t(key) : fallback;
    const statusMap = {
        'pending': { class: 'badge-pending', label: _t('common.status_pending', '待處理') },
        'processing': { class: 'badge-processing', label: _t('common.status_processing', '處理中') },
        'shipping': { class: 'badge-shipping', label: _t('common.status_shipping', '運送中') },
        'shipped': { class: 'badge-shipping', label: _t('common.status_shipped', '已發貨') },
        'delivered': { class: 'badge-delivered', label: _t('common.status_delivered', '已送達') },
        'completed': { class: 'badge-delivered', label: _t('common.status_completed', '已完成') },
        'cancelled': { class: 'badge-cancelled', label: _t('common.status_cancelled', '已取消') },
        'paid': { class: 'badge-paid', label: _t('common.status_paid', '已付款') },
        'unpaid': { class: 'badge-unpaid', label: _t('common.status_unpaid', '未付款') },
        'active': { class: 'badge-active', label: _t('common.status_enabled', '啟用') },
        'inactive': { class: 'badge-inactive', label: _t('common.status_disabled', '停用') },
        'received': { class: 'badge-delivered', label: _t('common.status_stored', '已入庫') },
        'in_warehouse': { class: 'badge-processing', label: _t('common.status_in_warehouse', '倉庫中') },
        'dispatched': { class: 'badge-shipping', label: _t('common.status_out', '已出庫') },
        'returned': { class: 'badge-cancelled', label: _t('common.status_returned', '已退回') },
        'approved': { class: 'badge-delivered', label: _t('common.status_reviewed', '已審核') },
        'rejected': { class: 'badge-cancelled', label: _t('common.status_rejected', '已拒絕') },
    };
    const s = statusMap[status] || { class: 'badge-inactive', label: status || _t('common.status_unknown', '未知') };
    return `<span class="badge ${s.class}">${s.label}</span>`;
}

/**
 * Logout user (clear both localStorage and cookies, matching vWork pattern)
 */
async function logout() {
    try {
        const token = localStorage.getItem('auth_token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        await fetch(API_BASE + '/auth/logout', {
            method: 'POST',
            headers: headers,
            credentials: 'same-origin',
        });
    } catch (e) {
        // ignore
    }
    // Clear localStorage
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
    localStorage.removeItem('tenant_id');
    // Clear cookie
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
    window.location.href = '/login';
}

/**
 * Open a Bootstrap modal
 * @param {string} modalId - Modal element ID
 */
function openModal(modalId) {
    const el = document.getElementById(modalId);
    if (el) {
        const modal = bootstrap.Modal.getOrCreateInstance(el);
        modal.show();
    }
}

/**
 * Close a Bootstrap modal
 * @param {string} modalId - Modal element ID
 */
function closeModal(modalId) {
    const el = document.getElementById(modalId);
    if (el) {
        const modal = bootstrap.Modal.getInstance(el);
        if (modal) modal.hide();
    }
}

/**
 * Reset a form
 * @param {string} formId - Form element ID
 */
function resetForm(formId) {
    const form = document.getElementById(formId);
    if (form) form.reset();
}

/**
 * Get form data as object
 * @param {string} formId - Form element ID
 * @returns {object} Form data
 */
function getFormData(formId) {
    const form = document.getElementById(formId);
    if (!form) return {};
    const formData = new FormData(form);
    const data = {};
    for (const [key, value] of formData.entries()) {
        // Handle numeric fields
        const input = form.querySelector(`[name="${key}"]`);
        if (input && (input.type === 'number' || input.dataset.type === 'number')) {
            data[key] = value === '' ? 0 : Number(value);
        } else {
            data[key] = value;
        }
    }
    return data;
}

/**
 * Populate form fields from data object
 * @param {string} formId - Form element ID
 * @param {object} data - Data object
 */
function populateForm(formId, data) {
    const form = document.getElementById(formId);
    if (!form || !data) return;
    Object.keys(data).forEach(key => {
        const input = form.querySelector(`[name="${key}"]`);
        if (input) {
            if (input.type === 'checkbox') {
                input.checked = !!data[key];
            } else if (input.type === 'select-one') {
                input.value = data[key] || '';
            } else {
                input.value = data[key] !== null && data[key] !== undefined ? data[key] : '';
            }
        }
    });
}

/**
 * Generate query string from filter object
 * @param {object} filters - Filter key-value pairs
 * @returns {string} Query string without leading ?
 */
function buildQueryParams(filters) {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
        if (value !== '' && value !== null && value !== undefined) {
            params.append(key, value);
        }
    });
    return params.toString();
}
