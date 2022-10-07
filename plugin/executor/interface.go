package plugin

type Executor interface {
	Execute(command string) (string, error)
}
