package version

import "fmt"

const AppName = "pixelc"

var (
	Version   = "0.0.0-dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func FullVersion() string {
	return fmt.Sprintf("%s %s (commit=%s build_date=%s)", AppName, Version, Commit, BuildDate)
}
