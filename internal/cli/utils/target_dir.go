package utils

import (
	"path/filepath"
)

func GetTargetDir(
	base, structure, serviceDir, flagDir string,
	fallbackFunc func(string, string, string) (string, error),
) (string, error) {
	if flagDir != "" {
		return filepath.Join(serviceDir, flagDir), nil
	}
	return fallbackFunc(base, structure, serviceDir)
}

func GetTargetDirWithDomainModule(
	serviceDir, domainDir, modularDir, flagDir string,
	fallbackFunc func(string, string, string, string, string) (string, error),
	base, structure string,
) (string, error) {
	if flagDir != "" {
		p := serviceDir
		if domainDir != "" {
			p = filepath.Join(p, domainDir)
		}
		if modularDir != "" {
			p = filepath.Join(p, "modules", modularDir)
		}
		p = filepath.Join(p, flagDir)
		return p, nil
	}
	return fallbackFunc(base, structure, serviceDir, domainDir, modularDir)
}
