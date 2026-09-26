# Cutover wizard (#812)

**One-off.** Delete this directory once the two-week rollback window after
cutover closes (ADR-0015). It is kept in the repo, rather than as a throwaway
script, only because #812 requires a rehearsal on a copy of production before
the real night. Both runs must use the same script.

`cutover.sh` walks the cutover night stage by stage (the `wizard` skill,
`.agents/skills/wizard`). Every stage ends in a check, either a command or a
question you answer, and a failed check stops the wizard. Re-running is safe.
The database steps are idempotent, and non-secret state (chart versions, the
legacy Helm revision to roll back to, duplicate counts) is saved in
`~/.local/state/hexlet-basics-cutover*.env`, outside the repo.

## Run

```sh
# 1. Rehearsal: the database half against a COPY of production, then legacy and
#    Go booted locally on that copy. Needs psql, atlas (mise), yq, gh.
REHEARSAL=1 scripts/oneoff/cutover/cutover.sh

# 2. The night: production DB + cluster. Adds helm, kubectl, sops, age and the
#    helm-secrets plugin (make -C k8s k8s-utils-install).
scripts/oneoff/cutover/cutover.sh
```

## What it does, in order

1. Preflight: tools, a clean release checkout, and the state of the tickets
   the cutover depends on (#793 #765 #767 #768 #795 #797 #811).
2. Cluster context, plus the rollback point: the live legacy Helm revision.
   Kubernetes must be at least 1.30, and both GHCR images must exist.
3. Postbox domain and static key (#793).
4. Cluster secrets through `make -C k8s secrets-edit`. The check prints key
   names only. `JWT_SECRET` and `EMAIL_TOKEN_SECRET` are compared in memory
   against the public development values. `SITE_URL` and `PUBLIC_URL` are
   checked (#797, #808).
5. Database access. The URL is entered hidden, and it is never written or
   echoed.
6. Duplicate enrollments (#765): `count.sql`, then its decision rule, then
   `collapse.sql`.
7. Duplicate blog likes (#799): `blog-likes-count.sql`, then
   `blog-likes-collapse.sql`.
8. atlas: `migrate apply --baseline 20260727053619` by hand, then
   `schema-verify.sql`.
9. Blog bodies: ActionText goes into `rich_body` (`blog-bodies-*.sql`). Posts
   with `<action-text-attachment>` are listed for a hand fix.
10. Legacy is still healthy on the migrated schema. This is what a rollback
    relies on.
11. Upload `book.pdf` to the bucket (#804).
12. A server-side dry run of the rendered chart.
13. `make -C k8s helm-upgrade-app`, which is the switch, followed by the rollout
    status of each deployment.
14. Routing: `/`, `/api`, both feed URLs (#805) and the `/webhooks/github`
    alias. The legacy nginx 301s are recorded but not enforced, because #811
    left them as an open decision.
15. Password sign-in, Magic Link email (the `account_email` River job), and the
    book download.
16. A lesson check on `api-check`, which also warms the course images (#767).
17. GitHub webhook redelivery.
18. A lead reaching amoCRM (the `amocrm_lead` River job).
19. Sentry, the blog hand-fixes, and the rollback card.

## Why atlas runs by hand and not in the Helm hook

`make -C k8s helm-upgrade-app-baseline` would migrate and switch traffic in a
single command. The blog-body copy needs the `rich_body` column before the
switch, and legacy must be seen working on the migrated schema before anything
depends on it. Running atlas by hand separates those steps. After the hand run
the revisions table exists, so the deploy uses plain `helm-upgrade-app`, and its
migrate hook finds nothing to apply. **Never** run `helm-upgrade-app-baseline`
after this wizard has migrated.

## Deploy and rollback

The #811 chart upgrades the existing `codebasics` release in place: the same
Services, the same Secret and the same ALB group. The Go stack therefore cannot
run next to legacy on its own host. The pre-switch confidence comes from three
places: the rehearsal (legacy and Go both run on the migrated copy), stage 10,
and the server-side dry run. Rollback is
`helm rollback codebasics <recorded revision> -n codebasics`, which restores the
legacy workloads and ingress. Migrations stay applied. A rollback runs no
pre-upgrade hook, so legacy's `rails db:prepare` does not fire.
