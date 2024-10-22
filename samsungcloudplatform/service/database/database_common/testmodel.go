package database_common

type TestModel struct {
	ArchiveBackupScheduleFrequency string
	BackupRetentionPeriod          string
	BackupStartHour                int32
	DnsServerIps                   []string
	AuditEnabled                   bool
}
