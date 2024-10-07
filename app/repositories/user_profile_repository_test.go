package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestUserProfileRepository_GetProfileByUserID(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "user_profiles" WHERE "user_profiles"."user_id" = \$1 ORDER BY "user_profiles"."id" LIMIT \$2`).
			WithArgs(profile.UserID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "first_name", "last_name", "phone_number", "date_of_birth", "profile_picture_url", "bio", "created_at", "updated_at"}).
				AddRow(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, profile.CreatedAt, profile.UpdatedAt))

		result, err := NewUserProfileRepository(db).GetProfileByUserID(profile.UserID)

		assert.NotNil(t, profile)
		assert.Nil(t, err)
		assert.Equal(t, profile.ID, result.ID)
		assert.Equal(t, profile.UserID, result.UserID)
		assert.Equal(t, profile.FirstName, result.FirstName)
		assert.Equal(t, profile.LastName, result.LastName)
		assert.Equal(t, profile.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, profile.DateOfBirth, result.DateOfBirth)
		assert.Equal(t, profile.ProfilePictureUrl, result.ProfilePictureUrl)
		assert.Equal(t, profile.Bio, result.Bio)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "user_profiles" WHERE "user_profiles"."user_id" = \$1 ORDER BY "user_profiles"."id" LIMIT \$2`).
			WithArgs(profile.UserID, 1).
			WillReturnError(errors.New("db error"))

		profile, err := NewUserProfileRepository(db).GetProfileByUserID(profile.UserID)

		assert.Nil(t, profile)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserProfileRepository_CreateUserProfile(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "user_profiles"`).
			WithArgs(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewUserProfileRepository(db).CreateUserProfile(tx, profile)
		tx.Commit()

		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "user_profiles"`).
			WithArgs(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewUserProfileRepository(db).CreateUserProfile(tx, profile)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserProfileRepository_UpdateUserProfile(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "user_profiles" SET "user_id"=\$1,"first_name"=\$2,"last_name"=\$3,"phone_number"=\$4,"date_of_birth"=\$5,"profile_picture_url"=\$6,"bio"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10`).
			WithArgs(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg(), profile.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewUserProfileRepository(db).UpdateUserProfile(tx, profile)
		tx.Commit()

		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "user_profiles" SET "user_id"=\$1,"first_name"=\$2,"last_name"=\$3,"phone_number"=\$4,"date_of_birth"=\$5,"profile_picture_url"=\$6,"bio"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10`).
			WithArgs(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg(), profile.ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewUserProfileRepository(db).UpdateUserProfile(tx, profile)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_UpdateProfilePicture(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()
	url := utils.GenerateRandomURL()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "user_profiles" SET "profile_picture_url"=\$1,"updated_at"=\$2 WHERE "id" = \$3`).
			WithArgs(url, sqlmock.AnyArg(), profile.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		result, err := NewUserProfileRepository(db).UpdateProfilePicture(tx, profile, &url)
		tx.Commit()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, url, *result.ProfilePictureUrl)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "user_profiles" SET "profile_picture_url"=\$1,"updated_at"=\$2 WHERE "id" = \$3`).
			WithArgs(url, sqlmock.AnyArg(), profile.ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		result, err := NewUserProfileRepository(db).UpdateProfilePicture(tx, profile, &url)
		tx.Rollback()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
