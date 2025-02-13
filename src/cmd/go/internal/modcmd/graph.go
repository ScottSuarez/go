// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go mod graph

package modcmd

import (
	"bufio"
	"context"
	"os"

	"cmd/go/internal/base"
	"cmd/go/internal/modload"

	"golang.org/x/mod/module"
)

var cmdGraph = &base.Command{
	UsageLine: "go mod graph [-go=version]",
	Short:     "print module requirement graph",
	Long: `
Graph prints the module requirement graph (with replacements applied)
in text form. Each line in the output has two space-separated fields: a module
and one of its requirements. Each module is identified as a string of the form
path@version, except for the main module, which has no @version suffix.

The -go flag causes graph to report the module graph as loaded by the
given Go version, instead of the version indicated by the 'go' directive
in the go.mod file.

See https://golang.org/ref/mod#go-mod-graph for more about 'go mod graph'.
	`,
	Run: runGraph,
}

var (
	graphGo goVersionFlag
)

func init() {
	cmdGraph.Flag.Var(&graphGo, "go", "")
	base.AddModCommonFlags(&cmdGraph.Flag)
	base.AddWorkfileFlag(&cmdGraph.Flag)
}

func runGraph(ctx context.Context, cmd *base.Command, args []string) {
	modload.InitWorkfile()

	if len(args) > 0 {
		base.Fatalf("go mod graph: graph takes no arguments")
	}
	modload.ForceUseModules = true
	modload.RootMode = modload.NeedRoot
	mg := modload.LoadModGraph(ctx, graphGo.String())

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	format := func(m module.Version) {
		w.WriteString(m.Path)
		if m.Version != "" {
			w.WriteString("@")
			w.WriteString(m.Version)
		}
	}

	mg.WalkBreadthFirst(func(m module.Version) {
		reqs, _ := mg.RequiredBy(m)
		for _, r := range reqs {
			format(m)
			w.WriteByte(' ')
			format(r)
			w.WriteByte('\n')
		}
	})
}
