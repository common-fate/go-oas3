package tests

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"testing"

	"github.com/common-fate/go-oas3/configurator"
	"github.com/common-fate/go-oas3/generator"
	"github.com/common-fate/go-oas3/loader"
	"github.com/stretchr/testify/assert"
)

func TestGolden(t *testing.T) {
	files, err := filepath.Glob("./testdata/*/swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		generate(f, t)
	}
}

func generate(f string, t *testing.T) {
	out := t.TempDir()
	t.Logf("running %s", f)
	parent := filepath.Dir(f)

	cfg := configurator.Config{SwaggerAddr: f, Package: "output", ComponentsPackage: "output"}

	l := loader.New(cfg)
	spec, err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	g := generator.New(cfg)

	res := g.Generate(spec)

	err = res.RouterCode.Save(path.Join(out, "router.go"))
	if err != nil {
		t.Fatal(err)
	}

	gotRouter, err := ioutil.ReadFile(path.Join(out, "router.go"))
	if err != nil {
		t.Fatal(err)
	}

	wantRouter, err := ioutil.ReadFile(path.Join(parent, "router.go"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(wantRouter), string(gotRouter))

	err = res.ComponentsCode.Save(path.Join(out, "components.go"))
	if err != nil {
		t.Fatal(err)
	}

	gotComp, err := ioutil.ReadFile(path.Join(out, "components.go"))
	if err != nil {
		t.Fatal(err)
	}

	wantComp, err := ioutil.ReadFile(path.Join(parent, "components.go"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(wantComp), string(gotComp))
}
