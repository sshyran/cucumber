package cucumberexpressions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCucumberExpressionParser(t *testing.T) {
	var assertAst = func(t *testing.T, expression string, expected node) {
		ast, err := parse(expression)
		require.NoError(t, err)
		require.Equal(t, expected, ast)
	}

	t.Run("empty string", func(t *testing.T) {
		assertAst(t, "", node{
			expressionNode,
			-1, -1,
			"",
			[]node{},
		})
	})

	t.Run("phrase", func(t *testing.T) {
		assertAst(t, "three blind mice", node{
			expressionNode,
			-1, -1,
			"",
			[]node{
				{textNode, -1, -1, "three", []node{}},
				{textNode, -1, -1, " ", []node{}},
				{textNode, -1, -1, "blind", []node{}},
				{textNode, -1, -1, " ", []node{}},
				{textNode, -1, -1, "mice", []node{}},
			},
		})
	})

	t.Run("optional", func(t *testing.T) {
		assertAst(t, "(blind)", node{
			expressionNode,
			-1, -1,
			"",
			[]node{
				{optionalNode,
					-1,
					-1,
					"",
					[]node{
						{textNode, -1, -1, "blind", []node{}},
					},
				},
			},
		})
	})

	t.Run("parameter", func(t *testing.T) {
		assertAst(t, "{string}", node{
			expressionNode, -1, -1,
			"",
			[]node{
				{parameterNode, -1, -1,
					"",
					[]node{
						{textNode, -1, -1, "string", []node{}},
					},
				},
			},
		})
	})

	t.Run("anonymous parameter", func(t *testing.T) {
		assertAst(t, "{}", node{
			expressionNode, -1, -1,
			[]node{
				{parameterNode, -1, -1,
					[]node{},
					"",
				},
			},
			"",
		})
	})

	t.Run("optional phrase", func(t *testing.T) {
		assertAst(t, "three (blind) mice", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{optionalNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "blind"},
				}, ""},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
			},
			"",
		})
	})

	t.Run("slash", func(t *testing.T) {
		assertAst(t, "\\", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "\\"},
			},
			"",
		})
	})

	t.Run("opening brace", func(t *testing.T) {
		assertAst(t, "{", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "{"},
			},
			"",
		})
	})

	t.Run("unfinished parameter", func(t *testing.T) {
		assertAst(t, "{string", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "{"},
				{textNode, -1, -1 []node{}, "string"},
			},
			"",
		})
	})

	t.Run("opening parenthesis", func(t *testing.T) {
		assertAst(t, "(", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "("},
			},
			"",
		})
	})

	t.Run("escaped opening parenthesis", func(t *testing.T) {
		assertAst(t, "\\(", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "("},
			},
			"",
		})
	})

	t.Run("escaped optional phrase", func(t *testing.T) {
		assertAst(t, "three \\(blind) mice", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "("},
				{textNode, -1, -1 []node{}, "blind"},
				{textNode, -1, -1 []node{}, ")"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
			},
			"",
		})
	})

	t.Run("escaped optional followed by optional", func(t *testing.T) {
		assertAst(t, "three \\((very) blind) mice", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "("},
				{optionalNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "very"},
				}, ""},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "blind"},
				{textNode, -1, -1 []node{}, ")"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
			},
			"",
		})
	})

	t.Run("optional containing escaped optional", func(t *testing.T) {
		assertAst(t, "three ((very\\) blind) mice", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{optionalNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "("},
				{textNode, -1, -1 []node{}, "very"},
				{textNode, -1, -1 []node{}, ")"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "blind"},
				}, ""},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
			},
			"",
		})
	})

	t.Run("alternation", func(t *testing.T) {
		assertAst(t, "mice/rats", node{
			expressionNode, -1, -1,
			[]node{
				{alternationNode, -1, -1 []node{
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "mice"},
				}, ""},
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "rats"},
				}, ""},
				}, ""},
			},
			"",
		})
	})

	t.Run("escaped alternation", func(t *testing.T) {
		assertAst(t, "mice\\/rats", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "mice"},
				{textNode, -1, -1 []node{}, "/"},
				{textNode, -1, -1 []node{}, "rats"},
			},
			"",
		})
	})

	t.Run("alternation phrase", func(t *testing.T) {
		assertAst(t, "three hungry/blind mice", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{alternationNode, -1, -1 []node{
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "hungry"},
				}, ""},
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "blind"},
				}, ""},
				}, ""},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
			},
			"",
		})
	})

	t.Run("alternation with whitespace", func(t *testing.T) {
		assertAst(t, "\\ three\\ hungry/blind\\ mice\\ ", node{
			expressionNode, -1, -1,
			[]node{
				{alternationNode, -1, -1 []node{
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "hungry"},
				}, ""},
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "blind"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
				{textNode, -1, -1 []node{}, " "},
				}, ""},
				}, ""},
			},
			"",
		})
	})

	t.Run("alternation with unused end optional", func(t *testing.T) {
		assertAst(t, "three )blind\\ mice/rats", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{alternationNode, -1, -1 []node{
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, ")"},
				{textNode, -1, -1 []node{}, "blind"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
				}, ""},
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "rats"},
				}, ""},
				}, ""},
			},
			"",
		})
	})

	t.Run("alternation with unused start optional", func(t *testing.T) {
		assertAst(t, "three blind\\ mice/rats(", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{alternationNode, -1, -1 []node{
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "blind"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "mice"},
				}, ""},
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "rats"},
				{textNode, -1, -1 []node{}, "("},
				}, ""},
				}, ""},
			},
			"",
		})
	})

	t.Run("alternation with unused start optional", func(t *testing.T) {
		assertAst(t, "three blind\\ rat/cat(s)", node{
			expressionNode, -1, -1,
			[]node{
				{textNode, -1, -1 []node{}, "three"},
				{textNode, -1, -1 []node{}, " "},
				{alternationNode, -1, -1 []node{
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "blind"},
				{textNode, -1, -1 []node{}, " "},
				{textNode, -1, -1 []node{}, "rat"},
				}, ""},
				{alternativeNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "cat"},
				{optionalNode, -1, -1 []node{
				{textNode, -1, -1 []node{}, "s"},
				}, ""},
				}, ""},
				}, ""},
			},
			"",
		})
	})

}