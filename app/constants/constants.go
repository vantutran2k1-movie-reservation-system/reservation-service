package constants

const (
	// Gin mode
	GinReleaseMode = "release"
	GinDebugMode   = "debug"

	// Request headers
	ProfilePictureRequestFormKey = "Profile-Picture"
	UserPasswordResetToken       = "Reset-Token"
	UserVerificationToken        = "Verification-Token"
	RetryAfter                   = "Retry-After"

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

	// Redis key
	ClientRateLimit = "rateLimit"
)

type SeatType string

const (
	Regular SeatType = "REGULAR"
	Vip     SeatType = "VIP"
)

type ShowStatus string

const (
	Active    ShowStatus = "ACTIVE"
	Cancelled ShowStatus = "CANCELLED"
	Completed ShowStatus = "COMPLETED"
	Expired   ShowStatus = "EXPIRED"
	Scheduled ShowStatus = "SCHEDULED"
	OnHold    ShowStatus = "ON-HOLD"
)
