package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Analyze(cf Config) {
	for _, path := range cf.Paths {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("could not get file info for path %q: %s\n", path, err)
			continue
		}
		if info.IsDir() {
			analyzeDir(path, cf)
		} else {
			analyzeFile(path, cf)
		}
	}
}

func analyzeDir(path string, cf Config) {
	filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if !cf.SkipTest && strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		if !cf.SkipVendor && strings.Contains(path, "vendor") {
			return nil
		}
		if err == nil && isGoFile(entry) {
			analyzeFile(path, cf)
		}
		return err
	})
}

func analyzeFile(path string, cf Config) {
	file := path
	fset := token.NewFileSet()
	src, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var fundecl_map map[string]struct{} = make(map[string]struct{})
	var typ_name = ""
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Recv != nil {
				typ := fn.Recv.List[0].Type
				stn := strings.TrimPrefix(recvString(typ), "*")
				if _, ok := fundecl_map[stn]; !ok {
					fundecl_map[stn] = struct{}{}
					typ_name = stn
				} else {
					if typ_name != stn {
						fmt.Printf("%s:%v %s: bad order\n", file, fset.Position(decl.Pos()), funcName(fn))
						typ_name = stn
					}
				}
			}
		}
	}
}

func isGoFile(entry fs.DirEntry) bool {
	return !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go")
}

func Analyzebcx(path string, info fs.FileInfo, err error) error {

	if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {

	}
	return nil
}

func funcName(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		typ := fn.Recv.List[0].Type
		return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
	}
	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}
