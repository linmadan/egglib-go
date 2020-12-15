package hooks

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
)

type SqlGeneratePrintHook struct{}

func (hook SqlGeneratePrintHook) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (hook SqlGeneratePrintHook) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	data, _ := q.FormattedQuery()
	fmt.Println(string(data))
	return nil
}
