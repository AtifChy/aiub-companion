package database

import "database/sql"

func StringOrNull(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func BoolOrNull(b *bool) sql.NullBool {
	return sql.NullBool{Bool: b != nil && *b, Valid: b != nil}
}

func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
