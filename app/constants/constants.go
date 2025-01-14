package constants

const (
	// Gin mode
	GinReleaseMode = "release"
	GinDebugMode   = "debug"

	// Request headers
	ProfilePictureRequestFormKey = "profile_picture"
	ResetToken                   = "Reset-Token"
	UserVerificationToken        = "Verification-Token"

	// Request query params
	Limit                  = "limit"
	Offset                 = "offset"
	IncludeUserProfile     = "includeProfile"
	IncludeGenres          = "includeGenres"
	IncludeTheaterLocation = "includeLocation"
	MaxDistance            = "distance"
	Email                  = "email"

	// Content types
	ContentType     = "Content-Type"
	ImageJpeg       = "image/jpeg"
	ImagePng        = "image/png"
	ApplicationJson = "application/json"

	// Configcat feature flags
	CanModifyUsers     = "canModifyUsers"
	CanModifyMovies    = "canModifyMovies"
	CanModifyGenres    = "canModifyGenres"
	CanModifyLocations = "canModifyLocations"
	CanModifyTheaters  = "canModifyTheaters"
)
