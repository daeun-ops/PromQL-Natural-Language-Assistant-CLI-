package prom

import (
	"github.com/prometheus/prometheus/promql/parser"
)

func Validate(expr string) error {
	_, err := parser.ParseExpr(expr)
	return err
}

