package version

const AppName = "pixelc"

const Version = "0.0.0-dev"

func FullVersion() string {
	return AppName + " " + Version
}
