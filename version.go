package main

import "runtime"

// The current version of chef-runner. A ".dev" suffix denotes that the version
// is currently being developed.
const Version = "v0.9.0"

// GitVersion is the Git version that is being compiled. This string contains
// tag and commit information. It will be filled in by the compiler.
var GitVersion string

// VersionString returns the current program version, which is either the Git
// version if available or the static version defined above.
func VersionString() string {
	if GitVersion != "" {
		return GitVersion
	}
	return Version
}

// GoVersionString returns the Go version used to build the program.
func GoVersionString() string {
	return runtime.Version()
}

// TargetString returns the target operating system and architecture.
func TargetString() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
