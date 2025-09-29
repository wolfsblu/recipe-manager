package mapper

import (
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (m *DBMapper) FromUserForPasswordResetToken(user *domain.User, token string) database.CreatePasswordResetTokenParams {
	return database.CreatePasswordResetTokenParams{
		UserID: user.ID,
		Token:  token,
	}
}

func (m *DBMapper) FromUserDetails(userDetails domain.UserDetails, roleID int64) database.CreateUserParams {
	return database.CreateUserParams{
		Email:        userDetails.Email,
		PasswordHash: userDetails.PasswordHash,
		RoleID:       roleID,
		Locale:       userDetails.Locale,
	}
}

func (m *DBMapper) FromUserForRegistration(user *domain.User, token string) database.CreateUserRegistrationParams {
	return database.CreateUserRegistrationParams{
		UserID: user.ID,
		Token:  token,
	}
}

func (m *DBMapper) FromUserForPasswordUpdate(hashedPassword string, userID int64) database.UpdatePasswordByUserIdParams {
	return database.UpdatePasswordByUserIdParams{
		PasswordHash: hashedPassword,
		ID:           userID,
	}
}

func (m *DBMapper) FromUserForUpdate(user *domain.User) database.UpdateUserParams {
	return database.UpdateUserParams{
		Email:       user.Email,
		IsConfirmed: user.Confirmed,
		ID:          user.ID,
	}
}
