package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"promql-nlq-assistant/internal/llm"  // OpenAI LLM 클라이언트
	"promql-nlq-assistant/internal/prom" // Prometheus API 및 검증 유틸
)

// main 함수: CLI 프로그램의 진입점
// 사용자가 터미널에서 자연어 입력 → PromQL 쿼리 생성 → 문법 검증 → (옵션) Prometheus 쿼리 실행
func main() {
	fmt.Println(" PromQL NL Assistant")
	fmt.Println("Enter natural language (type 'exit' to quit).")

	// LLM 및 Prometheus 클라이언트 초기화
	lc := llm.New()      // OpenAI 또는 mock LLM 클라이언트
	pc := prom.NewClient() // Prometheus API 클라이언트

	// 표준 입력(터미널) 스캐너 생성
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ") // 사용자 입력 프롬프트
		if !scanner.Scan() {
			break // EOF나 입력 오류 발생 시 종료
		}
		in := strings.TrimSpace(scanner.Text())
		if in == "" {
			continue // 빈 줄은 무시
		}
		if strings.ToLower(in) == "exit" {
			fmt.Println(" Bye sophie log~~")
			return // 종료 명령
		}

		// 요청마다 컨텍스트 생성 (타임아웃 60초)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		//  Prometheus에서 스키마 힌트 가져오기
		//    - LLM이 실제 메트릭 이름을 참고할 수 있도록
		var hint string
		if metrics, err := pc.ListMetrics(ctx, 200); err == nil {
			// 너무 길면 프롬프트에 불필요한 토큰 낭비라 간단히 join
			hint = strings.Join(metrics, ", ")
		}

		//  LLM을 통해 PromQL 쿼리 생성
		promql, err := lc.GeneratePromQL(ctx, in, hint)
		if err != nil {
			fmt.Println("[LLM ERROR]", err)
			continue // LLM 오류 시 다음 입력으로
		}

		fmt.Println("\n[Generated PromQL]")
		fmt.Println(promql)

		//  PromQL 문법 검증
		if err := prom.Validate(promql); err != nil {
			fmt.Println("\n[Validation]  Invalid PromQL:", err)
			continue // 문법 오류 시 실행하지 않음
		}

		fmt.Println("\n[Validation]  OK")

		//  Prometheus에 즉시 쿼리 수행 (옵션)
		//   PROM_URL 환경변수가 설정된 경우에만 실행
		if os.Getenv("PROM_URL") != "" {
			fmt.Printf("\n[Query @ %s]\n", os.Getenv("PROM_URL"))

			// Prometheus instant query 실행
			if b, err := pc.InstantQuery(ctx, promql); err == nil {
				out := string(b)

				// 출력이 너무 길면 800자까지만 표시
				if len(out) > 800 {
					out = out[:800] + "..."
				}
				fmt.Println(out)
			} else {
				fmt.Println("Query error:", err)
			}
		} else {
			// PROM_URL이 없으면 힌트 메시지 출력
			fmt.Println("\n[Hint] Set PROM_URL to run the query against Prometheus (e.g., http://localhost:9090)")
		}
	}
}
