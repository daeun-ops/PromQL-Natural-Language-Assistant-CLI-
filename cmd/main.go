package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(" PromQL Natural Language Assistant (local test mode)")
	fmt.Println("Type your question below (or 'exit' to quit):")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if strings.ToLower(input) == "exit" {
			fmt.Println(" Bye!")
			break
		}

		if input == "" {
			continue
		}

		// 간단한 가짜 매핑 (LLM 대신)
		promql := fakePromQL(input)

		fmt.Println("\n[Generated PromQL]")
		fmt.Println(promql)
		fmt.Println("\n[Validation]")
		fmt.Println(" Syntax looks fine (mock check)")
	}
}

func fakePromQL(input string) string {
	switch {
	case strings.Contains(input, "에러") || strings.Contains(input, "error"):
		return `sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))`
	case strings.Contains(input, "지연") || strings.Contains(input, "latency"):
		return `histogram_quantile(0.95, sum by (le, endpoint)(rate(http_request_duration_seconds_bucket[5m])))`
	default:
		return `sum(rate(http_requests_total[5m]))`
	}
}

