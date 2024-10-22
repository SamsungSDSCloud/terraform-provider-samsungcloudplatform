package database_common

const (
	CreatingState    string = "CREATING"
	ReservedState    string = "RESERVED"
	ActiveState      string = "ACTIVE"
	InActiveState    string = "INACTIVE"
	DeployingState   string = "DEPLOYING"
	DeletedState     string = "DELETED"
	TerminatingState string = "TERMINATING"
	RunningState     string = "RUNNING"
	AvailableState   string = "AVAILABLE"
	UnavailableState string = "UNAVAILABLE"
	UnknownState     string = "UNKNOWN"
	ErrorState       string = "ERROR"
	EditingState     string = "EDITING"
	StartingState    string = "STARTING"
	StoppingState    string = "STOPPING"
	StoppedState     string = "STOPPED"
	RestartingState  string = "RESTARTING"
	SoftDeletedState string = "SOFT_DELETED"
	UpgradingState   string = "UPGRADING"
)

const (
	None      string = "None"
	OneYear   string = "1 Year"
	ThreeYear string = "3 Year"
)

func DatabaseProcessingAndStoppedStates() []string {
	return []string{CreatingState, EditingState, StartingState, RestartingState, StoppingState, TerminatingState, UpgradingState, StoppedState}
}
