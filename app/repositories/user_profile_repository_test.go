package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"gorm.io/gorm"
)

func TestUserProfileRepository_GetProfileByUserID_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	expectedProfile := utils.GenerateRandomUserProfile()

	mock.ExpectQuery(`SELECT \* FROM "user_profiles" WHERE "user_profiles"."user_id" = \$1 ORDER BY "user_profiles"."id" LIMIT \$2`).
		WithArgs(expectedProfile.UserID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "first_name", "last_name", "phone_number", "date_of_birth", "profile_picture_url", "bio", "created_at", "updated_at"}).
			AddRow(expectedProfile.ID, expectedProfile.UserID, expectedProfile.FirstName, expectedProfile.LastName, expectedProfile.PhoneNumber, expectedProfile.DateOfBirth, expectedProfile.ProfilePictureUrl, expectedProfile.Bio, expectedProfile.CreatedAt, expectedProfile.UpdatedAt))

	profile, err := NewUserProfileRepository(db).GetProfileByUserID(expectedProfile.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, expectedProfile.ID, profile.ID)
	assert.Equal(t, expectedProfile.UserID, profile.UserID)
	assert.Equal(t, expectedProfile.FirstName, profile.FirstName)
	assert.Equal(t, expectedProfile.LastName, profile.LastName)
	assert.Equal(t, expectedProfile.PhoneNumber, profile.PhoneNumber)
	assert.Equal(t, expectedProfile.DateOfBirth, profile.DateOfBirth)
	assert.Equal(t, expectedProfile.ProfilePictureUrl, profile.ProfilePictureUrl)
	assert.Equal(t, expectedProfile.Bio, profile.Bio)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepository_GetProfileByUserID_NotFound(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "user_profiles" WHERE "user_profiles"."user_id" = \$1 ORDER BY "user_profiles"."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	profile, err := NewUserProfileRepository(db).GetProfileByUserID(userID)

	assert.Nil(t, profile)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepository_CreateUserProfile_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "user_profiles"`).
		WithArgs(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := db.Begin()
	err := NewUserProfileRepository(db).CreateUserProfile(tx, profile)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepository_CreateUserProfile_Failure(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "user_profiles"`).
		WithArgs(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	tx := db.Begin()
	err := NewUserProfileRepository(db).CreateUserProfile(tx, profile)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidData, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepository_UpdateUserProfile_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "user_profiles" SET "user_id"=\$1,"first_name"=\$2,"last_name"=\$3,"phone_number"=\$4,"date_of_birth"=\$5,"profile_picture_url"=\$6,"bio"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10`).
		WithArgs(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg(), profile.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := db.Begin()
	err := NewUserProfileRepository(db).UpdateUserProfile(tx, profile)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepository_UpdateUserProfile_Failure(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	profile := utils.GenerateRandomUserProfile()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "user_profiles" SET "user_id"=\$1,"first_name"=\$2,"last_name"=\$3,"phone_number"=\$4,"date_of_birth"=\$5,"profile_picture_url"=\$6,"bio"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10`).
		WithArgs(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg(), profile.ID).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	tx := db.Begin()
	err := NewUserProfileRepository(db).UpdateUserProfile(tx, profile)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
