package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"go/parser"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type opt struct {
	pkg               string
	json              bool
	ignoreStdPkg      bool
	ignoreInternalPkg bool
	disableShowID     bool
}
type s struct {
	prog    *loader.Program
	opt     *opt
	arrived map[string]int
}

func isStdPackage(s *s, pkg *types.Package) bool {
	files := s.prog.Package(pkg.Path()).Files
	if len(files) > 0 {
		filepath := s.prog.Fset.Position(files[0].Package).Filename
		return strings.HasPrefix(filepath, build.Default.GOROOT)
	}
	// fmt.Println("!!!", pkg.Path()) // xxx (e.g. unsafe)
	return true
}

func isInternalPackage(s *s, pkg *types.Package) bool {
	return strings.Contains(pkg.Path(), "/internal/")
}

func guessPkg() (string, error) {
	curdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(curdir)
	if err != nil {
		return "", err
	}
	for _, srcdir := range build.Default.SrcDirs() {
		if strings.HasPrefix(path, srcdir) {
			pkgname := strings.TrimLeft(strings.Replace(path, srcdir, "", 1), "/")
			return pkgname, nil
		}
	}

	return "", errors.Errorf("%q is not subdir of srcdirs(%q)", path, build.Default.SrcDirs())
}

func load(pkg string) (*loader.Program, error) {
	conf := loader.Config{
		ParserMode: parser.ImportsOnly,
		TypeCheckFuncBodies: func(path string) bool {
			return false
		},
	}
	conf.Import(pkg)

	prog, err := conf.Load()
	if err != nil {
		return nil, errors.Wrap(err, "load")
	}

	return prog, nil
}

func dump(pkg *types.Package, s *s, depth int) error {
	id, arrived := s.arrived[pkg.Path()]
	if !arrived {
		id = len(s.arrived)
		s.arrived[pkg.Path()] = id
	}
	if (!s.opt.ignoreStdPkg || !isStdPackage(s, pkg)) && (!s.opt.ignoreInternalPkg || !isInternalPackage(s, pkg)) {
		if s.opt.disableShowID {
			fmt.Printf("%s%s\n", strings.Repeat("  ", depth), pkg.Path())
		} else {
			fmt.Printf("%s%s #=%d\n", strings.Repeat("  ", depth), pkg.Path(), id)
		}
	}
	if arrived {
		return nil
	}

	deps := pkg.Imports()
	sort.Slice(deps, func(i int, j int) bool {
		return deps[i].Name() < deps[j].Name()
	})
	for _, deppkg := range deps {
		dump(deppkg, s, depth+1)
	}
	return nil
}

type tree struct {
	ID           int     `json:"id"`
	Pkg          string  `json:"pkg"`
	Dependencies []*tree `json:"dependencies"`
	Depth        int     `json:"-"`
}

func buildtree(pkg *types.Package, s *s, depth int) *tree {
	id, arrived := s.arrived[pkg.Path()]
	if !arrived {
		id = len(s.arrived)
		s.arrived[pkg.Path()] = id
	}

	t := &tree{Pkg: pkg.Path(), Depth: depth, Dependencies: []*tree{}}
	if !s.opt.disableShowID {
		t.ID = id
	}

	if !arrived {
		deps := pkg.Imports()
		sort.Slice(deps, func(i int, j int) bool {
			return deps[i].Name() < deps[j].Name()
		})
		for _, deppkg := range deps {
			child := buildtree(deppkg, s, depth+1)
			if child != nil {
				t.Dependencies = append(t.Dependencies, child)
			}
		}
	}

	if s.opt.ignoreStdPkg && isStdPackage(s, pkg) {
		return nil
	}
	if s.opt.ignoreInternalPkg && isInternalPackage(s, pkg) {
		return nil
	}
	return t
}

func dumpJSON(pkg *types.Package, s *s, depth int) error {
	t := buildtree(pkg, s, depth)
	if t == nil {
		return errors.New("not found")
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(t)
}

func run(opt *opt) error {
	prog, err := load(opt.pkg)
	if err != nil {
		return err
	}
	s := &s{
		arrived: map[string]int{},
		opt:     opt,
		prog:    prog,
	}
	dumpfn := dump
	if opt.json {
		dumpfn = dumpJSON
	}
	if err := dumpfn(prog.Package(opt.pkg).Pkg, s, 0); err != nil {
		return err
	}
	return nil
}

func main() {
	app := kingpin.New("pkgtree", "dump pkg dependencies")
	var opt opt
	app.Arg("pkg", "pkg").Required().StringVar(&opt.pkg)
	app.Flag("json", "").BoolVar(&opt.json)
	app.Flag("ignore-std-pkg", "").BoolVar(&opt.ignoreStdPkg)
	app.Flag("ignore-internal-pkg", "").BoolVar(&opt.ignoreInternalPkg)
	app.Flag("disable-show-id", "").BoolVar(&opt.disableShowID)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(fmt.Sprintf("%v", err))
	}

	if opt.pkg == "" || opt.pkg == "." {
		pkg, err := guessPkg()
		if err != nil {
			app.FatalUsage(fmt.Sprintf("%v", err))
		}
		opt.pkg = pkg
		log.Printf("guess pkg name .. %q\n", opt.pkg)
	}
	if err := run(&opt); err != nil {
		log.Fatalf("!!%+v", err)
	}
}
