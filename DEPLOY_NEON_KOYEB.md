# vShip Deploy Guide

## 1. Create Neon Postgres

1. Create a free Neon project.
2. Create a database named `vship`.
3. Copy the connection string.
4. Ensure the connection string uses SSL.

Expected env on the app side:

```text
DATABASE_URL=postgresql://USER:PASSWORD@HOST/DB?sslmode=require
```

## 2. Upload This Project To GitHub

1. Create a new private GitHub repository.
2. Upload the full `vship` project.
3. Do not commit your local `.env` file.

## 3. Deploy Backend To Koyeb

1. Create a Koyeb account.
2. Click `Create Web Service`.
3. Choose `GitHub` and select the repository.
4. Keep the service root at the repository root.
5. Koyeb should detect the `Dockerfile` automatically.
6. Set these environment variables:

```text
APP_ENV=production
DATABASE_URL=<your Neon connection string>
JWT_SECRET=<generate-a-long-random-secret>
BASE_URL=https://<your-koyeb-service>.koyeb.app
SERVER_HOST=0.0.0.0
```

7. Set the health check path to `/healthz`.
8. Deploy.

After deployment, note the service URL:

```text
https://<your-koyeb-service>.koyeb.app
```

## 4. Point The Mini Program At The Deployed API

Create `miniprogram/.env.production` with:

```text
VITE_API_BASE=https://<your-koyeb-service>.koyeb.app/api/v1/mp
```

## 5. Build The WeChat Mini Program

From `miniprogram/` run:

```text
npm install
npm run build:mp-weixin
```

The built mini program will be in:

```text
miniprogram/dist/build/mp-weixin
```

## 6. Configure WeChat Backend

If the mini program `AppID` and `AppSecret` are not stored in the account or tenant `extra_fields`, set these fallback environment variables on the backend service:

```text
WECHAT_MINIPROGRAM_APPID=<your-mini-program-appid>
WECHAT_MINIPROGRAM_SECRET=<your-mini-program-secret>
```

Then configure the mini program domain in WeChat:

In `mp.weixin.qq.com` for your mini program:

1. Go to `开发管理` -> `开发设置`.
2. Add the request legal domain:

```text
https://<your-koyeb-service>.koyeb.app
```

3. Save the change.

## 7. Upload From WeChat DevTools

1. Open WeChat DevTools.
2. Import:

```text
miniprogram/dist/build/mp-weixin
```

3. Confirm the AppID is:

```text
wx3fe977cf21ca15e9
```

4. Test login and API calls.
5. Click `上传`.
6. Go back to `mp.weixin.qq.com` and submit for review.

## Notes

- Koyeb free tier is fine for early testing, but it may sleep or throttle.
- Neon free tier is fine for early testing, but it has limits.
- If you later buy a domain, update both `BASE_URL` and `VITE_API_BASE`.
