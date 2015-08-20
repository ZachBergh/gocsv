### Gocsv

A data to csv file productor tool.

##### Installation
````
go get github.com/ZachBergh/gocsv
````

##### Usage

````
import "github.com/ZachBergh/gocsv"

/* 
Params:
    fileName string: csv file's name 
    path string: csv file's path
    header []string: csv file's header
        if len(header) <= 0 , will create default header
    data []interface{}: source data
*/

err := gocsv.NewCsvFile(fileName, path, header, data)
if err != nil {
    return err
}
````