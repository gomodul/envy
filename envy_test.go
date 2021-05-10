package envy_test

import (
	"os"
	"testing"

	"github.com/gomodul/envy"
	"github.com/stretchr/testify/require"
)

const GOPATH = "GOPATH"

func TestGet(t *testing.T) {
	r := require.New(t)

	r.NotZero(os.Getenv(GOPATH))
	r.Equal("", envy.Get(envy.Version))
	r.Equal(os.Getenv(GOPATH), envy.Get(GOPATH, "foo"))
}

func TestSet(t *testing.T) {
	r := require.New(t)

	before := os.Getenv(envy.Version)

	err := envy.Set(envy.Version, "foo")
	r.NoError(err)

	after := os.Getenv(envy.Version)
	r.NotEqual(before, after)

	err = envy.Set(envy.Version, before)
	r.NoError(err)

	after = os.Getenv(envy.Version)
	r.Equal(before, after)
}

func TestCurrentPkgName(t *testing.T) {
	r := require.New(t)

	r.Equal("github.com/gomodul/envy", envy.CurrentPkgName())
	r.Equal("envy", envy.CurrentFolderName())
}

func TestGoPath(t *testing.T) {
	r := require.New(t)

	r.Equal(os.Getenv("GOPATH"), envy.GoPath())
}

func TestCurrentFolderName(t *testing.T) {
	actual := envy.CurrentFolderName()
	if actual != "envy" {
		t.Fatalf("expected (envy), got (%v)", actual)
	}
}

func TestList(t *testing.T) {
	r := require.New(t)

	list := envy.List()

	r.NotZero(os.Getenv(GOPATH))
	r.Equal("", list[envy.Version])
	r.Equal(os.Getenv(GOPATH), list[GOPATH])
}
