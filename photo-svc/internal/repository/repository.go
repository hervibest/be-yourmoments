package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Helper untuk membuat placeholder dinamis ($2, $3, $4, ...)
func generatePlaceholders(count, start int) string {
	placeholders := []string{}
	for i := 0; i < count; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", start+i))
	}
	return strings.Join(placeholders, ", ")
}

// Helper untuk membuat values dinamis untuk batch INSERT
func generateInsertValues(count int) string {
	values := []string{}
	for i := 0; i < count; i++ {
		n := i + 2 // karena $1 untuk photo_id
		values = append(values, fmt.Sprintf("($1, $%d, NOW(), NOW())", n))
	}
	return strings.Join(values, ", ")
}

// Helper untuk konversi []string ke []interface{}
func convertToInterface(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
