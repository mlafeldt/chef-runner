package main

// The current version of chef-runner.
const Version = "v0.4.0.dev"

// The Git version that is being compiled. This string contains tag and commit
// information. It will be filled in by the compiler.
var GitVersion string

func VersionString() string {
	if GitVersion != "" {
		return GitVersion
	}
	return Version
}
