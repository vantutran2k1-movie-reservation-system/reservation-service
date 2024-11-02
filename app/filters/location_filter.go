package filters

import "gorm.io/gorm"

type CountryFilter struct {
	Filter
	ID   *Condition
	Name *Condition
	Code *Condition
}

func (f *CountryFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.Name != nil {
		conditions = append(conditions, f.Name.ToFilterCondition("name"))
	}

	if f.Code != nil {
		conditions = append(conditions, f.Code.ToFilterCondition("code"))
	}

	return conditions
}

func (f *CountryFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}

type StateFilter struct {
	Filter
	ID        *Condition
	CountryID *Condition
	Name      *Condition
}

func (f *StateFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.CountryID != nil {
		conditions = append(conditions, f.CountryID.ToFilterCondition("country_id"))
	}

	if f.Name != nil {
		conditions = append(conditions, f.Name.ToFilterCondition("name"))
	}

	return conditions
}

func (f *StateFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}

type CityFilter struct {
	Filter
	ID      *Condition
	StateID *Condition
	Name    *Condition
}

func (f *CityFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.StateID != nil {
		conditions = append(conditions, f.StateID.ToFilterCondition("state_id"))
	}

	if f.Name != nil {
		conditions = append(conditions, f.Name.ToFilterCondition("name"))
	}

	return conditions
}

func (f *CityFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
