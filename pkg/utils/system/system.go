package system

import (
	"os"
	"path/filepath"
)

const defaultVelaHome = ".vela"

const (
	// VelaHomeEnv defines vela home system env
	VelaHomeEnv = "VELA_HOME"
	// StorageDriverEnv defines vela storage driver env
	StorageDriverEnv = "STORAGE_DRIVER"
)

// GetVelaHomeDir return vela home dir
func GetVelaHomeDir() (string, error) {
	if custom := os.Getenv(VelaHomeEnv); custom != "" {
		return custom, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultVelaHome), nil
}

// GetDefaultFrontendDir return default vela frontend dir
func GetDefaultFrontendDir() (string, error) {
	home, err := GetVelaHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "frontend"), nil
}

// GetCapCenterDir return cap center dir
func GetCapCenterDir() (string, error) {
	home, err := GetVelaHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "centers"), nil
}

// GetRepoConfig return repo config
func GetRepoConfig() (string, error) {
	home, err := GetCapCenterDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "config.yaml"), nil
}

// GetCapabilityDir return capability dirs including workloads and traits
func GetCapabilityDir() (string, error) {
	home, err := GetVelaHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "capabilities"), nil
}

// GetEnvDirByName will get env dir from name
func GetEnvDirByName(name string) string {
	homedir, err := GetVelaHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homedir, "envs", name)
}

// InitDirs create dir if not exits
func InitDirs() error {
	if err := InitCapabilityDir(); err != nil {
		return err
	}
	if err := InitCapCenterDir(); err != nil {
		return err
	}
	return nil
}

// InitCapCenterDir create dir if not exits
func InitCapCenterDir() error {
	home, err := GetCapCenterDir()
	if err != nil {
		return err
	}
	_, err = CreateIfNotExist(filepath.Join(home, ".tmp"))
	return err
}

// InitCapabilityDir create dir if not exits
func InitCapabilityDir() error {
	dir, err := GetCapabilityDir()
	if err != nil {
		return err
	}
	_, err = CreateIfNotExist(dir)
	return err
}

// CreateIfNotExist create dir if not exist
func CreateIfNotExist(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// nolint:gosec
			return false, os.MkdirAll(dir, 0755)
		}
		return false, err
	}
	return true, nil
}
