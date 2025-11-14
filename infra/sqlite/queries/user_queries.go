package queries

import (
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectUserByID returns a query to fetch a user with their role
func SelectUserByID(id int64) SelectStatement {
	return SELECT(
		Users.ID,
		Users.Email,
		Users.PasswordHash,
		Users.IsConfirmed,
		Users.RoleID,
		Users.Locale,
		Users.CreatedAt,
		Roles.Name.AS("role.name"),
	).FROM(
		Users.INNER_JOIN(Roles, Users.RoleID.EQ(Roles.ID)),
	).WHERE(Users.ID.EQ(Int(id))).
		LIMIT(1)
}

// SelectUserByEmail returns a query to fetch a user by email
func SelectUserByEmail(email string) SelectStatement {
	return SELECT(
		Users.ID,
		Users.Email,
		Users.PasswordHash,
		Users.IsConfirmed,
		Users.RoleID,
		Users.Locale,
		Users.CreatedAt,
	).FROM(Users).
		WHERE(Users.Email.EQ(String(email))).
		LIMIT(1)
}

// SelectPermissionsByRole returns a query to fetch all permissions for a role
func SelectPermissionsByRole(roleID int64) SelectStatement {
	return SELECT(
		Permissions.ID,
		Permissions.Slug,
		Permissions.Name,
	).FROM(
		Permissions.INNER_JOIN(RolePermissions, Permissions.ID.EQ(RolePermissions.PermissionID)),
	).WHERE(RolePermissions.RoleID.EQ(Int(roleID)))
}

// SelectPasswordResetByToken returns a query to fetch a password reset with user
func SelectPasswordResetByToken(token string) SelectStatement {
	return SELECT(
		PasswordResets.UserID,
		PasswordResets.Token,
		PasswordResets.CreatedAt,
		Users.ID,
		Users.Email,
		Users.PasswordHash,
		Users.IsConfirmed,
		Users.RoleID,
		Users.Locale,
		Users.CreatedAt,
	).FROM(
		PasswordResets.INNER_JOIN(Users, Users.ID.EQ(PasswordResets.UserID)),
	).WHERE(PasswordResets.Token.EQ(String(token))).
		LIMIT(1)
}

// SelectPasswordResetByUser returns a query to fetch a password reset by user ID
func SelectPasswordResetByUser(userID int64) SelectStatement {
	return SELECT(
		PasswordResets.UserID,
		PasswordResets.Token,
		PasswordResets.CreatedAt,
	).FROM(PasswordResets).
		WHERE(PasswordResets.UserID.EQ(Int(userID))).
		LIMIT(1)
}

// SelectRegistrationByToken returns a query to fetch a user registration with user
func SelectRegistrationByToken(token string) SelectStatement {
	return SELECT(
		UserRegistrations.UserID,
		UserRegistrations.Token,
		UserRegistrations.CreatedAt,
		Users.ID,
		Users.Email,
		Users.PasswordHash,
		Users.IsConfirmed,
		Users.RoleID,
		Users.Locale,
		Users.CreatedAt,
	).FROM(
		UserRegistrations.INNER_JOIN(Users, Users.ID.EQ(UserRegistrations.UserID)),
	).WHERE(UserRegistrations.Token.EQ(String(token))).
		LIMIT(1)
}

// InsertUser returns an insert statement for a new user
func InsertUser(email, passwordHash string, roleID int64, locale string) InsertStatement {
	stmt := Users.INSERT(Users.Email, Users.PasswordHash, Users.RoleID, Users.Locale).
		VALUES(email, passwordHash, roleID, locale).
		RETURNING(Users.AllColumns()...)
	stmt.DebugSql()
	return stmt
}

// InsertUserRegistration returns an insert statement for a user registration
func InsertUserRegistration(userID int64, token string) InsertStatement {
	return UserRegistrations.INSERT(UserRegistrations.UserID, UserRegistrations.Token).
		VALUES(userID, token).
		RETURNING(UserRegistrations.AllColumns()...)
}

// InsertPasswordReset returns an insert statement for a password reset
func InsertPasswordReset(userID int64, token string) InsertStatement {
	return PasswordResets.INSERT(PasswordResets.UserID, PasswordResets.Token).
		VALUES(userID, token).
		RETURNING(PasswordResets.AllColumns()...)
}

// UpdateUser returns an update statement for user details
func UpdateUser(userID int64, email string, isConfirmed bool) UpdateStatement {
	return Users.UPDATE(Users.Email, Users.IsConfirmed).
		SET(email, isConfirmed).
		WHERE(Users.ID.EQ(Int(userID)))
}

// UpdatePassword returns an update statement for user password
func UpdatePassword(userID int64, passwordHash string) UpdateStatement {
	return Users.UPDATE(Users.PasswordHash).
		SET(passwordHash).
		WHERE(Users.ID.EQ(Int(userID)))
}

// DeletePasswordReset returns a delete statement for password reset by user
func DeletePasswordReset(userID int64) DeleteStatement {
	return PasswordResets.DELETE().WHERE(PasswordResets.UserID.EQ(Int(userID)))
}

// DeleteRegistration returns a delete statement for user registration
func DeleteRegistration(userID int64) DeleteStatement {
	return UserRegistrations.DELETE().WHERE(UserRegistrations.UserID.EQ(Int(userID)))
}

// DeletePasswordResetsBefore returns a delete statement for old password resets
func DeletePasswordResetsBefore(before time.Time) DeleteStatement {
	// SQLite stores timestamps as strings, so we can compare them directly
	return PasswordResets.DELETE().WHERE(
		RawBool("created_at < #ts", RawArgs{"#ts": before.Format(time.RFC3339)}),
	)
}

// DeleteRegistrationsBefore returns a delete statement for old registrations
func DeleteRegistrationsBefore(before time.Time) DeleteStatement {
	return UserRegistrations.DELETE().WHERE(
		RawBool("created_at < #ts", RawArgs{"#ts": before.Format(time.RFC3339)}),
	)
}
