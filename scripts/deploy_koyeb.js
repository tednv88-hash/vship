// scripts/deploy_koyeb.js
// Usage:
//   node scripts/deploy_koyeb.js
// Requires env vars:
//   KOYEB_TOKEN       - Koyeb API token (https://app.koyeb.com/account/api)
//   GITHUB_REPO       - e.g. "tednv/vship" (must be public or Koyeb GitHub-app installed)
//   DATABASE_URL      - Neon connection string (postgres://...?sslmode=require)
//   JWT_SECRET        - any long random string
// Optional:
//   KOYEB_APP_NAME    - default "vship"
//   KOYEB_REGION      - default "was" (Washington)
//   KOYEB_INSTANCE    - default "free"
//   GITHUB_BRANCH     - default "main"

const https = require('https');

const TOKEN = process.env.KOYEB_TOKEN;
const REPO = process.env.GITHUB_REPO;
const DATABASE_URL = process.env.DATABASE_URL;
const JWT_SECRET = process.env.JWT_SECRET || 'change-me-' + Math.random().toString(36).slice(2);
const APP_NAME = process.env.KOYEB_APP_NAME || 'vship';
const REGION = process.env.KOYEB_REGION || 'was';
const INSTANCE = process.env.KOYEB_INSTANCE || 'free';
const BRANCH = process.env.GITHUB_BRANCH || 'main';

if (!TOKEN || !REPO || !DATABASE_URL) {
  console.error('Missing required env: KOYEB_TOKEN, GITHUB_REPO, DATABASE_URL');
  process.exit(1);
}

function api(method, path, body) {
  return new Promise((resolve, reject) => {
    const data = body ? JSON.stringify(body) : null;
    const req = https.request({
      method,
      hostname: 'app.koyeb.com',
      path: '/v1' + path,
      headers: {
        'Authorization': 'Bearer ' + TOKEN,
        'Content-Type': 'application/json',
        ...(data ? { 'Content-Length': Buffer.byteLength(data) } : {}),
      },
    }, (res) => {
      let buf = '';
      res.on('data', c => buf += c);
      res.on('end', () => {
        let parsed;
        try { parsed = JSON.parse(buf); } catch { parsed = buf; }
        if (res.statusCode >= 400) {
          reject(new Error(`HTTP ${res.statusCode} ${method} ${path}: ${buf}`));
        } else resolve(parsed);
      });
    });
    req.on('error', reject);
    if (data) req.write(data);
    req.end();
  });
}

(async () => {
  console.log('=> Looking up app:', APP_NAME);
  let app;
  const apps = await api('GET', `/apps?name=${APP_NAME}`);
  if (apps.apps && apps.apps.length) {
    app = apps.apps[0];
    console.log('   exists, id=', app.id);
  } else {
    console.log('=> Creating app');
    const r = await api('POST', '/apps', { name: APP_NAME });
    app = r.app;
    console.log('   created, id=', app.id);
  }

  const serviceDef = {
    app_id: app.id,
    definition: {
      name: 'api',
      type: 'WEB',
      git: {
        repository: 'github.com/' + REPO,
        branch: BRANCH,
        no_deploy_on_push: false,
        buildpack: { build_command: '', run_command: '' },
        // use Dockerfile in repo root
        docker: { dockerfile: 'Dockerfile' },
      },
      instance_types: [{ type: INSTANCE }],
      regions: [REGION],
      scalings: [{ min: 1, max: 1 }],
      ports: [{ port: 3002, protocol: 'http' }],
      routes: [{ path: '/', port: 3002 }],
      health_checks: [{ http: { path: '/healthz', port: 3002 }, grace_period: 30 }],
      env: [
        { key: 'DATABASE_URL', value: DATABASE_URL },
        { key: 'JWT_SECRET', value: JWT_SECRET },
        { key: 'APP_ENV', value: 'production' },
        { key: 'PORT', value: '3002' },
      ],
    },
  };

  console.log('=> Looking up service: api');
  const services = await api('GET', `/services?app_id=${app.id}`);
  let svc = (services.services || []).find(s => s.name === 'api');
  if (svc) {
    console.log('   exists, updating, id=', svc.id);
    await api('PUT', `/services/${svc.id}`, serviceDef);
  } else {
    console.log('=> Creating service');
    const r = await api('POST', '/services', serviceDef);
    svc = r.service;
    console.log('   created, id=', svc.id);
  }

  console.log('=> Service deploying. Domain will appear at:');
  console.log(`   https://${APP_NAME}-<org>.koyeb.app`);
  console.log('=> Poll status:');
  for (let i = 0; i < 60; i++) {
    await new Promise(r => setTimeout(r, 5000));
    const s = await api('GET', `/services/${svc.id}`);
    const status = s.service.status;
    console.log(`   [${i}] status=${status}`);
    if (status === 'HEALTHY' || status === 'DEGRADED') break;
    if (status === 'ERROR') { console.error('Service errored'); process.exit(1); }
  }

  // fetch domain
  const domains = await api('GET', `/domains?app_id=${app.id}`);
  const d = (domains.domains || []).find(x => x.type === 'AUTOASSIGNED');
  if (d) console.log('PUBLIC_DOMAIN=https://' + d.name);
})().catch(e => { console.error(e); process.exit(1); });
