/*
Copyright 2018 The Kubernetes Authors.

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

package templates

import (
	"fmt"
	"path/filepath"
	"time"

	"sigs.k8s.io/kubebuilder/pkg/model/file"
)

var _ file.Template = &Boilerplate{}

// Boilerplate scaffolds a boilerplate header file.
type Boilerplate struct {
	file.TemplateMixin
	file.BoilerplateMixin

	// License is the License type to write
	License string

	// Licenses maps License types to their actual string
	Licenses map[string]string

	// Owner is the copyright owner - e.g. "The Kubernetes Authors"
	Owner string

	// Year is the copyright year
	Year string
}

// Validate implements file.RequiresValidation
func (f Boilerplate) Validate() error {
	if f.License == "" {
		// A default license will be set later
	} else if _, found := knownLicenses[f.License]; found {
		// One of the know licenses
	} else if _, found := f.Licenses[f.License]; found {
		// A map containing the requested license was also provided
	} else {
		return fmt.Errorf("unknown specified license %s", f.License)
	}

	return nil
}

// SetTemplateDefaults implements input.Template
func (f *Boilerplate) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("hack", "boilerplate.go.txt")
	}

	if f.License == "" {
		f.License = "apache2"
	}

	if f.Licenses == nil {
		f.Licenses = make(map[string]string, len(knownLicenses))
	}

	for key, value := range knownLicenses {
		if _, hasLicense := f.Licenses[key]; !hasLicense {
			f.Licenses[key] = value
		}
	}

	if f.Year == "" {
		f.Year = fmt.Sprintf("%v", time.Now().Year())
	}

	// Boilerplate given
	if len(f.Boilerplate) > 0 {
		f.TemplateBody = f.Boilerplate
		return nil
	}

	f.TemplateBody = boilerplateTemplate

	return nil
}

const boilerplateTemplate = `/*
{{ if .Owner -}}
Copyright {{ .Year }} {{ .Owner }}.
{{- else -}}
Copyright {{ .Year }}.
{{- end }}
{{ index .Licenses .License }}*/`

var knownLicenses = map[string]string{
	"apache2": apache2,
	"none":    "",
}

const apache2 = `
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`
