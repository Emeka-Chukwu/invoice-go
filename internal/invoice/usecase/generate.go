package invoice_usecase

// import (
// 	"bytes"
// 	"fmt"
// 	"strconv"

// 	"github.com/jung-kurt/gofpdf"
// )

// // func (i invoiceUsecase) GenerateInvoicePDF() ([]byte, error) {
// // pdf := gofpdf.New("P", "mm", "A4", "")
// // pdf.AddPage()

// // // Set font
// // pdf.SetFont("Arial", "B", 16)

// // // Add sender information
// // pdf.SetXY(10, 10)
// // pdf.SetFont("Arial", "B", 16)
// // pdf.Cell(0, 10, "Fabulous Enterprise")

// // pdf.SetFont("Arial", "", 12)
// // pdf.SetXY(10, 20)
// // pdf.Cell(0, 10, "+386 989 271 3115")
// // pdf.SetXY(10, 25)
// // pdf.Cell(0, 10, "1331 Hart Ridge Road 48436 Gaines, MI")
// // pdf.SetXY(10, 30)
// // pdf.Cell(0, 10, "info@fabulousenterise.co")

// // // Add recipient information
// // pdf.SetXY(140, 10)
// // pdf.SetFont("Arial", "B", 12)
// // pdf.Cell(0, 10, "CUSTOMER")
// // pdf.SetFont("Arial", "", 12)
// // pdf.SetXY(140, 20)
// // pdf.Cell(0, 10, "Olaniyi Ojo Adewale")
// // pdf.SetXY(140, 25)
// // pdf.Cell(0, 10, "+386 989 271 3115")
// // pdf.SetXY(140, 30)
// // pdf.Cell(0, 10, "info@fabulousenterise.co")

// // // Add invoice details
// // pdf.SetFont("Arial", "B", 12)
// // pdf.SetXY(10, 40)
// // pdf.Cell(0, 10, "INVOICE DETAILS")
// // pdf.SetFont("Arial", "", 12)

// // pdf.SetXY(10, 50)
// // pdf.Cell(0, 10, fmt.Sprintf("Invoice No: %s", "1023902390"))
// // pdf.SetXY(10, 55)
// // pdf.Cell(0, 10, fmt.Sprintf("Issue Date: %s", "March 30th, 2023"))
// // pdf.SetXY(10, 60)
// // pdf.Cell(0, 10, fmt.Sprintf("Due Date: %s", "May 19th, 2023"))
// // pdf.SetXY(10, 65)
// // pdf.Cell(0, 10, "Billing Currency: USD ($)")

// // // Add items header
// // pdf.SetFont("Arial", "B", 12)
// // pdf.SetXY(10, 75)
// // pdf.Cell(40, 10, "Items")
// // pdf.Cell(40, 10, "Quantity")
// // pdf.Cell(40, 10, "Unit Price")
// // pdf.Cell(40, 10, "Amount")

// // // Add items
// // items := []struct {
// // 	Description string
// // 	Quantity    int
// // 	UnitPrice   float64
// // 	Amount      float64
// // }{
// // 	{"Email Marketing", 10, 1500, 15000},
// // 	{"Video looping effect", 6, 1110500, 6663000},
// // 	{"Graphic design for emails", 7, 2750, 19250},
// // 	{"Video looping effect", 6, 1110500, 6663000},
// // }

// // // Print items
// // pdf.SetFont("Arial", "", 12)
// // startY := 85
// // for _, item := range items {
// // 	pdf.SetXY(10, float64(startY))
// // 	pdf.Cell(40, 10, item.Description)
// // 	pdf.Cell(40, 10, fmt.Sprintf("%d", item.Quantity))
// // 	pdf.Cell(40, 10, fmt.Sprintf("$%.2f", item.UnitPrice))
// // 	pdf.Cell(40, 10, fmt.Sprintf("$%.2f", item.Amount))
// // 	startY += 10
// // }

// // // Add totals
// // pdf.SetFont("Arial", "B", 12)
// // pdf.SetXY(10, float64(startY+10))
// // pdf.Cell(40, 10, "Subtotal:")
// // pdf.Cell(40, 10, "$6,697,200.00")

// // pdf.SetXY(10, float64(startY+20))
// // pdf.Cell(40, 10, "Discount (2.5%):")
// // pdf.Cell(40, 10, "$167,430.00")

// // pdf.SetFont("Arial", "B", 16)
// // pdf.SetXY(10, float64(startY+30))
// // pdf.Cell(40, 10, "TOTAL AMOUNT DUE:")
// // pdf.Cell(40, 10, "$6,529,770.00")

// // // Add payment information
// // pdf.SetFont("Arial", "B", 12)
// // pdf.SetXY(10, float64(startY+50))
// // pdf.Cell(40, 10, "PAYMENT INFORMATION")

