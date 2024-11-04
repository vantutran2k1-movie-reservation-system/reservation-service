package filters

import "gorm.io/gorm"

type UserFilter struct {
	Filter
	ID    *Condition
	Email *Condition
}

func (f *UserFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.Email != nil {
		conditions = append(conditions, f.Email.ToFilterCondition("email"))
	}

	return conditions
}

func (f *UserFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
