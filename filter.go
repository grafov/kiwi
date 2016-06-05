package kiwi

type filter interface {
	Check(string, string) bool
}

type keyFilter struct {
	Key string
}

func (k *keyFilter) Check(key, val string) bool {
	return k.Key == key
}

type valsFilter struct {
	Key  string
	Vals []string
}

func (k *valsFilter) Check(key, val string) bool {
	if key != k.Key {
		return false
	}
	for _, v := range k.Vals {
		if v == val {
			return true
		}
	}
	return false
}