// // pdf.SetFont("Arial", "", 12)
// // pdf.SetXY(10, float64(startY+60))
// // pdf.Cell(40, 10, "Account Name: 1023902390")
// // pdf.SetXY(10, float64(startY+65))
// // pdf.Cell(40, 10, "Account Number: March 30th, 2023")
// // pdf.SetXY(10, float64(startY+70))
// // pdf.Cell(40, 10, "ACH Routing No: May 19th, 2023")
// // pdf.SetXY(10, float64(startY+75))
// // pdf.Cell(40, 10, "Bank Name: USD ($)")

// // // Output to PDF file
// // err := pdf.OutputFileAndClose("invoice.pdf")
// // if err != nil {
// // 	log.Fatal(err)
// // }

// // Buffer to hold the PDF
// // var buf bytes.Buffer

// // // Output the PDF to the buffer
// // err := pdf.Output(&buf)
// // if err != nil {
// // 	return nil, err
// // }
// // return buf.Bytes(), nil
// // fmt.Println(pdf)
// // base64PDF := base64.StdEncoding.EncodeToString(buf.Bytes())
// // fmt.Println("PDF Base64 Output:\n", base64PDF)
// // // Set headers to return PDF in the response
// // w.Header().Set("Content-Type", "application/pdf")
// // w.Header().Set("Content-Disposition", "attachment; filename=\"invoice.pdf\"")
// // w.WriteHeader(http.StatusOK)

// // // Write the PDF bytes to the response
// // _, err = w.Write(pdfData)
// // if err != nil {
// // 	log.Println("Failed to send PDF:", err)
// // }

// type InvoiceItem struct {
// 	Description string
// 	Quantity    int
// 	UnitPrice   float64
// 	Total       float64
// }

// type Invoice struct {
// 	SenderName    string
// 	SenderPhone   string
// 	SenderAddress string
// 	SenderEmail   string

// 	CustomerName  string
// 	CustomerPhone string
// 	CustomerEmail string

// 	InvoiceNumber string
// 	IssueDate     string
// 	DueDate       string
// 	Currency      string

// 	Items          []InvoiceItem
// 	Subtotal       float64
// 	DiscountRate   float64
// 	DiscountAmount float64
// 	TotalDue       float64

// 	AccountName   string
// 	AccountNumber string
// 	RoutingNumber string
// 	BankName      string
// 	BankAddress   string

// 	Note string
// }

// func generateInvoicePDF(invoice Invoice, filename string) ([]byte, error) {
// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	pdf.AddPage()

// 	// Set up fonts
// 	pdf.SetFont("Arial", "B", 16)
// 	pdf.SetTextColor(100, 100, 100)

// 	// Header
// 	pdf.SetFillColor(253, 233, 242)
// 	pdf.Rect(10, 10, 190, 50, "F")

// 	// Sender details
// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.SetXY(15, 15)
// 	pdf.Cell(100, 10, invoice.SenderName)

// 	pdf.SetFont("Arial", "", 10)
// 	pdf.SetXY(15, 22)
// 	pdf.Cell(100, 6, invoice.SenderPhone)
// 	pdf.SetXY(15, 28)
// 	pdf.Cell(100, 6, invoice.SenderAddress)
// 	pdf.SetXY(15, 34)
// 	pdf.Cell(100, 6, invoice.SenderEmail)

// 	// Customer details
// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.SetXY(120, 15)
// 	pdf.Cell(80, 10, "CUSTOMER")

// 	pdf.SetFont("Arial", "", 10)
// 	pdf.SetXY(120, 22)
// 	pdf.Cell(80, 6, invoice.CustomerName)
// 	pdf.SetXY(120, 28)
// 	pdf.Cell(80, 6, invoice.CustomerPhone)
// 	pdf.SetXY(120, 34)
// 	pdf.Cell(80, 6, invoice.CustomerEmail)

// 	// Invoice details
// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.SetXY(15, 45)
// 	pdf.Cell(50, 10, "INVOICE DETAILS")

// 	pdf.SetFont("Arial", "", 10)
// 	pdf.SetXY(15, 52)
// 	pdf.Cell(40, 6, "INVOICE NO")
// 	pdf.SetXY(55, 52)
// 	pdf.Cell(40, 6, invoice.InvoiceNumber)

// 	pdf.SetXY(95, 52)
// 	pdf.Cell(30, 6, "ISSUE DATE")
// 	pdf.SetXY(125, 52)
// 	pdf.Cell(30, 6, invoice.IssueDate)

// 	pdf.SetXY(155, 52)
// 	pdf.Cell(20, 6, "DUE DATE")
// 	pdf.SetXY(175, 52)
// 	pdf.Cell(25, 6, invoice.DueDate)

// 	// Items table
// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.SetXY(15, 70)
// 	pdf.Cell(100, 10, "Items")

// 	pdf.SetFont("Arial", "B", 10)
// 	pdf.SetXY(15, 80)
// 	pdf.Cell(80, 8, "Description")
// 	pdf.Cell(20, 8, "Quantity")
// 	pdf.Cell(40, 8, "Unit Price")
// 	pdf.Cell(40, 8, "Total")

