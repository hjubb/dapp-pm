package main

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mitchellh/go-homedir"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

func GetDependencies(s *Sol) (*hashset.Set, error) {
	allDeps := hashset.New()
  allDeps.Add(s)
  if len(s.deps) == 0 {
    return allDeps, nil
  }

  for _, dep := range s.deps {
    result, err := GetDependencies(dep)

    if err != nil {
      panic(err)
    }

    for _, some_dep := range result.Values(){
      allDeps.Add(some_dep)
    }
  }

  return allDeps, nil
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
			relPath, err := filepath.Rel(gitRoot, path)
			if err != nil {
				panic(err)
			}
			sols = append(sols, &Sol{
				name: info.Name(),
				path: relPath,
				deps: make([]*Sol, 0),
			})
		}
		return nil
	}
	filepath.Walk(path.Join(gitRoot, "src"), checker)
	filepath.Walk(path.Join(gitRoot, "contracts"), checker)
	filepath.Walk(path.Join(gitRoot, "libs"), checker)

	// check imports and sort
	err := s.ResolveImports(gitRoot, sols)
	if err != nil {
		return nil, err
	}

	return &Extractor{
		dep:  dep,
		sols: sols,
	}, nil
}

func (s *Storage) ResolveImports(root string, sols []*Sol) error {

	reg, err := regexp.Compile(ImportRegex)
	if err != nil {
		return err
	}

	for _, sol := range sols {
		file, err := os.Open(path.Join(root, sol.path))
		if err != nil {
			return err
		}
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

				solDir := path.Dir(path.Join(root, sol.path))

				depPath := path.Join(solDir, match)
				for ii, searchSol := range sols {
					toCheck := path.Join(root, sols[ii].path)
					if toCheck == depPath {
						sol.deps = append(sol.deps, searchSol)
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		file.Close()
	}

	return nil
}

func (s *Storage) ResolveDepPath(dep *Dependency) string {
	return path.Join(s.root, dep.path)
}

func (s *Storage) Commit(e *Extractor, set *hashset.Set) {
	working, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	p := path.Join(working, "src", "lib", e.dep.path) // extract "lib" into a config object or read dapp file
	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		panic(err)
	}

	inclDeps := hashset.New()
	for _, inter := range set.Values() {
		sol := inter.(*Sol)
		inclDeps.Add(sol)
    deps, err := GetDependencies(sol)
    if err != nil {
      panic(err)
    }
		for _, dep := range deps.Values() {
			inclDeps.Add(dep)
		}
	}

	for _, toStore := range inclDeps.Values() {
		sol := toStore.(*Sol)
		s.WriteTo(p, e, sol)
	}
}

func (s *Storage) WriteTo(libPath string, e *Extractor, sol *Sol) {
	fromPath := path.Join(s.ResolveDepPath(e.dep), sol.path)
	fromFile, err := os.Open(fromPath)
	if err != nil {
		panic(err)
	}
	defer fromFile.Close()
	err = os.MkdirAll(path.Dir(path.Join(libPath, sol.path)), os.ModePerm)
	if err != nil {
		panic(err)
	}
	toFile, err := os.Create(path.Join(libPath, sol.path))
	if err != nil {
		panic(err)
	}
	defer toFile.Close()
	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		panic(err)
	}
}
