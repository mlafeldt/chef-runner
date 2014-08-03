package main

import "runtime"

// The current version of chef-runner.
const Version = "v0.4.0"

// The Git version that is being compiled. This string contains tag and commit
// information. It will be filled in by the compiler.
var GitVersion string

// The current program version, which is either the Git version if available or
// the static version defined above.
func VersionString() string {
	if GitVersion != "" {
		return GitVersion
	}
	return Version
}

// The target operating system and architecture.
func TargetString() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
