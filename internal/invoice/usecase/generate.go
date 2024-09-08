package invoice_usecase

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func GenerateInvoicePDF() ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)

	// Add sender information
	pdf.SetXY(10, 10)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Fabulous Enterprise")

	pdf.SetFont("Arial", "", 12)
	pdf.SetXY(10, 20)
	pdf.Cell(0, 10, "+386 989 271 3115")
	pdf.SetXY(10, 25)
	pdf.Cell(0, 10, "1331 Hart Ridge Road 48436 Gaines, MI")
	pdf.SetXY(10, 30)
	pdf.Cell(0, 10, "info@fabulousenterise.co")

	// Add recipient information
	pdf.SetXY(140, 10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "CUSTOMER")
	pdf.SetFont("Arial", "", 12)
	pdf.SetXY(140, 20)
	pdf.Cell(0, 10, "Olaniyi Ojo Adewale")
	pdf.SetXY(140, 25)
	pdf.Cell(0, 10, "+386 989 271 3115")
	pdf.SetXY(140, 30)
	pdf.Cell(0, 10, "info@fabulousenterise.co")

	// Add invoice details
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(10, 40)
	pdf.Cell(0, 10, "INVOICE DETAILS")
	pdf.SetFont("Arial", "", 12)

	pdf.SetXY(10, 50)
	pdf.Cell(0, 10, fmt.Sprintf("Invoice No: %s", "1023902390"))
	pdf.SetXY(10, 55)
	pdf.Cell(0, 10, fmt.Sprintf("Issue Date: %s", "March 30th, 2023"))
	pdf.SetXY(10, 60)
	pdf.Cell(0, 10, fmt.Sprintf("Due Date: %s", "May 19th, 2023"))
	pdf.SetXY(10, 65)
	pdf.Cell(0, 10, "Billing Currency: USD ($)")

	// Add items header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(10, 75)
	pdf.Cell(40, 10, "Items")
	pdf.Cell(40, 10, "Quantity")
	pdf.Cell(40, 10, "Unit Price")
	pdf.Cell(40, 10, "Amount")

	// Add items
	items := []struct {
		Description string
		Quantity    int
		UnitPrice   float64
		Amount      float64
	}{
		{"Email Marketing", 10, 1500, 15000},
		{"Video looping effect", 6, 1110500, 6663000},
		{"Graphic design for emails", 7, 2750, 19250},
		{"Video looping effect", 6, 1110500, 6663000},
	}

	// Print items
	pdf.SetFont("Arial", "", 12)
	startY := 85
	for _, item := range items {
		pdf.SetXY(10, float64(startY))
		pdf.Cell(40, 10, item.Description)
		pdf.Cell(40, 10, fmt.Sprintf("%d", item.Quantity))
		pdf.Cell(40, 10, fmt.Sprintf("$%.2f", item.UnitPrice))
		pdf.Cell(40, 10, fmt.Sprintf("$%.2f", item.Amount))
		startY += 10
	}

	// Add totals
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(10, float64(startY+10))
	pdf.Cell(40, 10, "Subtotal:")
	pdf.Cell(40, 10, "$6,697,200.00")

	pdf.SetXY(10, float64(startY+20))
	pdf.Cell(40, 10, "Discount (2.5%):")
	pdf.Cell(40, 10, "$167,430.00")

	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(10, float64(startY+30))
	pdf.Cell(40, 10, "TOTAL AMOUNT DUE:")
	pdf.Cell(40, 10, "$6,529,770.00")

	// Add payment information
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(10, float64(startY+50))
	pdf.Cell(40, 10, "PAYMENT INFORMATION")

	pdf.SetFont("Arial", "", 12)
	pdf.SetXY(10, float64(startY+60))
	pdf.Cell(40, 10, "Account Name: 1023902390")
	pdf.SetXY(10, float64(startY+65))
	pdf.Cell(40, 10, "Account Number: March 30th, 2023")
	pdf.SetXY(10, float64(startY+70))
	pdf.Cell(40, 10, "ACH Routing No: May 19th, 2023")
	pdf.SetXY(10, float64(startY+75))
	pdf.Cell(40, 10, "Bank Name: USD ($)")

	// // Output to PDF file
	// err := pdf.OutputFileAndClose("invoice.pdf")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Buffer to hold the PDF
	var buf bytes.Buffer

	// Output the PDF to the buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
	// fmt.Println(pdf)
	// base64PDF := base64.StdEncoding.EncodeToString(buf.Bytes())
	// fmt.Println("PDF Base64 Output:\n", base64PDF)
	// // Set headers to return PDF in the response
	// w.Header().Set("Content-Type", "application/pdf")
	// w.Header().Set("Content-Disposition", "attachment; filename=\"invoice.pdf\"")
	// w.WriteHeader(http.StatusOK)

	// // Write the PDF bytes to the response
	// _, err = w.Write(pdfData)
	// if err != nil {
	// 	log.Println("Failed to send PDF:", err)
	// }
}
