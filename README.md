````
import "github.com/ZachBergh/gocsv"

...
// fileName, path string
// header []string, if len(header) <= 0 , will create default header
// data []interface{}
err := gocsv.NewCsvFile(fileName, path, header, data)
````