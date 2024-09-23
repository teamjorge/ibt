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
	// Number of times the setup has been modified
	UPDATE_COUNT_FIELD_NAME string = "UpdateCount"
)

// CarSetup is the overarching structure to represent the setup of a car.
type CarSetup struct {
	Name   string
	Update int
	Values CarSetupDetails
}

// CarSetupDetails is the structure for storing the individual items that make up a car setup.
//
// Example: The front right tyre pressure for the given setup will be 128 PSI.
type CarSetupDetails map[CarSetupKey]*CarSetupItem

// CarSetupItem is the raw and parsed values of a single car setup item.
type CarSetupItem struct {
	RawValue string
	Parsed   []CarSetupItemParsedValue
}

// CarSetupItemParsedValue is a detailed numerical value of a single car setup item.
//
// Parsed items will refer only to numerical values. Additionally, the sign of the numeric value will
// be preserved for cases where + / - is applicable. For example, some cars might have +3 clicks of
// wing, rather than just 3 clicks.
type CarSetupItemParsedValue struct {
	NumericalValue  float64
	NumericalSign   int
	MeasurementUnit string
}

// CarSetupKey refers to the map key used when storing CarSetupItems.
//
// A CarSetupKey should ideally be created using the NewCarSetupKey function. This
// eliminates any risk with deconstructing it to find specific values, such as
// category, subcategory, and item name.
//
// This key consists of the category, subcategory and item name.
//
// For example:
//
//	Category |    SubCategory  |  Item Name
//
// DriveBrake|BrakeSystemConfig|BaseBrakeBias
type CarSetupKey string

// CarSetupKeys are multiple CarSetupKey values
//
// This structure primarily exists for sorting purposes
type CarSetupKeys []CarSetupKey

func (a CarSetupKeys) Len() int           { return len(a) }
func (a CarSetupKeys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CarSetupKeys) Less(i, j int) bool { return a[i] < a[j] }

// NewCarSetupKey initialises a new CarSetupKey value
func NewCarSetupKey(category, subcategory, itemName string) CarSetupKey {
	return CarSetupKey(strings.Join([]string{category, subcategory, itemName}, "|"))
}

// Category part of the given CarSetupKey
func (csk CarSetupKey) Category() string { return strings.Split(string(csk), "|")[0] }

// SubCategory part of the given CarSetupKey
func (csk CarSetupKey) SubCategory() string { return strings.Split(string(csk), "|")[1] }

// ItemName part of the given CarSetupKey
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

// IsParsed determines if the given CarSetupItem has any parsed numerical values
func (c *CarSetupItem) IsParsed() bool { return len(c.Parsed) > 0 }

// Storage for preserving numerical signs in parsed numerical items
var carSetupNumericalSigns = map[string]int{
	"-": -1,
	"":  0,
	"+": 1,
}

// ParseSetupItem attempts to parse a numerical value from the given setup item.
//
// This function returns a slice due to some setup items consisting of multiple values.
//
// For example: Tire pressures will contain 3 independent values for inner, outer, and carcass
// temperatures.
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

// Patterns used for finding and parsing numerical values in setup items
var (
	downforcePattern = regexp.MustCompile(`([\d|\.]*):([\d.*])`)
	numericalPattern = regexp.MustCompile(`((?P<numSign>[-|+])|)(?P<numValue>[\d|\.|]*)`)
)

// parseCarSetupItemFromInput attempts to parse numerical values from the given raw value
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

// CarSetupComparison is the overarching structure for storing the difference between two
// CarSetups.
type CarSetupComparison map[CarSetupKey]*CarSetupComparisonItem

// Differences between the two CarSetups that were compared
func (c CarSetupComparison) Differences() map[CarSetupKey]*CarSetupComparisonItem {
	differences := make(map[CarSetupKey]*CarSetupComparisonItem)

	for setupItemName, setupItemValue := range c {
		if setupItemValue.different {
			differences[setupItemName] = setupItemValue
		}
	}

	return differences
}

// CarSetupComparisonItem is the comparison of a single item from the two compared setups
type CarSetupComparisonItem struct {
	I1 *CarSetupItem
	I2 *CarSetupItem

	NumericalDifferences []float64
	RawDifference        string
	different            bool
}

// CompareSetups will compare the differences between the two given setups
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

// CompareSetupItemParsedValue compares the parsed values of each setup item.
//
// If a parsed item is not found in one or either of the items, the raw value comparison
// will be the only populated field
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

// compareSetupNumericalItem performs the comparison when parsed values are populated
// and requires numerical values (with signs) to be compared
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

// SetupFilter is used to filter a CarSetup or CarSetupComparison for specific categories,
// subcategories and/or items.
//
// Prefixes can be used instead of specifying full categories, subcategories, and/or item names.
//
// Fields that are not populated will be ignored during equality checks.
type SetupFilter struct {
	Category    string
	SubCategory string
	ItemName    string
}

// Compare determines if SetupFilter contains all of the populated filter fields
func (s SetupFilter) Compare(key CarSetupKey) bool {
	return strings.Contains(key.Category(), s.Category) &&
		strings.Contains(key.SubCategory(), s.SubCategory) &&
		strings.Contains(key.ItemName(), s.ItemName)
}

// FilterSetupItems for only the given categories, subcategories, and items.
func FilterSetupItems[k comparable](setupItems map[CarSetupKey]k, filters ...SetupFilter) map[CarSetupKey]k {
	filteredItems := make(map[CarSetupKey]k)

	for item, value := range setupItems {
		for _, filter := range filters {
			if filter.Compare(item) {
				filteredItems[item] = value
				break
			}
		}
	}

	return filteredItems
}

// DiscardSetupItems removes the given categories, subcategories, and/or items.
func DiscardSetupItems[k comparable](setupItems map[CarSetupKey]k, filters ...SetupFilter) map[CarSetupKey]k {
	filteredItems := make(map[CarSetupKey]k)

	for item, value := range setupItems {
		triggered := false
		for _, filter := range filters {
			if filter.Compare(item) {
				triggered = true
				break
			}
		}
		if !triggered {
			filteredItems[item] = value
		}
	}

	return filteredItems
}
