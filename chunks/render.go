package chunks

import (
	"fmt"
	"io"

	"github.com/osteele/liquid/generics"
)

// Render evaluates an AST node and writes the result to an io.Writer.
func (n *ASTSeq) Render(w io.Writer, ctx Context) error {
	for _, c := range n.Children {
		if err := c.Render(w, ctx); err != nil {
			return err
		}
	}
	return nil
}

// Render evaluates an AST node and writes the result to an io.Writer.
// TODO probably safe to remove this type and method, once the test suite is larger
func (n *ASTChunks) Render(w io.Writer, _ Context) error {
	fmt.Println(MustYAML(n))
	return fmt.Errorf("unimplemented: ASTChunks.Render")
}

// Render evaluates an AST node and writes the result to an io.Writer.
func (n *ASTGenericTag) Render(w io.Writer, ctx Context) error {
	return n.render(w, ctx)
}

// Render evaluates an AST node and writes the result to an io.Writer.
func (n *ASTText) Render(w io.Writer, _ Context) error {
	_, err := w.Write([]byte(n.Source))
	return err
}

// Render evaluates an AST node and writes the result to an io.Writer.
func (n *ASTRaw) Render(w io.Writer, _ Context) error {
	for _, s := range n.slices {
		_, err := w.Write([]byte(s))
		if err != nil {
			return err
		}
	}
	return nil
}

// RenderASTSequence renders a sequence of nodes.
func (ctx Context) RenderASTSequence(w io.Writer, seq []ASTNode) error {
	for _, n := range seq {
		if err := n.Render(w, ctx); err != nil {
			return err
		}
	}
	return nil
}

// Render evaluates an AST node and writes the result to an io.Writer.
func (n *ASTControlTag) Render(w io.Writer, ctx Context) error {
	cd, ok := findControlTagDefinition(n.Tag)
	if !ok || cd.action == nil {
		return fmt.Errorf("unimplemented tag: %s", n.Tag)
	}
	f := cd.action(*n)
	return f(w, ctx)
}

// Render evaluates an AST node and writes the result to an io.Writer.
func (n *ASTObject) Render(w io.Writer, ctx Context) error {
	// TODO separate this into parse and evaluate stages.
	value, err := ctx.EvaluateExpr(n.Args)
	if err != nil {
		return err
	}
	if generics.IsEmpty(value) {
		return nil
	}
	_, err = w.Write([]byte(fmt.Sprint(value)))
	return err
}
