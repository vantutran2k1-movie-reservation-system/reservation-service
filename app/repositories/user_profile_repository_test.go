package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestUserProfileRepository_GetProfileByUserID(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserProfileRepository(db)

	profile := utils.GenerateUserProfile()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "first_name", "last_name", "phone_number", "date_of_birth", "profile_picture_url", "bio", "created_at", "updated_at"}).
			AddRow(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, profile.CreatedAt, profile.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_profiles" WHERE user_id = $1 ORDER BY "user_profiles"."id" LIMIT $2`)).
			WithArgs(profile.UserID, 1).
			WillReturnRows(rows)

		result, err := repo.GetProfileByUserID(profile.UserID)

		assert.NotNil(t, profile)
		assert.Nil(t, err)
		assert.Equal(t, profile, result)
	})

	t.Run("profile not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_profiles" WHERE user_id = $1 ORDER BY "user_profiles"."id" LIMIT $2`)).
			WithArgs(profile.UserID, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetProfileByUserID(profile.UserID)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting profile", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_profiles" WHERE user_id = $1 ORDER BY "user_profiles"."id" LIMIT $2`)).
			WithArgs(profile.UserID, 1).
			WillReturnError(errors.New("error getting profile"))

		result, err := NewUserProfileRepository(db).GetProfileByUserID(profile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting profile", err.Error())
	})
}

func TestUserProfileRepository_CreateUserProfile(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserProfileRepository(db)

	profile := utils.GenerateUserProfile()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "user_profiles" ("id","user_id","first_name","last_name","phone_number","date_of_birth","profile_picture_url","bio","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`)).
			WithArgs(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateUserProfile(tx, profile)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating profile", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "user_profiles" ("id","user_id","first_name","last_name","phone_number","date_of_birth","profile_picture_url","bio","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`)).
			WithArgs(profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("error creating profile"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateUserProfile(tx, profile)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating profile", err.Error())
	})
}

func TestUserProfileRepository_UpdateUserProfile(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserProfileRepository(db)

	profile := utils.GenerateUserProfile()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_profiles" SET "user_id"=$1,"first_name"=$2,"last_name"=$3,"phone_number"=$4,"date_of_birth"=$5,"profile_picture_url"=$6,"bio"=$7,"created_at"=$8,"updated_at"=$9 WHERE "id" = $10`)).
			WithArgs(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg(), profile.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UpdateUserProfile(tx, profile)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error updating profile", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_profiles" SET "user_id"=$1,"first_name"=$2,"last_name"=$3,"phone_number"=$4,"date_of_birth"=$5,"profile_picture_url"=$6,"bio"=$7,"created_at"=$8,"updated_at"=$9 WHERE "id" = $10`)).
			WithArgs(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth, profile.ProfilePictureUrl, profile.Bio, sqlmock.AnyArg(), sqlmock.AnyArg(), profile.ID).
			WillReturnError(errors.New("error updating profile"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateUserProfile(tx, profile)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error updating profile", err.Error())
	})
}

func TestUserRepository_UpdateProfilePicture(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserProfileRepository(db)

	profile := utils.GenerateUserProfile()
	url := utils.GenerateURL()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_profiles" SET "profile_picture_url"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs(url, sqlmock.AnyArg(), profile.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		result, err := repo.UpdateProfilePicture(tx, profile, &url)
		tx.Commit()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, url, *result.ProfilePictureUrl)
	})

	t.Run("error updating profile", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "user_profiles" SET "profile_picture_url"=\$1,"updated_at"=\$2 WHERE "id" = \$3`).
			WithArgs(url, sqlmock.AnyArg(), profile.ID).
			WillReturnError(errors.New("error updating profile"))
		mock.ExpectRollback()

		tx := db.Begin()
		result, err := repo.UpdateProfilePicture(tx, profile, &url)
		tx.Rollback()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error updating profile", err.Error())
	})
}
