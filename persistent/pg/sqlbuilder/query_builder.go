package sqlbuilder

import (
	"fmt"
	"github.com/go-pg/pg/v10/orm"
	"strconv"
	"time"
)

type Query struct {
	*orm.Query
	queryOptions map[string]interface{}
	AffectRow    int
}

func BuildQuery(query *orm.Query, queryOptions map[string]interface{}) *Query {
	return &Query{
		query,
		queryOptions,
		0,
	}
}

func (query *Query) SetWhereByQueryOption(condition string, optionKey string) *Query {
	if v, ok := query.queryOptions[optionKey]; ok {
		// time.Time 零值特殊处理
		if _, ok := v.(time.Time); ok {
			if t, e := time.Parse(time.RFC3339, fmt.Sprintf("%v", v)); e == nil {
				if t.IsZero() {
					return query
				}
			}
		}
		query.Where(condition, v)
	}
	return query
}

func (query *Query) SetUpdateByQueryOption(condition string, optionKey string) *Query {
	if v, ok := query.queryOptions[optionKey]; ok {
		query.Set(condition, v)
	}
	return query
}

func (query *Query) SetOffsetAndLimit(defaultLimit int) *Query {
	if offset, ok := query.queryOptions["offset"]; ok {
		offset, _ := strconv.ParseInt(fmt.Sprintf("%v", offset), 10, 64)
		if offset >= 0 {
			query.Offset(int(offset))
		}
	} else {
		query.Offset(0)
	}
	if limit, ok := query.queryOptions["limit"]; ok {
		limit, _ := strconv.ParseInt(fmt.Sprintf("%v", limit), 10, 64)
		if limit == 0 {
			query.Limit(defaultLimit)
		} else {
			if limit > 0 {
				query.Limit(int(limit))
			}
		}
	} else {
		query.Limit(defaultLimit)
	}
	return query
}

func (query *Query) SetOrderByQueryOption(orderColumn string, optionKey string) *Query {
	if v, ok := query.queryOptions[optionKey]; ok && (v == "ASC" || v == "DESC") {
		query.Order(fmt.Sprintf("%s %s", orderColumn, v))
	}
	return query
}

func (query *Query) SetOrderDirect(orderColumn string, sortOrder string) *Query {
	query.Order(fmt.Sprintf("%s %s", orderColumn, sortOrder))
	return query
}
