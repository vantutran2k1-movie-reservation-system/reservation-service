package filters

import "gorm.io/gorm"

type GenreFilter struct {
	Filter
	ID   *Condition
	Name *Condition
}

func (f *GenreFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.Name != nil {
		conditions = append(conditions, f.Name.ToFilterCondition("name"))
	}

	return conditions
}

func (f *GenreFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