// 	pdf.SetFont("Arial", "", 10)
// 	y := 88.0
// 	for _, item := range invoice.Items {
// 		pdf.SetXY(15, y)
// 		pdf.Cell(80, 6, item.Description)
// 		pdf.SetXY(95, y)
// 		pdf.Cell(20, 6, strconv.Itoa(item.Quantity))
// 		pdf.SetXY(115, y)
// 		pdf.Cell(40, 6, fmt.Sprintf("$%.2f", item.UnitPrice))
// 		pdf.SetXY(155, y)
// 		pdf.Cell(40, 6, fmt.Sprintf("$%.2f", item.Total))
// 		y += 6
// 	}

// 	// Totals
// 	y += 10
// 	pdf.SetFont("Arial", "B", 10)
// 	pdf.SetXY(135, y)
// 	pdf.Cell(30, 6, "SUBTOTAL")
// 	pdf.SetXY(165, y)
// 	pdf.Cell(30, 6, fmt.Sprintf("$%.2f", invoice.Subtotal))

// 	y += 6
// 	pdf.SetXY(135, y)
// 	pdf.Cell(30, 6, fmt.Sprintf("DISCOUNT (%.1f%%)", invoice.DiscountRate*100))
// 	pdf.SetXY(165, y)
// 	pdf.Cell(30, 6, fmt.Sprintf("$%.2f", invoice.DiscountAmount))

// 	y += 6
// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.SetXY(135, y)
// 	pdf.Cell(30, 6, "TOTAL AMOUNT DUE")
// 	pdf.SetXY(165, y)
// 	pdf.Cell(30, 6, fmt.Sprintf("$%.2f", invoice.TotalDue))

// 	// Payment information
// 	y += 15
// 	pdf.SetFillColor(240, 240, 240)
// 	pdf.Rect(15, y, 180, 30, "F")

// 	pdf.SetFont("Arial", "B", 10)
// 	pdf.SetXY(20, y+5)
// 	pdf.Cell(50, 6, "PAYMENT INFORMATION")

// 	pdf.SetFont("Arial", "", 8)
// 	pdf.SetXY(20, y+11)
// 	pdf.Cell(40, 5, "ACCOUNT NAME")
// 	pdf.SetXY(60, y+11)
// 	pdf.Cell(50, 5, invoice.AccountName)

// 	pdf.SetXY(110, y+11)
// 	pdf.Cell(40, 5, "ACCOUNT NUMBER")
// 	pdf.SetXY(150, y+11)
// 	pdf.Cell(40, 5, invoice.AccountNumber)

// 	pdf.SetXY(20, y+16)
// 	pdf.Cell(40, 5, "ACH ROUTING NO")
// 	pdf.SetXY(60, y+16)
// 	pdf.Cell(50, 5, invoice.RoutingNumber)

// 	pdf.SetXY(110, y+16)
// 	pdf.Cell(40, 5, "BANK NAME")
// 	pdf.SetXY(150, y+16)
// 	pdf.Cell(40, 5, invoice.BankName)

// 	pdf.SetXY(20, y+21)
// 	pdf.Cell(40, 5, "BANK ADDRESS")
// 	pdf.SetXY(60, y+21)
// 	pdf.Cell(130, 5, invoice.BankAddress)

// 	// Note
// 	pdf.SetFont("Arial", "I", 10)
// 	pdf.SetXY(15, y+40)
// 	pdf.Cell(180, 6, invoice.Note)

// 	var buf bytes.Buffer

// 	// Output the PDF to the buffer
// 	err := pdf.Output(&buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func main() {
// 	invoice := Invoice{
// 		SenderName:    "Fabulous Enterprise",
// 		SenderPhone:   "+386 989 271 3115",
// 		SenderAddress: "1331 Hart Ridge Road 48436 Gaines, MI",
// 		SenderEmail:   "info@fabulousenterise.co",

// 		CustomerName:  "Olaniyi Ojo Adewale",
// 		CustomerPhone: "+386 989 271 3115",
// 		CustomerEmail: "info@fabulousenterise.co",

// 		InvoiceNumber: "1023902390",
// 		IssueDate:     "March 30th, 2023",
// 		DueDate:       "May 19th, 2023",
// 		Currency:      "USD ($)",

// 		Items: []InvoiceItem{
// 			{"Email Marketing", 10, 1500, 15000},
// 			{"Video looping effect", 6, 1110500, 6663000},
// 			{"Graphic design for emails", 7, 2750, 19250},
// 			{"Video looping effect", 6, 1110500, 6663000},
// 		},
// 		Subtotal:       6697200.00,
// 		DiscountRate:   0.025,
// 		DiscountAmount: 167430.00,
// 		TotalDue:       6529770.00,

// 		AccountName:   "1023902390",
// 		AccountNumber: "March 30th, 2023",
// 		RoutingNumber: "May 19th, 2023",
// 		BankName:      "USD ($)",
// 		BankAddress:   "1023902390",

// 		Note: "Thank you for your patronage",
// 	}

// 	_, err := generateInvoicePDF(invoice, "invoice.pdf")
// 	if err != nil {
// 		fmt.Println("Error generating PDF:", err)
// 	} else {
// 		fmt.Println("PDF generated successfully")
// 	}
// }
