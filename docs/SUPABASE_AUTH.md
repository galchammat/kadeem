# Supabase Authentication

Auth is proxy-driven. oauth2-proxy handles the full OIDC flow against Supabase -- the frontend never touches Supabase auth directly. The Go API reads trusted headers set by the proxy.

**Flow:** User hits protected route on `api.cyanlab.cc` → oauth2-proxy redirects to Supabase login → user authenticates → Supabase callbacks to `/oauth2/callback` → oauth2-proxy creates session cookie → requests forwarded to Go API with `X-Auth-User-Id`, `X-Auth-Email`, `X-Auth-Role` headers.

Auth middleware: `packages/server/internal/api/middleware/auth.go`

## Supabase Project Setup

1. Create project at [supabase.com](https://supabase.com)
2. From **Settings > API > API Keys** tab, grab:
   - **Project URL** → `SUPABASE_URL`
   - **Publishable key** (`sb_publishable_...`) → `SUPABASE_PUBLISHABLE_KEY`
3. In **Settings > JWT Keys**, migrate to asymmetric signing keys (ES256)
4. In **Authentication > URL Configuration**, set:
   - Site URL: `https://cyanlab.cc`
   - Redirect URLs: `https://cyanlab.cc/auth/callback`, `https://api.cyanlab.cc/oauth2/callback`
5. Enable providers in **Authentication > Providers** (email, Google, GitHub, etc.)

oauth2-proxy uses PKCE (`code_challenge_method = "S256"`), so no `client_secret` is needed.

## User Roles

Two roles: **user** (default) and **admin**.

- **user** -- can only access their own resources (filtered by `user_id` column)
- **admin** -- bypasses `user_id` filtering, sees all resources

### Setting a User as Admin

**SQL** (in Supabase SQL Editor):

```sql
UPDATE auth.users
SET raw_user_meta_data = jsonb_set(
    COALESCE(raw_user_meta_data, '{}'::jsonb),
    '{role}',
    '"admin"'
)
WHERE email = 'admin@example.com';
```

**Dashboard:** Authentication > Users > click user > User Metadata > add `{"role": "admin"}` > Save.

User must log out and back in to refresh their JWT after a role change.

## Environment Variables

```bash
SUPABASE_URL=https://xxxxxxxxxxxxx.supabase.co
SUPABASE_PUBLISHABLE_KEY=sb_publishable_...
OAUTH2_PROXY_COOKIE_SECRET=$(openssl rand -base64 32)
DOMAIN=api.cyanlab.cc
FRONTEND_DOMAIN=cyanlab.cc
```

oauth2-proxy is provisioned by Ansible (`ansible/roles/nginx/tasks/oauth2_proxy.yml`). Config template: `ansible/roles/nginx/templates/oauth2-proxy.cfg.j2`.

## Troubleshooting

**401 Unauthorized** -- Token invalid or expired. User needs to re-login. Verify `SUPABASE_PUBLISHABLE_KEY` is correct and JWT signing keys are properly rotated in Supabase dashboard.

**403 Forbidden** -- User trying to access another user's resource. Check `user_id` in DB matches the JWT `sub` claim. If admin role isn't working, user needs to re-login after the role change.

**CORS errors** -- Check nginx config allows `https://cyanlab.cc`. Verify frontend isn't using `localhost` against prod.

**Logs:**

```bash
journalctl -u oauth2-proxy -f
journalctl -u kadeem-daemon -f
```
