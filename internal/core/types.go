package core

type CommandType string

const (
	CmdSync         CommandType = "SYNC"
	CmdKick         CommandType = "KICK"
	CmdPing         CommandType = "PING"
	CmdCreateSecret CommandType = "CREATE_SECRET"
	CmdUpdateSecret CommandType = "UPDATE_SECRET"
	CmdIsolate      CommandType = "ISOLATE"
	CmdGetTraffic   CommandType = "GET_TRAFFIC"
	CmdRefreshMetrics CommandType = "REFRESH_METRICS"
	CmdBackup       CommandType = "BACKUP"
)

type Command struct {
	Type    CommandType
	Payload interface{}
	Result  chan interface{}
	Error   chan error
}
