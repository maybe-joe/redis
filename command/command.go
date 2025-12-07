package command

type Type string

const (
	UNKNOWN Type = "UNKNOWN"
	PING    Type = "PING"
	ECHO    Type = "ECHO"
)

type Command struct {
	Type Type
	Args []string
}

func Unknown() Command {
	return Command{Type: UNKNOWN}
}

func Ping() Command {
	return Command{Type: PING}
}

func Echo(message string) Command {
	return Command{Type: ECHO, Args: []string{message}}
}
