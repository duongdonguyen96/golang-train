package file

import (
	"bytes"
	"context"
	"encoding/csv"
	"strconv"
	"time"

	"golang-train/internal/payment/application/dto"
)

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter { return &CSVExporter{} }

func (e *CSVExporter) ExportPaymentsCSV(ctx context.Context, rows []dto.PaymentListItemDTO) ([]byte, error) {
	_ = ctx
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)

	_ = w.Write([]string{"payment_id", "user_id", "user_name", "amount", "currency", "created_at"})
	for _, r := range rows {
		_ = w.Write([]string{
			strconv.FormatUint(r.PaymentID, 10),
			strconv.FormatUint(r.UserID, 10),
			r.UserName,
			strconv.FormatInt(r.Amount, 10),
			r.Currency,
			r.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
