package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"regexp"
)

const (
	InputRegex  = "(?P<repo>[^\\r\\n\\t\\f\\v@ ]+)@v?(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?"
	ImportRegex = "import\\s+?(?:(?:(?:[\\w*\\s{},]*)\\s+from\\s+?)|)(?:(?:\"(?P<ImportPath>.*?)\")|(?:'(?P<ImportPath>.*?)'))[\\s]*?(?:;|$|)"
)

type Sol struct {
	name, path string
	deps       []*Sol
}

type Dependency struct {
	name, version, path string
}

type Extractor struct {
	dep  *Dependency
	sols []*Sol
}

func NewDependency(dependency string) (*Dependency, error) {
	match, err := regexp.MatchString(InputRegex, dependency)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf("Bad repo format, try [user]/[repo]@[version]")
	}
	splitter, err := regexp.Compile("@")
	if err != nil {
		return nil, err
	}

	split := splitter.Split(dependency, 2)

	return &Dependency{
		name:    split[0],
		version: split[1],
		path:    dependency,
	}, nil

}

func (e *Extractor) GetFileList(selected *hashset.Set) []string {
	list := make([]string, len(e.sols))
	for i, sol := range e.sols {
		isSelected := " "
		if selected.Contains(sol) {
			isSelected = "X"
		}
		list[i] = fmt.Sprintf("[%s] %s (%s)", isSelected, sol.name, sol.path)
	}
	return list
}

func (e *Extractor) GetSol(n int) *Sol {
	return e.sols[n]
}
