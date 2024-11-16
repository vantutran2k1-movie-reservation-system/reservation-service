package filters

import "gorm.io/gorm"

type MovieFilter struct {
	Filter
	ID        *Condition
	IsActive  *Condition
	IsDeleted *Condition
}

func (f *MovieFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.IsActive != nil {
		conditions = append(conditions, f.IsActive.ToFilterCondition("is_active"))
	}

	if f.IsDeleted != nil {
		conditions = append(conditions, f.IsDeleted.ToFilterCondition("is_deleted"))
	}

	return conditions
}

func (f *MovieFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
