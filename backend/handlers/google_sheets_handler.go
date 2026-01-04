package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// GoogleSheetsHandler handles Google Sheets operations
type GoogleSheetsHandler struct{}

// NewGoogleSheetsHandler creates a new GoogleSheetsHandler
func NewGoogleSheetsHandler() *GoogleSheetsHandler {
	return &GoogleSheetsHandler{}
}

// GetPartNameByOrderNumber fetches part name from Google Sheets based on order number
func (h *GoogleSheetsHandler) GetPartNameByOrderNumber(c *gin.Context) {
	orderNumber := c.Param("orderNumber")

	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Order number is required",
		})
		return
	}

	// Google Sheets configuration
	spreadsheetID := "1ONL4NWJSPZZ4BIYVNTidSVpZh7yKlBKeMpDUXXBmT8s"
	sheetName := "Data Transfer"

	// Build the Google Sheets CSV export URL
	sheetsURL := fmt.Sprintf(
		"https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&sheet=%s",
		spreadsheetID,
		url.QueryEscape(sheetName),
	)

	// Fetch the CSV data from Google Sheets
	resp, err := http.Get(sheetsURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch Google Sheets data",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch Google Sheets data",
			"details": fmt.Sprintf("HTTP status: %d", resp.StatusCode),
		})
		return
	}

	// Parse CSV
	reader := csv.NewReader(resp.Body)

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse CSV headers",
			"details": err.Error(),
		})
		return
	}

	// Find column indices
	noOrderIndex := -1
	partNameIndex := -1

	for i, header := range headers {
		cleanHeader := strings.TrimSpace(strings.Trim(header, "\""))
		if cleanHeader == "No Order" {
			noOrderIndex = i
		} else if cleanHeader == "Part Name" {
			partNameIndex = i
		}
	}

	if noOrderIndex == -1 || partNameIndex == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Required columns not found in Google Sheets",
			"details": fmt.Sprintf("No Order column: %d, Part Name column: %d", noOrderIndex, partNameIndex),
		})
		return
	}

	// Search for the order number from all rows (don't skip)
	partName := ""
	found := false
	searchedCount := 0
	rowNumber := 1 // Header is row 1

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to read CSV row",
				"details": err.Error(),
			})
			return
		}

		rowNumber++
		searchedCount++

		// Check if we have enough columns
		if len(row) > noOrderIndex && len(row) > partNameIndex {
			currentOrderNumber := strings.TrimSpace(strings.Trim(row[noOrderIndex], "\""))

			// Skip empty order numbers
			if currentOrderNumber == "" {
				continue
			}

			// Match the order number
			if currentOrderNumber == orderNumber {
				partName = strings.TrimSpace(strings.Trim(row[partNameIndex], "\""))
				found = true
				fmt.Printf("Found order number '%s' at row %d with part name '%s'\n", orderNumber, rowNumber, partName)
				break
			}
		}
	}

	// Log search results for debugging
	fmt.Printf("Searched %d rows total\n", searchedCount)
	fmt.Printf("Looking for order number: '%s'\n", orderNumber)
	fmt.Printf("Found: %v\n", found)

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Order number '%s' not found in Google Sheets", orderNumber),
		})
		return
	}

	// Return the part name
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"order_number": orderNumber,
			"part_name":    partName,
		},
	})
}

// GetAllSheetData fetches all data from Google Sheets
func (h *GoogleSheetsHandler) GetAllSheetData(c *gin.Context) {
	// Google Sheets configuration
	spreadsheetID := "1ONL4NWJSPZZ4BIYVNTidSVpZh7yKlBKeMpDUXXBmT8s"
	sheetName := "Data Transfer"

	// Build the Google Sheets CSV export URL
	sheetsURL := fmt.Sprintf(
		"https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&sheet=%s",
		spreadsheetID,
		url.QueryEscape(sheetName),
	)

	// Fetch the CSV data from Google Sheets
	resp, err := http.Get(sheetsURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch Google Sheets data",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch Google Sheets data",
			"details": fmt.Sprintf("HTTP status: %d", resp.StatusCode),
		})
		return
	}

	// Parse CSV
	reader := csv.NewReader(resp.Body)

	// Read all rows
	allRows, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse CSV data",
			"details": err.Error(),
		})
		return
	}

	if len(allRows) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"headers": []string{},
				"rows":    [][]string{},
				"total":   0,
			},
		})
		return
	}

	// First row is headers
	headers := make([]string, len(allRows[0]))
	for i, header := range allRows[0] {
		headers[i] = strings.TrimSpace(strings.Trim(header, "\""))
	}

	// Rest are data rows
	dataRows := make([][]string, 0)
	for i := 1; i < len(allRows); i++ {
		row := make([]string, len(allRows[i]))
		for j, cell := range allRows[i] {
			row[j] = strings.TrimSpace(strings.Trim(cell, "\""))
		}
		dataRows = append(dataRows, row)
	}

	// Reverse the order so newest data appears first
	for i, j := 0, len(dataRows)-1; i < j; i, j = i+1, j-1 {
		dataRows[i], dataRows[j] = dataRows[j], dataRows[i]
	}

	// Return all data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"headers": headers,
			"rows":    dataRows,
			"total":   len(dataRows),
		},
	})
}
