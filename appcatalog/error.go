package appcatalog

import "fmt"

type catalogError struct {
	operation  string
	statusCode int
	response   *appCatalogResponse
}

func (err *catalogError) Error() string {
	return fmt.Sprintf("%s failed\n%s", err.operation, err.response.getDetails(err.statusCode))
}

func (err *catalogError) isBadRequest() bool {
	return err.statusCode == 400
}

func (err *catalogError) authenticationIssue() bool {
	return err.statusCode == 401 || err.statusCode == 403
}

func (err *catalogError) serverError() bool {
	return err.statusCode == 500 // SERVER ERROR
}

func (err *catalogError) canRetry() bool {
	return !err.isBadRequest() && !err.serverError() && !err.authenticationIssue()
}
