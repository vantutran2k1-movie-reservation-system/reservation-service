package filters

import (
	"fmt"
	"gorm.io/gorm"
)

type LogicOperator string
type ComparisonOperator string
type SortDirection string

const (
	// Logic operators
	And LogicOperator = "AND"
	Or  LogicOperator = "OR"

	// Comparison operators
	OpEqual        ComparisonOperator = "="
	OpNotEqual     ComparisonOperator = "!="
	OpGreater      ComparisonOperator = ">"
	OpGreaterEqual ComparisonOperator = ">="
	OpLess         ComparisonOperator = "<"
	OpLessEqual    ComparisonOperator = "<="
	OpLike         ComparisonOperator = "LIKE"
	OpIn           ComparisonOperator = "IN"
	OpNotIn        ComparisonOperator = "NOT IN"

	// Sort directions
	Asc  SortDirection = "ASC"
	Desc SortDirection = "DESC"
)

type Condition struct {
	Operator ComparisonOperator
	Value    any
}

func (c *Condition) ToFilterCondition(field string) FilterCondition {
	return FilterCondition{
		Field: field,
		Condition: Condition{
			Operator: c.Operator,
			Value:    c.Value,
		},
	}
}

type FilterCondition struct {
	Field string
	Condition
}

type SortOption struct {
	Field     string
	Direction SortDirection
}

type Filter interface {
	GetConditions() []FilterCondition
	GetFilterQuery(query *gorm.DB, conditions []FilterCondition) *gorm.DB
}

type SingleFilter struct {
	Logic LogicOperator
}

type MultiFilter struct {
	Logic  LogicOperator
	Limit  *int
	Offset *int
	Sort   []SortOption
}

func (f *SingleFilter) GetConditions() []FilterCondition {
	return []FilterCondition{}
}

func (f *SingleFilter) GetFilterQuery(query *gorm.DB, conditions []FilterCondition) *gorm.DB {
	return applyConditions(query, conditions, f.Logic)
}

func (f *MultiFilter) GetConditions() []FilterCondition {
	return []FilterCondition{}
}

func (f *MultiFilter) GetFilterQuery(query *gorm.DB, conditions []FilterCondition) *gorm.DB {
	query = applyConditions(query, conditions, f.Logic)

	for _, sort := range f.Sort {
		query = query.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
	}

	if f.Limit != nil {
		query = query.Limit(*f.Limit)
	}

	if f.Offset != nil {
		query = query.Offset(*f.Offset)
	}

	return query
}

func applyConditions(query *gorm.DB, conditions []FilterCondition, logic LogicOperator) *gorm.DB {
	for i, condition := range conditions {
		clause := fmt.Sprintf("%s %s ?", condition.Field, condition.Operator)
		if i == 0 {
			query = query.Where(clause, condition.Value)
		} else {
			if logic == And {
				query = query.Where(clause, condition.Value)
			} else {
				query = query.Or(clause, condition.Value)
			}
		}
	}
	return query
}
