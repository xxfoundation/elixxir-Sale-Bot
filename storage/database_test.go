package storage

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	s, err := NewStorage(Params{
		Username: "",
		Password: "",
		DBName:   "",
		Address:  "0.0.0.0",
		Port:     "5432",
	})
	if err != nil {
		t.Errorf("Failed to initialize storage: %+v", err)
	}
	err = s.UpsertMember("test")
	if err != nil {
		t.Errorf("Failed to insert test member: %+v", err)
	}
}
