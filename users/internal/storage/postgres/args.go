package postgres

import "strings"

const (
	userTable  = "users"
	tokenTable = "tokens"
)

const (
	returning = "RETURNING "
	separator = ","
)

const (
	fieldID     = "id"
	fieldUserID = "user_id"

	fieldUsername       = "username"
	fieldHashedPassword = "hashed_password"
	fieldIsAdmin        = "is_admin"

	fieldTokenValue     = "value"
	fieldCreatedAt      = "created_at"
	fieldExpirationDate = "expiration_date"
)

var (
	userFields = []string{
		fieldID, fieldUsername, fieldHashedPassword, fieldIsAdmin,
	}
	tokenFields = []string{
		fieldID, fieldUserID, fieldTokenValue, fieldCreatedAt, fieldExpirationDate,
	}
	returningUser  = returning + strings.Join(userFields, separator)
	returningToken = returning + strings.Join(tokenFields, separator)
)
