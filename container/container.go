package container

import (
	"fmt"
	"io"

	"github.com/clipperhouse/typewriter"
)

func init() {
	err := typewriter.Register(NewContainerWriter())
	if err != nil {
		panic(err)
	}
}

type ContainerWriter struct {
	tagsByType map[string]typewriter.Tag // typewriter.Type is not comparable, key by .String()
}

func NewContainerWriter() ContainerWriter {
	return ContainerWriter{
		tagsByType: make(map[string]typewriter.Tag),
	}
}

func (c ContainerWriter) Name() string {
	return "container"
}

func (c ContainerWriter) Validate(t typewriter.Type) (bool, error) {
	tag, found, err := t.Tags.ByName("containers")
	if found && err == nil {
		c.tagsByType[t.String()] = tag
	}
	return found, err
}

func (c ContainerWriter) WriteHeader(w io.Writer, t typewriter.Type) {
	tag := c.tagsByType[t.String()] // validated above

	var list, ring, set bool

	for _, s := range tag.Items {
		if s == "List" {
			list = true
		}
		if s == "Ring" {
			ring = true
		}
		if s == "Set" {
			set = true
		}
	}

	if list {
		license := `// List is a modification of http://golang.org/pkg/container/list/
`
		w.Write([]byte(license))
	}

	if ring {
		license := `// Ring is a modification of http://golang.org/pkg/container/ring/
`
		w.Write([]byte(license))
	}

	if list || ring {
		license := `// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at http://golang.org/LICENSE

`
		w.Write([]byte(license))
	}

	if set {
		license := `// Set is a modification of https://github.com/deckarep/golang-set
// The MIT License (MIT)
// Copyright (c) 2013 Ralph Caraveo (deckarep@gmail.com)
`
		w.Write([]byte(license))
	}

	return
}

func (c ContainerWriter) Imports(t typewriter.Type) (result []string) {
	return result
}

func (c ContainerWriter) Write(w io.Writer, t typewriter.Type) {
	tag := c.tagsByType[t.String()] // validated above

	for _, s := range tag.Items {
		tmpl, err := containerTemplates.Get(s) // validate above to avoid err check here?
		if err != nil {
			continue
		}
		err = tmpl.Execute(w, t)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return
}
