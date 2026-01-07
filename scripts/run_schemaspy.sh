#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
OUTPUT_DIR="${ROOT_DIR}/api/docs/schemaspy"

mkdir -p "${OUTPUT_DIR}"

# Docker Compose経由でSchemaSpyを実行するラッパー。
# 接続先はPostgreSQL（サービス名 db）を前提にしている。
cd "${ROOT_DIR}"
docker compose run --rm schemaspy "$@"

