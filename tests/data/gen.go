package main

import (
	"log"
	"path"
	"path/filepath"

	"github.com/common-fate/go-oas3/configurator"
	"github.com/common-fate/go-oas3/generator"
	"github.com/common-fate/go-oas3/loader"
)

//go:generate go run gen.go

func main() {
	err := generate()
	if err != nil {
		log.Fatal(err)
	}
}

func generate() error {
	files, err := filepath.Glob("./*/swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		cfg := configurator.Config{SwaggerAddr: f, Package: "output", ComponentsPackage: "output"}

		l := loader.New(cfg)
		spec, err := l.Load()
		if err != nil {
			return err
		}
		g := generator.New(cfg)

		res := g.Generate(spec)

		folder := filepath.Dir(f)

		err = res.RouterCode.Save(path.Join(folder, "router.go"))
		if err != nil {
			return err
		}

		err = res.ComponentsCode.Save(path.Join(folder, "components.go"))
		if err != nil {
			return err
		}
	}

	// fmt.Fprintln(os.Stderr, files) // contains a list of all files in the current directory
	// return errors.New("ala")
	return nil
}
