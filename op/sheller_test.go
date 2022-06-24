package op

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSheller(t *testing.T) {
	// only run on not-windows, and explicitly set SHELL env
	if runtime.GOOS == "windows" {
		return
	}
	curShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", curShell)
	os.Setenv("SHELL", "/bin/sh")

	// features pipe "|" and file redirection ">"
	c := NewSheller("echo hiii | grep hi > cooltestfile")
	defer os.Remove("cooltestfile")
	o := c.Run()
	assert.Equal(t, "", o.Result)
	assert.NoError(t, o.Error)

	bts, err := ioutil.ReadFile("cooltestfile")
	assert.Equal(t, "hiii\n", string(bts))
	assert.NoError(t, err)
}
