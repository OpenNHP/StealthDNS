package version

var (
	// Version is the version number, set at build time via ldflags
	Version = "1.0.0"
	// BuildNumber is the build number, set at build time
	BuildNumber = ""
	// CommitID is the git commit hash, set at build time
	CommitID = ""
	// BuildTime is the build timestamp, set at build time
	BuildTime = ""
)

