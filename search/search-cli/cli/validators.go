package cli

import (
    "net/url"
    "errors"
    "fmt"
)

func ValidateMandatory(value string, paramName string, zeroValue string) error {
    if value == zeroValue {
        return errors.New(fmt.Sprintf("%s is missing", paramName))
    }
    return nil
}

func ValidateNotEmpty(value []string, paramName string) error {
    if len(value)==0 {
        return errors.New(fmt.Sprintf("%s is empty, at least one value must be provided", paramName))
    }
    return nil
}

func ValidateURL(urlStr string, paramName string) error {
    _, err := url.ParseRequestURI(urlStr)
    if err != nil {
        return errors.New(fmt.Sprintf("%s is invalid, must be a valid URL", paramName))
    }
    return nil
}
