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

package client

import "fmt"

// FakeClient implements GraviteeClient with injectable functions for testing.
type FakeClient struct {
	GetFunc    func(path string) ([]byte, error)
	PostFunc   func(path string, body any) ([]byte, error)
	PutFunc    func(path string, body any) ([]byte, error)
	PatchFunc  func(path string, body any) ([]byte, error)
	DeleteFunc func(path string) error
}

func (f *FakeClient) Get(path string) ([]byte, error) {
	if f.GetFunc == nil {
		return nil, fmt.Errorf("unexpected Get call: %s", path)
	}

	return f.GetFunc(path)
}

func (f *FakeClient) Post(path string, body any) ([]byte, error) {
	if f.PostFunc == nil {
		return nil, fmt.Errorf("unexpected Post call: %s", path)
	}

	return f.PostFunc(path, body)
}

func (f *FakeClient) Put(path string, body any) ([]byte, error) {
	if f.PutFunc == nil {
		return nil, fmt.Errorf("unexpected Put call: %s", path)
	}

	return f.PutFunc(path, body)
}

func (f *FakeClient) Patch(path string, body any) ([]byte, error) {
	if f.PatchFunc == nil {
		return nil, fmt.Errorf("unexpected Patch call: %s", path)
	}

	return f.PatchFunc(path, body)
}

func (f *FakeClient) Delete(path string) error {
	if f.DeleteFunc == nil {
		return fmt.Errorf("unexpected Delete call: %s", path)
	}

	return f.DeleteFunc(path)
}
