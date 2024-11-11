package ldap

type Stack struct {
	data []byte
}

func (s *Stack) Len() int { return len(s.data) }
func (s *Stack) Pop() byte {
	if len(s.data) == 0 {
		return 0 // 栈为空
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}
func (s *Stack) Push(v byte) {
	s.data = append(s.data, v)
}
