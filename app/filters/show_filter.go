package filters

import "gorm.io/gorm"

type ShowFilter struct {
	Filter
	Id        *Condition
	MovieId   *Condition
	TheaterId *Condition
	StartTime *Condition
	EndTime   *Condition
	Status    *Condition
}

func (f *ShowFilter) GetConditions() []FilterCondition {
	var conditions []FilterCondition

	if f.Id != nil {
		conditions = append(conditions, f.Id.ToFilterCondition("id"))
	}

	if f.MovieId != nil {
		conditions = append(conditions, f.MovieId.ToFilterCondition("movie_id"))
	}

	if f.TheaterId != nil {
		conditions = append(conditions, f.TheaterId.ToFilterCondition("theater_id"))
	}

	if f.StartTime != nil {
		conditions = append(conditions, f.StartTime.ToFilterCondition("start_time"))
	}

	if f.EndTime != nil {
		conditions = append(conditions, f.EndTime.ToFilterCondition("end_time"))
	}

	if f.Status != nil {
		conditions = append(conditions, f.Status.ToFilterCondition("status"))
	}

	return conditions
}

func (f *ShowFilter) GetFilterQuery(query *gorm.DB) *gorm.DB {
	return f.Filter.GetFilterQuery(query, f.GetConditions())
}
