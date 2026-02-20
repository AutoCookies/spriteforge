package compiler

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"runtime"

	"pixelc/internal/imageutil"
	"pixelc/internal/version"
)

func VersionString() string {
	return version.FullVersion()
}

func Doctor() (string, error) {
	tmp, err := os.MkdirTemp("", "pixelc-doctor-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.SetRGBA(0, 0, color.RGBA{R: 255, A: 255})
	p := filepath.Join(tmp, "doctor.png")
	if err := imageutil.SavePNG(p, img); err != nil {
		return "", err
	}
	if _, err := imageutil.LoadPNG(p); err != nil {
		return "", err
	}
	return fmt.Sprintf("doctor ok os=%s arch=%s temp=%s", runtime.GOOS, runtime.GOARCH, tmp), nil
}
