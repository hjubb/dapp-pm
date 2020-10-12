package main

import (
	"reflect"
	"testing"
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