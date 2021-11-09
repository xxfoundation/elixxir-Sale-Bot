package storage

func (m *MapImpl) UpsertMember(email, id string) error {
	m.salemems[id] = email
	return nil
}
