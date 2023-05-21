package tck_test

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/mburbidg/cypher"
	"github.com/mburbidg/cypher/parser"
	"github.com/mburbidg/cypher/scanner"
	"os"
	"testing"
)

type graphFeature struct{}

type syntaxErrKey struct{}

type reporter struct{}

func (r reporter) Error(line int, msg string) error {
	return fmt.Errorf("Error: %s (line %d)", msg, line)
}

func (g *graphFeature) anyGraph(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) beforeScenario(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) afterScenario(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) anEmptyGraph(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) executingQuery(ctx context.Context, query *godog.DocString) (context.Context, error) {
	reporter := &reporter{}
	s := scanner.New([]byte(query.Content), reporter)
	p := parser.New(s, reporter)
	stmt, err := p.Parse()
	if err != nil {
		return context.WithValue(ctx, syntaxErrKey{}, err), nil
	}
	r := &astRuntime{}
	err = r.eval(stmt)
	if err != nil {
		return context.WithValue(ctx, syntaxErrKey{}, err), nil
	}
	return ctx, nil
}

func (g *graphFeature) executingControlQuery(ctx context.Context, query *godog.DocString) (context.Context, error) {
	reporter := &reporter{}
	s := scanner.New([]byte(query.Content), reporter)
	p := parser.New(s, reporter)
	r := &astRuntime{}
	stmt, err := p.Parse()
	if err != nil {
		return context.WithValue(ctx, syntaxErrKey{}, err), nil
	}
	err = r.eval(stmt)
	if err != nil {
		return context.WithValue(ctx, syntaxErrKey{}, err), nil
	}
	return ctx, nil
}

func (g *graphFeature) havingExecutedQuery(ctx context.Context, query *godog.DocString) (context.Context, error) {
	reporter := &reporter{}
	s := scanner.New([]byte(query.Content), reporter)
	p := parser.New(s, reporter)
	r := &astRuntime{}
	stmt, err := p.Parse()
	if err != nil {
		return context.WithValue(ctx, syntaxErrKey{}, err), nil
	}
	err = r.eval(stmt)
	if err != nil {
		return context.WithValue(ctx, syntaxErrKey{}, err), nil
	}
	return ctx, nil
}

func (g *graphFeature) theResultShouldBeEmpty(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) theResultShouldBeInAnyOrder(ctx context.Context, table *godog.Table) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) theSideEffectsShouldBe(ctx context.Context, values *godog.Table) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) noSideEffects(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (g *graphFeature) syntaxErrorRaised(ctx context.Context, errStr string) (context.Context, error) {
	if err, ok := ctx.Value(syntaxErrKey{}).(error); ok {
		if cypherErr, ok := err.(*cypher.CypherErr); ok {
			if cypherErr.Code != errStr {
				return ctx, fmt.Errorf("expecting syntax error: expected=%s, actual=%s", errStr, cypherErr.Code)
			}
		}
		return ctx, nil
	}
	return ctx, fmt.Errorf("expecting syntax error: %s", errStr)
}

func TestCypherFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeCypherScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths: []string{
				"tck/features/clauses/create/Create2.feature",
				//"tck/features/clauses/match/Match1.feature",
			},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeCypherScenario(sc *godog.ScenarioContext) {
	g := &graphFeature{}
	sc.Before(g.beforeScenario)
	sc.After(g.afterScenario)
	sc.Step(`^any graph$`, g.anyGraph)
	sc.Step(`^an empty graph$`, g.anEmptyGraph)
	sc.Step(`^executing query:$`, g.executingQuery)
	sc.Step(`^the result should be empty$`, g.theResultShouldBeEmpty)
	sc.Step(`^the result should be, in any order:$`, g.theResultShouldBeInAnyOrder)
	sc.Step(`^the side effects should be:$`, g.theSideEffectsShouldBe)
	sc.Step(`^no side effects$`, g.noSideEffects)
	sc.Step(`^a SyntaxError should be raised at compile time: ([a-zA-Z]+)$`, g.syntaxErrorRaised)
	sc.Step(`^executing control query:$`, g.executingControlQuery)
	sc.Step(`^having executed:$`, g.havingExecutedQuery)

}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
