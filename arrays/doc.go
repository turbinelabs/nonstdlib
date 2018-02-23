/*
Copyright 2018 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package arrays includes several sub-packages allowing type-safe execution
// of tasks commonly applied to slices.
package arrays

//go:generate codegen --output=equal.go --source=$GOFILE equal.template types[]=int,int64,float64,string
//go:generate codegen --output=equal_test.go --source=$GOFILE equal_test.template types[]=int,int64,float64,string
