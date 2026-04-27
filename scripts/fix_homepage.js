// scripts/fix_homepage.js
// Replaces homepage navbar items with clean emoji/iconfont and uses better banners

const https = require('https');
const BASE = 'https://vship-api.onrender.com';
const EMAIL = 'admin@guoyun.com';
const PWD = 'Admin123!';

function req(method, path, body, token) {
  return new Promise((resolve, reject) => {
    const url = new URL(BASE + path);
    const data = body ? JSON.stringify(body) : null;
    const headers = { Accept: 'application/json' };
    if (data) { headers['Content-Type'] = 'application/json'; headers['Content-Length'] = Buffer.byteLength(data); }
    if (token) headers.Authorization = 'Bearer ' + token;
    const r = https.request({ method, hostname: url.hostname, path: url.pathname + url.search, headers }, (res) => {
      let buf = ''; res.on('data', c => buf += c);
      res.on('end', () => { let p; try { p = JSON.parse(buf); } catch { p = buf; }
        if (res.statusCode >= 400) reject(new Error(`HTTP ${res.statusCode}: ${buf}`));
        else resolve(p);
      });
    });
    r.on('error', reject); if (data) r.write(data); r.end();
  });
}

// Use iconify CDN — clean monochrome SVG icons
const ic = (name, color = 'fbbf24') =>
  `https://api.iconify.design/material-symbols/${name}.svg?color=%23${color}&width=64`;

(async () => {
  const lr = await req('POST', '/api/v1/auth/login', { email: EMAIL, password: PWD });
  const token = lr.token;

  const list = await req('GET', '/api/v1/page-designs', null, token);
  const designs = list.data || list;
  const home = (designs || []).find(d => d.type === 'home' && d.is_default);
  if (!home) { console.error('no home page'); process.exit(1); }

  const newPageData = {
    page: {
      id: 'page', type: 'page', name: '页面设置',
      params: { name: '国韵好运仓储', title: '国韵好运 · 仓储集运', share_title: '国韵好运仓储集运' },
      style: { titleBackgroundColor: '#0f3a57', titleTextColor: 'white' },
    },
    items: [
      {
        type: 'banner', name: '轮播广告',
        data: [
          { imgUrl: 'https://api.iconify.design/material-symbols/local-shipping-rounded.svg?color=%23ffffff&width=750&height=300', linkUrl: '' },
        ],
        params: { interval: '3000' },
        style: { btnColor: '#fbbf24', btnShape: 'round' },
      },
      {
        type: 'notice', name: '公告通知',
        params: { icon: '', text: '欢迎使用 国韵好运 — 海外仓储 · 集运转运一站式服务' },
        style: { background: '#ffffff', paddingTop: '8', textColor: '#333333' },
      },
      {
        type: 'navBar', name: '导航菜单',
        data: [
          { color: '#333', imgUrl: ic('inventory-2-outline-rounded'), linkUrl: '/pages/forecast/index', text: '入仓预报' },
          { color: '#333', imgUrl: ic('inbox-outline-rounded'),       linkUrl: '/pages/package/index',  text: '仓内包裹' },
          { color: '#333', imgUrl: ic('receipt-long-outline-rounded'),linkUrl: '/pages/order/list',    text: '集运订单' },
          { color: '#333', imgUrl: ic('account-balance-wallet-outline-rounded'), linkUrl: '/pages/recharge/index', text: '充值中心' },
          { color: '#333', imgUrl: ic('warehouse-outline-rounded'),   linkUrl: '/pages/warehouse/address', text: '海外仓地址' },
          { color: '#333', imgUrl: ic('help-outline-rounded'),        linkUrl: '/pages/help/index',     text: '帮助中心' },
          { color: '#333', imgUrl: ic('confirmation-number-outline-rounded'), linkUrl: '/pages/coupon/my', text: '优惠券' },
          { color: '#333', imgUrl: ic('calculate-outline-rounded'),   linkUrl: '/pages/estimate/index',    text: '仓储计费' },
        ],
        style: { background: '#ffffff', rowsNum: '4' },
      },
      {
        type: 'goods', name: '推荐服务',
        data: [],
        params: { auto: { category: 0, goodsSort: 'all', showNum: 8 }, source: 'auto' },
        style: { background: '#f5f6f7', column: '2', display: 'list',
                 show: { goodsName: '1', goodsPrice: '1', goodsSales: '1', linePrice: '1', sellingPoint: '0' } },
      },
    ],
  };

  await req('PUT', `/api/v1/page-designs/${home.id}`, { page_data: newPageData }, token);
  console.log('homepage updated');

  // Replace banner image URLs with cleaner gradient placeholders too
  const bannersResp = await req('GET', '/api/v1/banners', null, token);
  const bs = (bannersResp.data || bannersResp || []);
  const cleanImgs = [
    'https://placehold.co/750x300/0f3a57/fbbf24/png?text=%E5%9B%BD%E9%9F%B5%E5%A5%BD%E8%BF%90%E4%BB%93%E5%82%A8+%C2%B7+%E9%9B%86%E8%BF%90%E8%BD%AC%E8%BF%90',
    'https://placehold.co/750x300/d97706/ffffff/png?text=%E6%96%B0%E4%BA%BA%E9%A6%96%E5%8D%95%E4%BB%93%E5%82%A8%E8%B4%B9%E7%AB%8B%E5%87%8F50',
    'https://placehold.co/750x300/059669/ffffff/png?text=%E4%BC%9A%E5%91%98%E6%97%A5+%E5%8F%8C%E5%80%8D%E7%A7%AF%E5%88%86',
  ];
  let i = 0;
  for (const b of bs) {
    if (i >= cleanImgs.length) break;
    try {
      await req('PUT', `/api/v1/banners/${b.id}`, { image_url: cleanImgs[i] }, token);
      console.log('banner updated:', b.title);
    } catch (e) { console.log('skip', b.title, e.message); }
    i++;
  }

  console.log('DONE');
})().catch(e => { console.error(e); process.exit(1); });
