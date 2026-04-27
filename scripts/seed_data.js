// scripts/seed_data.js
// Seeds banner, sample goods, and replaces the default home page design
// with one that uses real public image URLs (so the home page isn't blank).
//
// Required env:
//   API_BASE   default https://vship-api.onrender.com
//   ADMIN_EMAIL / ADMIN_PASSWORD

const https = require('https');

const BASE = process.env.API_BASE || 'https://vship-api.onrender.com';
const EMAIL = process.env.ADMIN_EMAIL || 'admin@guoyun.com';
const PWD = process.env.ADMIN_PASSWORD || 'Admin123!';

function req(method, path, body, token) {
  return new Promise((resolve, reject) => {
    const url = new URL(BASE + path);
    const data = body ? JSON.stringify(body) : null;
    const headers = { 'Accept': 'application/json' };
    if (data) {
      headers['Content-Type'] = 'application/json';
      headers['Content-Length'] = Buffer.byteLength(data);
    }
    if (token) headers['Authorization'] = 'Bearer ' + token;
    const r = https.request({
      method, hostname: url.hostname, path: url.pathname + url.search, headers,
    }, (res) => {
      let buf = '';
      res.on('data', c => buf += c);
      res.on('end', () => {
        let p; try { p = JSON.parse(buf); } catch { p = buf; }
        if (res.statusCode >= 400) reject(new Error(`HTTP ${res.statusCode} ${method} ${path}: ${buf}`));
        else resolve(p);
      });
    });
    r.on('error', reject);
    if (data) r.write(data);
    r.end();
  });
}

(async () => {
  console.log('=> Login');
  const lr = await req('POST', '/api/v1/auth/login', { email: EMAIL, password: PWD });
  const token = lr.token;
  const tenantId = lr.user.tenant_id;
  console.log('   tenant=', tenantId);

  // -- Banner --
  console.log('=> Create banner');
  const banners = [
    { title: '新人专享 国韵好运首单立减50',
      image_url: 'https://picsum.photos/seed/guoyun1/750/360',
      link_url: '', position: 'home', sort_order: 1, status: 'active' },
    { title: '集运优惠 美国直邮5折起',
      image_url: 'https://picsum.photos/seed/guoyun2/750/360',
      link_url: '', position: 'home', sort_order: 2, status: 'active' },
    { title: '会员日 双倍积分',
      image_url: 'https://picsum.photos/seed/guoyun3/750/360',
      link_url: '', position: 'home', sort_order: 3, status: 'active' },
  ];
  for (const b of banners) {
    try { const r = await req('POST', '/api/v1/banners', b, token); console.log('   +', r.title); }
    catch (e) { console.log('   skip', b.title, e.message); }
  }

  // -- Goods --
  console.log('=> Create goods');
  const goods = [
    { name: '美国转运 标准服务', description: '美国华盛顿仓库 7-12天到港', image_url: 'https://picsum.photos/seed/g1/400/400', price: 25.00, original_price: 35.00, stock: 999, sales_count: 128, unit: 'kg', weight: 1.0, status: 'active' },
    { name: '日本转运 海运空运', description: '日本东京仓库 3-5天到港', image_url: 'https://picsum.photos/seed/g2/400/400', price: 18.00, original_price: 22.00, stock: 999, sales_count: 256, unit: 'kg', weight: 1.0, status: 'active' },
    { name: '韩国转运 美妆专线', description: '韩国首尔仓库 2-4天到港', image_url: 'https://picsum.photos/seed/g3/400/400', price: 16.00, original_price: 20.00, stock: 999, sales_count: 312, unit: 'kg', weight: 1.0, status: 'active' },
    { name: '欧洲转运 德国仓', description: '欧洲德国仓库 10-15天到港', image_url: 'https://picsum.photos/seed/g4/400/400', price: 32.00, original_price: 40.00, stock: 999, sales_count: 88, unit: 'kg', weight: 1.0, status: 'active' },
  ];
  for (const g of goods) {
    try { const r = await req('POST', '/api/v1/goods', g, token); console.log('   +', r.name); }
    catch (e) { console.log('   skip', g.name, e.message); }
  }

  // -- Replace home page design --
  console.log('=> Find home page design');
  const list = await req('GET', '/api/v1/page-designs', null, token);
  const designs = list.data || list;
  const home = (designs || []).find(d => d.type === 'home' && d.is_default);
  if (!home) { console.log('   no home page design found, skipping'); }
  else {
    console.log('   id=', home.id);
    const newPageData = {
      page: {
        id: 'page', type: 'page', name: '页面设置',
        params: { name: '国韵好运首页', title: '国韵好运 - 全球集运', share_title: '国韵好运' },
        style: { titleBackgroundColor: '#0f3a57', titleTextColor: 'white' },
      },
      items: [
        {
          type: 'banner', name: '轮播广告',
          data: banners.map(b => ({ imgUrl: b.image_url, linkUrl: '' })),
          params: { interval: '3000' },
          style: { btnColor: '#fbbf24', btnShape: 'round' },
        },
        {
          type: 'notice', name: '公告通知',
          params: { icon: '', text: '欢迎使用国韵好运 — 全球集运 一站式服务' },
          style: { background: '#ffffff', paddingTop: '8', textColor: '#333333' },
        },
        {
          type: 'navBar', name: '导航菜单',
          data: [
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n1/80/80', linkUrl: '/pages/forecast/index', text: '我的预报' },
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n2/80/80', linkUrl: '/pages/package/list', text: '我的包裹' },
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n3/80/80', linkUrl: '/pages/order/list', text: '我的订单' },
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n4/80/80', linkUrl: '/pages/recharge/index', text: '充值中心' },
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n5/80/80', linkUrl: '/pages/address/list', text: '收货地址' },
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n6/80/80', linkUrl: '/pages/help/list', text: '帮助中心' },
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n7/80/80', linkUrl: '/pages/coupon/list', text: '我的优惠' },
            { color: '#333', imgUrl: 'https://picsum.photos/seed/n8/80/80', linkUrl: '/pages/route/list', text: '运费查询' },
          ],
          style: { background: '#ffffff', rowsNum: '4' },
        },
        {
          type: 'goods', name: '推荐商品',
          data: [],
          params: { auto: { category: 0, goodsSort: 'all', showNum: 8 }, source: 'auto' },
          style: { background: '#f5f6f7', column: '2', display: 'list',
                   show: { goodsName: '1', goodsPrice: '1', goodsSales: '1', linePrice: '1', sellingPoint: '0' } },
        },
      ],
    };
    await req('PUT', `/api/v1/page-designs/${home.id}`, { page_data: newPageData }, token);
    console.log('   updated');
  }

  console.log('DONE');
})().catch(e => { console.error(e); process.exit(1); });
