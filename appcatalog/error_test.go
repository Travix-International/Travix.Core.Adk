package appcatalog

import "testing"

func TestIsBadRequest(t *testing.T) {
	err := &catalogError{
		operation:  "test operation",
		statusCode: 400,
		response:   nil,
	}

	if err.isBadRequest() == true {
		t.Log("isBadRequest returns the expected value\n")
	} else {
		t.Fatal("isBadRequest was expecting to return `true`\n")
	}
}

func TestAuthenticationIssue401(t *testing.T) {
	err := &catalogError{
		operation:  "test operation",
		statusCode: 401,
		response:   nil,
	}

	if err.authenticationIssue() == true {
		t.Log("authenticationIssue returns the expected value\n")
	} else {
		t.Fatal("isBadRequest was expecting to return `true`\n")
	}
}

func TestAuthenticationIssue403(t *testing.T) {
	err := &catalogError{
		operation:  "test operation",
		statusCode: 403,
		response:   nil,
	}

	if err.authenticationIssue() == true {
		t.Log("authenticationIssue returns the expected value\n")
	} else {
		t.Fatal("isBadRequest was expecting to return `true`\n")
	}
}

func TestServerError(t *testing.T) {
	err := &catalogError{
		operation:  "test operation",
		statusCode: 500,
		response:   nil,
	}

	if err.serverError() == true {
		t.Log("serverError returns the expected value\n")
	} else {
		t.Fatal("serverError was expecting to return `true`\n")
	}
}

func TestCanRetry(t *testing.T) {
	err := &catalogError{
		operation:  "test operation",
		statusCode: 503,
		response:   nil,
	}

	if err.canRetry() == true {
		t.Log("canRetry returns the expected value\n")
	} else {
		t.Fatal("canRetry was expecting to return `true`\n")
	}
}

func TestCanRetryNot(t *testing.T) {
	err := &catalogError{
		operation:  "test operation",
		statusCode: 500,
		response:   nil,
	}

	if err.canRetry() == false {
		t.Log("canRetry returns the expected value\n")
	} else {
		t.Fatal("canRetry was expecting to return `false`\n")
	}
}
