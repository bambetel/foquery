package main

import (
	"fmt"
	"log"
	"strings"
)

type Query struct {
	Keys   []string
	Op     []string
	Values []any
}

// TODO: insert value placeholders
// $1, $2, $3..$N
func (q *Query) Select() string {
	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM tablename")

	if len(q.Op) > 0 {
		sb.WriteString(" WHERE ")
	}

ops:
	for i := range q.Op {
		switch q.Op[i] {
		case "=":
			sb.WriteString(fmt.Sprintf("%s=%v", q.Keys[i], q.Values[i]))
		case "r":
			r := q.Values[i].([]int)
			sb.WriteString(fmt.Sprintf("%s BETWEEN %d AND %d", q.Keys[i], r[0], r[1]))
		default:
			log.Printf("Invalid op: %q\n", q.Op[i])
		}
		if i == len(q.Op)-1 {
			break ops
		}
		sb.WriteString(" AND ")
	}

	return sb.String()
}

type QueryOption func(*Query)

func WithStringEq(col string, val string) QueryOption {
	return func(q *Query) {
		q.Keys = append(q.Keys, col)
		q.Op = append(q.Op, "=")
		q.Values = append(q.Values, fmt.Sprintf("'%s'", val))
	}
}

func WithIntEq(col string, val int) QueryOption {
	return func(q *Query) {
		q.Keys = append(q.Keys, col)
		q.Op = append(q.Op, "=")
		q.Values = append(q.Values, val)
	}
}

func WithIntRange(col string, min, max int) QueryOption {
	return func(q *Query) {
		q.Keys = append(q.Keys, col)
		q.Op = append(q.Op, "r")
		q.Values = append(q.Values, []int{min, max})
	}
}

func NewQuery(options ...QueryOption) *Query {
	q := &Query{}
	for _, o := range options {
		o(q)
	}
	return q
}

func main() {
	q := NewQuery(WithIntEq("priority", 10))
	fmt.Printf("%+v\n", q)
	fmt.Printf("SQL Select: %s\n", q.Select())
	fmt.Println("---")
	q1 := NewQuery()
	fmt.Printf("%+v\n", q1)
	fmt.Printf("SQL Select: %s\n", q1.Select())
	fmt.Println("---")
	q2 := NewQuery(WithIntRange("val", 2, 5))
	fmt.Printf("%+v\n", q2)
	fmt.Printf("SQL Select: %s\n", q2.Select())
	fmt.Println("---")
	q3 := NewQuery(WithIntRange("val", 2, 5), WithStringEq("name", "test"))
	fmt.Printf("%+v\n", q3)
	fmt.Printf("SQL Select: %s\n", q3.Select())
}
