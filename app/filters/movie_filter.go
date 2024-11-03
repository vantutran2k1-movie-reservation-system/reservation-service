package filters

import "gorm.io/gorm"

type MovieFilter struct {
	Filter
	ID *Condition
}

func (f *MovieFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	return conditions
}

func (f *MovieFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
