package constants

const (
	UserSession = "user_session"

	// Gin mode
	GinReleaseMode = "release"
	GinDebugMode   = "debug"

	// Request headers
	ProfilePictureRequestFormKey = "profile_picture"
	ResetToken                   = "Reset-Token"

	// Request query params
	IncludeUserProfile     = "includeProfile"
	IncludeGenres          = "includeGenres"
	IncludeTheaterLocation = "includeLocation"

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
