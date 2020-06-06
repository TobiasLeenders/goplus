package golang

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/qiniu/goplus/ast"
	"github.com/qiniu/goplus/cl"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"
	"github.com/qiniu/x/log"
)

func saveGoFile(dir string, code *Code) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	b, err := code.Bytes(nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dir+"/gop_autogen.go", b, 0666)
}

// -----------------------------------------------------------------------------

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

func testGenGo(t *testing.T, pkgDir, sel, exclude string) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	if exclude != "" && strings.Contains(pkgDir, exclude) {
		return
	}
	log.Debug("Compiling", pkgDir)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseDir failed:", err, len(pkgs))
	}

	bar := getPkg(pkgs)
	b := NewBuilder(bar.Name, nil, fset)
	_, err = cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClAll)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()
	err = saveGoFile(pkgDir, code)
	if err != nil {
		t.Fatal("saveGoFile failed:", err)
	}
}

func TestGenGofile(t *testing.T) {
	sel, exclude := "", "17"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir += "/testdata"
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		testGenGo(t, dir+"/"+fi.Name(), sel, exclude)
	}
}

// -----------------------------------------------------------------------------
