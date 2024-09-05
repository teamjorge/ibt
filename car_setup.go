package ibt

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/teamjorge/ibt/headers"
)

const (
	UPDATE_COUNT_FIELD_NAME string = "UpdateCount"
)

type CarSetup struct {
	Name   string
	Update int
	Values CarSetupDetails
}

type CarSetupDetails map[CarSetupKey]*CarSetupItem

type CarSetupItem struct {
	RawValue string
	Parsed   []CarSetupItemParsedValue
}

type CarSetupItemParsedValue struct {
	NumericalValue  float64
	NumericalSign   int
	MeasurementUnit string
}

type CarSetupKey string

type CarSetupKeys []CarSetupKey

func (a CarSetupKeys) Len() int           { return len(a) }
func (a CarSetupKeys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CarSetupKeys) Less(i, j int) bool { return a[i] < a[j] }

func NewCarSetupKey(category, subcategory, itemName string) CarSetupKey {
	return CarSetupKey(strings.Join([]string{category, subcategory, itemName}, "|"))
}

func (csk CarSetupKey) Category() string { return strings.Split(string(csk), "|")[0] }

func (csk CarSetupKey) SubCategory() string { return strings.Split(string(csk), "|")[1] }

func (csk CarSetupKey) ItemName() string { return strings.Split(string(csk), "|")[2] }

// ParseCarSetup from the given session info
func ParseCarSetup(sessionInfo *headers.Session) *CarSetup {
	setup := new(CarSetup)
	if update, ok := sessionInfo.CarSetup[UPDATE_COUNT_FIELD_NAME]; ok {
		setup.Update = update.(int)
	} else {
		setup.Update = 0
	}
	setup.Name = sessionInfo.DriverInfo.DriverSetupName

	setupItems := make(CarSetupDetails)
	for categoryName, category := range sessionInfo.CarSetup {
		if categoryName == UPDATE_COUNT_FIELD_NAME {
			continue
		}
		for subCategoryName, subCategory := range category.(map[string]interface{}) {
			for setupItemName, setupItemValue := range subCategory.(map[string]interface{}) {
				setupItem := CarSetupItem{
					RawValue: setupItemValue.(string),
					Parsed:   ParseSetupItem(setupItemValue.(string)),
				}
				setupItems.Add(categoryName, subCategoryName, setupItemName, &setupItem)
			}
		}
	}

	setup.Values = setupItems

	return setup
}

// CarSetup for the given stub
func (stub Stub) CarSetup() *CarSetup { return ParseCarSetup(stub.header.SessionInfo) }

// Add a new setup item for the specified category and subcategory
func (s CarSetupDetails) Add(category, subcategory, itemName string, value *CarSetupItem) {
	storeKey := NewCarSetupKey(category, subcategory, itemName)

	s[storeKey] = value
}

func (c *CarSetupItem) IsParsed() bool { return len(c.Parsed) > 0 }

var carSetupNumericalSigns = map[string]int{
	"-": -1,
	"":  0,
	"+": 1,
}

func ParseSetupItem(input string) []CarSetupItemParsedValue {
	// Handle multiple values - such as tyre wear/temperature etc
	valueParts := strings.Split(input, ", ")

	parsedValues := make([]CarSetupItemParsedValue, 0)

	for _, part := range valueParts {
		partItem, err := parseCarSetupItemFromInput(part)
		if err != nil {
			continue
		}
		parsedValues = append(parsedValues, partItem)
	}

	return parsedValues
}

var (
	downforcePattern = regexp.MustCompile(`([\d|\.]*):([\d.*])`)
	numericalPattern = regexp.MustCompile(`((?P<numSign>[-|+])|)(?P<numValue>[\d|\.|]*)`)
)

func parseCarSetupItemFromInput(input string) (CarSetupItemParsedValue, error) {
	var parsedValue CarSetupItemParsedValue
	input = strings.TrimSpace(input)

	if downforcePattern.MatchString(input) {
		return parsedValue, errors.New("downforce pattern detected")
	}

	matches := numericalPattern.FindStringSubmatch(input)
	if matches[numericalPattern.SubexpIndex("numValue")] == "" {
		return parsedValue, errors.New("no numerical values detected")
	}

	foundNumValue := strings.TrimSpace(matches[numericalPattern.SubexpIndex("numValue")])
	numVal, err := strconv.ParseFloat(foundNumValue, 64)
	if err != nil {
		return parsedValue, errors.New("failed to parse expected float")
	}

	numSign := matches[numericalPattern.SubexpIndex("numSign")]

	parsedValue.NumericalValue = numVal
	parsedValue.NumericalSign = carSetupNumericalSigns[numSign]
	parsedValue.MeasurementUnit = strings.TrimSpace(strings.Replace(input, numSign+foundNumValue, "", 1))

	return parsedValue, nil
}

