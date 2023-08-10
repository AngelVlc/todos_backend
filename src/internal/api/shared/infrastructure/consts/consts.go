package consts

type contextKey string

const (
	ReqContextUserIDKey      contextKey = "userID"
	ReqContextUserNameKey    contextKey = "userName"
	ReqContextUserIsAdminKey contextKey = "userIsAdmin"
	ReqContextRequestKey     contextKey = "requestID"
	ReqContextStartTime      contextKey = "startTime"
)
