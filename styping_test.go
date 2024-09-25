package main

import (
	"fmt"
	"testing"

	"github.com/mbivert/ftests"
)

func TestSTypingSInferType(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"Unbounded variable",
			sInferType,
			[]any{&VarExpr{expr{}, "A"},},
			[]any{
				nil,
				fmt.Errorf("'A' isn't bounded!"),
			},
		},
		{
			"λx:bool . *",
			sInferType,
			[]any{mustParse("λx:bool.*")},
			[]any{
				&AbsExpr{expr{&ArrowType{typ{},
						&BoolType{typ{}},
						&UnitType{typ{}},
					}},
					&BoolType{typ{}},
					"x",
					&UnitExpr{expr{&UnitType{typ{}}}},
				},
				nil,
			},
		},
	})
}
