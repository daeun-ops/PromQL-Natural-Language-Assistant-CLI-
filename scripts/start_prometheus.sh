#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT_DIR"

CONF="./ops/prometheus.yml"
BIN="./prometheus"

if [[ ! -f "$CONF" ]]; then
  echo "Missing $CONF"
  exit 1
fi

if ! command -v "$BIN" >/dev/null 2>&1 && ! [[ -x "$BIN" ]]; then
  echo "Prometheus binary not found at $BIN"
  echo "Copy the downloaded 'prometheus' binary here or put it in PATH."
  exit 1
fi

echo "Starting Prometheus on 0.0.0.0:9090 with $CONF"
"$BIN" --config.file="$CONF" --web.listen-address="0.0.0.0:9090"
