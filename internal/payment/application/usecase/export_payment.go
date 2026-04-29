package usecase

import (
	"context"
	"fmt"
	"time"

	"golang-train/internal/payment/application/port"
)

type ExportPaymentsUsecase struct {
	q        port.PaymentQuery
	exporter port.FileExporter
	storage  port.Storage
}

func NewExportPaymentsUsecase(q port.PaymentQuery, exporter port.FileExporter, storage port.Storage) *ExportPaymentsUsecase {
	return &ExportPaymentsUsecase{q: q, exporter: exporter, storage: storage}
}

func (uc *ExportPaymentsUsecase) Execute(ctx context.Context) (string, error) {
	rows, err := uc.q.ListWithUserName(ctx, 1000)
	if err != nil {
		return "", err
	}
	csvBytes, err := uc.exporter.ExportPaymentsCSV(ctx, rows)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("exports/payments_%d.csv", time.Now().Unix())
	return uc.storage.PutObject(ctx, key, "text/csv", csvBytes)
}
