package models

func (m *UserAdd) ToUser(e User) User {
	e.Username = m.Username
	e.Password = m.Password
	e.Role = m.Role

	return e
}
