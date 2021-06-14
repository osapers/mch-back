package pgutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
)

func IsSqlNoRows(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, pg.ErrNoRows) {
		return true
	}

	return false
}

func HandleSqlIntegrityViolation(err error) error {
	if err == nil {
		return nil
	}

	pgErr, ok := err.(pg.Error)
	if !ok {
		return err
	}

	if pgErr.IntegrityViolation() {
		if strings.Contains(pgErr.Error(), "unique constraint") {
			return fmt.Errorf("entity already exists")
		} else {
			return fmt.Errorf("integrity violation error")
		}
	}

	return err
}
