package utils

import (
	"encoding/csv"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"sort"
	"strconv"
)

type CsvFile struct {
	FileName string
	Path     string
	Header   []string
	data     []interface{}
}

func (this *CsvFile) CreateCsvFile() error {

	var err error

	f, err := os.Create(this.Path + this.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)

	for k, v := range this.data {
		if k == 0 {
			if len(this.Header) <= 0 || this.Header == nil {
				this.Header = this.CreateHeader(v, "")
			}
			w.Write(this.Header)
		}
		row := this.CreateRow(v)
		w.Write(row)
		w.Flush()
	}
	return err
}

func (this *CsvFile) CreateHeader(data interface{}, prefix string) (result []string) {
	if data == nil {
		return
	}
	if prefix != "" {
		prefix = prefix + "."
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
				result = CombineStringSlice(result, this.CreateHeader(iv, k+"."+strconv.Itoa(ik)))
			}
		case reflect.Struct, reflect.Map:
			result = CombineStringSlice(result, this.CreateHeader(mapV[k], k))
		default:
			result = append(result, prefix+k)
		}
	}
	return
}

func (this *CsvFile) CreateRow(data interface{}) (result []string) {
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
			result = CombineStringSlice(result, this.CreateRow(mapV[k]))
		case reflect.Bool:
			boolVal := "false"
			if vtype.Bool() {
				boolVal = "true"
			}
			result = append(result, boolVal)
		case reflect.Slice, reflect.Array:
			for _, iv := range mapV[k].([]interface{}) {
				result = CombineStringSlice(result, this.CreateRow(iv))
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
