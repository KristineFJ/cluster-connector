/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"kubeops.dev/cluster-connector/pkg/cmds"

	"github.com/spf13/cobra/doc"
	"gomodules.xyz/runtime"
	"k8s.io/klog/v2"
)

const (
	version = "0.2.0"
)

var (
	tplFrontMatter = template.Must(template.New("index").Parse(`---
title: Reference | Vault Operator
description: Vault Operator CLI Reference
menu:
  docs_{{ "{{ .version }}" }}:
    identifier: reference-operator
    name: Vault Operator
    weight: 10
    parent: reference
menu_name: docs_{{ "{{ .version }}" }}
---
`))

	_ = template.Must(tplFrontMatter.New("cmd").Parse(`---
title: {{ .Name }}
menu:
  docs_{{ "{{ .version }}" }}:
    identifier: {{ .ID }}
    name: {{ .Name }}
    parent: reference-operator
{{- if .RootCmd }}
    weight: 0
{{ end }}
menu_name: docs_{{ "{{ .version }}" }}
section_menu_id: reference
{{- if .RootCmd }}
url: /docs/{{ "{{ .version }}" }}/reference/operator/
aliases:
- /docs/{{ "{{ .version }}" }}/reference/operator/{{ .ID }}/
{{- end }}
---
`))
)

// ref: https://github.com/spf13/cobra/blob/master/doc/md_docs.md
func main() {
	rootCmd := cmds.NewRootCmd()
	dir := runtime.GOPath() + "/src/kubeops.dev/cluster-connector/docs/reference/operator"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		klog.Fatalln(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		klog.Fatalln(err)
	}

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		data := struct {
			ID      string
			Name    string
			Version string
			RootCmd bool
		}{
			strings.Replace(base, "_", "-", -1),
			strings.Title(strings.Replace(base, "_", " ", -1)),
			version,
			!strings.ContainsRune(base, '_'),
		}
		var buf bytes.Buffer
		if err := tplFrontMatter.ExecuteTemplate(&buf, "cmd", data); err != nil {
			klog.Fatalln(err)
		}
		return buf.String()
	}

	linkHandler := func(name string) string {
		return "/docs/reference/operator/" + name
	}
	err = doc.GenMarkdownTreeCustom(rootCmd, dir, filePrepender, linkHandler)
	if err != nil {
		klog.Fatalln(err)
	}

	index := filepath.Join(dir, "_index.md")
	f, err := os.OpenFile(index, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		klog.Fatalln(err)
	}
	err = tplFrontMatter.ExecuteTemplate(f, "index", struct{ Version string }{version})
	if err != nil {
		klog.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		klog.Fatalln(err)
	}
}
