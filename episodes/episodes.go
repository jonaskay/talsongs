package episodes

type Episodes []string

func (e Episodes) Unique() Episodes {
	var r Episodes
	m := make(map[string]bool)

	for _, str := range e {
		if _, ok := m[str]; !ok {
			m[str] = true
			r = append(r, str)
		}
	}

	return r
}
