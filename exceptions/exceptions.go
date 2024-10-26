package exceptions

type ERROR_TYPE string

type PublicError struct {
	Type        ERROR_TYPE `json:"key"`
	Description string     `json:"description"`
}

type PrivateError struct {
	Type        ERROR_TYPE `json:"key"`
	Description string     `json:"description"`
}

var errors []Error

type Error struct {
	Private ErrorContent `json:"private_error"`
	Public  ErrorContent `json:"public_error"`
}

type ErrorContent struct {
	Description string     `json:"description"`
	Type        ERROR_TYPE `json:"key"`
}

func AddError(public_description string, public_type ERROR_TYPE, private_description string, private_type ERROR_TYPE) {
	error := Error{}
	error.Public = ErrorContent{
		Description: public_description,
		Type:        public_type,
	}
	error.Private = ErrorContent{
		Description: private_description,
		Type:        private_type,
	}
	errors = append(errors, error)
}

func GetErrors() []Error {
	return errors
}

const (
	GENERIC_ERROR                        ERROR_TYPE = "GenericException"
	PREVIOUS_STAGE_FAILED                ERROR_TYPE = "PreviousStageFailed"
	FAILED_TO_READ_PREVIOUS_STAGE_OUTPUT ERROR_TYPE = "FailedToReadPreviousStageOutput"
	UNSUPPORTED_LANGUAGE_REQUESTED       ERROR_TYPE = "UnsupportedLanguageRequested"
)
