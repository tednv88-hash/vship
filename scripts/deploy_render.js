// scripts/deploy_render.js
// Deploys the vship backend to Render using the Render REST API.
// Required env:
//   RENDER_API_KEY    - https://dashboard.render.com/account/api-keys
//   GITHUB_REPO_URL   - e.g. https://github.com/tednv88-hash/vship
//   DATABASE_URL      - Neon Postgres URL
//   JWT_SECRET        - any long random string
// Optional:
//   RENDER_SERVICE_NAME (default "vship-api")
//   RENDER_REGION       (default "oregon")
//   RENDER_OWNER_ID     (auto-detected from /v1/owners)
//   GITHUB_BRANCH       (default "main")

const https = require('https');

const KEY = process.env.RENDER_API_KEY;
const REPO = process.env.GITHUB_REPO_URL;
const DB = process.env.DATABASE_URL;
const JWT = process.env.JWT_SECRET || 'jwt-' + require('crypto').randomBytes(24).toString('hex');
const NAME = process.env.RENDER_SERVICE_NAME || 'vship-api';
const REGION = process.env.RENDER_REGION || 'oregon';
const BRANCH = process.env.GITHUB_BRANCH || 'main';

if (!KEY || !REPO || !DB) {
  console.error('Missing RENDER_API_KEY / GITHUB_REPO_URL / DATABASE_URL');
  process.exit(1);
}

function api(method, path, body) {
  return new Promise((resolve, reject) => {
    const data = body ? JSON.stringify(body) : null;
    const req = https.request({
      method,
      hostname: 'api.render.com',
      path: '/v1' + path,
      headers: {
        'Authorization': 'Bearer ' + KEY,
        'Accept': 'application/json',
        ...(data ? { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(data) } : {}),
      },
    }, (res) => {
      let buf = '';
      res.on('data', c => buf += c);
      res.on('end', () => {
        let parsed;
        try { parsed = JSON.parse(buf); } catch { parsed = buf; }
        if (res.statusCode >= 400) reject(new Error(`HTTP ${res.statusCode} ${method} ${path}: ${buf}`));
        else resolve(parsed);
      });
    });
    req.on('error', reject);
    if (data) req.write(data);
    req.end();
  });
}

(async () => {
  let ownerId = process.env.RENDER_OWNER_ID;
  if (!ownerId) {
    console.log('=> Fetching owner');
    const owners = await api('GET', '/owners');
    ownerId = (owners[0] && (owners[0].owner.id || owners[0].id));
    console.log('   owner=', ownerId);
  }

  console.log('=> Looking up service:', NAME);
  const list = await api('GET', `/services?name=${NAME}&limit=20`);
  let svc = (list || []).map(x => x.service || x).find(s => s.name === NAME);

  const envVars = [
    { key: 'DATABASE_URL', value: DB },
    { key: 'JWT_SECRET', value: JWT },
    { key: 'APP_ENV', value: 'production' },
    { key: 'PORT', value: '3002' },
  ];

  if (!svc) {
    console.log('=> Creating service');
    const body = {
      type: 'web_service',
      name: NAME,
      ownerId,
      repo: REPO,
      branch: BRANCH,
      autoDeploy: 'yes',
      serviceDetails: {
        env: 'docker',
        region: REGION,
        plan: 'free',
        dockerfilePath: './Dockerfile',
        envSpecificDetails: { dockerCommand: '' },
        healthCheckPath: '/healthz',
      },
      envVars,
    };
    const r = await api('POST', '/services', body);
    svc = r.service || r;
    console.log('   created id=', svc.id);
  } else {
    console.log('   exists id=', svc.id);
    console.log('=> Updating env vars');
    await api('PUT', `/services/${svc.id}/env-vars`, envVars);
    console.log('=> Triggering deploy');
    await api('POST', `/services/${svc.id}/deploys`, {});
  }

  console.log('=> Polling deploy status');
  for (let i = 0; i < 80; i++) {
    await new Promise(r => setTimeout(r, 6000));
    const deploys = await api('GET', `/services/${svc.id}/deploys?limit=1`);
    const d = (deploys[0] && (deploys[0].deploy || deploys[0]));
    const st = d ? d.status : '?';
    console.log(`   [${i}] ${st}`);
    if (st === 'live' || st === 'build_failed' || st === 'update_failed' || st === 'canceled' || st === 'deactivated') {
      if (st !== 'live') process.exit(2);
      break;
    }
  }

  const fresh = await api('GET', `/services/${svc.id}`);
  const url = (fresh.service && fresh.service.serviceDetails && fresh.service.serviceDetails.url) || fresh.serviceDetails && fresh.serviceDetails.url;
  console.log('PUBLIC_URL=' + url);
})().catch(e => { console.error(e); process.exit(1); });
