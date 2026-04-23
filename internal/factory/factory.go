// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package factory

import (
	"io"
	"os"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/apim"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
)

// IOStreams abstracts standard I/O for testability.
type IOStreams struct {
	Out io.Writer
	Err io.Writer
	In  io.Reader
}

// DefaultIOStreams returns IOStreams connected to os stdin/stdout/stderr.
func DefaultIOStreams() IOStreams {
	return IOStreams{
		Out: os.Stdout,
		Err: os.Stderr,
		In:  os.Stdin,
	}
}

// Factory is the central dependency injection container passed to all commands.
type Factory struct {
	Config            *config.Config
	Resolved          *config.ResolvedContext
	Overrides         config.Overrides
	Client            client.GraviteeClient
	apimService       apim.Service
	amService         am.Service
	IOStreams         IOStreams
	ContextResolveErr error
	ConfigPath        string
	Product           string
	OutputFormat      string
	NoHeaders         bool
	Quiet             bool
	Debug             bool
}

// APIM returns the APIM service, creating it lazily from Client + Resolved if needed.
func (f *Factory) APIM() apim.Service {
	if f.apimService != nil {
		return f.apimService
	}

	if f.Client != nil && f.Resolved != nil {
		f.apimService = apim.NewService(f.Client, f.Resolved)
	}

	return f.apimService
}

// SetAPIMService sets the APIM service (used in tests).
func (f *Factory) SetAPIMService(s apim.Service) {
	f.apimService = s
}

// AM returns the AM service, creating it lazily from Client + Resolved if needed.
func (f *Factory) AM() am.Service {
	if f.amService != nil {
		return f.amService
	}

	if f.Client != nil && f.Resolved != nil {
		f.amService = am.NewService(f.Client, f.Resolved)
	}

	return f.amService
}

// SetAMService sets the AM service (used in tests).
func (f *Factory) SetAMService(s am.Service) {
	f.amService = s
}
