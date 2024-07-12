package bbva

import (
	"errors"
	"fmt"
	"strings"

	excelize "github.com/xuri/excelize/v2"
)

const (
	SheetName = "Informe BBVA"

	headerConcept = "Concepto"
)

var ErrCouldNotFindHeader error = errors.New("could not find header")

type XLSXItem struct {
	Cells [][2]string
}

type XLSX struct {
	Rows [][]string

	HeaderRowIndex    int
	HeaderColumnIndex int

	// HeaderKeys maps a column index to its header label.
	HeaderKeys map[int]string

	Items []XLSXItem
}

func (x *XLSX) findHeaderRowIndex() error {
	search := strings.ToLower(headerConcept)
	for i, row := range x.Rows {
		for _, cell := range row {
			if strings.ToLower(cell) == search {
				x.HeaderRowIndex = i

				return nil
			}
		}
	}

	x.HeaderRowIndex = -1

	return ErrCouldNotFindHeader
}

func (x *XLSX) readHeaderKeys() error {
	if x.HeaderRowIndex == -1 {
		return ErrCouldNotFindHeader
	}

	columnIndex := -1
	row := x.Rows[x.HeaderRowIndex]
	headers := map[int]string{}
	for i, cell := range row {
		cell = strings.TrimSpace(cell)
		if cell == "" {
			if columnIndex == -1 {
				// Empty space to the left.
				continue
			} else {
				// Empty space to the right.
				//
				// Stop reading headers.
				break
			}
		}

		if columnIndex == -1 {
			columnIndex = i

			x.HeaderColumnIndex = columnIndex
		}

		headers[i] = cell
	}

	x.HeaderKeys = headers

	return nil
}

func (x *XLSX) parseItems() error {
	var items []XLSXItem

	columnOffset := x.HeaderColumnIndex
	rowOffset := x.HeaderRowIndex
	for i := rowOffset + 1; i < len(x.Rows); i++ {
		row := x.Rows[i]
		firstValue := strings.TrimSpace(row[columnOffset])
		if firstValue == "" {
			// Assume the value under the leftmost header always has a value,
			// since it should be date.
			break
		}

		var cells [][2]string

		for j := columnOffset; j < len(row); j++ {
			key, ok := x.HeaderKeys[j]
			if !ok {
				break
			}

			value := strings.TrimSpace(row[j])

			cell := [2]string{
				key,
				value,
			}

			cells = append(cells, cell)
		}

		item := XLSXItem{
			Cells: cells,
		}

		items = append(items, item)
	}

	x.Items = items

	return nil
}

func ParseXLSXFile(file string) (*XLSX, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(SheetName)
	if err != nil {
		return nil, fmt.Errorf("could not get rows: %w", err)
	}

	xlsx := &XLSX{
		Rows: rows,
	}

	err = xlsx.findHeaderRowIndex()
	if err != nil {
		return nil, fmt.Errorf("could not find sheet header: %w", err)
	}

	err = xlsx.readHeaderKeys()
	if err != nil {
		return nil, fmt.Errorf("could not read sheet header: %w", err)
	}

	err = xlsx.parseItems()
	if err != nil {
		return nil, fmt.Errorf("could not parse items: %w", err)
	}

	return xlsx, nil
}
