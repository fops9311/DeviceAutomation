package entangled

func Entangle() (in func(interface{}), out func() interface{}) {
	var data interface{}
	in = func(d interface{}) {
		data = d
	}
	out = func() interface{} {
		return data
	}
	return in, out
}
