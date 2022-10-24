package loggingaudit

type UserResponse struct {
	Email       string
	UserId      string
	UserLoginId string
	UserName    string
}

type CreateTrailRequest struct {
	TrailName                  string
	ObsBucketId                string
	IsLoggingTargetAllUser     string
	LoggingTargetUsers         []UserResponse
	IsLoggingTargetAllResource string
	LoggingTargetResourceIds   []string
	TrailSaveType              string
	TrailDescription           string
}
