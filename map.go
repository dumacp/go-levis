package levis

type conf struct {
	buttons_start int
	buttons_end   int
	read_start    int
	read_end      int
	write_start   int
	write_end     int
}

type Conf interface {
	SetReadMem(start, end int) Conf
	SetWriteMem(start, end int) Conf
	SetButtonMem(start, end int) Conf
}

func (m *conf) SetReadMem(start, end int) Conf {

	m.read_start = start
	m.read_end = end
	return m
}

func (m *conf) SetWriteMem(start, end int) Conf {

	m.write_start = start
	m.write_end = end
	return m
}

func (m *conf) SetButtonMem(start, end int) Conf {

	m.buttons_start = start
	m.buttons_end = end
	return m
}
