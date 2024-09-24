package main

import (
	"fmt"
	"testing"

	"github.com/mbivert/ftests"
)

func TestTypingApplySubst(t *testing.T) {
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

func TestTypingOccursIn(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"VarType, no match",
			occursIn,
			[]any{
				&VarType{typ{}, "A"},
				"B",
			},
			[]any{false},
		},
		{
			"VarType, match",
			occursIn,
			[]any{
				&VarType{typ{}, "A"},
				"A",
			},
			[]any{true},
		},
		{
			"ArrowType, match",
			occursIn,
			[]any{
				&ArrowType{typ{},
					&VarType{typ{}, "A"},
					&VarType{typ{}, "B"},
				},
				"A",
			},
			[]any{true},
		},
		{
			"ArrowType, no match",
			occursIn,
			[]any{
				&ArrowType{typ{},
					&VarType{typ{}, "A"},
					&VarType{typ{}, "B"},
				},
				"C",
			},
			[]any{false},
		},
		{
			"ProductType, match",
			occursIn,
			[]any{
				&ProductType{typ{},
					&VarType{typ{}, "A"},
					&VarType{typ{}, "B"},
				},
				"A",
			},
			[]any{true},
		},
		{
			"ProductType, no match",
			occursIn,
			[]any{
				&ProductType{typ{},
					&VarType{typ{}, "A"},
					&VarType{typ{}, "B"},
				},
				"C",
			},
			[]any{false},
		},
		{
			"UnitType: never matched",
			occursIn,
			[]any{
				&UnitType{typ{}},
				"A",
			},
			[]any{false},
		},
		{
			"BoolType: never matched",
			occursIn,
			[]any{
				&BoolType{typ{}},
				"A",
			},
			[]any{false},
		},
		{
			"IntType: never matched",
			occursIn,
			[]any{
				&IntType{typ{}},
				"A",
			},
			[]any{false},
		},
		{
			"FloatType: never matched",
			occursIn,
			[]any{
				&FloatType{typ{}},
				"A",
			},
			[]any{false},
		},
	})
}

func TestTypingMgu(t *testing.T) {
	var nilSubst Subst

	ftests.Run(t, []ftests.Test{
		{
			"Empty input",
			mgu,
			[]any{[]Type{}, []Type{}},
			[]any{Subst{}, nil},
		},
		{
			"case 1: mgu(X; X) = id",
			mgu,
			[]any{[]Type{&VarType{typ{}, "X"}}, []Type{&VarType{typ{}, "X"}}},
			[]any{Subst{}, nil},
		},
		{
			"case 2: mgu(X; B) = [X ↦ B] if X ∉ B (B is →)",
			mgu,
			[]any{
				[]Type{&VarType{typ{}, "X"}},
				[]Type{
					&ArrowType{typ{},
						&VarType{typ{}, "Y"},
						&VarType{typ{}, "Z"},
					},
				},
			},
			[]any{Subst{
				"X": &ArrowType{typ{},
					&VarType{typ{}, "Y"},
					&VarType{typ{}, "Z"},
				},
			}, nil},
		},
		{
			"case 2: mgu(X; B) = [X ↦ B] if X ∉ B (B is ×, contains →)",
			mgu,
			[]any{
				[]Type{&VarType{typ{}, "X"}},
				[]Type{
					&ProductType{typ{},
						&VarType{typ{}, "Y"},
						&ArrowType{typ{},
							&VarType{typ{}, "Z"},
							&VarType{typ{}, "Z"},
						},
					},
				},
			},
			[]any{Subst{
				"X": &ProductType{typ{},
					&VarType{typ{}, "Y"},
					&ArrowType{typ{},
						&VarType{typ{}, "Z"},
						&VarType{typ{}, "Z"},
					},
				},
			}, nil},
		},
		{
			"case 2: mgu(X; B) = [X ↦ B] if X ∉ B (B is ι)",
			mgu,
			[]any{
				[]Type{&VarType{typ{}, "X"}},
				[]Type{&BoolType{typ{}}},
			},
			[]any{Subst{
				"X": &BoolType{typ{}},
			}, nil},
		},
		{
			"case 3: mgu(X; B) fails if X ∈ B (B is →)",
			mgu,
			[]any{
				[]Type{&VarType{typ{}, "X"}},
				[]Type{
					&ArrowType{typ{},
						&VarType{typ{}, "Y"},
						&VarType{typ{}, "X"},
					},
				},
			},
			[]any{nilSubst, fmt.Errorf("X occurs in Y → X")},
		},
		{
			"case 3: mgu(X; B) fails if X ∈ B (B is ×, contains →)",
			mgu,
			[]any{
				[]Type{&VarType{typ{}, "X"}},
				[]Type{
					&ProductType{typ{},
						&VarType{typ{}, "Y"},
						&ArrowType{typ{},
							&VarType{typ{}, "Z"},
							&VarType{typ{}, "X"},
						},
					},
				},
			},
			[]any{nilSubst, fmt.Errorf("X occurs in Y × (Z → X)")},
		},
	})
}
