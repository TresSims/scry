package facts

import "fmt"

type Priority int

const (
	Emergency Priority = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Informational
	Debug
)

func (p Priority) String() string {
	switch p {
	case Emergency:
		return "EMERG"
	case Alert:
		return "ALERT"
	case Critical:
		return "CRIT"
	case Error:
		return "ERR"
	case Warning:
		return "WARN"
	case Notice:
		return "NOTICE"
	case Informational:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		return "Unrecognized Facility"
	}
}

type Facility int

const (
	Kern Facility = iota
	User
	Mail
	Daemon
	SysLog
	LPR
	News
	UUCP
	CRON
	AuthPriv
	FTP
	NTP
	Security
	Console
	SolarisCron
	Local0
	Local1
	Local2
	Local3
	Local4
	Local5
	Local6
	Local7
)

func (f Facility) String() string {
	switch f {
	case Kern:
		return "kern"
	case User:
		return "User"
	case Mail:
		return "Mail"
	case Daemon:
		return "Daemon"
	case SysLog:
		return "SysLog"
	case LPR:
		return "LPR"
	case News:
		return "News"
	case UUCP:
		return "UUCP"
	case CRON:
		return "CRON"
	case AuthPriv:
		return "AuthPriv"
	case FTP:
		return "FTP"
	case NTP:
		return "NTP"
	case Security:
		return "Security"
	case Console:
		return "Console"
	case SolarisCron:
		return "SolarisCron"
	case Local0:
		return "Local0"
	case Local1:
		return "Local1"
	case Local2:
		return "Local2"
	case Local3:
		return "Local3"
	case Local4:
		return "Local4"
	case Local5:
		return "Local5"
	case Local6:
		return "Local6"
	case Local7:
		return "Local7"
	default:
		return "Unrecognized Facility"
	}
}

type SyslogLine struct {
	Message    string   `json:"MESSAGE"`
	Priority   Priority `json:"PRIORITY"`
	Identifier string   `json:"SYSLOG_IDENTIFIER"`
	Facility   Facility `json:"SYSLOG_FACILITY"`
}

func (s SyslogLine) String() string {
	return fmt.Sprintf("[%s] %s (%s): %s", s.Priority, s.Identifier, s.Facility, s.Message)
}
