package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/services"
	"github.com/jung-kurt/gofpdf"
)

type PayslipHandler struct {
	EmployeeService *services.EmployeeService
	PayrollService  *services.PayrollService
}

func (h *PayslipHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/employees/")
	id = strings.TrimSuffix(id, "/payslip")
	id = strings.TrimSuffix(id, "/")

	emp, ok := h.EmployeeService.Get(id)
	if !ok {
		h.writeError(w, http.StatusNotFound, "employee not found")
		return
	}

	payroll := h.PayrollService.Run(emp.BasicPay, nil)
	buffer, err := h.renderPayslipPDF(emp, payroll)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Payslip-%s.pdf", emp.ID))
	w.Write(buffer.Bytes())
}

func (h *PayslipHandler) GenerateBulkZip(w http.ResponseWriter, r *http.Request) {
	employees := h.EmployeeService.List()
	if len(employees) == 0 {
		h.writeError(w, http.StatusNotFound, "no employees available")
		return
	}

	buffer := &bytes.Buffer{}
	zipWriter := zip.NewWriter(buffer)

	for _, emp := range employees {
		payroll := h.PayrollService.Run(emp.BasicPay, nil)
		pdfBytes, err := h.renderPayslipPDF(emp, payroll)
		if err != nil {
			zipWriter.Close()
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fileName := fmt.Sprintf("Payslip-%s.pdf", emp.ID)
		zipFile, err := zipWriter.Create(fileName)
		if err != nil {
			zipWriter.Close()
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		_, _ = zipFile.Write(pdfBytes.Bytes())
	}

	zipWriter.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=Payslips.zip")
	w.Write(buffer.Bytes())
}

func (h *PayslipHandler) renderPayslipPDF(emp models.Employee, payroll models.PayrollResult) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Payslip", false)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Payslip")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, "Employee ID:")
	pdf.Cell(0, 8, emp.ID)
	pdf.Ln(8)
	pdf.Cell(40, 8, "Name:")
	pdf.Cell(0, 8, emp.Name)
	pdf.Ln(8)
	pdf.Cell(40, 8, "KRA PIN:")
	pdf.Cell(0, 8, emp.KraPin)
	pdf.Ln(8)
	pdf.Cell(40, 8, "Position:")
	pdf.Cell(0, 8, emp.Position)
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "Summary")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Gross Pay:")
	pdf.Cell(0, 8, fmt.Sprintf("%.2f", payroll.GrossPay))
	pdf.Ln(8)
	pdf.Cell(60, 8, "Net Pay:")
	pdf.Cell(0, 8, fmt.Sprintf("%.2f", payroll.NetPay))
	pdf.Ln(12)

	pdf.Cell(0, 8, "Deductions")
	pdf.Ln(8)
	for _, deduction := range payroll.Deductions {
		pdf.Cell(60, 8, deduction.Description+":")
		pdf.Cell(0, 8, fmt.Sprintf("%.2f", deduction.Amount))
		pdf.Ln(8)
	}

	buffer := &bytes.Buffer{}
	if err := pdf.Output(buffer); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (h *PayslipHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, message)))
}
