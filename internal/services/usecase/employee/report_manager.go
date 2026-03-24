package employee

import (
	"fmt"
	"os"
	"time"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type reportManager struct {
	report *excelize.File
}

func (r reportManager) createDeletedEmployeeReport(worker entity.Employee) (string, error) {
	const mainSheetName = "Deleted Employee common information"

	r.report = excelize.NewFile()
	defer r.report.Close()

	_, err := r.report.NewSheet(mainSheetName)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrReportCreating, err)
	}

	r.setDeletedReportTemplate(mainSheetName)
	r.setDeletedReportEmployeeValues(worker, mainSheetName)

	path := fmt.Sprintf("./.tmp/%s.xlsx", uuid.New().String())
	if err = r.report.SaveAs(path); err != nil {
		return "", fmt.Errorf("%w: %s", ErrReportCreating, err)
	}

	go func() {
		time.Sleep(30 * time.Minute)
		os.Remove(path)
	}()
	return path, nil
}

func (r reportManager) setDeletedReportTemplate(sheetName string) {
	r.report.SetCellValue(sheetName, "A1", "TIN")
	r.report.SetCellValue(sheetName, "B1", "Snils")
	r.report.SetCellValue(sheetName, "C1", "Passport")
	r.report.SetCellValue(sheetName, "D1", "Phone")
	r.report.SetCellValue(sheetName, "E1", "First Name")
	r.report.SetCellValue(sheetName, "F1", "Last Name")
	r.report.SetCellValue(sheetName, "G1", "Patronymic")
	r.report.SetCellValue(sheetName, "H1", "Address")
	r.report.SetCellValue(sheetName, "I1", "Title")
	r.report.SetCellValue(sheetName, "J1", "Citizenship")
	r.report.SetCellValue(sheetName, "K1", "Hiring Date")
	r.report.SetCellValue(sheetName, "L1", "Unit Type")
	r.report.SetCellValue(sheetName, "M1", "Unit Name")
	r.report.SetCellValue(sheetName, "N1", "Education")
	r.report.SetCellValue(sheetName, "O1", "Salary")
}

func (r reportManager) setDeletedReportEmployeeValues(worker entity.Employee, sheetName string) {
	r.report.SetCellValue(sheetName, "A2", worker.TinNum)
	r.report.SetCellValue(sheetName, "B2", worker.Snils)
	r.report.SetCellValue(sheetName, "C2", worker.PassportData)
	r.report.SetCellValue(sheetName, "D2", worker.PhoneNum)
	r.report.SetCellValue(sheetName, "E2", worker.FirstName)
	r.report.SetCellValue(sheetName, "F2", worker.LastName)
	r.report.SetCellValue(sheetName, "G2", worker.Patronymic)
	r.report.SetCellValue(sheetName, "H2", worker.Address)
	r.report.SetCellValue(sheetName, "I2", worker.Title.Name)
	r.report.SetCellValue(sheetName, "J2", worker.Citizenship.Name)
	r.report.SetCellValue(sheetName, "K2", worker.HiringDate)
	r.report.SetCellValue(sheetName, "L2", worker.Unit.Type)
	r.report.SetCellValue(sheetName, "M2", worker.Unit.Name)
	r.report.SetCellValue(sheetName, "N2", worker.Education)
	r.report.SetCellValue(sheetName, "O2", worker.Salary)
}
