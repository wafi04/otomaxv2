package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	validate := validator.New()
	
	// Register custom validations
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("operator", validateOperator)
	validate.RegisterValidation("game_type", validateGameType)
	validate.RegisterValidation("provider", validateProvider)
	validate.RegisterValidation("amount", validateAmount)
	
	return &Validator{validate: validate}
}

// ValidateStruct validates a struct
func (v *Validator) ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError
	
	err := v.validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Message: getErrorMessage(err),
			})
		}
	}
	
	return errors
}

// Custom validation functions
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	
	// Remove all non-digit characters
	phone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	// Check Indonesian phone number patterns
	patterns := []string{
		`^08[0-9]{8,11}$`,     // 08xxxxxxxxxx
		`^628[0-9]{8,11}$`,    // 628xxxxxxxxxx
		`^\+628[0-9]{8,11}$`,  // +628xxxxxxxxxx
	}
	
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, phone)
		if matched {
			return true
		}
	}
	
	return false
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	

	if len(password) < 8 {
		return false
	}
	
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func validateOperator(fl validator.FieldLevel) bool {
	operator := strings.ToLower(fl.Field().String())
	validOperators := []string{"telkomsel", "indosat", "xl", "axis", "three", "smartfren", "by.u"}
	
	for _, valid := range validOperators {
		if operator == valid {
			return true
		}
	}
	
	return false
}

func validateGameType(fl validator.FieldLevel) bool {
	gameType := strings.ToLower(fl.Field().String())
	validGameTypes := []string{"mobile_legends", "pubg", "free_fire", "valorant", "steam", "garena", "codashop"}
	
	for _, valid := range validGameTypes {
		if gameType == valid {
			return true
		}
	}
	
	return false
}

func validateProvider(fl validator.FieldLevel) bool {
	provider := strings.ToLower(fl.Field().String())
	validProviders := []string{"midtrans", "xendit", "gopay", "ovo", "dana", "bank_transfer", "virtual_account"}
	
	for _, valid := range validProviders {
		if provider == valid {
			return true
		}
	}
	
	return false
}

func validateAmount(fl validator.FieldLevel) bool {
	amount := fl.Field().Float()
	return amount >= 1000 && amount <= 10000000 // Min 1k, Max 10M
}


func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "phone":
		return fmt.Sprintf("%s must be a valid Indonesian phone number", err.Field())
	case "password":
		return fmt.Sprintf("%s must contain at least 8 characters with uppercase, lowercase, digit, and special character", err.Field())
	case "operator":
		return fmt.Sprintf("%s must be a valid operator", err.Field())
	case "game_type":
		return fmt.Sprintf("%s must be a valid game type", err.Field())
	case "provider":
		return fmt.Sprintf("%s must be a valid payment provider", err.Field())
	case "amount":
		return fmt.Sprintf("%s must be between 1,000 and 10,000,000", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// Utility validation functions
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func IsValidPhoneNumber(phone string) bool {
	phone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	patterns := []string{
		`^08[0-9]{8,11}$`,
		`^628[0-9]{8,11}$`,
		`^\+628[0-9]{8,11}$`,
	}
	
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, phone)
		if matched {
			return true
		}
	}
	
	return false
}

func NormalizePhoneNumber(phone string) string {
	phone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	if strings.HasPrefix(phone, "08") {
		return "628" + phone[2:]
	} else if strings.HasPrefix(phone, "+628") {
		return phone[1:]
	} else if strings.HasPrefix(phone, "628") {
		return phone
	}
	
	return phone
}