# PromQL-Natural-Language-Assistant-CLI-
 Convert natural language into **PromQL** using GPT models — and validate/test it against a local Prometheus.

---
```bash
promql-nlq-assistant/
├── README.md
├── .env.example                 # OPENAI_API_KEY, PROM_URL 등 환경변수
├── go.mod
├── go.sum
├── cmd/
│   └── main.go                  # CLI 
├── internal/
│   ├── llm/
│   │   └── client.go            # OpenAI 호출
│   ├── prom/
│   │   ├── fetcher.go           # Prometheus 메트릭/라벨 수집 
│   │   └── validator.go         # promql/parser 검증 (
│   └── utils/
│       └── formatter.go         # 출력 포맷/하이라이트 
├── ops/
│   ├── prometheus.yml           # 로컬 실행용 Prometheus 설정 (자체 스크래핑+기본 타겟)
│   └── wsl-setup.md             # Windows 11 + WSL2 설치/실행 가이드
├── scripts/
│   ├── start_prometheus.sh      # WSL에서 프로메테우스 실행 스크립트
│   └── smoke_queries.txt        # 터미널 테스트용 자연어 예시 모음
└── examples/
    └── sample_queries.txt       # NL → 생성된 PromQL 예시
```



 ##  What this repo demonstrates
- **Terminal-first** workflow:
  1) Enter a natural language question
  2) Get a generated **PromQL**
  3) Validate with `promql/parser`
  4) (Optional) Run the query against **local Prometheus** and preview results

- **No dashboard needed** — just CLI output.

---

##  Local Test Environment (Windows 11)
We recommend **WSL2 + Ubuntu** for the simplest local setup.  
(Alternative: Ubuntu VM. See notes at the end.)

### Prerequisites
- Windows 11 with **WSL2** enabled
- Ubuntu 22.04 (or 24.04) installed on WSL
- Go 1.22+ in WSL
- Local Prometheus (single binary)

> Detailed step-by-step: see [`ops/wsl-setup.md`](./ops/wsl-setup.md)

---

##  Quickstart (WSL2)

### 1) Clone & env
```bash
git clone https://github.com/yourname/promql-nlq-assistant.git
cd promql-nlq-assistant

# Copy environment template
cp .env.example .env
# Edit .env with your values
#   OPENAI_API_KEY=sk-xxxx
#   PROM_URL=http://localhost:9090

### 2) Start local Prometheus

- Prometheus config is at `ops/prometheus.yml`.
    
    It scrapes itself (`localhost:9090`) so you can test queries immediately.
    

```bash
chmod +x ./scripts/start_prometheus.sh
./scripts/start_prometheus.sh
# Prometheus UI: http://localhost:9090  (from Windows browser too)
```

> Tip: To generate some traffic/metrics, just open the UI and click around;
> 
> 
> Prometheus exposes its own HTTP metrics like `prometheus_http_request_duration_seconds_bucket`.
> 

### 3) Run the CLI (once code is added)

```bash
go mod tidy
go run ./cmd/main.go
```

You’ll be prompted:

```
Enter your question: 최근 5분 동안 p95 HTTP latency가 1초 넘는 엔드포인트 상위 5개ㅔ
```

Expected CLI output flow:

```
[LLM] Generated PromQL:
histogram_quantile(0.95, sum by (le, handler) (rate(prometheus_http_request_duration_seconds_bucket[5m])))

[Validate] OK (promql/parser)

[Query @ http://localhost:9090/api/v1/query]
<json result preview ...>

```

---

## Terminal Smoke Tests

You can copy/paste from `scripts/smoke_queries.txt`:

- "5분 평균 에러율(5xx) 보여줘"
- "지난 10분 p95 응답시간 엔드포인트별 상위 5개"
- "현재 업/다운 타겟 요약"
- "메모리 사용률 평균"

When the CLI is implemented, each line should return a PromQL string, pass validation, and (optionally) return a quick preview from Prometheus.

 Environment Variables

---

See `.env.example`

- `OPENAI_API_KEY` — required for LLM generation
- `PROM_URL` — e.g. `http://localhost:9090`

---

## Repo Layout

```
promql-nlq-assistant/
├── cmd/                 # CLI entry (main.go)
├── internal/llm/        # OpenAI client wrappers
├── internal/prom/       # Prometheus API + validator
├── internal/utils/      # CLI formatting
├── ops/                 # Local ops: prometheus.yml, setup docs
├── scripts/             # Shell scripts & test prompts
└── examples/            # Saved input/output snapshots
```

---

## Notes (Alternative: Ubuntu VM)

- You can do the same with a **Ubuntu VM** (VirtualBox/Hyper-V):
    - Download Prometheus binary inside the VM
    - Use the same `ops/prometheus.yml`
    - Expose VM port `9090 -> host 9090`
- WSL2 is lighter for a quick local check, so we recommend it first.

---

## Roadmap

- Self-refine loop (auto-correct invalid PromQL)
- Range queries + sample graph export
- Alert rule skeleton generation

```


### `ops/wsl-setup.md`
- Windows 11에서 WSL2 활성화 → Ubuntu 설치
- WSL 우분투에서 Go 설치(`sudo apt-get update && sudo apt-get install -y golang`)
- Prometheus 설치 (압축 해제 후 바이너리 실행)
- `ops/prometheus.yml` 설명 + `scripts/start_prometheus.sh` 실행법
- 방화벽/포트 9090 접근 확인 
- Node Exporter 추가 방법

### `scripts/start_prometheus.sh`
- 현재 디렉토리 기준으로 `./prometheus` 바이너리 실행
- `--config.file=./ops/prometheus.yml --web.listen-address="0.0.0.0:9090"`
- 로그를 터미널로 출력하도록 (서비스 등록 대신 간단 실행)

### `ops/prometheus.yml` (로컬 전용 샘플)
- `global.scrape_interval: 5s`
- `scrape_configs`:
  - job_name: `prometheus`
    static_configs: `targets: ["localhost:9090"]`

### `.env.example`

```

OPENAI_API_KEY=sk-your-key

PROM_URL=[http://localhost:9090](http://localhost:9090/)


### `scripts/smoke_queries.txt`
- 자연어 입력 예시 여러 줄 (한 줄당 한 테스트)
  - “5분 평균 에러율(5xx) 보여줘”
  - “지난 10분 p95 응답시간 엔드포인트별 상위 5개”
  - “현재 업/다운 타겟 요약”
  - “메모리 사용률 평균”
  - “CPU 사용률 상위 5개 인스턴스”

### `examples/sample_queries.txt`
- (나중에) 실제 실행 스냅샷을 붙여둘 아카이브 파일  
  - NL 입력 → 생성된 PromQL → 검증 결과 → 짧은 결과 요약
















