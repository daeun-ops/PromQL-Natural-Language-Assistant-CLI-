#!/bin/bash
set -euo pipefail

# 다운로드 위치 또는 PATH에 prometheus가 있다고 가정하에 진햇ㅇ
# 현재  루트에서 실행한다고 가정
if [[ ! -f "./ops/prometheus.yml" ]]; then
  echo "Run from repo root. Missing ./ops/prometheus.yml"
  exit 1
fi

echo "Starting Prometheus on 0.0.0.0:9090 ..."
./prometheus --config.file=./ops/prometheus.yml --web.listen-address="0.0.0.0:9090"

