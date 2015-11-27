package cassandra

type counterStmt bool

const (
	UP   counterStmt = true
	DOWN counterStmt = false
)

func (c counterStmt) String() string {
	sign := ""
	if bool(c) {
		sign = "+"
	} else {
		sign = "-"
	}
	return "UPDATE " + TABLE_NAME + " SET version = version " + sign + " 1 where versionRow = ?"
}
