// Package driver defines the interface that all drivers need to implement. A
// driver is responsible for running commands on an instance using whatever
// mechanism is available.
package driver

type Driver interface {
	RunCommand(command string) error
	String() string
}
