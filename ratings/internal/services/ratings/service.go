package ratings

import (
	"database/sql"
	"fmt"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/ratings"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetAllRatingsForAMusic(musicID uuid.UUID) ([]models.Rating, error) {
	var err error
	// calling repository
	ratings, err := repository.GetAllRatingsForAMusic(musicID)
	// managing errors
	if err != nil {
		logrus.Errorf("error retrieving ratings : %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    500,
		}
	}

	return ratings, nil
}

func GetMusicRating(musicID uuid.UUID, ratingID uuid.UUID) (*models.Rating, error) {
	rating, err := repository.GetMusicRating(musicID, ratingID)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.CustomError{
				Message: "rating not found",
				Code:    http.StatusNotFound,
			}
		}
		logrus.Errorf("error retrieving rating : %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    500,
		}
	}

	return rating, err
}

func AddMusicRating(musicID uuid.UUID, ratingRequest models.RatingRequest) (*models.Rating, error) {
	if ratingRequest.Comment == nil || ratingRequest.Rating == nil || ratingRequest.UserID == nil {
		return nil, &models.CustomError{
			Message: "missing fields",
			Code:    http.StatusUnprocessableEntity,
		}
	}

	if *ratingRequest.Rating < 0 || *ratingRequest.Rating > 5 {
		return nil, &models.CustomError{
			Message: "rating must be between 0 and 5",
			Code:    http.StatusUnprocessableEntity,
		}
	}

	userID, err := uuid.FromString(*ratingRequest.UserID)
	if err != nil {
		logrus.Errorf("parsing error : %s", err.Error())
		return nil, &models.CustomError{
			Message: fmt.Sprintf("cannot parse id (%s) as UUID", ratingRequest.UserID),
			Code:    http.StatusUnprocessableEntity,
		}
	}

	id, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("error creating uuid : %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	newRating := models.Rating{
		Id:         id,
		Comment:    *ratingRequest.Comment,
		Rating:     *ratingRequest.Rating,
		RatingDate: time.Now(),
		MusicID:    musicID,
		UserID:     userID,
	}

	err = repository.AddMusicRating(newRating)
	if err != nil {
		logrus.Errorf("Error adding rating : %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return &newRating, err
}

func ModifyMusicRating(musicID uuid.UUID, ratingID uuid.UUID, ratingRequest models.RatingRequest) (*models.Rating, error) {
	if ratingRequest.Rating != nil {
		if *ratingRequest.Rating < 0 || *ratingRequest.Rating > 5 {
			return nil, &models.CustomError{
				Message: "rating must be between 0 and 5",
				Code:    http.StatusUnprocessableEntity,
			}
		}
	}

	if ratingRequest.UserID != nil {
		_, err := uuid.FromString(*ratingRequest.UserID)
		if err != nil {
			logrus.Errorf("parsing error : %s", err.Error())
			return nil, &models.CustomError{
				Message: fmt.Sprintf("cannot parse id (%s) as UUID", ratingRequest.UserID),
				Code:    http.StatusUnprocessableEntity,
			}
		}
	}

	rating, err := repository.ModifyMusicRating(musicID, ratingID, ratingRequest)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.CustomError{
				Message: "rating not found",
				Code:    http.StatusNotFound,
			}
		}
		logrus.Errorf("Error adding rating : %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return rating, err
}

func DeleteMusicRating(musicID uuid.UUID, ratingID uuid.UUID) error {
	err := repository.DeleteMusicRating(musicID, ratingID)
	if err != nil {
		logrus.Errorf("Error deleting music : %s", err.Error())
		return &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return err
}
