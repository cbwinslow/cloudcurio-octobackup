# Supabase Integration in CloudCurio Monorepo

## Setup
- Ansible installs self-hosted Supabase (/opt/supabase).
- Env: Add SUPABASE_URL, SUPABASE_ANON_KEY to .env.local (from Supabase dashboard post-install).
- Client: Added @supabase/supabase-js to package.json.

## Usage
- Storage: Upload scripts/docs to Supabase buckets (e.g., in /api/scripts: supabase.storage.from('scripts').upload()).
- Auth: Optional alongside NextAuth (e.g., email/password via Supabase Auth).
- Realtime: For chat/reviews (supabase.channel() for subscriptions).
- DB: Use Supabase Postgres if replacing Prisma (migrate schema).

Extend worker.py for embeddings upload to Supabase vector extension.

See https://supabase.com/docs/guides/self-hosting for config.
