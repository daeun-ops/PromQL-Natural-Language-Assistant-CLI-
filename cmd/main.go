package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"promql-nlq-assistant/internal/llm"
	"promql-nlq-assistant/internal/prom"
)

func main() {
	fmt.Println(" PromQL NL Assistant")
	fmt.Println("Enter natural language (type 'exit' to quit).")

	lc := llm.New()
	pc := prom.NewClient()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}
		in := strings.TrimSpace(scanner.Text())
		if in == "" {
			continue
		}
		if strings.ToLower(in) == "exit" {
			fmt.Println(" Bye")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// 1) (Optional) Schema hint from Prometheus
		var hint string
		if metrics, err := pc.ListMetrics(ctx, 200); err == nil {
			// Keep it brief to avoid long prompts
			hint = strings.Join(metrics, ", ")
		}

		// 2) LLM → PromQL
		promql, err := lc.GeneratePromQL(ctx, in, hint)
		if err != nil {
			fmt.Println("[LLM ERROR]", err)
			continue
		}
		fmt.Println("\n[Generated PromQL]")
		fmt.Println(promql)

		// 3) Validate syntax
		if err := prom.Validate(promql); err != nil {
			fmt.Println("\n[Validation]  Invalid PromQL:", err)
			continue
		}
		fmt.Println("\n[Validation]  OK")

		// 4) Run instant query (optional)
		if os.Getenv("PROM_URL") != "" {
			fmt.Printf("\n[Query @ %s]\n", os.Getenv("PROM_URL"))
			if b, err := pc.InstantQuery(ctx, promql); err == nil {
				// print first ~800 chars for brevity
				out := string(b)
				if len(out) > 800 {
					out = out[:800] + "..."
				}
				fmt.Println(out)
			} else {
				fmt.Println("Query error:", err)
			}
		} else {
			fmt.Println("\n[Hint] Set PROM_URL to run the query against Prometheus (e.g., http://localhost:9090)")
		}
	}
}
