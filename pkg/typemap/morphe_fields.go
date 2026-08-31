package typemap

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kalo-build/morphe-go/pkg/yaml"
)

var knownFormats = map[string]func(*gofakeit.Faker) any{
	"email":     func(f *gofakeit.Faker) any { return f.Email() },
	"phone":     func(f *gofakeit.Faker) any { return f.Phone() },
	"url":       func(f *gofakeit.Faker) any { return f.URL() },
	"ipv4":      func(f *gofakeit.Faker) any { return f.IPv4Address() },
	"ipv6":      func(f *gofakeit.Faker) any { return f.IPv6Address() },
	"name":      func(f *gofakeit.Faker) any { return f.Name() },
	"firstName": func(f *gofakeit.Faker) any { return f.FirstName() },
	"lastName":  func(f *gofakeit.Faker) any { return f.LastName() },
	"city":      func(f *gofakeit.Faker) any { return f.City() },
	"country":   func(f *gofakeit.Faker) any { return f.Country() },
	"street":    func(f *gofakeit.Faker) any { return f.Street() },
	"zip":       func(f *gofakeit.Faker) any { return f.Zip() },
	"state":     func(f *gofakeit.Faker) any { return f.State() },
	"company":   func(f *gofakeit.Faker) any { return f.Company() },
	"sentence":  func(f *gofakeit.Faker) any { return f.Sentence(6) },
	"paragraph": func(f *gofakeit.Faker) any { return f.Paragraph(1, 3, 5, " ") },
	"username":  func(f *gofakeit.Faker) any { return f.Username() },
	"color":     func(f *gofakeit.Faker) any { return f.Color() },
	"hexColor":  func(f *gofakeit.Faker) any { return f.HexColor() },
	"currency":  func(f *gofakeit.Faker) any { return f.CurrencyShort() },
}

func GenerateFieldValue(faker *gofakeit.Faker, fieldName string, fieldType yaml.ModelFieldType, hints FormatHints) any {
	if hints.Format != "" {
		if fn, ok := knownFormats[hints.Format]; ok {
			return fn(faker)
		}
	}

	if hints.Regex != "" {
		if val, err := generateFromRegex(faker, hints.Regex); err == nil {
			return val
		}
	}

	switch fieldType {
	case yaml.ModelFieldTypeUUID:
		return faker.UUID()

	case yaml.ModelFieldTypeString:
		return generateString(faker, fieldName, hints)

	case yaml.ModelFieldTypeInteger:
		return generateInteger(faker, hints)

	case yaml.ModelFieldTypeFloat:
		return generateFloat(faker, hints)

	case yaml.ModelFieldTypeBoolean:
		return faker.Bool()

	case yaml.ModelFieldTypeTime:
		t := faker.DateRange(
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
		)
		return t.UTC().Format(time.RFC3339)

	case yaml.ModelFieldTypeDate:
		t := faker.DateRange(
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		)
		return t.UTC().Format("2006-01-02")

	case yaml.ModelFieldTypeProtected:
		return faker.Password(true, true, true, true, false, 24)

	case yaml.ModelFieldTypeSealed:
		return "$2a$10$" + faker.LetterN(53)

	default:
		return faker.Word()
	}
}

func generateString(faker *gofakeit.Faker, fieldName string, hints FormatHints) string {
	val := inferStringFromFieldName(faker, fieldName)

	if hints.MaxLength > 0 && len(val) > hints.MaxLength {
		val = val[:hints.MaxLength]
	}
	if hints.MinLength > 0 && len(val) < hints.MinLength {
		val = val + faker.LetterN(uint(hints.MinLength-len(val)))
	}
	return val
}

func inferStringFromFieldName(faker *gofakeit.Faker, fieldName string) string {
	lower := strings.ToLower(fieldName)
	switch {
	case lower == "email" || strings.HasSuffix(lower, "email"):
		return faker.Email()
	case lower == "phone" || strings.HasSuffix(lower, "phone"):
		return faker.Phone()
	case lower == "firstname" || lower == "first_name":
		return faker.FirstName()
	case lower == "lastname" || lower == "last_name":
		return faker.LastName()
	case lower == "name" || lower == "displayname" || lower == "display_name":
		return faker.Name()
	case lower == "title":
		return faker.JobTitle()
	case lower == "description" || lower == "body" || lower == "content" || lower == "text":
		return faker.Sentence(8)
	case lower == "url" || lower == "website":
		return faker.URL()
	case lower == "city":
		return faker.City()
	case lower == "state":
		return faker.State()
	case lower == "country":
		return faker.Country()
	case lower == "street" || lower == "address" || lower == "addressline1" || lower == "address_line_1":
		return faker.Street()
	case lower == "zip" || lower == "zipcode" || lower == "zip_code" || lower == "postalcode" || lower == "postal_code":
		return faker.Zip()
	case lower == "taxid" || lower == "tax_id":
		return fmt.Sprintf("%d-%d", faker.Number(10, 99), faker.Number(1000000, 9999999))
	default:
		return faker.Sentence(3)
	}
}

func generateInteger(faker *gofakeit.Faker, hints FormatHints) int {
	minVal := 1
	maxVal := 10000
	if hints.HasMin {
		minVal = int(hints.Min)
	}
	if hints.HasMax {
		maxVal = int(hints.Max)
	}
	return faker.Number(minVal, maxVal)
}

func generateFloat(faker *gofakeit.Faker, hints FormatHints) float64 {
	minVal := 0.0
	maxVal := 10000.0
	if hints.HasMin {
		minVal = hints.Min
	}
	if hints.HasMax {
		maxVal = hints.Max
	}
	return faker.Float64Range(minVal, maxVal)
}

func generateFromRegex(faker *gofakeit.Faker, pattern string) (string, error) {
	if _, err := regexp.Compile(pattern); err != nil {
		return "", err
	}
	return faker.Regex(pattern), nil
}
