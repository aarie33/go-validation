package govalidation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidation(t *testing.T) {
	var validate *validator.Validate = validator.New()
	if validate == nil {
		t.Error("validator is nil")
	}
}

func TestValidationVariable(t *testing.T) {
	validate := validator.New()
	user := "joko"

	err := validate.Var(user, "required")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestValidateTwoVariables(t *testing.T) {
	validate := validator.New()

	password := "rahasia"
	confirmPassword := "rahasia"

	err := validate.VarWithValue(password, confirmPassword, "eqfield")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestMultipleTag(t *testing.T) {
	validate := validator.New()
	user := "123343"

	err := validate.Var(user, "required,numeric,min=5,max=10")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestStruct(t *testing.T) {
	type LoginRequest struct {
		Username string `validate:"required,email"`
		Password string `validate:"required,min=5"`
	}

	validate := validator.New()
	loginRequest := LoginRequest{
		Username: "santoso@gmail.com",
		Password: "rahasia",
	}

	err := validate.Struct(loginRequest)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestValidationErrors(t *testing.T) {
	type LoginRequest struct {
		Username string `validate:"required,email"`
		Password string `validate:"required,min=5"`
	}

	validate := validator.New()
	loginRequest := LoginRequest{
		Username: "santoso@gmail.com",
		Password: "rahasia",
	}

	err := validate.Struct(loginRequest)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestStructCrossField(t *testing.T) {
	type RegisterRequest struct {
		Username        string `validate:"required,email"`
		Password        string `validate:"required,min=5"`
		ConfirmPassword string `validate:"required,eqfield=Password"`
	}

	validate := validator.New()
	registerRequest := RegisterRequest{
		Username:        "santoso@gmail.com",
		Password:        "rahasia",
		ConfirmPassword: "rahasisa",
	}

	err := validate.Struct(registerRequest)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestNestedStruct(t *testing.T) {
	type Address struct {
		City    string `validate:"required,max=255"`
		Country string `validate:"required,max=255"`
	}

	type User struct {
		Id      int     `validate:"required"`
		Name    string  `validate:"required,max=255"`
		Address Address `validate:"required"`
	}

	validate := validator.New()
	request := User{
		Id:   1,
		Name: "",
		Address: Address{
			City:    "Jakarta",
			Country: "Indonesia",
		},
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestCollection(t *testing.T) {
	type Address struct {
		City    string `validate:"required,max=255"`
		Country string `validate:"required,max=255"`
	}

	type User struct {
		Id        int       `validate:"required"`
		Name      string    `validate:"required,max=255"`
		Addresses []Address `validate:"required,dive"`
	}

	validate := validator.New()
	request := User{
		Id:   1,
		Name: "Joko",
		Addresses: []Address{
			{
				City:    "Jakarta",
				Country: "Indonesia",
			},
			{
				City:    "Bandung",
				Country: "",
			},
		},
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestBasicCollection(t *testing.T) {
	type Address struct {
		City    string `validate:"required,max=255"`
		Country string `validate:"required,max=255"`
	}

	type User struct {
		Id        int       `validate:"required"`
		Name      string    `validate:"required,max=255"`
		Addresses []Address `validate:"required,dive"`
		Hobbies   []string  `validate:"required,dive,required,min=1,oneof=football basketball"`
	}

	validate := validator.New()
	request := User{
		Id:   1,
		Name: "Joko",
		Addresses: []Address{
			{
				City:    "Jakarta",
				Country: "Indonesia",
			},
			{
				City:    "Bandung",
				Country: "Indonesia",
			},
		},
		Hobbies: []string{
			"football",
			"",
		},
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestMapCollection(t *testing.T) {
	type Address struct {
		City    string `validate:"required,max=255"`
		Country string `validate:"required,max=255"`
	}

	type School struct {
		Name string `validate:"required,max=255"`
	}

	type User struct {
		Id        int               `validate:"required"`
		Name      string            `validate:"required,max=255"`
		Addresses []Address         `validate:"required,dive"`
		Hobbies   []string          `validate:"required,dive,required,min=1,oneof=football basketball"`
		Schools   map[string]School `validate:"required,dive,keys,required,min=1,endkeys,dive"`
	}

	validate := validator.New()
	request := User{
		Id:   1,
		Name: "Joko",
		Addresses: []Address{
			{
				City:    "Jakarta",
				Country: "Indonesia",
			},
			{
				City:    "Bandung",
				Country: "Indonesia",
			},
		},
		Hobbies: []string{
			"football",
		},
		Schools: map[string]School{
			"SD": {
				Name: "",
			},
			"SMP": {
				Name: "SMP Negeri 1",
			},
		},
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestBasicMapCollection(t *testing.T) {
	type Address struct {
		City    string `validate:"required,max=255"`
		Country string `validate:"required,max=255"`
	}

	type School struct {
		Name string `validate:"required,max=255"`
	}

	type User struct {
		Id        int               `validate:"required"`
		Name      string            `validate:"required,max=255"`
		Addresses []Address         `validate:"required,dive"`
		Hobbies   []string          `validate:"required,dive,required,min=1,oneof=football basketball"`
		Schools   map[string]School `validate:"required,dive,keys,required,min=1,endkeys,dive"`
		Wallet    map[string]int    `validate:"required,dive,keys,required,endkeys,required,gt=0"`
	}

	validate := validator.New()
	request := User{
		Id:   1,
		Name: "Joko",
		Addresses: []Address{
			{
				City:    "Jakarta",
				Country: "Indonesia",
			},
			{
				City:    "Bandung",
				Country: "Indonesia",
			},
		},
		Hobbies: []string{
			"football",
		},
		Schools: map[string]School{
			"SD": {
				Name: "SD Negeri 1",
			},
			"SMP": {
				Name: "SMP Negeri 1",
			},
		},
		Wallet: map[string]int{
			"USD": 100,
			"IDR": 0,
		},
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestAlias(t *testing.T) {
	validate := validator.New()
	validate.RegisterAlias("varchar", "required,max=255")

	type User struct {
		Id       int    `validate:"required"`
		Name     string `validate:"varchar,min=5"`
		Position string `validate:"required"`
	}

	request := User{
		Id:       1,
		Name:     "Joko Wicaksono",
		Position: "Software Engineer",
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func MustValidUsername(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)

	if ok {
		if value != strings.ToUpper(value) {
			return false
		}

		if len(value) < 5 {
			return false
		}
	}

	return true
}

func TestCustomValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("valid_username", MustValidUsername)

	type User struct {
		Username string `validate:"required,valid_username"`
		Password string `validate:"required"`
	}

	request := User{
		Username: "JOKOSANTOSO",
		Password: "123456",
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

var regexNumber = regexp.MustCompile(`^[0-9]+$`)

func MustValidPin(field validator.FieldLevel) bool {
	length, err := strconv.Atoi(field.Param())

	if err != nil {
		panic(err)
	}

	value := field.Field().String()
	if !regexNumber.MatchString(value) {
		return false
	}

	return len(value) == length
}

func TestCustomValidationWithParam(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("valid_pin", MustValidPin)

	type Login struct {
		Phone string `validate:"required,numeric"`
		Pin   string `validate:"required,numeric,valid_pin=6"`
	}

	request := Login{
		Phone: "081234567890",
		Pin:   "123456",
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func TestOrRule(t *testing.T) {
	type Login struct {
		Username string `validate:"required,email|numeric"`
		Password string `validate:"required"`
	}

	request := Login{
		Username: "joko@gmail.com",
		Password: "123456",
	}

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

func MustEqualsIgnoreCase(field validator.FieldLevel) bool {
	value, _, _, ok := field.GetStructFieldOK2()

	if !ok {
		return false
	}

	firstValue := strings.ToUpper(field.Field().String())
	secondValue := strings.ToUpper(value.String())

	return firstValue == secondValue
}

func TestCrossFieldValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("equals_ignore_case", MustEqualsIgnoreCase)

	type User struct {
		Username string `validate:"required,equals_ignore_case=Email|equals_ignore_case=Phone"`
		Email    string `validate:"required,email"`
		Phone    string `validate:"required,numeric"`
		Name     string `validate:"required"`
	}

	user := User{
		Username: "santoso@example.com",
		Email:    "santoso@example.com",
		Phone:    "23242",
		Name:     "santoso",
	}

	err := validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}

type RegisterRequest struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required,numeric"`
	Password string `validate:"required"`
}

func MustValidRegisterSuccess(level validator.StructLevel) {
	request := level.Current().Interface().(RegisterRequest)

	if request.Username == request.Email || request.Username == request.Phone {

	} else {
		level.ReportError(request.Username, "Username", "Username", "username_not_equals_email_or_phone", "")
	}
}

func TestStructLevelValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(MustValidRegisterSuccess, RegisterRequest{})

	request := RegisterRequest{
		Username: "joko@example.com",
		Email:    "joko@example.com",
		Phone:    "039383939",
		Password: "123456",
	}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Field(), " on tag ", validationError.Tag(), " with value ", validationError.Value())
		}
	}
}
