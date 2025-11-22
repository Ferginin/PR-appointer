package handler

type ErrorCode string

const (
	ErrCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrCodePRExists    ErrorCode = "PR_EXISTS"
	ErrCodePRMerged    ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
)

type APIError struct {
	Error struct {
		Code    ErrorCode `json:"code"`
		Message string    `json:"message"`
	} `json:"error"`
}

func newAPIError(code ErrorCode, message string) APIError {
	var apiErr APIError
	apiErr.Error.Code = code
	apiErr.Error.Message = message
	return apiErr
}
