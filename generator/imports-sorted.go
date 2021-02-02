package genutils

import (
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/kucjac/gentools/types"
)

var standardPackages = make(map[string]struct{})

func init() {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		standardPackages[p.PkgPath] = struct{}{}
	}
}

func isStandardPackage(pkg string) bool {
	_, ok := standardPackages[pkg]
	return ok
}

// Imports is a structure that contains groups of imports.
type Imports struct {
	StdLibs       []*types.Package
	External      []*types.Package
	Local         []*types.Package
	LocalInternal []*types.Package
}

// WriteTo writes the sorted imports into provided writer. The imports are used in a golang like style.
func (im *Imports) WriteTo(w io.Writer) (int64, error) {
	if len(im.StdLibs) == 0 && len(im.External) == 0 && len(im.Local) == 0 && len(im.LocalInternal) == 0 {
		return 0, nil
	}

	var nt int
	n, err := fmt.Fprintln(w, "import (")
	if err != nil {
		return 0, err
	}

	var needBreak bool
	if len(im.StdLibs) > 0 {
		needBreak = true
		for _, imp := range im.StdLibs {
			nt, err = fmt.Fprintf(w, "\t%q\n", imp.Path)
			if err != nil {
				return 0, err
			}
			n += nt
		}
	}
	if len(im.External) > 0 {
		if needBreak {
			nt, err = fmt.Fprintln(w, "")
			if err != nil {
				return 0, err
			}
			n += nt
		}
		needBreak = true
		for _, imp := range im.External {
			nt, err = fmt.Fprintf(w, "\t%q\n", imp.Path)
			if err != nil {
				return 0, err
			}
			n += nt
		}
	}
	if len(im.Local) > 0 {
		if needBreak {
			nt, err = fmt.Fprintln(w, "")
			if err != nil {
				return 0, err
			}
			n += nt
		}
		needBreak = true
		for _, imp := range im.Local {
			nt, err = fmt.Fprintf(w, "\t%q\n", imp.Path)
			if err != nil {
				return 0, err
			}
			n += nt
		}
	}
	if len(im.LocalInternal) > 0 {
		if needBreak {
			nt, err = fmt.Fprintln(w, "")
			if err != nil {
				return 0, err
			}
			n += nt
		}
		for _, imp := range im.LocalInternal {
			nt, err = fmt.Fprintf(w, "\t%q\n", imp.Path)
			if err != nil {
				return 0, err
			}
			n += nt
		}
	}
	nt, err = fmt.Fprint(w, ")")
	if err != nil {
		return 0, err
	}
	n += nt
	return int64(n), nil
}

// String prints sorted imports.
func (im *Imports) String() string {
	sb := &strings.Builder{}
	im.WriteTo(sb)
	return sb.String()
}

// SortImports sorts and groups imports slice.
func SortImports(contextPackage string, imports []*types.Package) *Imports {
	var im Imports
	localInternal := path.Join(contextPackage, "internal")
	for _, imp := range imports {
		if isStandardPackage(imp.Path) {
			im.StdLibs = append(im.StdLibs, imp)
			continue
		}
		if !strings.HasPrefix(imp.Path, contextPackage) {
			im.External = append(im.External, imp)
			continue
		}

		if !strings.HasPrefix(imp.Path, localInternal) {
			im.Local = append(im.Local, imp)
			continue
		}

		im.LocalInternal = append(im.LocalInternal, imp)
	}
	sort.Slice(im.StdLibs, func(i, j int) bool { return im.StdLibs[i].Path < im.StdLibs[j].Path })
	sort.Slice(im.External, func(i, j int) bool { return im.External[i].Path < im.External[j].Path })
	sort.Slice(im.Local, func(i, j int) bool { return im.Local[i].Path < im.Local[j].Path })
	sort.Slice(im.LocalInternal, func(i, j int) bool { return im.LocalInternal[i].Path < im.LocalInternal[j].Path })
	return &im
}
