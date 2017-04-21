package modules

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"
)

type TemplatedBlock struct {
	Template string `yaml:"template"`

	tpl *template.Template
	buf bytes.Buffer
}

func (t *TemplatedBlock) ParseTemplate() {
	tplName := fmt.Sprintf("%p", t)
	t.tpl = template.Must(
		template.
			New(tplName).
			Funcs(tplFuncs).
			Parse(t.Template),
	)
}

func (t *TemplatedBlock) Format(data interface{}) string {
	t.buf.Reset()
	if err := t.tpl.Execute(&t.buf, data); err != nil {
		return err.Error()
	}
	return t.buf.String()
}

var tplFuncs = template.FuncMap{
	"float": toFloat,
	"sum": func(x0 interface{}, xs ...interface{}) (r float64, err error) {
		if r, err = toFloat(x0); err != nil {
			return
		}
		var xn float64
		for _, x := range xs {
			if xn, err = toFloat(x); err != nil {
				return
			}
			r += xn
		}
		return
	},
	"diff": func(x0 interface{}, xs ...interface{}) (r float64, err error) {
		if r, err = toFloat(x0); err != nil {
			return
		}
		var xn float64
		for _, x := range xs {
			if xn, err = toFloat(x); err != nil {
				return
			}
			r -= xn
		}
		return
	},
	"percent": func(v0 interface{}, v1 interface{}) (r float64, err error) {
		var x0, x1 float64
		if x0, err = toFloat(v0); err == nil {
			if x1, err = toFloat(v1); err == nil {
				r = 100.0 * x0 / x1
			}
		}
		return
	},
	"unit": func(u Unit, v interface{}) (r float64, err error) {
		var x float64
		if x, err = toFloat(v); err == nil {
			r = u.Apply(x)
		}
		return
	},
}

func toFloat(x interface{}) (float64, error) {
	switch v := x.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case int:
		return float64(v), nil
	default:
		return 0.0, fmt.Errorf("toFloat: cannot convert %s to float64", reflect.TypeOf(x).Name())
	}
}
