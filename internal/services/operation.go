package services

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
)

type operationService struct {
	operationRepository domain.OperationRepository
	contextTimeout      time.Duration
}

func NewOperationService(operationRepository domain.OperationRepository, timeout time.Duration) domain.OperationService {
	return &operationService{
		operationRepository: operationRepository,
		contextTimeout:      timeout,
	}
}

func (osr *operationService) GetUserOperations(ctx context.Context, userID int, fromTS sql.NullTime, toTS sql.NullTime) (fileName string, err error) {
	c, cancel := context.WithTimeout(ctx, osr.contextTimeout)
	defer cancel()

	operations, err := osr.operationRepository.GetUserOperations(c, userID, fromTS, toTS)
	if err != nil {
		return
	}

	fileName, err = buildFileName(userID, fromTS, toTS)
	if err != nil {
		return
	}

	file, err := os.Create("./web/static/" + fileName)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"N", "Operation", "Segment", "Time"})
	if err != nil {
		return
	}
	for n, operation := range operations {
		err = writer.Write([]string{strconv.FormatInt(int64(n), 10), string(operation.OpType), operation.DoneAt.String()})
		return
	}
	return
}

func buildFileName(userID int, fromTS sql.NullTime, toTS sql.NullTime) (fileName string, err error) {
	var userIDStr, fromTSStr, toTSStr string
	userIDStr = strconv.FormatInt(int64(userID), 10)
	if !fromTS.Valid {
		err = errors.New("fromTS must be specified")
		return
	}
	if !toTS.Valid {
		toTS.Time = time.Now()
		toTS.Valid = true
	}
	fromTSStr = fromTS.Time.Format("2006-01-02")
	toTSStr = toTS.Time.Format("2006-01-02")

	fileName = userIDStr + "_" + fromTSStr + "_" + toTSStr + ".csv"
	return
}
