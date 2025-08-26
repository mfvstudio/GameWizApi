package helpers

import (
	"log"
	"regexp"
	"time"

	"github.com/mfvstudio/gamewizapi/cmd/gen"
)

// At least one of lower, upper case
// At least one digit
// At least one special char from the allowed list
// At least 8 chars long
var passwordPattern string = "^[A-Za-z\\d@$!%*?&_+]{8,}$"
var atleastOneLowercase string = ".*[a-z]"
var atleastOneUpperCase string = ".*[A-Z]"
var atLeastOneSpecial string = ".*[@$!%*?&_+]"
var atLeastOneDigit string = "\\d*"

// var emailPattern string = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+[a-z]{3}$"
var userNamePattern string = "^[^<>'\";]{3,12}$"

// Checks for blank inputs
func IsSignUpFormAccepted(form *gen.CreateAccount) (bool, error) {
	if form.Username == "" ||
		form.Email == "" ||
		form.Birthdate == 0 ||
		form.Password == "" {
		return false, nil
	}
	if !IsAtLeastEighteen(form) {
		log.Printf("user does not meet age requirement.")
		return false, nil
	}
	res, err := IsRegexAccepted(userNamePattern, form.Username)
	if err != nil {
		return false, err
	}
	if !res {
		log.Printf("Username is in incorrect format")
		return false, nil
	}
	patterns := [7]string{passwordPattern, atLeastOneDigit, atLeastOneSpecial, atleastOneLowercase, atleastOneUpperCase}
	for i := range len(patterns) {
		res, err := IsRegexAccepted(patterns[i], form.Password)
		if err != nil || !res {
			log.Printf("Auth failed for no.%v. Result is %v", i, res)
			return false, err
		}
	}
	return true, nil
}

// Users must be 18 or older
func IsAtLeastEighteen(form *gen.CreateAccount) bool {
	bdate := time.Unix(form.Birthdate, 0)
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return bdate.Before(eighteenYearsAgo)
}

func IsRegexAccepted(reg string, p string) (bool, error) {
	pattern, err := regexp.Compile(reg)
	if err != nil {
		log.Printf("Error while compiling regex pattern.")
		return false, err
	}
	match := pattern.Match([]byte(p))
	return match, nil
}
