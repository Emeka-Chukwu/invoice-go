package helper

import (
	"fmt"
	"go-invoice/domain"

	"github.com/jung-kurt/gofpdf"
)

func AddSenderCustomerInfo(pdf *gofpdf.Fpdf, customer domain.CustomerResponse) {
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(150, 150, 150)

	pdf.SetXY(20, 20)
	pdf.Cell(40, 5, "SENDER")
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(20, 26)
	pdf.Cell(40, 5, "Fabulous Enterprise")
	pdf.SetXY(20, 31)
	pdf.Cell(40, 5, "+386 969 271 3115")
	pdf.SetXY(20, 36)
	pdf.Cell(40, 5, "1331 Hart Ridge Road 48436 Gaines, MI")
	pdf.SetXY(20, 41)
	pdf.Cell(40, 5, "info@fabulousenterise.co")

	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(130, 20)
	pdf.Cell(40, 5, "CUSTOMER")
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(130, 26)
	pdf.Cell(40, 5, customer.Name)
	pdf.SetXY(130, 31)
	pdf.Cell(40, 5, customer.Phone)
	pdf.SetXY(130, 36)
	pdf.Cell(40, 5, customer.Email)
}

func AddInvoiceDetails(pdf *gofpdf.Fpdf, invoice domain.InvoiceResponse) {
	pdf.SetFillColor(252, 242, 244) // Light pink
	pdf.Rect(10, 50, 190, 25, "F")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(150, 150, 150)

	pdf.SetXY(20, 55)
	pdf.Cell(40, 5, "INVOICE NO")
	pdf.SetXY(20, 60)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, fmt.Sprintf("%d", invoice.ID))

	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(70, 55)
	pdf.Cell(40, 5, "ISSUE DATE")
	pdf.SetXY(70, 60)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, invoice.IssueDate.String())

	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(120, 55)
	pdf.Cell(40, 5, "DUE DATE")
	pdf.SetXY(120, 60)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, invoice.DueDate.GoString())

	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(170, 55)
	pdf.Cell(40, 5, "BILLING CURRENCY")
	pdf.SetXY(170, 60)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, "USD ($)")
}

func AddItemsTable(pdf *gofpdf.Fpdf, items []domain.InvoiceItem) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(20, 85)
	pdf.Cell(40, 10, "Items")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(20, 95)
	pdf.Cell(90, 10, "ITEM")
	pdf.Cell(20, 10, "QTY")
	pdf.Cell(30, 10, "PRICE")
	pdf.Cell(30, 10, "TOTAL")
	for index, item := range items {
		AddItemRow(pdf, float64((105 + (index * 22))), item.Title, item.Description, fmt.Sprintf("%d", item.Quantity), fmt.Sprintf("$%f", item.UnitPrice), fmt.Sprintf("$%f", item.UnitPrice*float64(item.Quantity)))
	}
}

func AddItemRow(pdf *gofpdf.Fpdf, y float64, item, description, quantity, price, total string) {
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(20, y)
	pdf.Cell(90, 5, item)
	pdf.SetXY(110, y)
	pdf.Cell(20, 5, quantity)
	pdf.SetXY(130, y)
	pdf.Cell(30, 5, price)
	pdf.SetXY(160, y)
	pdf.Cell(30, 5, total)

	if description != "" {
		pdf.SetTextColor(150, 150, 150)
		pdf.SetXY(20, y+5)
		pdf.Cell(90, 5, description)
	}
}

func AddTotals(pdf *gofpdf.Fpdf, subTotal, discount, discountedAmount, totalAmount string) {
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(130, 180)
	pdf.Cell(30, 5, "SUBTOTAL")
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(160, 180)
	pdf.Cell(30, 5, subTotal)

	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(130, 185)
	pdf.Cell(30, 5, discount)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(160, 185)
	pdf.Cell(30, 5, discountedAmount)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(120, 195)
	pdf.Cell(30, 5, "TOTAL AMOUNT DUE")
	pdf.SetXY(160, 195)
	pdf.Cell(30, 5, totalAmount)
}

func AddPaymentInfo(pdf *gofpdf.Fpdf, bankInfo domain.BankInformation) {
	pdf.SetFillColor(252, 242, 244) // Light pink
	pdf.Rect(10, 210, 190, 50, "F")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(20, 215)
	pdf.Cell(40, 5, "PAYMENT INFORMATION")

	AddPaymentInfoRow(pdf, 225, "ACCOUNT NAME", bankInfo.AccountName)
	AddPaymentInfoRow(pdf, 235, "ACCOUNT NUMBER", bankInfo.AccountNumber)
	AddPaymentInfoRow(pdf, 245, "BANK ADDRESS", bankInfo.BankAddress)

	pdf.SetXY(110, 225)
	pdf.SetTextColor(150, 150, 150)
	pdf.Cell(40, 5, "ACH ROUTING NO")
	pdf.SetXY(110, 230)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, bankInfo.ACHRoutingNo)

	pdf.SetXY(160, 225)
	pdf.SetTextColor(150, 150, 150)
	pdf.Cell(40, 5, "BANK NAME")
	pdf.SetXY(160, 230)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, "USD ($)")

	pdf.SetXY(110, 235)
	pdf.SetTextColor(150, 150, 150)
	pdf.Cell(40, 5, "ACCOUNT NUMBER")
	pdf.SetXY(110, 240)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, bankInfo.BankName)
}

func AddPaymentInfoRow(pdf *gofpdf.Fpdf, y float64, label, value string) {
	pdf.SetXY(20, y)
	pdf.SetTextColor(150, 150, 150)
	pdf.Cell(40, 5, label)
	pdf.SetXY(20, y+5)
	pdf.SetTextColor(80, 80, 80)
	pdf.Cell(40, 5, value)
}

func AddNote(pdf *gofpdf.Fpdf, note string) {
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(20, 270)
	pdf.Cell(40, 5, "NOTE")
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(20, 275)
	pdf.Cell(40, 5, note)
}
