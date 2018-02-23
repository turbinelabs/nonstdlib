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

// Package tabwriter provides a set of sane defaults for converting tab separated
// values into a pretty column formatted output.
package tabwriter

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

type Config struct {
	// MinWidth is the minimal cell width including any padding
	MinWidth int

	// TabWidth is the width of a tab characters (equivalent number of spaces)
	TabWidth int

	// Padding is added to a cell before computing its width
	Padding int

	// PadChar is an ASCII char used for padding; see more detail in
	// https://golang.org/pkg/text/tabwriter/#NewWriter
	PadChar byte

	// Flags specify formatting controls; see more detail in
	// https://golang.org/pkg/text/tabwriter/#pkg-constants
	Flags uint
}

// DefaultConfig is the tabwriter config that will be used if none is specified.
var DefaultConfig = Config{
	MinWidth: 5,
	TabWidth: 0,
	Padding:  2,
	PadChar:  ' ',
	Flags:    0,
}

type TbnTabWriter interface {
	// Format returns the input with the content aligned by tabs.
	Format(string) string

	// FormatWithHeader returns the input with the content aligned by tabs and
	// includes a header even if no data is provided.
	FormatWithHeader(string, string) string
}

// New constructs a new TbnTabWriter with the provided output config.
func New(cfg Config) TbnTabWriter {
	return tbnTabWriter{cfg}
}

var defaultWriter = New(DefaultConfig)

type tbnTabWriter struct {
	cfg Config
}

func (t tbnTabWriter) Format(in string) string {
	return t.FormatWithHeader("", in)
}

func (t tbnTabWriter) FormatWithHeader(hdr, in string) string {
	c := t.cfg

	w := &tabwriter.Writer{}
	buf := bytes.NewBuffer(nil)
	w.Init(buf, c.MinWidth, c.TabWidth, c.Padding, c.PadChar, c.Flags)

	if hdr != "" {
		fmt.Fprintln(w, hdr)
	}

	for _, l := range strings.Split(in, "\n") {
		fmt.Fprintln(w, l)
	}
	w.Flush()

	return buf.String()
}

// Format returns the input with the content aligned by tabs. If a config is
// specified it will be used otherwise DefaultConfig will be used.
func Format(in string) string {
	return defaultWriter.Format(in)
}

// FormatWithHeader returns the input with the content aligned by tabs and
// includes a header even if no data is provided. If a config is specified it
// will be used otherwise DefaultConfig will be used.
func FormatWithHeader(hdr, in string) string {
	return defaultWriter.FormatWithHeader(hdr, in)
}
