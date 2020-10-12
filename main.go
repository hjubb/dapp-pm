package main

import (
	"fmt"
)

func main() {

	storage, err := NewStorageEngine()
	if err != nil {
		panic(err)
	}

	dep, err := storage.GetDependency("openzeppelin/openzeppelin-contracts@v3.2.0")
	if err != nil {
		panic(err)
	}

	extractor, err := storage.ExtractPaths(dep)

	if err != nil {
		panic(err)
	}

	list := extractor.Render()
}