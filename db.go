package chord

type Db map[string]string

func newDb() Db {
	var s Db = make(map[string]string)

	return s
}
