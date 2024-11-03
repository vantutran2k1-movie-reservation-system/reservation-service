package filters

import "gorm.io/gorm"

type LoginTokenFilter struct {
	Filter
	UserID     *Condition
	TokenValue *Condition
	ExpiresAt  *Condition
}

func (f *LoginTokenFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.UserID != nil {
		conditions = append(conditions, f.UserID.ToFilterCondition("user_id"))
	}

	if f.TokenValue != nil {
		conditions = append(conditions, f.TokenValue.ToFilterCondition("token_value"))
	}

	if f.ExpiresAt != nil {
		conditions = append(conditions, f.ExpiresAt.ToFilterCondition("expires_at"))
	}

	return conditions
}

func (f *LoginTokenFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
