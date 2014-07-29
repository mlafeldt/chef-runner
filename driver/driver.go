package driver

type Driver interface {
	RunCommand(command string) error
	String() string
}
