package filters

import "gorm.io/gorm"

type UserProfileFilter struct {
	Filter
	UserID *Condition
}

func (f *UserProfileFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.UserID != nil {
		conditions = append(conditions, f.UserID.ToFilterCondition("user_id"))
	}

	return conditions
}

func (f *UserProfileFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
