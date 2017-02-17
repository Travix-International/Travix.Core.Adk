package appcatalog

import "fmt"

type catalogError struct {
	operation  string
	statusCode int
}

func (err *catalogError) Error() string {
	return fmt.Sprintf("%s failed, App Catalog returned status code %v\n", err.operation, err.statusCode)
}

func (err *catalogError) canRetry() bool {
	return err.statusCode == 500 // SERVER ERROR
}
