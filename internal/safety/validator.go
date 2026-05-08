package safety

import (
	"fmt"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v6"
)

// ValidateSQL uses the real PostgreSQL parser to check
// that the given SQL is a pure SELECT statement

func ValidateSQL(sql string) error {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return fmt.Errorf("SQL is empty")
	}

	// Parse using the actual Postgres parser
	tree, err := pg_query.Parse(sql)
	if err != nil {
		return fmt.Errorf("SQL parse error: %w", err)
	}

	if len(tree.Stmts) == 0 {
		return fmt.Errorf("no statements found in SQL")
	}

	// Walk every statement in the parse tree
	for i, stmt := range tree.Stmts {
		if stmt.Stmt == nil {
			return fmt.Errorf("statement %d is nil", i)
		}

		// The only allowed statement type is SelectStmt
		if stmt.Stmt.GetSelectStmt() == nil {
			return fmt.Errorf(
				"rejected: only SELECT statements are allowed — "+
					"got a non-SELECT statement at position %d", i+1,
			)
		}
	}

	return nil
}