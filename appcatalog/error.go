package appcatalog

import "fmt"

type catalogError struct {
	operation  string
	statusCode int
}

func (err *catalogError) Error() string {
	return fmt.Sprintf("%s failed, App Catalog returned status code %v\n", err.operation, err.statusCode)
}

func (err *catalogError) versionAlreadySubmitted() bool {
	return err.statusCode == 400
}

func (err *catalogError) authenticationIssue() bool {
	return err.statusCode == 401 || err.statusCode == 403
}

func (err *catalogError) serverError() bool {
	return err.statusCode == 500 // SERVER ERROR
}

func (err *catalogError) canRetry() bool {
	return !err.versionAlreadySubmitted() && !err.serverError() && !err.authenticationIssue()
}
