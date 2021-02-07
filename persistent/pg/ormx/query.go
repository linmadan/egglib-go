package ormx

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

func NewQuery(query *orm.Query, queryOptions map[string]interface{}) *Query {
	return &Query{
		query,
		queryOptions,
		0,
	}
}

func (query *Query) SetWhere(condition, key string) *Query {
	if v, ok := query.queryOptions[key]; ok {
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

func (query *Query) SetUpdate(condition, key string) *Query {
	if v, ok := query.queryOptions[key]; ok {
		query.Set(condition, v)
	}
	return query
}

func (query *Query) SetLimit() *Query {
	if offset, ok := query.queryOptions["offset"]; ok {
		offset, _ := strconv.ParseInt(fmt.Sprintf("%v", offset), 10, 64)
		if offset > -1 {
			query.Offset(int(offset))
		}
	} else {
		query.Offset(0)
	}
	if limit, ok := query.queryOptions["limit"]; ok {
		limit, _ := strconv.ParseInt(fmt.Sprintf("%v", limit), 10, 64)
		if limit > 0 {
			query.Limit(int(limit))
		} else {
			query.Limit(20)
		}
	}
	return query
}

func (query *Query) SetOrder(orderColumn string, key string) *Query {
	if v, ok := query.queryOptions[key]; ok {
		query.Order(fmt.Sprintf("%v %v", orderColumn, v))
	}
	return query
}
