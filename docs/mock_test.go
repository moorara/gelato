package lookup

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type ServiceMocker struct {
	t            *testing.T
	spew         *spew.ConfigState
	expectations *ServiceExpectations
}

func MockService(t *testing.T) *ServiceMocker {
	return &ServiceMocker{
		t: t,
		spew: &spew.ConfigState{
			Indent:                  "  ",
			DisablePointerAddresses: true,
			DisableCapacities:       true,
			SortKeys:                true,
		},
		expectations: new(ServiceExpectations),
	}
}

func (m *ServiceMocker) Expect() *ServiceExpectations {
	return m.expectations
}

func (m *ServiceMocker) Assert() {
	// Repeat for all methods
	for _, e := range m.expectations.lookupExpectations {
		if e.recorded == nil {
			m.t.Errorf("\nExpected Lookup method be called with %s", m.spew.Sdump(e.inputs))
		}
	}
}

func (m *ServiceMocker) Impl() Service {
	return &ServiceImpl{
		t:            m.t,
		spew:         m.spew,
		expectations: m.expectations,
	}
}

type ServiceExpectations struct {
	lookupExpectations []*LookupExpectation
}

func (e *ServiceExpectations) Lookup() *LookupExpectation {
	expectation := new(LookupExpectation)
	e.lookupExpectations = append(e.lookupExpectations, expectation)
	return expectation
}

type LookupExpectation struct {
	inputs   *lookupInputs
	outputs  *lookupOutputs
	callback func(*Request) (*Response, error)
	recorded *lookupInputs
}

type lookupInputs struct {
	request *Request
}

type lookupOutputs struct {
	response *Response
	err      error
}

func (e *LookupExpectation) WithArgs(request *Request) *LookupExpectation {
	e.inputs = &lookupInputs{
		request: request,
	}
	return e
}

func (e *LookupExpectation) Return(response *Response, err error) *LookupExpectation {
	e.outputs = &lookupOutputs{
		response: response,
		err:      err,
	}
	return e
}

func (e *LookupExpectation) Call(callback func(*Request) (*Response, error)) *LookupExpectation {
	e.callback = callback
	return e
}

type ServiceImpl struct {
	t            *testing.T
	spew         *spew.ConfigState
	expectations *ServiceExpectations
}

func (i *ServiceImpl) Lookup(request *Request) (*Response, error) {
	inputs := &lookupInputs{
		request: request,
	}

	for _, e := range i.expectations.lookupExpectations {
		if e.inputs == nil || reflect.DeepEqual(e.inputs, inputs) {
			e.recorded = inputs
			if e.callback != nil {
				return e.callback(e.inputs.request)
			}
			return e.outputs.response, e.outputs.err
		}
	}

	i.t.Errorf("Expectation missing: Lookup method called with %s", i.spew.Sdump(inputs))

	return nil, nil
}
