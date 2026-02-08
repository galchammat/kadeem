# Supabase Authentication

Direct JWT authentication. Frontend handles OAuth via Supabase JS. JWT sent in `Authorization` header to Go API, validated via JWKS (ES256).

## Setup

1. Create Supabase project at [supabase.com](https://supabase.com)
2. From **Settings > API**, grab:
   - Project URL → `VITE_SUPABASE_URL`
   - Publishable key → `VITE_SUPABASE_PUBLISHABLE_KEY`
   - JWKS URL (Settings > API > JWT Settings > JWKS URL) → `SUPABASE_JWKS_URL`
3. In **Authentication > URL Configuration**, set Site URL and Redirect URLs
4. Enable providers in **Authentication > Providers**

## Environment Variables

```bash
# Frontend (Vite)
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_PUBLISHABLE_KEY=sb_publishable_...

# Backend (JWKS URL for JWT validation)
SUPABASE_JWKS_URL=https://your-project.supabase.co/auth/v1/.well-known/jwks.json
```

## User Roles

**user** (default): own resources only. **admin**: all resources.

Set admin in SQL:
```sql
UPDATE auth.users SET raw_user_meta_data = jsonb_set(
  COALESCE(raw_user_meta_data, '{}'::jsonb), '{role}', '"admin"'
) WHERE email = 'admin@example.com';
```

User must re-login for JWT refresh after role change.

## Troubleshooting

- **401 Unauthorized**: Token expired/invalid, or wrong issuer/audience. User must re-login. Verify `SUPABASE_JWKS_URL` is correct.
- **403 Forbidden**: User accessing another's resource. Check `user_id` matches JWT `sub` claim.
