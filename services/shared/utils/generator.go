package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// IDOptions berisi konfigurasi untuk generate ID
type IDOptions struct {
	// Prefix yang akan ditambahkan di awal ID (default: kosong)
	Prefix string

	// Separator antara prefix dan random number (default: "-")
	Separator string

	// Panjang angka random yang diinginkan (default: 12)
	NumberLength int

	// Format custom untuk ID
	// Contoh: "{prefix}_{rand:6}_{timestamp}"
	CustomFormat string
}

func GenerateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func GenerateRandomId(folder string) string {
	return fmt.Sprintf("%s-%012d", folder, rand.Intn(1000000000000))
}

// DefaultIDOptions memberikan default values untuk IDOptions
func DefaultIDOptions() IDOptions {
	return IDOptions{
		Separator:    "-",
		NumberLength: 12,
	}
}

// GenerateID membuat random ID dengan opsi default
func GenerateID(prefix string) string {
	opts := DefaultIDOptions()
	opts.Prefix = prefix
	return GenerateCustomID(opts)
}

// GenerateCustomID membuat random ID dengan opsi kustom
func GenerateCustomID(opts IDOptions) string {
	// Jika ada custom format, gunakan itu
	if opts.CustomFormat != "" {
		return generateWithFormat(opts)
	}

	// Gunakan default values jika tidak diset
	if opts.Separator == "" {
		opts.Separator = "-"
	}
	if opts.NumberLength <= 0 {
		opts.NumberLength = 12
	}

	// Generate random number dengan panjang yang ditentukan
	maxNum := int64(1)
	for i := 0; i < opts.NumberLength; i++ {
		maxNum *= 10
	}
	randomNum := rand.Int63n(maxNum)

	// Format number dengan leading zeros
	numberFormat := fmt.Sprintf("%%0%dd", opts.NumberLength)

	// Jika tidak ada prefix, return number saja
	if opts.Prefix == "" {
		return fmt.Sprintf(numberFormat, randomNum)
	}

	return fmt.Sprintf("%s%s"+numberFormat,
		opts.Prefix,
		opts.Separator,
		randomNum,
	)
}

// generateWithFormat membuat ID dengan format kustom
func generateWithFormat(opts IDOptions) string {
	result := opts.CustomFormat

	// Replace placeholder untuk prefix
	result = strings.ReplaceAll(result, "{prefix}", opts.Prefix)

	result = strings.ReplaceAll(result, "{timestamp}",
		fmt.Sprintf("%d", time.Now().Unix()))

	for i := 1; i <= 20; i++ {
		placeholder := fmt.Sprintf("{rand:%d}", i)
		if strings.Contains(result, placeholder) {
			maxNum := int64(1)
			for j := 0; j < i; j++ {
				maxNum *= 10
			}
			randomNum := rand.Int63n(maxNum)
			numberFormat := fmt.Sprintf("%%0%dd", i)
			result = strings.ReplaceAll(result, placeholder,
				fmt.Sprintf(numberFormat, randomNum))
		}
	}

	return result
}
