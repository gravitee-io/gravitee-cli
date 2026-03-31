package client

import "fmt"

// FakeClient implements GraviteeClient with injectable functions for testing.
type FakeClient struct {
	GetFunc    func(path string) ([]byte, error)
	PostFunc   func(path string, body interface{}) ([]byte, error)
	PutFunc    func(path string, body interface{}) ([]byte, error)
	DeleteFunc func(path string) error
}

func (f *FakeClient) Get(path string) ([]byte, error) {
	if f.GetFunc == nil {
		return nil, fmt.Errorf("unexpected Get call: %s", path)
	}

	return f.GetFunc(path)
}

func (f *FakeClient) Post(path string, body interface{}) ([]byte, error) {
	if f.PostFunc == nil {
		return nil, fmt.Errorf("unexpected Post call: %s", path)
	}

	return f.PostFunc(path, body)
}

func (f *FakeClient) Put(path string, body interface{}) ([]byte, error) {
	if f.PutFunc == nil {
		return nil, fmt.Errorf("unexpected Put call: %s", path)
	}

	return f.PutFunc(path, body)
}

func (f *FakeClient) Delete(path string) error {
	if f.DeleteFunc == nil {
		return fmt.Errorf("unexpected Delete call: %s", path)
	}

	return f.DeleteFunc(path)
}
