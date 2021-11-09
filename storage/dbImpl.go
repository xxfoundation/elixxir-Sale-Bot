package storage

import "gorm.io/gorm/clause"

func (db *DatabaseImpl) UpsertMember(email, id string) error {
	m := SaleMem{Email: email, ID: id}
	return db.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(m).Error
}
