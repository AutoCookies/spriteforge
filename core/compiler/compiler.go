package compiler

import "pixelc/pkg/model"

func Compile(inputPath string, cfg model.Config) (*model.Atlas, []byte, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, err
	}
	return nil, nil, ErrNotImplemented
}
