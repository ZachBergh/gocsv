package gocsv

import (
	"encoding/csv"
	"os"
	"reflect"
	"sort"
	"strconv"
)

// CSV struction
type CsvFile struct {
	FileName string
	Path     string
	Header   []string
	data     []interface{}
}

func (c *CsvFile) CreateCsvFile() error {

	var err error

	f, err := os.Create(c.Path + c.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)

	for k, v := range c.data {
		if k == 0 {
			if len(c.Header) <= 0 || c.Header == nil {
				c.Header = c.CreateHeader(v, "")
			}
			w.Write(c.Header)
		}
		row := c.CreateRow(v)
		w.Write(row)
		w.Flush()
	}
	return err
}

func (c *CsvFile) CreateHeader(data interface{}, prefix string) (result []string) {
	if data == nil {
		return
	}
	if prefix != "" {
		prefix = prefix + "."
	}
	if reflect.ValueOf(data).Kind() == reflect.String {
		result = append(result, reflect.ValueOf(data).String())
		return
	}
	mapV := data.(map[string]interface{})
	var keys []string
	for k := range mapV {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if k == "_id" {
			continue
		}

		vtype := reflect.ValueOf(mapV[k])

		switch vtype.Kind() {
		case reflect.Slice, reflect.Array:
			for ik, iv := range mapV[k].([]interface{}) {
				result = CombineStringSlice(result, c.CreateHeader(iv, k+"."+strconv.Itoa(ik)))
			}
		case reflect.Struct, reflect.Map:
			result = CombineStringSlice(result, c.CreateHeader(mapV[k], k))
		default:
			result = append(result, prefix+k)
		}
	}
	return
}

func (c *CsvFile) CreateRow(data interface{}) (result []string) {
	if data == nil {
		return
	}
	if reflect.ValueOf(data).Kind() == reflect.String {
		result = append(result, reflect.ValueOf(data).String())
		return
	}
	mapV := data.(map[string]interface{})

	var keys []string
	for k := range mapV {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if k == "_id" {
			continue
		}

		vtype := reflect.ValueOf(mapV[k])

		switch vtype.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result = append(result, strconv.Itoa(int(vtype.Int())))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = append(result, strconv.Itoa(int(vtype.Uint())))
		case reflect.Float32, reflect.Float64:
			result = append(result, strconv.Itoa(int(vtype.Float())))
		case reflect.String:
			result = append(result, vtype.String())
		case reflect.Struct, reflect.Map:
			result = CombineStringSlice(result, c.CreateRow(mapV[k]))
		case reflect.Bool:
			boolVal := "false"
			if vtype.Bool() {
				boolVal = "true"
			}
			result = append(result, boolVal)
		case reflect.Slice, reflect.Array:
			for _, iv := range mapV[k].([]interface{}) {
				result = CombineStringSlice(result, c.CreateRow(iv))
			}
		}
	}
	return
}

func CombineStringSlice(a []string, b []string) (c []string) {
	c = make([]string, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return
}

func NewCsvFile(filename, path string, header []string, data []interface{}) (err error) {
	csv := &CsvFile{filename, path, header, data}
	err = csv.CreateCsvFile()
	return
}