type CarSetupComparison map[CarSetupKey]*CarSetupComparisonItem

func (c CarSetupComparison) Differences() map[CarSetupKey]*CarSetupComparisonItem {
	differences := make(map[CarSetupKey]*CarSetupComparisonItem)

	for setupItemName, setupItemValue := range c {
		if setupItemValue.different {
			differences[setupItemName] = setupItemValue
		}
	}

	return differences
}

type CarSetupComparisonItem struct {
	I1 *CarSetupItem
	I2 *CarSetupItem

	NumericalDifferences []float64
	RawDifference        string
	different            bool
}

func CompareSetups(s1, s2 *CarSetup) CarSetupComparison {
	comparisons := make(CarSetupComparison)

	for itemName, initialItem := range s1.Values {
		if targetItem, ok := s2.Values[itemName]; ok {
			difference := CompareSetupItemParsedValue(initialItem, targetItem)
			comparisons[itemName] = &difference
		}
	}

	return comparisons
}

func CompareSetupItemParsedValue(i1, i2 *CarSetupItem) CarSetupComparisonItem {
	diff := CarSetupComparisonItem{I1: i1, I2: i2}

	// No difference
	if i1.RawValue == i2.RawValue {
		return diff
	}

	diff.different = true
	diff.RawDifference = fmt.Sprintf("%s -> %s", i1.RawValue, i2.RawValue)

	if i1.IsParsed() && i2.IsParsed() && len(i1.Parsed) == len(i2.Parsed) {
		numericalDifferences := make([]float64, 0)

		for idx, parsedValue := range i1.Parsed {
			numericalDifferences = append(numericalDifferences, compareSetupNumericalItem(parsedValue, i2.Parsed[idx]))
		}
		diff.NumericalDifferences = numericalDifferences
	}

	return diff
}

func compareSetupNumericalItem(i1, i2 CarSetupItemParsedValue) float64 {
	sign1 := i1.NumericalSign
	sign2 := i2.NumericalSign

	if sign1 == 0 {
		sign1 = 1
	}
	if sign2 == 0 {
		sign2 = 1
	}

	value1 := i1.NumericalValue * float64(sign1)
	value2 := i2.NumericalValue * float64(sign2)

	return value2 - value1
}

func FilterSetupItems[k comparable](setupItems map[CarSetupKey]k, category, subCategory, itemName string) map[CarSetupKey]k {
	var flag int
	if category != "" {
		flag += 1
	}
	if subCategory != "" {
		flag += 2
	}
	if itemName != "" {
		flag += 4
	}

	var keepFunc func(CarSetupKey) bool

	switch flag {
	case 0:
		return setupItems
	case 1:
		keepFunc = func(csk CarSetupKey) bool {
			return strings.HasPrefix(string(csk), category)
		}
	case 2:
		keepFunc = func(csk CarSetupKey) bool {
			return strings.Split(string(csk), "|")[1] == subCategory
		}
	case 3:
		keepFunc = func(csk CarSetupKey) bool {
			return strings.HasPrefix(string(csk), strings.Join([]string{category, subCategory}, "|"))
		}
	case 5:
		keepFunc = func(csk CarSetupKey) bool {
			return strings.HasPrefix(string(csk), category) && strings.HasSuffix(string(csk), itemName)
		}
	case 6:
		keepFunc = func(csk CarSetupKey) bool {
			return strings.HasSuffix(string(csk), strings.Join([]string{subCategory, itemName}, "|"))
		}
	case 7:
		keepFunc = func(csk CarSetupKey) bool {
			return csk == NewCarSetupKey(category, subCategory, itemName)
		}
	}

	filteredItems := make(map[CarSetupKey]k)

	for item, value := range setupItems {
		if keepFunc(item) {
			filteredItems[item] = value
		}
	}

	return filteredItems
}

type SetupDiscardFunc func(string) func(CarSetupKey) bool

func HasCategory(category string) func(csk CarSetupKey) bool {
	return func(csk CarSetupKey) bool {
		return strings.HasPrefix(string(csk), category)
	}
}

func HasSubCategory(subCategory string) func(csk CarSetupKey) bool {
	return func(csk CarSetupKey) bool {
		extractedSubCategory := strings.Split(string(csk), "|")[1]
		return strings.HasPrefix(extractedSubCategory, subCategory)
	}
}

func HasItemName(itemName string) func(csk CarSetupKey) bool {
	return func(csk CarSetupKey) bool {
		extractedSubCategory := strings.Split(string(csk), "|")[1]
		return strings.HasPrefix(extractedSubCategory, itemName)
	}
}
