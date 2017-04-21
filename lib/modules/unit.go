package modules

import "fmt"

// The Unit type
type Unit float64

const (
	Byte Unit = 1.0
	Kilo Unit = 1 << 10
	Mega Unit = 1 << 20
	Giga Unit = 1 << 30
	Tera Unit = 1 << 40
)

var unitToString = map[Unit]string{Byte: "B", Kilo: "K", Mega: "M", Giga: "G", Tera: "T"}
var stringToUnit = map[string]Unit{"B": Byte, "K": Kilo, "M": Mega, "G": Giga, "T": Tera}

func (u *Unit) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var name string
	if err = unmarshal(&name); err == nil {
		if u1, found := stringToUnit[name]; found {
			*u = u1
		} else {
			err = fmt.Errorf("Unknown unit %q", name)
		}
	}
	return
}

func (u Unit) String() string {
	return unitToString[u]
}

func (u Unit) Apply(v float64) float64 {
	return v / float64(u)
}
