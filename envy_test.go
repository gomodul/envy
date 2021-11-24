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
}

func TestGoPath(t *testing.T) {
	r := require.New(t)

	r.Equal(os.Getenv("GOPATH"), envy.GoPath())
}

func TestList(t *testing.T) {
	r := require.New(t)

	list := envy.List()

	r.NotZero(os.Getenv(GOPATH))
	r.Equal("", list[envy.Version])
	r.Equal(os.Getenv(GOPATH), list[GOPATH])
}

func TestStage(t *testing.T) {
	key := "ENV"
	stage := envy.Get(key)

	var err error

	tests := []struct {
		name string
		env  string
		want string
	}{
		{
			name: "should have stage dev",
			env:  "dev  ",
			want: "dev",
		},
		{
			name: "should have empty stage",
			env:  "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = envy.Set(key, tt.env)
			require.NoError(t, err)

			if got := envy.Stage(); got != tt.want {
				t.Errorf("Stage() = %v, want %v", got, tt.want)
			}
		})
	}

	err = envy.Set(key, stage)
	require.NoError(t, err)
}
