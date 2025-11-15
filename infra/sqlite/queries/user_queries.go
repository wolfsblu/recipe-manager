package queries

import (
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectUserByID returns a query to fetch a user with their role
func SelectUserByID(id int32) SelectStatement {
	return SELECT(
		User.ID,
		User.Email,
		User.PasswordHash,
		User.IsConfirmed,
		User.RoleID,
		User.Locale,
		User.CreatedAt,
		Role.Name.AS("role.name"),
	).FROM(
		User.INNER_JOIN(Role, User.RoleID.EQ(Role.ID)),
	).WHERE(User.ID.EQ(Int32(id))).
		LIMIT(1)
}

// SelectUserByEmail returns a query to fetch a user by email
func SelectUserByEmail(email string) SelectStatement {
	return SELECT(User.AllColumns).
		FROM(User).
		WHERE(User.Email.EQ(String(email)))
}

// SelectPermissionsByRole returns a query to fetch all permissions for a role
func SelectPermissionsByRole(roleID int32) SelectStatement {
	return SELECT(
		Permission.ID,
		Permission.Slug,
		Permission.Name,
	).FROM(
		Permission.INNER_JOIN(RolePermission, Permission.ID.EQ(RolePermission.PermissionID)),
	).WHERE(RolePermission.RoleID.EQ(Int32(roleID)))
}

// SelectPasswordResetByToken returns a query to fetch a password reset with user
func SelectPasswordResetByToken(token string) SelectStatement {
	return SELECT(
		PasswordReset.UserID,
		PasswordReset.Token,
		PasswordReset.CreatedAt,
		User.ID,
		User.Email,
		User.PasswordHash,
		User.IsConfirmed,
		User.RoleID,
		User.Locale,
		User.CreatedAt,
	).FROM(
		PasswordReset.INNER_JOIN(User, User.ID.EQ(PasswordReset.UserID)),
	).WHERE(PasswordReset.Token.EQ(String(token))).
		LIMIT(1)
}

// SelectPasswordResetByUser returns a query to fetch a password reset by user ID
func SelectPasswordResetByUser(userID int32) SelectStatement {
	return SELECT(
		PasswordReset.UserID,
		PasswordReset.Token,
		PasswordReset.CreatedAt,
	).FROM(PasswordReset).
		WHERE(PasswordReset.UserID.EQ(Int32(userID))).
		LIMIT(1)
}

// SelectRegistrationByToken returns a query to fetch a user registration with user
func SelectRegistrationByToken(token string) SelectStatement {
	return SELECT(
		UserRegistration.UserID,
		UserRegistration.Token,
		UserRegistration.CreatedAt,
		User.ID,
		User.Email,
		User.PasswordHash,
		User.IsConfirmed,
		User.RoleID,
		User.Locale,
		User.CreatedAt,
	).FROM(
		UserRegistration.INNER_JOIN(User, User.ID.EQ(UserRegistration.UserID)),
	).WHERE(UserRegistration.Token.EQ(String(token))).
		LIMIT(1)
}

// InsertUser returns an insert statement for a new user
func InsertUser(email, passwordHash string, roleID int32, locale string) InsertStatement {
	stmt := User.INSERT(User.Email, User.PasswordHash, User.RoleID, User.Locale).
		VALUES(email, passwordHash, roleID, locale).
		RETURNING(User.AllColumns)
	stmt.DebugSql()
	return stmt
}

// InsertUserRegistration returns an insert statement for a user registration
func InsertUserRegistration(userID int32, token string) InsertStatement {
	return UserRegistration.INSERT(UserRegistration.UserID, UserRegistration.Token).
		VALUES(userID, token).
		RETURNING(UserRegistration.AllColumns)
}

// InsertPasswordReset returns an insert statement for a password reset
func InsertPasswordReset(userID int32, token string) InsertStatement {
	return PasswordReset.INSERT(PasswordReset.UserID, PasswordReset.Token).
		VALUES(userID, token).
		RETURNING(PasswordReset.AllColumns)
}

// UpdateUser returns an update statement for user details
func UpdateUser(userID int32, email string, isConfirmed bool) UpdateStatement {
	return User.UPDATE(User.Email, User.IsConfirmed).
		SET(email, isConfirmed).
		WHERE(User.ID.EQ(Int32(userID)))
}

// UpdatePassword returns an update statement for user password
func UpdatePassword(userID int32, passwordHash string) UpdateStatement {
	return User.UPDATE(User.PasswordHash).
		SET(passwordHash).
		WHERE(User.ID.EQ(Int32(userID)))
}

// DeletePasswordReset returns a delete statement for password reset by user
func DeletePasswordReset(userID int32) DeleteStatement {
	return PasswordReset.DELETE().WHERE(PasswordReset.UserID.EQ(Int32(userID)))
}

// DeleteRegistration returns a delete statement for user registration
func DeleteRegistration(userID int32) DeleteStatement {
	return UserRegistration.DELETE().WHERE(UserRegistration.UserID.EQ(Int32(userID)))
}

// DeletePasswordResetsBefore returns a delete statement for old password resets
func DeletePasswordResetsBefore(before time.Time) DeleteStatement {
	// SQLite stores timestamps as strings, so we can compare them directly
	return PasswordReset.DELETE().WHERE(
		RawBool("created_at < #ts", RawArgs{"#ts": before.Format(time.RFC3339)}),
	)
}

// DeleteRegistrationsBefore returns a delete statement for old registrations
func DeleteRegistrationsBefore(before time.Time) DeleteStatement {
	return UserRegistration.DELETE().WHERE(
		RawBool("created_at < #ts", RawArgs{"#ts": before.Format(time.RFC3339)}),
	)
}
