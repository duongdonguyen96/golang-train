package port

import (
	"context"

	"golang-train/internal/payment/application/dto"
)

type FileExporter interface {
	ExportPaymentsCSV(ctx context.Context, rows []dto.PaymentListItemDTO) ([]byte, error)
}
