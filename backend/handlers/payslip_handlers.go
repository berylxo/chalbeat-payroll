package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/services"
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Payslip-%s.pdf", emp.EmployeeID))
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

		fileName := fmt.Sprintf("Payslip-%s.pdf", emp.EmployeeID)
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
	html, err := h.renderPayslipHTML(emp, payroll)
	if err != nil {
		return nil, err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(html))
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	return bytes.NewBuffer(pdfg.Bytes()), nil
}

func (h *PayslipHandler) renderPayslipHTML(emp models.Employee, payroll models.PayrollResult) (string, error) {
	templateDir, err := findTemplatesDir()
	if err != nil {
		return "", err
	}

	templatePath := filepath.Join(templateDir, "payslip.html")
	logoPath := filepath.Join(templateDir, "logo.png")

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	logoBytes, err := os.ReadFile(logoPath)
	if err != nil {
		return "", err
	}
	logoDataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(logoBytes)

	type employeeView struct {
		Name       string
		ID         string
		EmployeeID string
		KraPin     string
		Department string
	}
	type deductionItemView struct {
		Code        string
		Description string
		Amount      string // Changed to formatted string
	}
	type earningsView struct {
		BasicSalary string // Changed to string
		GrossPay    string // Changed to string
		Allowances  []deductionItemView
	}
	type totalsView struct {
		GrossPay        string // Changed to string
		TotalDeductions string // Changed to string
		NetPay          string // Changed to string
	}
	type viewData struct {
		CompanyLogo string
		LogoPath    string
		Period      string
		Employee    employeeView
		Earnings    earningsView
		Deductions  []deductionItemView
		Totals      totalsView
		GeneratedAt string
	}

	var formattedDeductions []deductionItemView
	var totalDeductionsCents int64

	for _, deduction := range payroll.Deductions {
		totalDeductionsCents += deduction.Amount
		formattedDeductions = append(formattedDeductions, deductionItemView{
			Code:        deduction.Code,
			Description: deduction.Description,
			Amount:      engine.FormatCentsToKES(deduction.Amount),
		})
	}

	// Calculate and format values safely using integer cent foundations
	data := viewData{
		CompanyLogo: logoDataURI,
		LogoPath:    "file://" + logoPath,
		Period:      time.Now().Format("January 2006"),
		Employee: employeeView{
			Name:       emp.Name,
			ID:         emp.EmployeeID,
			EmployeeID: emp.EmployeeID,
			KraPin:     emp.KraPin,
			Department: emp.Position,
		},
		Earnings: earningsView{
			BasicSalary: engine.FormatCentsToKES(int64(math.Round(emp.BasicPay * 100))),
			GrossPay:    engine.FormatCentsToKES(payroll.GrossPay),
			Allowances:  nil,
		},
		Deductions: formattedDeductions,
		Totals: totalsView{
			GrossPay:        engine.FormatCentsToKES(payroll.GrossPay),
			TotalDeductions: engine.FormatCentsToKES(totalDeductionsCents),
			NetPay:          engine.FormatCentsToKES(payroll.NetPay),
		},
		GeneratedAt: time.Now().Format("02 Jan 2006 15:04"),
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func findTemplatesDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	tryPath := filepath.Join(cwd, "templates")
	if fi, err := os.Stat(tryPath); err == nil && fi.IsDir() {
		return tryPath, nil
	}

	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(executable), "templates"), nil
}

func (h *PayslipHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, message)))
}
