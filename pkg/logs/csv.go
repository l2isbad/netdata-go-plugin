package logs

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

type (
	CSVConfig struct {
		Delimiter        rune                             `yaml:"delimiter"`
		TrimLeadingSpace bool                             `yaml:"trim_leading_space"`
		Format           string                           `yaml:"format"`
		CheckField       func(string) (string, int, bool) `yaml:"-"`
	}

	CSVParser struct {
		config CSVConfig
		reader *csv.Reader
		format *csvFormat
	}

	csvFormat struct {
		raw      string
		maxIndex int
		fields   []csvField
	}

	csvField struct {
		name string
		idx  int
	}
)

func NewCSVParser(config CSVConfig, in io.Reader) (*CSVParser, error) {
	if config.Format == "" {
		return nil, errors.New("empty csv format")
	}

	format, err := newCSVFormat(config)
	if err != nil {
		return nil, fmt.Errorf("bad csv format '%s': %v", config.Format, err)
	}

	p := &CSVParser{
		config: config,
		reader: newCSVReader(in, config),
		format: format,
	}
	return p, nil
}

func (p *CSVParser) ReadLine(line LogLine) error {
	record, err := p.reader.Read()
	if err != nil {
		if isCSVParseError(err) {
			return &ParseError{msg: fmt.Sprintf("csv parse: %v", err), err: err}
		}
		return err
	}
	return p.format.parse(record, line)
}

func (p *CSVParser) Parse(row []byte, line LogLine) error {
	r := newCSVReader(bytes.NewBuffer(row), p.config)
	record, err := r.Read()
	if err != nil {
		if isCSVParseError(err) {
			return &ParseError{msg: fmt.Sprintf("csv parse: %v", err), err: err}
		}
		return err
	}
	return p.format.parse(record, line)
}

func (p CSVParser) Info() string {
	return fmt.Sprintf("csv: %s", p.format.raw)
}

func (f *csvFormat) parse(record []string, line LogLine) error {
	if len(record) <= f.maxIndex {
		return &ParseError{msg: "csv parse: unmatched line"}
	}

	for _, v := range f.fields {
		if err := line.Assign(v.name, record[v.idx]); err != nil {
			return &ParseError{msg: fmt.Sprintf("csv parse: %v", err), err: err}
		}
	}
	return nil
}

func newCSVReader(in io.Reader, config CSVConfig) *csv.Reader {
	r := csv.NewReader(in)
	r.Comma = config.Delimiter
	r.TrimLeadingSpace = config.TrimLeadingSpace
	r.ReuseRecord = true
	r.FieldsPerRecord = -1
	return r
}

func newCSVFormat(config CSVConfig) (*csvFormat, error) {
	r := csv.NewReader(strings.NewReader(config.Format))
	r.Comma = config.Delimiter
	r.TrimLeadingSpace = config.TrimLeadingSpace

	record, err := r.Read()
	if err != nil {
		return nil, err
	}

	fields, err := createCSVFields(record, config.CheckField)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return nil, errors.New("zero fields")
	}

	format := &csvFormat{
		raw:      config.Format,
		maxIndex: fields[len(fields)-1].idx,
		fields:   fields,
	}
	return format, nil
}

func createCSVFields(format []string, check func(string) (string, int, bool)) ([]csvField, error) {
	if check == nil {
		check = checkCSVFormatField
	}
	var fields []csvField
	var offset int
	seen := make(map[string]bool)

	for i, name := range format {
		name = strings.Trim(name, `"`)

		name, addOffset, valid := check(name)
		offset += addOffset
		if !valid {
			continue
		}
		if seen[name] {
			return nil, fmt.Errorf("duplicate field: %s", name)
		}
		seen[name] = true

		idx := i + offset
		fields = append(fields, csvField{name, idx})
	}
	return fields, nil
}

func isCSVParseError(err error) bool {
	return errors.Is(err, csv.ErrBareQuote) || errors.Is(err, csv.ErrFieldCount) || errors.Is(err, csv.ErrQuote)
}

func checkCSVFormatField(name string) (newName string, offset int, valid bool) {
	if len(name) < 2 || !strings.HasPrefix(name, "$") {
		return "", 0, false
	}
	return name, 0, true
}