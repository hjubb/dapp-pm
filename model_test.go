package main

import (
	"reflect"
	"testing"
	"github.com/emirpasic/gods/sets/hashset"
)

func TestNewDependency(t *testing.T) {
	type args struct {
		dependency string
	}
	tests := []struct {
		name    string
		args    args
		want    *Dependency
		wantErr bool
	}{
		{name: "test works", args: args{dependency: "openzeppelin/openzeppelin-contracts@3.0.0"}, want: &Dependency{
			name:    "openzeppelin/openzeppelin-contracts",
			version: "3.0.0",
			path:    "openzeppelin/openzeppelin-contracts@3.0.0",
		}, wantErr: false},
		{name: "test bad format", args: args{dependency: "openzeppelin/openzeppelin-contracts"}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDependency(tt.args.dependency)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDependency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDependency() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDependencies(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "test works", want: 8, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
      // level 2
      d := &Sol{name: "D", path: "./D", deps: make([]*Sol, 0)}
      e := &Sol{name: "E", path: "./E", deps: make([]*Sol, 0)}
      f := &Sol{name: "F", path: "./F", deps: make([]*Sol, 0)}
      g := &Sol{name: "G", path: "./G", deps: make([]*Sol, 0)}

      // level 1
      a := &Sol{name: "A", path: "./A", deps: []*Sol{d, e}}
      b := &Sol{name: "B", path: "./B", deps: []*Sol{f}}
      c := &Sol{name: "C", path: "./C", deps: []*Sol{g}}

      root := &Sol{
        name: "X", path: "./X", deps: []*Sol{a, b, c},
      }
      GetDependenciesOld := func(selected *Sol) (*hashset.Set){
        inclDeps := hashset.New()
        selections := hashset.New()
        selections.Add(selected)
        for _, inter := range selections.Values() {
          sol := inter.(*Sol)
          inclDeps.Add(sol)
          for _, dep := range sol.deps {
            inclDeps.Add(dep)
          }
        }
        return inclDeps
      }
      oldResults := len(GetDependenciesOld(root).Values())

      results, err := GetDependencies(root)
      uniqueDeps := len(results.Values())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
      if uniqueDeps != tt.want {
				t.Errorf("GetDependencies() got = %v, want %v", uniqueDeps, tt.want)
      }
      if uniqueDeps < oldResults {
				t.Errorf("GetDependencies() got = %v, OldWay got %v", uniqueDeps, oldResults)
      }
		})
	}
}
