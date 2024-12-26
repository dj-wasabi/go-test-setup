package utils

import (
	"errors"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

func IsWritable(path string) (bool, error) {
	tmpFile := "tmpfile"

	file, err := os.CreateTemp(path, tmpFile)
	if err != nil {
		return false, err
	}

	defer os.Remove(file.Name())
	defer file.Close()

	return true, nil
}

func HandleAuthError(errorCode, logID string, logger *slog.Logger) error {
	err := model.GetError(errorCode, logID)
	logger.Error("log_id", logID, err.Error)
	return errors.New(err.Error)
}

func HandleHTTPError(c *gin.Context, status int, err error) {
	error := model.NewError(err.Error())
	c.JSON(status, error)
}
