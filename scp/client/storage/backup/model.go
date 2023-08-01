package backup

type CreateBackupRequest struct {
	AzCode                     string
	BackupDrZoneId             string
	BackupName                 string
	BackupPolicyTypeCategory   string
	BackupRepository           string
	DrAzCode                   string
	FileSystemBackupSelections []string
	IsBackupDrEnabled          string
	ObjectId                   string
	ObjectType                 string
	PolicyType                 string
	ProductNames               []string
	RetentionPeriod            string
	Schedules                  []BackupScheduleInfo
	ServiceZoneId              string
	Tags                       []TagRequest
}

type BackupScheduleInfo struct {
	ScheduleFrequency       string
	ScheduleFrequencyDetail string
	ScheduleId              string
	ScheduleName            string
	ScheduleType            string
	StartTime               string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}

type UpdateBackupScheduleRequest struct {
	Schedules []BackupScheduleInfo
}
