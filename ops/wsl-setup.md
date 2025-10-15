## WSL 설치/로컬 실행 가이드 — `ops/wsl-setup.md`

```markdown
# Windows 11 + WSL2 + Prometheus 로컬 테스트 가이드

## 1) WSL2 & Ubuntu 설치
1. PowerShell(관리자)에서:
   ```powershell
   wsl --install -d Ubuntu-22.04
```

1. 설치 후 재부팅, Ubuntu 사용자/비번 설정.

## 필수 도구 설치 (WSL Ubuntu 터미널)

```bash
sudo apt update
sudo apt install -y build-essential curl wget unzip git
sudo apt install -y golang  # Go 1.22+ 권장
go version
```

## Prometheus 다운로드

```bash
cd ~
VER=$(curl -s https://api.github.com/repos/prometheus/prometheus/releases/latest \
 | grep tag_name | cut -d '"' -f4)
wget https://github.com/prometheus/prometheus/releases/download/${VER}/prometheus-${VER#v}.linux-amd64.tar.gz
tar xzf prometheus-*.tar.gz
cd prometheus-*.linux-amd64
```

바이너리를 레포 루트에서 쓰고 싶다면:

```bash
# 레포 루트가 ~/promql-nlq-assistant 라고 가정
cp prometheus ~/promql-nlq-assistant/

```

## 레포 클론 & 환경변수

```bash
cd ~
git clone https://github.com/yourname/promql-nlq-assistant.git
cd promql-nlq-assistant

cp .env.example .env
# .env 편집: OPENAI_API_KEY, PROM_URL=http://localhost:9090
```

## Prometheus 기동

```bash
chmod +x ./scripts/start_prometheus.sh
./scripts/start_prometheus.sh
# 브라우저에서 http://localhost:9090 접속
```

## CLI 실행

새 터미널에서:

```bash
cd ~/promql-nlq-assistant
export PROM_URL=http://localhost:9090
# 선택: export OPENAI_API_KEY=sk-xxxx

go mod init promql-nlq-assistant  # 최초 1회
# PromQL 파서 의존성
go get github.com/prometheus/prometheus@latest
go mod tidy

go run ./cmd/main.go
```

## 문제 해결

- 포트 충돌: 9090 사용 중이면 `scripts/start_prometheus.sh`에서 포트 변경.
- 인증 에러: `.env`의 키 오타/만료 확인.
- WSL 네트워크: `curl http://localhost:9090/-/healthy` 로 확인.

```
---

## 3) 스모크 테스트 프롬프트 — `scripts/smoke_queries.txt`

```text
5분 평균 에러율(5xx) 보여줘
지난 10분 p95 응답시간 엔드포인트별 상위 5개
현재 업/다운 타겟 요약
메모리 사용률 평균
CPU 사용률 상위 5개 인스턴스
```

나중에 cmd/main.go에 “파일에서 한 줄씩 읽어 자동 테스트” 옵션을 붙이면 금방 테스트 가능해짐

---
