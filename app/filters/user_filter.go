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

type UserRegistrationTokenFilter struct {
	Filter
	TokenValue *Condition
	IsUsed     *Condition
	ExpiresAt  *Condition
}

func (f *UserRegistrationTokenFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.TokenValue != nil {
		conditions = append(conditions, f.TokenValue.ToFilterCondition("token_value"))
	}

	if f.IsUsed != nil {
		conditions = append(conditions, f.IsUsed.ToFilterCondition("is_used"))
	}

	if f.ExpiresAt != nil {
		conditions = append(conditions, f.ExpiresAt.ToFilterCondition("expires_at"))
	}

	return conditions
}

func (f *UserRegistrationTokenFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}

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

type PasswordResetTokenFilter struct {
	Filter
	ID         *Condition
	UserID     *Condition
	TokenValue *Condition
	IsUsed     *Condition
	ExpiresAt  *Condition
}

func (f *PasswordResetTokenFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.UserID != nil {
		conditions = append(conditions, f.UserID.ToFilterCondition("user_id"))
	}

	if f.TokenValue != nil {
		conditions = append(conditions, f.TokenValue.ToFilterCondition("token_value"))
	}

	if f.IsUsed != nil {
		conditions = append(conditions, f.IsUsed.ToFilterCondition("is_used"))
	}

	if f.ExpiresAt != nil {
		conditions = append(conditions, f.ExpiresAt.ToFilterCondition("expires_at"))
	}

	return conditions
}

func (f *PasswordResetTokenFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
