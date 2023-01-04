package inflate

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNilTarget          = errors.New("nil target")
	ErrNotPointerReceiver = errors.New("not a pointer receiver")
)

// Reactor is an object configured with values / value providers
type Reactor interface {
	Get(interface{}) error
	Put(...interface{}) error
}

// New constructs new reactor
func New() Reactor {
	return &reactor{
		providers: make(map[reflect.Type]func() (interface{}, error)),
		values:    make(map[reflect.Type]interface{}),
	}
}

// NewWith constructs new reactor with given value/value providers
func NewWith(f ...interface{}) (Reactor, error) {
	reactor := New()
	if err := reactor.Put(f...); err != nil {
		return nil, err
	}

	return reactor, nil
}

// reactor is an object configured with value providers
type reactor struct {
	providers map[reflect.Type]func() (interface{}, error)
	values    map[reflect.Type]interface{}
}

// Get fills target (pointer to value) with value from reactor.
// If value is not calculated before it will be produced by
// corresponding provider
func (r *reactor) Get(target interface{}) error {
	if target == nil {
		return ErrNilTarget
	}

	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return ErrNotPointerReceiver
	}
	targetElemType := targetType.Elem()

	var toSave reflect.Value
	if v, ok := r.values[targetElemType]; ok {
		toSave = reflect.ValueOf(v)
	} else if p, ok := r.providers[targetElemType]; ok {
		if v, err := p(); err != nil {
			return err
		} else {
			r.values[targetElemType] = v
			toSave = reflect.ValueOf(v)
		}
	} else {
		return fmt.Errorf("provider for %s not registered", targetElemType.String())
	}

	reflect.ValueOf(target).Elem().Set(toSave)
	return nil
}

// Put places value or value provider function into reactor
func (r *reactor) Put(funcs ...interface{}) error {
	for i, f := range funcs {
		typeOf := reflect.TypeOf(f)
		if typeOf.Kind() != reflect.Func {
			r.values[typeOf] = f
			continue
		}
		numOut := typeOf.NumOut()
		if numOut == 0 {
			return fmt.Errorf("%d-th argument had no return arguments", i)
		} else if numOut > 2 {
			return fmt.Errorf("%d-th argument had %d return arguments", i, numOut)
		}

		funcValue := reflect.ValueOf(f)
		registeringType := typeOf.Out(0)

		var inReaders []func() (reflect.Value, error)
		numIn := typeOf.NumIn()
		for j := 0; j < numIn; j++ {
			argType := typeOf.In(j)
			inReaders = append(inReaders, func() (reflect.Value, error) {
				value := reflect.New(argType)
				err := r.Get(value.Interface())
				if err != nil {
					return reflect.Value{}, err
				}
				return value.Elem(), nil
			})
		}

		if numOut == 1 {
			if registeringType.Kind() == reflect.Ptr {
				r.providers[registeringType.Elem()] = func() (interface{}, error) {
					var in []reflect.Value
					for _, reader := range inReaders {
						v, err := reader()
						if err != nil {
							return nil, err
						}
						in = append(in, v)
					}
					values := funcValue.Call(in)
					return values[0].Elem().Interface(), nil
				}
			} else {
				r.providers[registeringType] = func() (interface{}, error) {
					var in []reflect.Value
					for _, reader := range inReaders {
						v, err := reader()
						if err != nil {
							return nil, err
						}
						in = append(in, v)
					}
					values := funcValue.Call(in)
					return values[0].Interface(), nil
				}
			}
		} else {
			if registeringType.Kind() == reflect.Ptr {
				r.providers[registeringType.Elem()] = func() (interface{}, error) {
					var in []reflect.Value
					for _, reader := range inReaders {
						v, err := reader()
						if err != nil {
							return nil, err
						}
						in = append(in, v)
					}
					values := funcValue.Call(in)
					if values[1].IsNil() {
						return values[0].Elem().Interface(), nil
					}
					return nil, values[1].Interface().(error)
				}
			} else {
				r.providers[registeringType] = func() (interface{}, error) {
					var in []reflect.Value
					for _, reader := range inReaders {
						v, err := reader()
						if err != nil {
							return nil, err
						}
						in = append(in, v)
					}
					values := funcValue.Call(in)
					if values[1].IsNil() {
						return values[0].Interface(), nil
					}
					return nil, values[1].Interface().(error)
				}
			}
		}
	}
	return nil
}
