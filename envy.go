package envy

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gomodul/godotenv"
	"github.com/rogpeppe/go-internal/modfile"
	"github.com/spf13/cast"
)

var splitter = "/"

func init() {
	if runtime.GOOS == "windows" {
		splitter = "\\"
	}
}

// Get args[0] = "Key Name", args[1] = "Default Value", args[2] = "file name or dir location filename".
func Get(args ...string) string {
	if len(args) < 1 {
		return ""
	}

	key := args[0]
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	var fileName string
	if len(args) > 2 {
		fileName = args[2]
	}

	Load(fileName)

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	var defaultValue string
	if len(args) > 1 {
		defaultValue = args[1]
	}
	return defaultValue
}

// GetBool return bool.
func GetBool(key ...string) bool {
	return cast.ToBool(Get(key...))
}

// GetInt return int.
func GetInt(key ...string) int {
	return cast.ToInt(Get(key...))
}

// GetInt32 return int32.
func GetInt32(key ...string) int32 {
	return cast.ToInt32(Get(key...))
}

// GetInt64 return int64.
func GetInt64(key ...string) int64 {
	return cast.ToInt64(Get(key...))
}

// GetUint return uint.
func GetUint(key ...string) uint {
	return cast.ToUint(Get(key...))
}

// GetUInt32 return uint32.
func GetUInt32(key ...string) uint32 {
	return cast.ToUint32(Get(key...))
}

// GetUInt64 return uint64.
func GetUInt64(key ...string) uint64 {
	return cast.ToUint64(Get(key...))
}

// GetTime return time.Time
func GetTime(key ...string) time.Time {
	return cast.ToTime(Get(key...))
}

// GetDuration return time.Duration
func GetDuration(key ...string) time.Duration {
	return cast.ToDuration(Get(key...))
}

// Set godoc.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Del unset env.
func Del(key string) error {
	return os.Unsetenv(key)
}

// Flush clear all env
func Flush() {
	os.Clearenv()
}

// List godoc.
func List() map[string]string {
	source := os.Environ()
	list := make(map[string]string, len(source))

	for _, e := range source {
		pair := strings.SplitN(e, "=", 2)
		list[pair[0]] = pair[1]
	}

	return list
}

// Stage get stage from GO_ENV or APP_ENV or ENV
func Stage() string {
	for _, v := range []string{"GO_ENV", "APP_ENV", "ENV"} {
		if value, exist := os.LookupEnv(v); exist {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// Load godoc.
func Load(args ...string) {
	cwdSplit := make([]string, 0)
	cwdSplitLength := 0

	if len(args) > 1 {
		log.Println("just need 1 arg")
	}

	var fileName string
	if len(args) > 0 && len(args[0]) > 0 {
		if strings.Index(args[0], "\\") > -1 {
			cwdSplit = strings.Split(args[0], "\\")
		} else {
			cwdSplit = strings.Split(args[0], "/")
		}
		cwdSplit = strings.Split(strings.Join(cwdSplit, splitter), splitter)
	} else {
		fileName = ".env"
		if stage := Stage(); stage != "" {
			fileName += "." + strings.ToLower(stage)
		}
	}

	if len(cwdSplit) > 0 {
		tmp := cwdSplit[:0]
		for _, v := range cwdSplit {
			if len(strings.TrimSpace(v)) > 0 {
				tmp = append(tmp, v)
			}
		}
		cwdSplit = tmp

		lastIndex := len(cwdSplit) - 1
		if len(cwdSplit[lastIndex]) > 3 && string(cwdSplit[lastIndex][0]) == "." {
			fileName = cwdSplit[lastIndex]
			if len(cwdSplit) == 1 {
				cwdSplit = nil
			} else {
				cwdSplit = cwdSplit[:lastIndex]
			}
		}
	}

	if len(cwdSplit) <= 0 {
		cwd, _ := os.Getwd()
		cwdSplit = strings.Split(cwd, splitter)
	}

	cwdSplitLength = len(cwdSplit)

	var filePathENV string
	for i := cwdSplitLength; i >= 0; i-- {
		filePathENV = strings.Join(cwdSplit[:i], splitter) + splitter + fileName

		matches, _ := filepath.Glob(filePathENV)
		if len(matches) != 0 {
			_ = godotenv.Overload(filePathENV)
			break
		}
	}
}

// GoPath godoc.
func GoPath() string {
	return Get("GOPATH", "")
}

// GoPaths godoc.
func GoPaths() []string {
	gp := Get("GOPATH", "")
	if runtime.GOOS == "windows" {
		return strings.Split(gp, ";")
	}
	return strings.Split(gp, ":")
}

func importPath(path string) string {
	path = strings.TrimPrefix(path, "/private")
	for _, gopath := range GoPaths() {
		srcpath := filepath.Join(gopath, "src")
		rel, err := filepath.Rel(srcpath, path)
		if err == nil {
			return filepath.ToSlash(rel)
		}
	}

	rel := strings.TrimPrefix(path, filepath.Join(GoPath(), "src"))
	rel = strings.TrimPrefix(rel, string(filepath.Separator))
	return filepath.ToSlash(rel)
}

// CurrentPkgName godoc.
func CurrentPkgName() string {
	mod, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return currentPackage()
	}

	packagePath := modfile.ModulePath(mod)
	if packagePath == "" {
		return ""
	}

	return packagePath
}

func currentPackage() string {
	pwd, _ := os.Getwd()
	return importPath(pwd)
}
