package keystruct

//as an ordered data struct, the key should be able to compare
//currently supoort 3 kinds of data type
type KeyStruct interface {
	CompareBiggerThan(other KeyStruct) bool
	KeyString() string
	KeyInt32() int
	KeyInt64() int64
}

//Generally you should define an struct contains DefaultKey,
//Then, you should implement the key you want
//Basically string is enough for most of the situation
type DefaultKey struct{}

func (k DefaultKey) CompareBiggerThan(other KeyStruct) bool {
	return false
}
func (k DefaultKey) KeyString() string {
	return ""
}
func (k DefaultKey) KeyInt32() int {
	return 0
}
func (k DefaultKey) KeyInt64() int64 {
	return 0
}
