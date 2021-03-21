package envy

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gomodul/godotenv"
	"github.com/rogpeppe/go-internal/modfile"
)

var splitter = "/"

func init() {
	if runtime.GOOS == "windows" {
		splitter = "\\"
	}
}

// Get args[0] = "Key Name", args[1] = "Default Value", args[2] = "file name or dir location filename".
func Get(args ...string) string {
	var key, defaultValue string
	if len(args) < 1 {
		return ""
	}

	key = args[0]
	if len(args) > 1 {
		defaultValue = args[1]
	}

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

	return defaultValue
}

// Load godoc.
func Load(args ...string) {
	cwdSplit := make([]string, 0)
	cwdSplitLength := 0
	fileName := ".env"

	if len(args) > 1 {
		log.Println("just need 1 arg")
	}

	if len(args) > 0 {
		arg := args[0]
		if len(arg) > 0 {
			cwdSplit = strings.Split(arg, "\\")
			cwdSplit = strings.Split(strings.Join(cwdSplit, splitter), splitter)
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

// CurrentFolderName godoc.
func CurrentFolderName() string {
	pkgName := CurrentPkgName()
	if pkgName == "" {
		return ""
	}

	pkgNamrArr := strings.Split(pkgName, splitter)
	if len(pkgNamrArr) > 0 {
		return pkgNamrArr[len(pkgNamrArr)-1]
	}

	return ""
}

func currentPackage() string {
	pwd, _ := os.Getwd()
	return importPath(pwd)
}
