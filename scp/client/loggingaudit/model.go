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
	IsLoggingTargetAllUser     bool
	LoggingTargetUsers         []string
	IsLoggingTargetAllResource bool
	LoggingTargetResourceIds   []string
	TrailSaveType              string
	TrailDescription           string
	IsLoggingTargetAllRegion   bool
	LoggingTargetRegions       []string
	UseVerification            bool
}

type UpdateTrailRequest struct {
	TrailUpdateType            string
	ObsBucketId                string
	IsLoggingTargetAllUser     string
	LoggingTargetUsers         []string
	IsLoggingTargetAllResource string
	LoggingTargetResourceIds   []string
	TrailSaveType              string
	TrailDescription           string
	IsLoggingTargetAllRegion   string
	LoggingTargetRegions       []string
	UseVerification            bool
}
