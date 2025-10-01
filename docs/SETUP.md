# Setup Guide

## App
1. Copy `.env.example` → `.env.local` and fill values.
2. `pnpm i && pnpm prisma generate && pnpm db:push`
3. `pnpm dev`

## Webhooks
- GitHub: add webhook to `/api/github/webhook` with secret `GITHUB_WEBHOOK_SECRET`.
- GitLab: add webhook to `/api/gitlab/webhook` with token `GITLAB_WEBHOOK_TOKEN`.

## Billing
- Stripe product/price → set env values.
- Webhook: `/api/stripe/webhook` (set `STRIPE_WEBHOOK_SECRET`).

## Auth
- GitHub OAuth app callback: `/api/auth/callback/github`.

## Worker
- Run `scripts/cloudcurio_worker_install.sh` on cbwdellr720 with ZeroTier network ID.
- Set `CONTAINER_IMAGE=ghcr.io/<you>/cloudcurio-review:latest`.
