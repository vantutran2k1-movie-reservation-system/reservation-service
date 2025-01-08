package filters

import "gorm.io/gorm"

type UserFilter struct {
	Filter
	ID         *Condition
	Email      *Condition
	IsVerified *Condition
}

func (f *UserFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.Email != nil {
		conditions = append(conditions, f.Email.ToFilterCondition("email"))
	}

	if f.IsVerified != nil {
		conditions = append(conditions, f.IsVerified.ToFilterCondition("is_verified"))
	}

	return conditions
}

func (f *UserFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
