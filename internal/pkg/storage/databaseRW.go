package storage

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

func ReadFromDB(s Storage, db *gorm.DB) error {
	var data JsonStorage
	data.Data = make(map[string]Value)

	rows, err := db.Raw(`
        SELECT d.key, d.value_type, d.expires_at,
               s.scalar_value_type, s.scalar_value_int, s.scalar_value_string,
               sl.value AS slice_value,
               di.field AS dict_key, di.scalar_value_type AS dict_value_type, di.scalar_value_int AS dict_value_int, di.scalar_value_string AS dict_value_string
        FROM data d
        LEFT JOIN scalars s ON d.id = s.data_id
        LEFT JOIN slices sl ON d.id = sl.data_id
        LEFT JOIN dicts di ON d.id = di.data_id`).Rows()

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			dataKey           string
			valueType         string
			expiresAt         int64
			scalarValueType   sql.NullString
			scalarValueInt    sql.NullInt64
			scalarValueString sql.NullString
			sliceValue        sql.NullInt64
			dictKey           sql.NullString
			dictValueType     sql.NullString
			dictValueInt      sql.NullInt64
			dictValueString   sql.NullString
		)

		if err := rows.Scan(&dataKey, &valueType, &expiresAt,
			&scalarValueType, &scalarValueInt, &scalarValueString,
			&sliceValue, &dictKey, &dictValueType, &dictValueInt, &dictValueString); err != nil {
			return err
		}

		// Add Value
		value := Value{
			ValueType: Kind(valueType),
			ExpiresAt: expiresAt,
			Scalar: ScalarValue{
				ScalarValueType:   ScalarKind(scalarValueType.String),
				ScalarValueInt:    scalarValueInt.Int64,
				ScalarValueString: scalarValueString.String,
			},
			Slice: []int{},
			Dict:  make(map[string]ScalarValue),
		}

		// Add to Slice
		if sliceValue.Valid {
			value.Slice = append(data.Data[dataKey].Slice, int(sliceValue.Int64))
		}

		// Add to Dict
		if dictKey.Valid {
			value.Dict[dictKey.String] = ScalarValue{
				ScalarValueType:   ScalarKind(dictValueType.String),
				ScalarValueInt:    dictValueInt.Int64,
				ScalarValueString: dictValueString.String,
			}
		}
		data.Data[dataKey] = value
	}

	s.LoadData(data)

	defer s.WriteLog("STATE was readed from DB")

	return nil
}

