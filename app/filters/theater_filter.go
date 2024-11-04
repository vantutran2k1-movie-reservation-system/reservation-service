package filters

import "gorm.io/gorm"

type TheaterFilter struct {
	Filter
	ID   *Condition
	Name *Condition
}

type TheaterLocationFilter struct {
	Filter
	TheaterID *Condition
}

func (f *TheaterFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.ID != nil {
		conditions = append(conditions, f.ID.ToFilterCondition("id"))
	}

	if f.Name != nil {
		conditions = append(conditions, f.Name.ToFilterCondition("name"))
	}

	return conditions
}

func (f *TheaterFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}

func (f *TheaterLocationFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.TheaterID != nil {
		conditions = append(conditions, f.TheaterID.ToFilterCondition("theater_id"))
	}

	return conditions
}

func (f *TheaterLocationFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
