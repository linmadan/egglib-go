package hooks

import (
	"fmt"
	"github.com/go-pg/pg"
)

type SqlGeneratePrintHook struct{}

func (hook SqlGeneratePrintHook) BeforeQuery(q *pg.QueryEvent) {
}

func (hook SqlGeneratePrintHook) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}
