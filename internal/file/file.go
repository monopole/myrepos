package file

import (
	"fmt"
	"os"
	"path/filepath"
)

type Path string

const defaultConfigFileName = ".myrepos"

var extensions = []string{".yml", ".yaml"}

func DefaultConfigFileName() string {
	return filepath.Join(string(Home()), defaultConfigFileName+extensions[0])
}

func computeDefaultConfigFilePath() (Path, error) {
	var badFiles []Path
	for i := range extensions {
		p := Path(defaultConfigFileName + extensions[i])
		if exists, isDir := p.Exists(); exists && !isDir {
			return p, nil
		}
		badFiles = append(badFiles, p)
		if home := Home(); home != "" {
			p = home.Append(p)
			if exists, isDir := p.Exists(); exists && !isDir {
				return p, nil
			}
			badFiles = append(badFiles, p)
		}
	}
	return "", fmt.Errorf("unable to open any of these: %v", badFiles)
}

func GetFilePath(args []string) ([]Path, error) {
	if len(args) == 0 {
		p, err := computeDefaultConfigFilePath()
		return []Path{p}, err
	}
	var result []Path
	for i := range args {
		p := Path(args[i])
		if exists, isDir := p.Exists(); exists && !isDir {
			result = append(result, p)
		} else {
			if isDir {
				return nil, fmt.Errorf(
					"%q found, but it's a directory, not a config file", p)
			}
			return nil, fmt.Errorf("no config file found at %q", p)
		}
	}
	return result, nil
}

// Exists returns a boolean pair, the first is true if the path exists,
// the second is true if the path is a directory.
func (p Path) Exists() (bool, bool) {
	info, err := os.Stat(string(p))
	if os.IsNotExist(err) {
		return false, false
	}
	return true, info.IsDir()
}

func (p Path) IsAbs() bool {
	return filepath.IsAbs(string(p))
}

func (p Path) Append(n Path) Path {
	return Path(filepath.Join(string(p), string(n)))
}

func (p Path) MkDir() error {
	return os.MkdirAll(string(p), 0755)
}

func Home() Path {
	return Path(os.Getenv("HOME"))
}
