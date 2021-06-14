package server

import (
	"os"
	"sync"
	"testing"
)

func TestHandleMessage(t *testing.T) {
	s := &Server{
		m: &sync.RWMutex{},
		report: &Report{
			unique: make(map[string]struct{}),
		},
	}

	f, err := os.Create(fileName)
	if err != nil {
		t.Fatal("Error creating file: ", err.Error())
	}
	s.file = f
	defer os.Remove(fileName)
	defer f.Close()

	productSku := "ABCD-1234\n"
	invalid := "AA11!"

	// Test Unique Product SKU
	s.handleMessage(productSku)
	_, exists := s.report.unique[productSku]
	if !exists {
		t.Error("ProductSku has not been saved into Unique list")
		return
	}

	// Test Duplicated Product SKU
	s.handleMessage(productSku)
	if len(s.report.duplicated) == 0 {
		t.Error("Duplicated ProductSku has not been saved into Duplicated list")
	}

	// Test Invalid Product SKU
	s.handleMessage(invalid)
	if len(s.report.invalid) == 0 {
		t.Error("Invalid ProductSku has not been saved into Invalid list")
	}
}
