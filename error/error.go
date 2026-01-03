package errconst

import "errors"

// Validation errors
var (
	ErrQueryParamRequired = errors.New("query parameter is required")
	ErrInvalidLatitude    = errors.New("invalid latitude value")
	ErrInvalidLongitude   = errors.New("invalid longitude value")
	ErrInvalidRadius      = errors.New("invalid radius value")
	ErrInvalidPage        = errors.New("invalid page value")
	ErrInvalidLimit       = errors.New("invalid limit value")
	ErrLimitExceed        = errors.New("limit exceeds max allowed")
	ErrRadiusExceed       = errors.New("radius exceeds max allowed")
	ErrInvalidScoreValue  = errors.New("score must be a valid number")
	ErrScoreThreshold     = errors.New("score must be between 0.7 and 1.0")
	ErrLatLonRequired     = errors.New("lat and lon parameters are required")
)

// News retrieval errors
var (
	ErrGettingNews         = errors.New("failed to fetch news")
	ErrGettingNearbyNews   = errors.New("failed to fetch nearby news")
	ErrGettingCategoryNews = errors.New("failed to fetch category news")
	ErrGettingSourceNews   = errors.New("failed to fetch source news")
	ErrGettingScoreNews    = errors.New("failed to fetch score-based news")
)
