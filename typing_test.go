package main

import (
	"testing"

	"github.com/mbivert/ftests"
)

func TestTypingapplySubst(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"VarType, empty substitution (id)",
			applySubst,
			[]any{
				&VarType{typ{}, "A"},
				Subst{},
			},
			[]any{
				&VarType{typ{}, "A"},
			},
		},
		{
			"VarType, non-empty, non-matching substitution",
			applySubst,
			[]any{
				&VarType{typ{}, "A"},
				Subst{"B": &VarType{typ{}, "C"}},
			},
			[]any{
				&VarType{typ{}, "A"},
			},
		},
		{
			"VarType, non-empty, matching substitution",
			applySubst,
			[]any{
				&VarType{typ{}, "A"},
				Subst{"A": &VarType{typ{}, "C"}},
			},
			[]any{
				&VarType{typ{}, "C"},
			},
		},
		{
			"Multi-level ArrowType",
			applySubst,
			[]any{
				&ArrowType{
					typ{},
					&VarType{typ{}, "A"},
					&ArrowType{
						typ{},
						&VarType{typ{}, "A"},
						&VarType{typ{}, "B"},
					},
				},
				Subst{"A": &VarType{typ{}, "C"}},
			},
			[]any{
				&ArrowType{
					typ{},
					&VarType{typ{}, "C"},
					&ArrowType{
						typ{},
						&VarType{typ{}, "C"},
						&VarType{typ{}, "B"},
					},
				},
			},
		},
		{
			"Multi-level ProductType/ArrowType",
			applySubst,
			[]any{
				&ProductType{
					typ{},
					&VarType{typ{}, "A"},
					&ArrowType{
						typ{},
						&VarType{typ{}, "A"},
						&VarType{typ{}, "B"},
					},
				},
				Subst{"A": &VarType{typ{}, "C"}},
			},
			[]any{
				&ProductType{
					typ{},
					&VarType{typ{}, "C"},
					&ArrowType{
						typ{},
						&VarType{typ{}, "C"},
						&VarType{typ{}, "B"},
					},
				},
			},
		},
		{
			"Multi-level ProductType/ArrowType, double-substitution",
			applySubst,
			[]any{
				&ProductType{
					typ{},
					&VarType{typ{}, "A"},
					&ArrowType{
						typ{},
						&VarType{typ{}, "A"},
						&ProductType{
							typ{},
							&VarType{typ{}, "B"},
							&VarType{typ{}, "D"},
						},
					},
				},
				Subst{
					"A": &VarType{typ{}, "C"},
					"B": &VarType{typ{}, "E"},
				},
			},
			[]any{
				&ProductType{
					typ{},
					&VarType{typ{}, "C"},
					&ArrowType{
						typ{},
						&VarType{typ{}, "C"},
						&ProductType{
							typ{},
							&VarType{typ{}, "E"},
							&VarType{typ{}, "D"},
						},
					},
				},
			},
		},
		{
			"UnitType: unchanged",
			applySubst,
			[]any{
				&UnitType{typ{}},
				Subst{"A": &VarType{typ{}, "C"}},
			},
			[]any{
				&UnitType{typ{}},
			},
		},
		{
			"BoolType: unchanged",
			applySubst,
			[]any{
				&BoolType{typ{}},
				Subst{"A": &VarType{typ{}, "C"}},
			},
			[]any{
				&BoolType{typ{}},
			},
		},
		{
			"IntType: unchanged",
			applySubst,
			[]any{
				&IntType{typ{}},
				Subst{"A": &VarType{typ{}, "C"}},
			},
			[]any{
				&IntType{typ{}},
			},
		},
		{
			"FloatType: unchanged",
			applySubst,
			[]any{
				&FloatType{typ{}},
				Subst{"A": &VarType{typ{}, "C"}},
			},
			[]any{
				&FloatType{typ{}},
			},
		},
	})
}
