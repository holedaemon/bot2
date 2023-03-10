package version

var version string

// Version returns a version string.
func Version() string {
	if version == "" {
		return "devel"
	}

	return version
}
