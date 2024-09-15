package payloads

type CreateUserProfileRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=1,max=255"`
	LastName    string `json:"last_name" binding:"required,min=1,max=255"`
	PhoneNumber string `json:"phone_number" binding:"phoneNumber"`
	DateOfBirth string `json:"date_of_birth" binding:"date,beforeToday"`
}
