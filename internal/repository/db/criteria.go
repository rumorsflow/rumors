package db

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"sort"
	"strings"
)

const (
	QueryIndex = "index"
	QuerySize  = "size"
	QuerySort  = "sort"
	QueryCond  = "cond"
	QueryField = "field"
	QueryValue = "value"

	CondEmpty = ""
	CondEq    = "eq"
	CondNe    = "ne"
	CondGt    = "gt"
	CondGte   = "gte"
	CondLt    = "lt"
	CondLte   = "lte"
	CondIn    = "in"
	CondNin   = "nin"
	CondRegex = "regex"
	CondLike  = "like"
)

var conditions = map[string]string{
	CondEmpty: "$eq",
	CondEq:    "$eq",
	CondNe:    "$ne",
	CondGt:    "$gt",
	CondGte:   "$gte",
	CondLt:    "$lt",
	CondLte:   "$lte",
	CondIn:    "$in",
	CondNin:   "$nin",
	CondRegex: "$regex",
	CondLike:  "$regex",
}

type field struct {
	field string
	cond  string
	value string
}

func (f field) valid() bool {
	return f.field != ""
}

func (f field) m() bson.M {
	return bson.M{f.field: bson.M{f.cond: f.parse()}}
}

func (f field) parse() any {
	if f.value == "" || f.cond == "$regex" {
		return f.value
	}

	if strings.EqualFold(f.value, "null") || strings.EqualFold(f.value, "nil") {
		return nil
	}

	switch f.cond {
	case "$in", "$nin":
		data := strings.Split(f.value, ",")
		value := make([]any, 0, len(data))
		for _, item := range data {
			value = append(value, parse(item))
		}
		return value
	}

	return parse(f.value)
}

func parse(value string) any {
	if v, err := cast.ToInt64E(value); err == nil {
		return v
	}
	if v, err := uuid.Parse(value); err == nil {
		return v
	}
	if v, err := cast.ToTimeE(value); err == nil {
		return v
	}
	if v, err := cast.ToBoolE(value); err == nil {
		return v
	}
	return value
}

func BuildCriteria(query string, excludeFields ...string) *repository.Criteria {
	exclude := make(map[string]struct{}, len(excludeFields))
	for _, f := range excludeFields {
		exclude[f] = struct{}{}
	}

	criteria := &repository.Criteria{}

	fields := map[int]map[int]*field{}
	var order bson.D

	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.Contains(key, ";") {
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		key, err := url.QueryUnescape(key)
		if err != nil {
			continue
		}
		value, err = url.QueryUnescape(value)
		if err != nil {
			continue
		}

		switch key {
		case QueryIndex:
			criteria.SetIndex(cast.ToInt64(value))
		case QuerySize:
			criteria.SetSize(cast.ToInt64(value))
		case QuerySort:
			if len(value) > 0 {
				if '-' == value[0] {
					if len(value) > 1 {
						order = append(order, primitive.E{Key: value[1:], Value: -1})
					}
				} else {
					order = append(order, primitive.E{Key: value, Value: 1})
				}
			}
		default:
			parts := strings.SplitN(key, ".", 3)
			switch parts[0] {
			case QueryField, QueryCond, QueryValue:
			default:
				continue
			}

			i := cast.ToInt(parts[1])
			j := cast.ToInt(parts[2])

			if _, ok := fields[i]; !ok {
				fields[i] = map[int]*field{}
			}
			if _, ok := fields[i][j]; !ok {
				fields[i][j] = &field{cond: conditions[CondEq]}
			}

			switch parts[0] {
			case QueryField:
				if _, ok := exclude[strings.TrimSpace(value)]; ok {
					continue
				}
				fields[i][j].field = strings.TrimSpace(value)
			case QueryCond:
				if c, ok := conditions[value]; ok {
					fields[i][j].cond = c
				}
			default:
				fields[i][j].value = value
			}
		}
	}

	idx := make([]int, 0, len(fields))
	for k1, m := range fields {
		for k2, f := range m {
			if f.valid() {
				continue
			}
			delete(m, k2)
		}
		if len(m) == 0 {
			delete(fields, k1)
		} else {
			idx = append(idx, k1)
		}
	}
	sort.Ints(idx)

	filter := bson.M{}
	for _, i := range idx {
		if len(fields[i]) > 1 {
			a := make(bson.A, 0, len(fields[i]))
			for _, f := range fields[i] {
				a = append(a, f.m())
			}
			filter["$or"] = a
		} else {
			for k, v := range fields[i][0].m() {
				filter[k] = v
			}
		}
	}
	criteria.Filter = filter

	if len(order) > 0 {
		criteria.Sort = order
	}

	return criteria
}
