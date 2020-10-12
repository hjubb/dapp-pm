package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Storage struct {
	root     string
	fetchMut sync.RWMutex
}

func NewStorageEngine() (*Storage, error) {
	dir, err := homedir.Expand("~/.dapp-pm")

	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return nil, err
	}

	return &Storage{
		root: dir,
	}, nil
}

func (s *Storage) CheckExisting(dep *Dependency) (bool, error) {
	s.fetchMut.RLock()
	defer s.fetchMut.RUnlock()

	file, err := os.Open(path.Join(s.root, dep.path))
	defer file.Close()
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s *Storage) InitializeDependency(dep *Dependency) error {
	s.fetchMut.Lock()
	defer s.fetchMut.Unlock()

	depPath := path.Join(s.root, dep.path)

	repo, err := git.PlainClone(depPath, false, &git.CloneOptions{
		URL:               fmt.Sprintf("https://github.com/%s.git", dep.name),
		ReferenceName:     plumbing.NewTagReferenceName(dep.version),
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: 1,
		Progress:          os.Stdout,
	})

	if err != nil {
		return err
	}

	fmt.Println("Initialized dependency!", repo)

	return nil
}

func (s *Storage) GetDependency(dependency string) (*Dependency, error) {

	dep, err := NewDependency(dependency)
	if err != nil {
		return nil, err
	}

	existing, err := s.CheckExisting(dep)
	if err != nil {
		return nil, err
	}

	if existing {
		// return the existing dependency
		fmt.Println("Using existing dependency")
		return dep, nil
	}

	// otherwise initialise and pull
	err = s.InitializeDependency(dep)

	if err != nil {
		return nil, err
	}

	return dep, nil
}

func (s *Storage) ExtractPaths(dep *Dependency) (*Extractor, error) {
	gitRoot := s.ResolveDepPath(dep)
	sols := make([]*Sol, 0)

	checker := func(path string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sol") {
			sols = append(sols, &Sol{
				name: info.Name(),
				path: path,
				deps: make([]*Sol, 0),
			})
		}
		return nil
	}
	filepath.Walk(path.Join(gitRoot, "src"), checker)
	filepath.Walk(path.Join(gitRoot, "contracts"), checker)
	filepath.Walk(path.Join(gitRoot, "libs"), checker)

	// check imports and sort
	err := s.ResolveImports(sols)
	if err != nil {
		return nil, err
	}

	return &Extractor{
		dep: dep,
		sols: sols,
	}, nil
}

func (s *Storage) ResolveImports(sols []*Sol) error {

	reg, err := regexp.Compile(ImportRegex)
	if err != nil {
		return err
	}

	for _, sol := range sols {
		file, err := os.Open(sol.path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// test the regex for imports here
			text := scanner.Text()
			if reg.MatchString(text) {
				matches := reg.FindStringSubmatch(text)

				var match string
				if len(matches[1]) == 0 {
					match = matches[2]
				} else {
					match = matches[1]
				}

				solDir := path.Dir(sol.path)

				depPath := path.Join(solDir, match)
				if el := sort.Search(len(sols), func(i int) bool {
					return sols[i].path == depPath
				}); el < len(sols) {
					sol.deps = append(sol.deps, sols[el])
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}

	}

	return nil
}

func (s *Storage) ResolveDepPath(dep *Dependency) string {
	return path.Join(s.root, dep.path)
}
