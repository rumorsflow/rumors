package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

const (
	// comparison
	eq  = operator("$eq")
	ne  = operator("$ne")
	gt  = operator("$gt")
	gte = operator("$gte")
	lt  = operator("$lt")
	lte = operator("$lte")
	in  = operator("$in")
	nin = operator("$nin")

	// logical
	and = operator("$and")
	or  = operator("$or")
	nor = operator("$nor")
	not = operator("$not")

	// array
	all       = operator("$all")
	size      = operator("$size")
	elemMatch = operator("$elemMatch")

	// element
	exists = operator("$exists")

	// evaluation
	regex = operator("$regex")
)

type (
	Filter []Expr

	Expr interface {
		Build() any
	}

	operator string

	logical struct {
		op   operator
		data []Expr
	}

	field struct {
		op    operator
		name  string
		value any
	}
)

func And(expr ...Expr) Expr {
	return &logical{op: and, data: expr}
}

func Or(expr ...Expr) Expr {
	return &logical{op: or, data: expr}
}

func Nor(expr ...Expr) Expr {
	return &logical{op: nor, data: expr}
}

func Not(name string, expr Expr) Expr {
	return &field{name: name, op: not, value: expr}
}

func Eq(name string, value any) Expr {
	return &field{name: name, op: eq, value: value}
}

func Ne(name string, value any) Expr {
	return &field{name: name, op: ne, value: value}
}

func Gt(name string, value any) Expr {
	return &field{name: name, op: gt, value: value}
}

func Gte(name string, value any) Expr {
	return &field{name: name, op: gte, value: value}
}

func Lt(name string, value any) Expr {
	return &field{name: name, op: lt, value: value}
}

func Lte(name string, value any) Expr {
	return &field{name: name, op: lte, value: value}
}

func In(name string, value ...any) Expr {
	return &field{name: name, op: in, value: value}
}

func Nin(name string, value ...any) Expr {
	return &field{name: name, op: nin, value: value}
}

func Size(name string, value uint) Expr {
	return &field{name: name, op: size, value: value}
}

func Exists(name string, value bool) Expr {
	return &field{name: name, op: exists, value: value}
}

func Regex(name string, pattern string, opts ...string) Expr {
	return &field{name: name, op: regex, value: primitive.Regex{
		Pattern: pattern,
		Options: strings.Join(opts, ""),
	}}
}

func All(name string, value ...any) Expr {
	return &field{name: name, op: all, value: value}
}

func ElemMatch(name string, query ...Expr) Expr {
	value := bson.M{}
	for _, e := range query {
		if m, ok := e.Build().(bson.M); ok {
			if f, ok := e.(*field); ok && len(f.name) > 0 {
				value[f.name] = m
			} else {
				for k, v := range m {
					value[k] = v
				}
			}
		}

	}
	return &field{name: name, op: elemMatch, value: value}
}

func (l *logical) Build() any {
	result := make(bson.A, len(l.data))
	for i, item := range l.data {
		switch e := item.(type) {
		case *logical:
			result[i] = bson.M{string(e.op): e.Build()}
			continue
		case *field:
			if len(e.name) > 0 {
				result[i] = bson.M{e.name: e.Build()}
				continue
			}
		}
		result[i] = item.Build()
	}
	return result
}

func (f *field) Build() any {
	return bson.M{string(f.op): f.value}
}

func (f Filter) Build() any {
	result := bson.M{}
	for _, i := range f {
		switch e := i.(type) {
		case *logical:
			result[string(e.op)] = e.Build()
		case *field:
			if len(e.name) > 0 {
				result[e.name] = e.Build()
			} else {
				if m, ok := e.Build().(bson.M); ok {
					for k, v := range m {
						result[k] = v
					}
				}
			}
		}
	}
	return result
}
