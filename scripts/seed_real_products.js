// scripts/seed_real_products.js — seed realistic warehouse/集运 service products
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

const img = (label, bg = '0f3a57', fg = 'ffffff') =>
  `https://placehold.co/600x600/${bg}/${fg}/png?text=${encodeURIComponent(label).replace(/%20/g, '+')}&font=roboto`;

const products = [
  // ========== 集运服务套餐 ==========
  {
    name: '美国海运集运（普货）',
    description: '美国洛杉矶仓 → 中国，海运拼箱 35-45 天送达。基础运费按实重计算，含基础保险。',
    image_url: img('USA SEA'),
    price: 38.00, original_price: 48.00,
    stock: 9999, unit: 'kg', weight: 1.0,
    sort_order: 1,
  },
  {
    name: '美国空运集运（快线）',
    description: '美国洛杉矶仓 → 中国，空运专线 7-12 天送达。适合电子产品、化妆品、保健品。',
    image_url: img('USA AIR', 'd97706'),
    price: 78.00, original_price: 95.00,
    stock: 9999, unit: 'kg', weight: 1.0,
    sort_order: 2,
  },
  {
    name: '日本海运集运',
    description: '日本东京仓 → 中国，海运 15-25 天。日妆、母婴产品首选。',
    image_url: img('JP SEA', 'dc2626'),
    price: 28.00, original_price: 35.00,
    stock: 9999, unit: 'kg', weight: 1.0,
    sort_order: 3,
  },
  {
    name: '日本空运集运',
    description: '日本东京仓 → 中国，空运 5-8 天。适合急件、易碎品。',
    image_url: img('JP AIR', 'be185d'),
    price: 58.00, original_price: 72.00,
    stock: 9999, unit: 'kg', weight: 1.0,
    sort_order: 4,
  },
  {
    name: '英国海运集运',
    description: '英国伦敦仓 → 中国，海运 40-55 天。奢侈品、英伦品牌专线。',
    image_url: img('UK SEA', '1e40af'),
    price: 45.00, original_price: 58.00,
    stock: 9999, unit: 'kg', weight: 1.0,
    sort_order: 5,
  },
  {
    name: '韩国空运集运',
    description: '韩国首尔仓 → 中国，空运 4-7 天。韩妆、明星同款专线。',
    image_url: img('KR AIR', '7c3aed'),
    price: 48.00, original_price: 62.00,
    stock: 9999, unit: 'kg', weight: 1.0,
    sort_order: 6,
  },

  // ========== 仓储增值服务 ==========
  {
    name: '免费仓储 60 天',
    description: '所有海外仓提供 60 天免费仓储期，超期 0.5 元/件/天。',
    image_url: img('60 DAYS FREE', '059669'),
    price: 0.00, original_price: 0.00,
    stock: 9999, unit: '件', weight: 0,
    sort_order: 10,
  },
  {
    name: '专业打包加固',
    description: '木箱打包 / 气泡膜加固 / 防摔防潮，适合易碎品、高价值物品。',
    image_url: img('PACKING', '7c2d12'),
    price: 25.00, original_price: 35.00,
    stock: 9999, unit: '件', weight: 0.5,
    sort_order: 11,
  },
  {
    name: '商品质检拍照',
    description: '入仓时质检并拍摄 5 张高清照片，确认无损后转运。',
    image_url: img('QC PHOTO', '0e7490'),
    price: 8.00, original_price: 12.00,
    stock: 9999, unit: '件', weight: 0,
    sort_order: 12,
  },
  {
    name: '拆箱合并打包',
    description: '将多个包裹拆箱合并为一个，节省 30-50% 运费。',
    image_url: img('CONSOLIDATE', '9333ea'),
    price: 15.00, original_price: 20.00,
    stock: 9999, unit: '次', weight: 0,
    sort_order: 13,
  },

  // ========== 保险/会员 ==========
  {
    name: '高额运输保险',
    description: '包裹价值 5%，最高保 10000 元。损失 24 小时内全额理赔。',
    image_url: img('INSURANCE', '991b1b'),
    price: 50.00, original_price: 80.00,
    stock: 9999, unit: '份', weight: 0,
    sort_order: 20,
  },
  {
    name: 'VIP 会员（年卡）',
    description: '全线运费 8 折 + 免费仓储 90 天 + 优先打包 + 专属客服。',
    image_url: img('VIP MEMBER', 'fbbf24', '0f3a57'),
    price: 388.00, original_price: 588.00,
    stock: 9999, unit: '张', weight: 0,
    sort_order: 21,
  },
];

(async () => {
  const lr = await req('POST', '/api/v1/auth/login', { email: EMAIL, password: PWD });
  const token = lr.token;
  console.log('logged in');

  // Soft-delete existing goods (clear old picsum ones)
  const oldList = await req('GET', '/api/v1/goods?limit=100', null, token);
  const olds = oldList.data || oldList || [];
  for (const g of olds) {
    try { await req('DELETE', `/api/v1/goods/${g.id}`, null, token); } catch {}
  }
  console.log(`removed ${olds.length} old goods`);

  for (const p of products) {
    try {
      await req('POST', '/api/v1/goods', { ...p, status: 'active' }, token);
      console.log('created:', p.name);
    } catch (e) {
      console.log('SKIP', p.name, e.message);
    }
  }
  console.log('DONE -', products.length, 'products');
})().catch(e => { console.error(e); process.exit(1); });
