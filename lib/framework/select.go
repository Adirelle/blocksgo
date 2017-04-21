package framework

import "reflect"

type handlerFunc func(value interface{})

type Select struct {
	cases    []reflect.SelectCase
	handlers []handlerFunc
}

func (s *Select) OnRecv(ch interface{}, h handlerFunc) (i int) {
	return s.addCase(
		reflect.SelectCase{
			Chan: reflect.ValueOf(ch),
			Dir:  reflect.SelectRecv,
		},
		h,
	)
}

func (s *Select) addCase(c reflect.SelectCase, h handlerFunc) (i int) {
	i = len(s.cases)
	s.cases = append(s.cases, c)
	s.handlers = append(s.handlers, h)
	return
}

func (s *Select) Run() (ok bool) {
	var (
		i int
		v reflect.Value
	)
	if i, v, ok = reflect.Select(s.cases); ok {
		s.handlers[i](v.Interface())
	}
	return
}