func WriteToDB(s Storage, db *gorm.DB) error {
	jsonStorage := s.ExportData()

	// Delete old data
	newKeys := make([]string, 0)
	for key := range jsonStorage.Data {
		newKeys = append(newKeys, key)
	}
	oldKeys := make([]string, 0)
	err := db.Raw("SELECT key FROM data").Scan(&oldKeys).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	for _, value := range oldKeys {
		if !contains(newKeys, value) {
			result := db.Exec("DELETE FROM data WHERE key = ?", value)
			if result.Error != nil {
				return err
			}
		}
	}

	// Update or add data
	for key, value := range jsonStorage.Data {
		// Work with table Data
		var dataID int
		err := db.Raw("SELECT id FROM data WHERE key = $1", key).Scan(&dataID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if dataID > 0 {
			// Update
			result := db.Exec("UPDATE data SET value_type = $1, expires_at = $2 WHERE id = $3", value.ValueType, value.ExpiresAt, dataID)
			if result.Error != nil {
				return err
			}
		} else {
			// Add new
			err = db.Raw("INSERT INTO data (key, value_type, expires_at) VALUES ($1, $2, $3) RETURNING id", key, value.ValueType, value.ExpiresAt).Scan(&dataID).Error
			if err != nil {
				return err
			}
		}

		// Work with table Scalars
		scalar := value.Scalar
		if scalar.ScalarValueType != "" {
			var scalarID int
			err = db.Raw("SELECT id FROM scalars WHERE data_id = $1 AND scalar_value_type = $2", dataID, scalar.ScalarValueType).Scan(&scalarID).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if scalarID > 0 {
				// Update
				result := db.Exec("UPDATE scalars SET scalar_value_int = $1, scalar_value_string = $2 WHERE id = $3",
					scalar.ScalarValueInt, scalar.ScalarValueString, scalarID)
				if result.Error != nil {
					return err
				}
			} else {
				// Add new
				result := db.Exec("INSERT INTO scalars (data_id, scalar_value_type, scalar_value_int, scalar_value_string) VALUES ($1, $2, $3, $4)",
					dataID, scalar.ScalarValueType, scalar.ScalarValueInt, scalar.ScalarValueString)
				if result.Error != nil {
					return err
				}
			}
		}

		// Work with table Slices
		for _, sliceValue := range value.Slice {
			var sliceID int
			err = db.Raw("SELECT id FROM slices WHERE data_id = $1 AND value = $2", dataID, sliceValue).Scan(&sliceID).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if sliceID > 0 {
				// Update
				result := db.Exec("UPDATE slices SET value = $1 WHERE id = $2",
					sliceValue, sliceID)
				if result.Error != nil {
					return err
				}
			} else {
				// Add new
				result := db.Exec("INSERT INTO slices (data_id, value) VALUES ($1, $2)",
					dataID, sliceValue)
				if result.Error != nil {
					return err
				}
			}
		}

		// Work with table Dict
		for dictKey, dictValue := range value.Dict {
			var dictID int
			err = db.Raw("SELECT id FROM dicts WHERE data_id = $1 AND field = $2", dataID, dictKey).Scan(&dictID).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if dictID > 0 {
				// Update
				result := db.Exec("UPDATE dicts SET field = $1, scalar_value_type = $2, scalar_value_int = $3, scalar_value_string = $4 WHERE id = $5",
					dictKey, dictValue.ScalarValueType, dictValue.ScalarValueInt, dictValue.ScalarValueString, dictID)
				if result.Error != nil {
					return err
				}
			} else {
				// Add new
				result := db.Exec("INSERT INTO dicts (data_id, field, scalar_value_type, scalar_value_int, scalar_value_string) VALUES ($1, $2, $3, $4, $5)",
					dataID, dictKey, dictValue.ScalarValueType, dictValue.ScalarValueInt, dictValue.ScalarValueString)
				if result.Error != nil {
					return err
				}
			}
		}

	}

	defer s.WriteLog("STATE was writed to DB")

	return nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Old func - Only INSERT
// func WriteToDB(s Storage, db *gorm.DB) error {
// 	jsonStorage := s.ExportData()

// 	for key, value := range jsonStorage.Data {
// 		// Save to Data
// 		var dataID int
// 		db.Raw("INSERT INTO data (key, value_type, expires_at) VALUES ($1, $2, $3) RETURNING id", key, value.ValueType, value.ExpiresAt).Scan(&dataID)
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		// Save to Scalar
// 		scalar := value.Scalar
// 		if scalar.ScalarValueType != "" {
// 			db.Exec("INSERT INTO scalars (data_id, scalar_value_type, scalar_value_int, scalar_value_string) VALUES ($1, $2, $3, $4)",
// 				dataID, scalar.ScalarValueType, scalar.ScalarValueInt, scalar.ScalarValueString)
// 		}

// 		// Save to Slice
// 		for _, sliceValue := range value.Slice {
// 			db.Exec("INSERT INTO slices (data_id, value) VALUES ($1, $2)", dataID, sliceValue)
// 		}

// 		// Save to Dict
// 		for dictKey, dictValue := range value.Dict {
// 			db.Exec("INSERT INTO dicts (data_id, field, scalar_value_type, scalar_value_int, scalar_value_string) VALUES ($1, $2, $3, $4, $5)",
// 				dataID, dictKey, dictValue.ScalarValueType, dictValue.ScalarValueInt, dictValue.ScalarValueString)
// 		}
// 	}

// 	defer s.WriteLog("STATE was writed to DB")

// 	return nil
// }
