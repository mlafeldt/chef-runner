// Package provisioner defines the interface that all provisioners need to
// implement.
package provisioner

// A Provisioner is responsible for provisioning a machine with Chef.
type Provisioner interface {
	PrepareFiles() error
	ProvisionCommand() []string
}
