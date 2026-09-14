#!/usr/bin/env bash

set -euo pipefail

DEPLOY_BRANCH="$1"
GITHUB_RUN_ID="$2"
COMPOSE_FILE="$3"
DUMP_RETENTION_DAYS=7
EXPECTED_SERVICES="api admin mobile cloudflare"

# 1. Stash production changes
# 本番環境で未pushの変更があれば退避する
echo "=== 1. Stash production changes ==="

rm -f /tmp/seeft-stash-created

if test -n "$(git status --porcelain)"; then
  git stash push --include-untracked \
    -m "github-actions-${GITHUB_RUN_ID}"
  touch /tmp/seeft-stash-created
fi

# 2. Pull remote changes
echo "=== 2. Pull remote changes ==="

git check-ref-format --branch "$DEPLOY_BRANCH"
git fetch origin "$DEPLOY_BRANCH"

if git show-ref --verify --quiet "refs/heads/$DEPLOY_BRANCH"; then
  git switch "$DEPLOY_BRANCH"
else
  git switch --track -c "$DEPLOY_BRANCH" \
    "origin/$DEPLOY_BRANCH"
fi

git pull --ff-only origin "$DEPLOY_BRANCH"

# 3. Restore production changes
echo "=== 3. Restore production changes ==="

if test -f /tmp/seeft-stash-created; then
  git stash pop
fi

# 4. Build application images
echo "=== 4. Build application images ==="

docker compose -f "$COMPOSE_FILE" build --pull

# 5. Dump database
# 本番環境にある api/env/seeft.env から接続情報を取得
echo "=== 5. Dump database ==="

set -a
source ./api/env/seeft.env
set +a

mkdir -p /var/backups/pg

PGPASSWORD="$NUTMEG_DB_PASSWORD" \
  pg_dump \
    -h "$NUTMEG_DB_HOST" \
    -p "$NUTMEG_DB_PORT" \
    -U "$NUTMEG_DB_USER" \
    -d "$NUTMEG_DB_NAME" \
    -Fc \
    -f "/var/backups/pg/seeft-$(date '+%Y%m%d-%H%M%S').dump"

# 6. Stop application containers
echo "=== 6. Stop application containers ==="

docker compose -f "$COMPOSE_FILE" stop $EXPECTED_SERVICES

# 7. Run database migration
echo "=== 7. Run database migration ==="

make prod-migrate

# 8. Start application containers
echo "=== 8. Start application containers ==="

docker compose -f "$COMPOSE_FILE" up -d $EXPECTED_SERVICES

# 9. Health check
echo "=== 9. Health check ==="

sleep 5

RUNNING_SERVICES="$(docker compose -f "$COMPOSE_FILE" ps --services --status running)"

for service in $EXPECTED_SERVICES; do
  if ! echo "$RUNNING_SERVICES" | grep -qx "$service"; then
    echo "Service is not running: $service"
    docker compose -f "$COMPOSE_FILE" logs --tail=100 "$service"
    exit 1
  fi
done

docker compose -f "$COMPOSE_FILE" ps

# 10. Prune unused Docker resources
echo "=== 10. Prune unused Docker resources ==="

docker system prune -af

# 11. Remove expired database dumps
echo "=== 11. Remove expired database dumps ==="

find /var/backups/pg \
  -type f \
  -name 'seeft-*.dump' \
  -mtime +"$DUMP_RETENTION_DAYS" \
  -delete
