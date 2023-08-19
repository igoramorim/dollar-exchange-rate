package txtfile

import (
	"context"
	"fmt"
	"github.com/igoramorim/dollar-exchange-rate/internal/repository"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

func New(filename string) (*TxtFile, error) {
	log.Println("initializing txtfile repository:", filename)

	file, err := createFile(filename)
	if err != nil {
		return nil, errors.WithMessagef(err, "txtfile: creating file %s", filename)
	}

	return &TxtFile{
		file: file,
	}, nil
}

func createFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

var _ repository.Repository = (*TxtFile)(nil)

type TxtFile struct {
	file *os.File
}

func (f *TxtFile) Save(_ context.Context, exchangeRate float64) error {
	log.Println("txtfile: saving", exchangeRate)

	_, err := f.file.WriteString(f.formatOutput(exchangeRate))
	if err != nil {
		return errors.WithMessage(err, "txtfile")
	}

	return nil
}

func (f *TxtFile) Close() error {
	return f.file.Close()
}

func (f *TxtFile) formatOutput(exchangeRate float64) string {
	now := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf("%s Dólar: %.4f\n", now, exchangeRate)
}
