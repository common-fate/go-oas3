package loader

import (
	"testing"

	"github.com/common-fate/go-oas3/configurator"
)

func TestLoader(t *testing.T) {
	l := New(configurator.Config{SwaggerAddr: "../tests/testdata/simple/swagger.yaml"})
	_, err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
}
