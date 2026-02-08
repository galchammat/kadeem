# Supabase Authentication

Direct JWT authentication. Frontend handles OAuth via Supabase JS. JWT sent in `Authorization` header to Go API, validated with `SUPABASE_JWT_SECRET`.

## Setup

1. Create Supabase project at [supabase.com](https://supabase.com)
2. From **Settings > API**, grab Project URL, Publishable key, and JWT Secret
3. In **Authentication > URL Configuration**, set Site URL and Redirect URLs
4. Enable providers in **Authentication > Providers**

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

- **401 Unauthorized**: Token expired/invalid. User must re-login. Verify `SUPABASE_JWT_SECRET` correct.
- **403 Forbidden**: User accessing another's resource. Check `user_id` matches JWT `sub` claim.
