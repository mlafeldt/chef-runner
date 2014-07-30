// Package driver defines the interface that all drivers need to implement. A
// driver is responsible for running commands on and uploading files to a
// machine using whatever mechanism is available.
package driver

type Driver interface {
	RunCommand(command string) error
	Upload(dst string, src ...string) error
	String() string
}
