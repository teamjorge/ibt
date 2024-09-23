package ibt

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/teamjorge/ibt/headers"
)

func TestCarSetupKeysSort(t *testing.T) {
	keys := CarSetupKeys{
		CarSetupKey("DriveBrake|PowerUnitConfig|EngineBraking"),
		CarSetupKey("TiresAero|RightFrontTire|StartingPressure"),
		CarSetupKey("Chassis|Front|HeaveRate"),
		CarSetupKey("DriveBrake|Differential|Middle"),
	}

	sort.Sort(keys)

	expected := CarSetupKeys{
		CarSetupKey("Chassis|Front|HeaveRate"),
		CarSetupKey("DriveBrake|Differential|Middle"),
		CarSetupKey("DriveBrake|PowerUnitConfig|EngineBraking"),
		CarSetupKey("TiresAero|RightFrontTire|StartingPressure"),
	}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Sort() CarSetupKeys = %v, want %v", keys, expected)
	}
}

func TestNewCarSetupKey(t *testing.T) {
	type args struct {
		category    string
		subcategory string
		itemName    string
	}
	tests := []struct {
		name string
		args args
		want CarSetupKey
	}{
		{
			"normal new car setup key",
			args{"cat", "subcat", "item"},
			"cat|subcat|item",
		},
		{
			"missing value new setup key",
			args{category: "cat", itemName: "item"},
			"cat||item",
		},
		{
			"empty new setup key",
			args{},
			"||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCarSetupKey(tt.args.category, tt.args.subcategory, tt.args.itemName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCarSetupKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCarSetupKey_CategorySubCategoryItemName(t *testing.T) {
	tests := []struct {
		name string
		csk  CarSetupKey
		want string
	}{
		{
			"normal setup key category",
			CarSetupKey("cat|sub|item"),
			"cat",
		},
		{
			"empty setup key category",
			CarSetupKey("|sub|item"),
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.csk.Category(); got != tt.want {
				t.Errorf("CarSetupKey.Category() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCarSetupKey_SubCategory(t *testing.T) {
	tests := []struct {
		name string
		csk  CarSetupKey
		want string
	}{
		{
			"normal setup key sub category",
			CarSetupKey("cat|sub|item"),
			"sub",
		},
		{
			"empty setup key sub category",
			CarSetupKey("cat||item"),
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.csk.SubCategory(); got != tt.want {
				t.Errorf("CarSetupKey.SubCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCarSetupKey_ItemName(t *testing.T) {
	tests := []struct {
		name string
		csk  CarSetupKey
		want string
	}{
		{
			"normal setup key item name",
			CarSetupKey("cat|sub|item"),
			"item",
		},
		{
			"empty setup key item name",
			CarSetupKey("cat|sub|"),
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.csk.ItemName(); got != tt.want {
				t.Errorf("CarSetupKey.ItemName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCarSetup(t *testing.T) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %+v", err)
		return
	}
	defer f.Close()

	header, err := headers.ParseHeaders(f)
	if err != nil {
		t.Errorf("failed to parse session header for testing file - %+v", err)
		return
	}

	type args struct {
		sessionInfo *headers.Session
	}
	tests := []struct {
		name string
		args args
		want *CarSetup
	}{
		{
			"test parse setup from valid header",
			args{header.SessionInfo},
			expectedFullCarSetup,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCarSetup(tt.args.sessionInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCarSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStub_CarSetup(t *testing.T) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %+v", err)
		return
	}
	defer f.Close()

	header, err := headers.ParseHeaders(f)
	if err != nil {
		t.Errorf("failed to parse session header for testing file - %+v", err)
		return
	}

	stub := Stub{
		header: header,
	}

	tests := []struct {
		name string
		stub Stub
		want *CarSetup
	}{
		{
			"parse setup from stub",
			stub,
			expectedFullCarSetup,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stub.CarSetup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stub.CarSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCarSetupDetails_Add(t *testing.T) {
	carSetupItem1 := &CarSetupItem{RawValue: "medium"}
	carSetupItem2 := &CarSetupItem{
		RawValue: "+3 Clicks",
		Parsed: []CarSetupItemParsedValue{
			{
				MeasurementUnit: "Clicks",
				NumericalSign:   1,
				NumericalValue:  3,
			},
		},
	}

	type args struct {
		category    string
		subcategory string
		itemName    string
		value       *CarSetupItem
	}
	tests := []struct {
		name string
		s    CarSetupDetails
		args args
		want string
	}{
		{
			"test string setup",
			CarSetupDetails{},
			args{
				"cat",
				"sub",
				"item name",
				carSetupItem1,
			},
			"medium",
		},
		{
			"test string setup",
			CarSetupDetails{},
			args{
				"cat",
				"sub",
				"item name",
				carSetupItem2,
			},
			"+3 Clicks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Add(tt.args.category, tt.args.subcategory, tt.args.itemName, tt.args.value)

			key := NewCarSetupKey(tt.args.category, tt.args.subcategory, tt.args.itemName)
			item, ok := tt.s[key]
			if !ok {
				t.Errorf("expected key %s to be added to car setup", key)
			}

			if item.RawValue != tt.want {
				t.Errorf("expected raw value to be %s. received %s", tt.want, item.RawValue)
			}
		})
	}
}

func TestCarSetupItem_IsParsed(t *testing.T) {
	carSetupItem1 := &CarSetupItem{RawValue: "medium"}
	carSetupItem2 := &CarSetupItem{
		RawValue: "+3 Clicks",
		Parsed: []CarSetupItemParsedValue{
			{
				MeasurementUnit: "Clicks",
				NumericalSign:   1,
				NumericalValue:  3,
			},
		},
	}

	tests := []struct {
		name string
		c    *CarSetupItem
		want bool
	}{
		{
			"not parsed",
			carSetupItem1,
			false,
		},
		{
			"parsed",
			carSetupItem2,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsParsed(); got != tt.want {
				t.Errorf("CarSetupItem.IsParsed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSetupItem(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want []CarSetupItemParsedValue
	}{
		{
			"parse non-metric setup item",
			args{
				input: "medium",
			},
			[]CarSetupItemParsedValue{},
		},
		{
			"parse single metric setup item",
			args{
				input: "+3 Clicks",
			},
			[]CarSetupItemParsedValue{
				{
					MeasurementUnit: "Clicks",
					NumericalSign:   1,
					NumericalValue:  3,
				},
			},
		},
		{
			"parse multiple metric setup item",
			args{
				input: "95%, 97%, 96%",
			},
			[]CarSetupItemParsedValue{
				{
					MeasurementUnit: "%",
					NumericalSign:   0,
					NumericalValue:  95,
				},
				{
					MeasurementUnit: "%",
					NumericalSign:   0,
					NumericalValue:  97,
				},
				{
					MeasurementUnit: "%",
					NumericalSign:   0,
					NumericalValue:  96,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSetupItem(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSetupItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseCarSetupItemFromInput(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    CarSetupItemParsedValue
		wantErr bool
	}{
		{
			"parse downforce pattern",
			args{"4.501:1"},
			CarSetupItemParsedValue{},
			true,
		},
		{
			"parse no numerical value",
			args{"medium"},
			CarSetupItemParsedValue{},
			true,
		},
		{
			"parse incorrect float",
			args{"0....32221 %"},
			CarSetupItemParsedValue{},
			true,
		},
		{
			"parse normal value",
			args{"%"},
			CarSetupItemParsedValue{},
			true,
		},
		{
			"test ride height",
			args{"25.0 mm"},
			CarSetupItemParsedValue{
				NumericalValue:  25.0,
				NumericalSign:   0,
				MeasurementUnit: "mm",
			},
			false,
		},
		{
			"test heave rate",
			args{"750 N/mm"},
			CarSetupItemParsedValue{
				NumericalValue:  750.0,
				NumericalSign:   0,
				MeasurementUnit: "N/mm",
			},
			false,
		},
		{
			"test corner weight",
			args{"1902 N"},
			CarSetupItemParsedValue{
				NumericalValue:  1902.0,
				NumericalSign:   0,
				MeasurementUnit: "N",
			},
			false,
		},
		{
			"test camber",
			args{"-3.15 deg"},
			CarSetupItemParsedValue{
				NumericalValue:  3.15,
				NumericalSign:   -1,
				MeasurementUnit: "deg",
			},
			false,
		},
		{
			"test toe in",
			args{"+0.05 deg"},
			CarSetupItemParsedValue{
				NumericalValue:  0.05,
				NumericalSign:   1,
				MeasurementUnit: "deg",
			},
			false,
		},
		{
			"test fuel",
			args{"45 Kg"},
			CarSetupItemParsedValue{
				NumericalValue:  45.0,
				NumericalSign:   0,
				MeasurementUnit: "Kg",
			},
			false,
		},
		{
			"test bbal",
			args{"52.5% (BBAL)"},
			CarSetupItemParsedValue{
				NumericalValue:  52.5,
				NumericalSign:   0,
				MeasurementUnit: "% (BBAL)",
			},
			false,
		},
		{
			"test fine bbal",
			args{"-1.0 (BB+/BB-)"},
			CarSetupItemParsedValue{
				NumericalValue:  1.0,
				NumericalSign:   -1,
				MeasurementUnit: "(BB+/BB-)",
			},
			false,
		},
		{
			"test dynamic ramping",
			args{"20% pedal"},
			CarSetupItemParsedValue{
				NumericalValue:  20,
				NumericalSign:   0,
				MeasurementUnit: "% pedal",
			},
			false,
		},
		{
			"test brake migration",
			args{"6 (BMIG)"},
			CarSetupItemParsedValue{
				NumericalValue:  6,
				NumericalSign:   0,
				MeasurementUnit: "(BMIG)",
			},
			false,
		},
		{
			"test aero balance",
			args{"48.14%"},
			CarSetupItemParsedValue{
				NumericalValue:  48.14,
				NumericalSign:   0,
				MeasurementUnit: "%",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCarSetupItemFromInput(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCarSetupItemFromInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCarSetupItemFromInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCarSetupComparison_Differences(t *testing.T) {
	heaveRate1 := &CarSetupItem{
		"750 N/mm",
		[]CarSetupItemParsedValue{
			{
				NumericalValue:  750.0,
				NumericalSign:   0,
				MeasurementUnit: "N/mm",
			},
		},
	}

	heaveRate2 := &CarSetupItem{
		"780 N/mm",
		[]CarSetupItemParsedValue{
			{
				NumericalValue:  780.0,
				NumericalSign:   0,
				MeasurementUnit: "N/mm",
			},
		},
	}

	tests := []struct {
		name string
		c    CarSetupComparison
		want map[CarSetupKey]*CarSetupComparisonItem
	}{
		{
			"test has differences",
			CarSetupComparison{
				"Chassis|Front|HeaveRate": &CarSetupComparisonItem{
					heaveRate1,
					heaveRate2,
					[]float64{30},
					"750 N/mm -> 780 N/mm",
					true,
				},
				"TiresAero|TireCompound|TireCompound": &CarSetupComparisonItem{
					I1: &CarSetupItem{RawValue: "Medium"},
					I2: &CarSetupItem{RawValue: "Medium"},
				},
			},
			map[CarSetupKey]*CarSetupComparisonItem{
				"Chassis|Front|HeaveRate": {
					heaveRate1, heaveRate2,
					[]float64{30},
					"750 N/mm -> 780 N/mm",
					true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Differences(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CarSetupComparison.Differences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareSetups(t *testing.T) {
	type args struct {
		s1 *CarSetup
		s2 *CarSetup
	}
	tests := []struct {
		name string
		args args
		want CarSetupComparison
	}{
		{
			"full comparison",
			args{
				expectedFullCarSetup,
				comparisonCarSetup,
			},
			CarSetupComparison{
				"Chassis|Front|HeaveRate": {
					expectedFullCarSetup.Values["Chassis|Front|HeaveRate"], comparisonCarSetup.Values["Chassis|Front|HeaveRate"],
					[]float64{30},
					"750 N/mm -> 780 N/mm",
					true,
				},
				"TiresAero|TireCompound|TireCompound": {
					expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"], comparisonCarSetup.Values["TiresAero|TireCompound|TireCompound"],
					nil,
					"Medium -> Soft",
					true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareSetups(tt.args.s1, tt.args.s2).Differences()
			for key, comparisonValue := range got {
				if !reflect.DeepEqual(comparisonValue, tt.want[key]) {
					t.Errorf("CompareSetups() key %s = %v, want %v", key, comparisonValue, tt.want[key])
				}
			}
		})
	}
}

func TestCompareSetupItemParsedValue(t *testing.T) {
	type args struct {
		i1 *CarSetupItem
		i2 *CarSetupItem
	}
	tests := []struct {
		name string
		args args
		want CarSetupComparisonItem
	}{
		{
			"normal with metric values",
			args{
				expectedFullCarSetup.Values["Chassis|Front|HeaveRate"],
				comparisonCarSetup.Values["Chassis|Front|HeaveRate"],
			},
			CarSetupComparisonItem{
				expectedFullCarSetup.Values["Chassis|Front|HeaveRate"], comparisonCarSetup.Values["Chassis|Front|HeaveRate"],
				[]float64{30},
				"750 N/mm -> 780 N/mm",
				true,
			},
		},
		{
			"normal with string values",
			args{
				expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"],
				comparisonCarSetup.Values["TiresAero|TireCompound|TireCompound"],
			},
			CarSetupComparisonItem{
				expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"], comparisonCarSetup.Values["TiresAero|TireCompound|TireCompound"],
				nil,
				"Medium -> Soft",
				true,
			},
		},
		{
			"mismatched items",
			args{
				expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"],
				comparisonCarSetup.Values["Chassis|Front|HeaveRate"],
			},
			CarSetupComparisonItem{
				expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"], comparisonCarSetup.Values["Chassis|Front|HeaveRate"],
				nil,
				"Medium -> 780 N/mm",
				true,
			},
		},
		{
			"same value",
			args{
				expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"],
				expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"],
			},
			CarSetupComparisonItem{
				expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"], expectedFullCarSetup.Values["TiresAero|TireCompound|TireCompound"],
				nil,
				"",
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareSetupItemParsedValue(tt.args.i1, tt.args.i2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareSetupItemParsedValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareSetupNumericalItem(t *testing.T) {
	type args struct {
		i1 CarSetupItemParsedValue
		i2 CarSetupItemParsedValue
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"neg lower",
			args{
				CarSetupItemParsedValue{NumericalValue: 4, NumericalSign: -1},
				CarSetupItemParsedValue{NumericalValue: 6, NumericalSign: -1},
			},
			-2,
		},
		{
			"neg higher",
			args{
				CarSetupItemParsedValue{NumericalValue: 4, NumericalSign: -1},
				CarSetupItemParsedValue{NumericalValue: 1, NumericalSign: -1},
			},
			3,
		},
		{
			"neg to pos",
			args{
				CarSetupItemParsedValue{NumericalValue: 4, NumericalSign: -1},
				CarSetupItemParsedValue{NumericalValue: 2, NumericalSign: 0},
			},
			6,
		},
		{
			"pos to neg",
			args{
				CarSetupItemParsedValue{NumericalValue: 4, NumericalSign: 0},
				CarSetupItemParsedValue{NumericalValue: 2, NumericalSign: -1},
			},
			-6,
		},
		{
			"pos lower",
			args{
				CarSetupItemParsedValue{NumericalValue: 4, NumericalSign: 1},
				CarSetupItemParsedValue{NumericalValue: 2, NumericalSign: 1},
			},
			-2,
		},
		{
			"pos higher",
			args{
				CarSetupItemParsedValue{NumericalValue: 1, NumericalSign: 1},
				CarSetupItemParsedValue{NumericalValue: 3, NumericalSign: 1},
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareSetupNumericalItem(tt.args.i1, tt.args.i2); got != tt.want {
				t.Errorf("compareSetupNumericalItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetupFilter_Compare(t *testing.T) {
	type args struct {
		key CarSetupKey
	}
	tests := []struct {
		name string
		s    SetupFilter
		args args
		want bool
	}{
		{
			"prefix category exists",
			SetupFilter{Category: "cat"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"contains category exists",
			SetupFilter{Category: "egory"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"category full",
			SetupFilter{Category: "category"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"prefix category does not exist",
			SetupFilter{Category: "bork?"},
			args{"category|subcategory|item_name"},
			false,
		},
		{
			"prefix subcategory exists",
			SetupFilter{SubCategory: "subcat"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"subcategory contains",
			SetupFilter{SubCategory: "bcat"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"subcategory full",
			SetupFilter{SubCategory: "subcategory"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"prefix subcategory does not exist",
			SetupFilter{SubCategory: "bork?"},
			args{"category|subcategory|item_name"},
			false,
		},
		{
			"prefix item name exists",
			SetupFilter{ItemName: "item"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"item name full",
			SetupFilter{ItemName: "item_name"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"contains item name",
			SetupFilter{ItemName: "m_n"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"prefix item name does not exist",
			SetupFilter{ItemName: "bork?"},
			args{"category|subcategory|item_name"},
			false,
		},
		{
			"prefix multiple exists",
			SetupFilter{ItemName: "item", Category: "cat"},
			args{"category|subcategory|item_name"},
			true,
		},
		{
			"prefix multiple one doesnt exist",
			SetupFilter{ItemName: "item_name", Category: "bork?"},
			args{"category|subcategory|item_name"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Compare(tt.args.key); got != tt.want {
				t.Errorf("SetupFilter.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterSetupItems(t *testing.T) {
	inputCarSetup := map[CarSetupKey]*CarSetup{
		"TiresAero|LeftRearTire|StartingPressure":     {},
		"TiresAero|RightFrontTire|LastHotPressure":    {},
		"TiresAero|RightRearTire|LastTempsIMO":        {},
		"TiresAero|TireCompound|TireCompound":         {},
		"DriveBrake|BrakeSystemConfig|TotalBrakeBias": {},
	}

	type args struct {
		setupItems map[CarSetupKey]*CarSetup
		filters    []SetupFilter
	}
	tests := []struct {
		name string
		args args
		want map[CarSetupKey]*CarSetup
	}{
		{
			"filter for tiresaero pressures",
			args{
				setupItems: inputCarSetup,
				filters:    []SetupFilter{{"TiresAero", "", "Pressure"}},
			},
			map[CarSetupKey]*CarSetup{
				"TiresAero|LeftRearTire|StartingPressure":  inputCarSetup["TiresAero|LeftRearTire|StartingPressure"],
				"TiresAero|RightFrontTire|LastHotPressure": inputCarSetup["TiresAero|RightFrontTire|LastHotPressure"],
			},
		},
		{
			"filter for tire data",
			args{
				setupItems: inputCarSetup,
				filters:    []SetupFilter{{"", "Tire", ""}},
			},
			map[CarSetupKey]*CarSetup{
				"TiresAero|LeftRearTire|StartingPressure":  inputCarSetup["TiresAero|LeftRearTire|StartingPressure"],
				"TiresAero|RightFrontTire|LastHotPressure": inputCarSetup["TiresAero|RightFrontTire|LastHotPressure"],
				"TiresAero|RightRearTire|LastTempsIMO":     inputCarSetup["TiresAero|RightRearTire|LastTempsIMO"],
				"TiresAero|TireCompound|TireCompound":      inputCarSetup["TiresAero|TireCompound|TireCompound"],
			},
		},
		{
			"filter for brake bias",
			args{
				setupItems: inputCarSetup,
				filters:    []SetupFilter{{"Brake", "", "Total"}},
			},
			map[CarSetupKey]*CarSetup{
				"DriveBrake|BrakeSystemConfig|TotalBrakeBias": inputCarSetup["DriveBrake|BrakeSystemConfig|TotalBrakeBias"],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterSetupItems(tt.args.setupItems, tt.args.filters...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterSetupItems() = \n%v\n, want \n%v", got, tt.want)
			}
		})
	}
}

func TestDiscardSetupItems(t *testing.T) {
	inputCarSetup := map[CarSetupKey]*CarSetup{
		"TiresAero|LeftRearTire|StartingPressure":     {},
		"TiresAero|RightFrontTire|LastHotPressure":    {},
		"TiresAero|RightRearTire|LastTempsIMO":        {},
		"TiresAero|TireCompound|TireCompound":         {},
		"DriveBrake|BrakeSystemConfig|TotalBrakeBias": {},
	}

	type args struct {
		setupItems map[CarSetupKey]*CarSetup
		filters    []SetupFilter
	}
	tests := []struct {
		name string
		args args
		want map[CarSetupKey]*CarSetup
	}{
		{
			"discard tiresaero pressures",
			args{
				setupItems: inputCarSetup,
				filters:    []SetupFilter{{"TiresAero", "", "Pressure"}},
			},
			map[CarSetupKey]*CarSetup{
				"TiresAero|RightRearTire|LastTempsIMO":        inputCarSetup["TiresAero|RightRearTire|LastTempsIMO"],
				"TiresAero|TireCompound|TireCompound":         inputCarSetup["TiresAero|TireCompound|TireCompound"],
				"DriveBrake|BrakeSystemConfig|TotalBrakeBias": inputCarSetup["DriveBrake|BrakeSystemConfig|TotalBrakeBias"],
			},
		},
		{
			"discard tire data",
			args{
				setupItems: inputCarSetup,
				filters:    []SetupFilter{{"", "Tire", ""}},
			},
			map[CarSetupKey]*CarSetup{
				"DriveBrake|BrakeSystemConfig|TotalBrakeBias": inputCarSetup["DriveBrake|BrakeSystemConfig|TotalBrakeBias"],
			},
		},
		{
			"discard brake bias",
			args{
				setupItems: inputCarSetup,
				filters:    []SetupFilter{{"Brake", "", "Total"}},
			},
			map[CarSetupKey]*CarSetup{
				"TiresAero|LeftRearTire|StartingPressure":  inputCarSetup["TiresAero|LeftRearTire|StartingPressure"],
				"TiresAero|RightFrontTire|LastHotPressure": inputCarSetup["TiresAero|RightFrontTire|LastHotPressure"],
				"TiresAero|RightRearTire|LastTempsIMO":     inputCarSetup["TiresAero|RightRearTire|LastTempsIMO"],
				"TiresAero|TireCompound|TireCompound":      inputCarSetup["TiresAero|TireCompound|TireCompound"],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiscardSetupItems(tt.args.setupItems, tt.args.filters...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiscardSetupItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

var expectedFullCarSetup = &CarSetup{
	Name:   "ARA_23S1_W13_RBR_R_2.sto",
	Update: 11,
	Values: CarSetupDetails{
		"Chassis|Front|HeaveRate": &CarSetupItem{
			RawValue: "750 N/mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  750.0,
					NumericalSign:   0,
					MeasurementUnit: "N/mm",
				},
			},
		},
		"Chassis|Front|RideHeight": &CarSetupItem{
			RawValue: "25.0 mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  25.0,
					NumericalSign:   0,
					MeasurementUnit: "mm",
				},
			},
		},
		"Chassis|Front|RollRate": &CarSetupItem{
			RawValue: "400 N/mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  400.0,
					NumericalSign:   0,
					MeasurementUnit: "N/mm",
				},
			},
		},
		"Chassis|LeftFront|Camber": &CarSetupItem{
			RawValue: "-3.15 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  3.15,
					NumericalSign:   -1,
					MeasurementUnit: "deg",
				},
			},
		},
		"Chassis|LeftFront|CornerWeight": &CarSetupItem{
			RawValue: "1902 N",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  1902.0,
					NumericalSign:   0,
					MeasurementUnit: "N",
				},
			},
		},
		"Chassis|LeftFront|ToeIn": &CarSetupItem{
			RawValue: "-0.05 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  0.05,
					NumericalSign:   -1,
					MeasurementUnit: "deg",
				},
			},
		},
		"Chassis|LeftRear|Camber": &CarSetupItem{
			RawValue: "-1.57 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  1.57,
					NumericalSign:   -1,
					MeasurementUnit: "deg",
				},
			},
		},
		"Chassis|LeftRear|CornerWeight": &CarSetupItem{
			RawValue: "2231 N",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  2231.0,
					NumericalSign:   0,
					MeasurementUnit: "N",
				},
			},
		},
		"Chassis|LeftRear|ToeIn": &CarSetupItem{
			RawValue: "+0.05 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  0.05,
					NumericalSign:   1,
					MeasurementUnit: "deg",
				},
			},
		},
		"Chassis|Rear|FuelLevel": &CarSetupItem{
			RawValue: "45 Kg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  45.0,
					NumericalSign:   0,
					MeasurementUnit: "Kg",
				},
			},
		},
		"Chassis|Rear|HeaveRate": &CarSetupItem{
			RawValue: "180 N/mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  180.0,
					NumericalSign:   0,
					MeasurementUnit: "N/mm",
				},
			},
		},
		"Chassis|Rear|RideHeight": &CarSetupItem{
			RawValue: "60.0 mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  60.0,
					NumericalSign:   0,
					MeasurementUnit: "mm",
				},
			},
		},
		"Chassis|Rear|RollRate": &CarSetupItem{
			RawValue: "300 N/mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  300.0,
					NumericalSign:   0,
					MeasurementUnit: "N/mm",
				},
			},
		},
		"Chassis|RightFront|Camber": &CarSetupItem{
			RawValue: "-3.15 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  3.15,
					NumericalSign:   -1,
					MeasurementUnit: "deg",
				},
			},
		},
		"Chassis|RightFront|CornerWeight": &CarSetupItem{
			RawValue: "1902 N",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  1902.0,
					NumericalSign:   0,
					MeasurementUnit: "N",
				},
			},
		},
		"Chassis|RightFront|ToeIn": &CarSetupItem{
			RawValue: "-0.05 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  0.05,
					NumericalSign:   -1,
					MeasurementUnit: "deg",
				},
			},
		},
		"Chassis|RightRear|Camber": &CarSetupItem{
			RawValue: "-1.57 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  1.57,
					NumericalSign:   -1,
					MeasurementUnit: "deg",
				},
			},
		},
		"Chassis|RightRear|CornerWeight": &CarSetupItem{
			RawValue: "2231 N",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  2231.0,
					NumericalSign:   0,
					MeasurementUnit: "N",
				},
			},
		},
		"Chassis|RightRear|ToeIn": &CarSetupItem{
			RawValue: "+0.05 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  0.05,
					NumericalSign:   1,
					MeasurementUnit: "deg",
				},
			},
		},
		"DriveBrake|BrakeSystemConfig|BaseBrakeBias": &CarSetupItem{
			RawValue: "52.5% (BBAL)",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  52.5,
					NumericalSign:   0,
					MeasurementUnit: "% (BBAL)",
				},
			},
		},
		"DriveBrake|BrakeSystemConfig|BrakeMigration": &CarSetupItem{
			RawValue: "6 (BMIG)",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  6.0,
					NumericalSign:   0,
					MeasurementUnit: "(BMIG)",
				},
			},
		},
		"DriveBrake|BrakeSystemConfig|DynamicRamping": &CarSetupItem{
			RawValue: "20% pedal",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  20.0,
					NumericalSign:   0,
					MeasurementUnit: "% pedal",
				},
			},
		},
		"DriveBrake|BrakeSystemConfig|FineBrakeBias": &CarSetupItem{
			RawValue: "-1.0 (BB+/BB-)",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  1.0,
					NumericalSign:   -1,
					MeasurementUnit: "(BB+/BB-)",
				},
			},
		},
		"DriveBrake|BrakeSystemConfig|TotalBrakeBias": &CarSetupItem{
			RawValue: "56.5% front",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  56.5,
					NumericalSign:   0,
					MeasurementUnit: "% front",
				},
			},
		},
		"DriveBrake|Differential|Entry": &CarSetupItem{
			RawValue: "1 (ENTRY)",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  1.0,
					NumericalSign:   0,
					MeasurementUnit: "(ENTRY)",
				},
			},
		},
		"DriveBrake|Differential|HighSpeed": &CarSetupItem{
			RawValue: "6 (HISPD)",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  6.0,
					NumericalSign:   0,
					MeasurementUnit: "(HISPD)",
				},
			},
		},
		"DriveBrake|Differential|Middle": &CarSetupItem{
			RawValue: "4 (MID)",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  4.0,
					NumericalSign:   0,
					MeasurementUnit: "(MID)",
				},
			},
		},
		"DriveBrake|Differential|Preload": &CarSetupItem{
			RawValue: "0 Nm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  0.0,
					NumericalSign:   0,
					MeasurementUnit: "Nm",
				},
			},
		},
		"DriveBrake|PowerUnitConfig|EngineBraking": &CarSetupItem{
			RawValue: "5 (EB)",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  5.0,
					NumericalSign:   0,
					MeasurementUnit: "(EB)",
				},
			},
		},
		"DriveBrake|PowerUnitConfig|MguKDeployMode": &CarSetupItem{
			RawValue: "Balanced",
			Parsed:   []CarSetupItemParsedValue{}, // p0
		},
		"TiresAero|AeroCalculator|AeroBalance": &CarSetupItem{
			RawValue: "48.14%",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  48.14,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
			},
		},
		"TiresAero|AeroCalculator|DownforceToDrag": &CarSetupItem{
			RawValue: "4.501:1",
			Parsed:   make([]CarSetupItemParsedValue, 0),
		},
		"TiresAero|AeroCalculator|FrontRhAtSpeed": &CarSetupItem{
			RawValue: "18.0 mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  18.0,
					NumericalSign:   0,
					MeasurementUnit: "mm",
				},
			},
		},
		"TiresAero|AeroCalculator|RearRhAtSpeed": &CarSetupItem{
			RawValue: "33.0 mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  33.0,
					NumericalSign:   0,
					MeasurementUnit: "mm",
				},
			},
		},
		"TiresAero|AeroPackage|DownforceTrim": &CarSetupItem{
			RawValue: "High",
			Parsed:   make([]CarSetupItemParsedValue, 0),
		},
		"TiresAero|AeroPackage|FrontFlapOffset": &CarSetupItem{
			RawValue: "0.50 deg",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  0.5,
					NumericalSign:   0,
					MeasurementUnit: "deg",
				},
			},
		},
		"TiresAero|AeroPackage|RearWingGurney": &CarSetupItem{
			RawValue: "0 mm",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  0.0,
					NumericalSign:   0,
					MeasurementUnit: "mm",
				},
			},
		},
		"TiresAero|LeftFrontTire|LastHotPressure": &CarSetupItem{
			RawValue: "24.1 psi",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  24.1,
					NumericalSign:   0,
					MeasurementUnit: "psi",
				},
			},
		},
		"TiresAero|LeftFrontTire|LastTempsOMI": &CarSetupItem{
			RawValue: "89C, 96C, 101C",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  89.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  96.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  101.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
			},
		},
		"TiresAero|LeftFrontTire|StartingPressure": &CarSetupItem{
			RawValue: "165.5 kPa",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  165.5,
					NumericalSign:   0,
					MeasurementUnit: "kPa",
				},
			},
		},
		"TiresAero|LeftFrontTire|TreadRemaining": &CarSetupItem{
			RawValue: "99%, 98%, 98%",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  99.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
			},
		},
		"TiresAero|LeftRearTire|LastHotPressure": &CarSetupItem{
			RawValue: "21.2 psi",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  21.2,
					NumericalSign:   0,
					MeasurementUnit: "psi",
				},
			},
		},
		"TiresAero|LeftRearTire|LastTempsOMI": &CarSetupItem{
			RawValue: "88C, 93C, 95C",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  88.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  93.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  95.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
			},
		},
		"TiresAero|LeftRearTire|StartingPressure": &CarSetupItem{
			RawValue: "144.8 kPa",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  144.8,
					NumericalSign:   0,
					MeasurementUnit: "kPa",
				},
			},
		},
		"TiresAero|LeftRearTire|TreadRemaining": &CarSetupItem{
			RawValue: "99%, 98%, 98%",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  99.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
			},
		},
		"TiresAero|RightFrontTire|LastHotPressure": &CarSetupItem{
			RawValue: "23.9 psi",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  23.9,
					NumericalSign:   0,
					MeasurementUnit: "psi",
				},
			},
		},
		"TiresAero|RightFrontTire|LastTempsIMO": &CarSetupItem{
			RawValue: "98C, 93C, 86C",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  93.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  86.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
			},
		},
		"TiresAero|RightFrontTire|StartingPressure": &CarSetupItem{
			RawValue: "165.5 kPa",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  165.5,
					NumericalSign:   0,
					MeasurementUnit: "kPa",
				},
			},
		},
		"TiresAero|RightFrontTire|TreadRemaining": &CarSetupItem{
			RawValue: "97%, 97%, 98%",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  97.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  97.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
			},
		},
		"TiresAero|RightRearTire|LastHotPressure": &CarSetupItem{
			RawValue: "21.0 psi",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  21.0,
					NumericalSign:   0,
					MeasurementUnit: "psi",
				},
			},
		},
		"TiresAero|RightRearTire|LastTempsIMO": &CarSetupItem{
			RawValue: "93C, 90C, 84C",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  93.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  90.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
				{
					NumericalValue:  84.0,
					NumericalSign:   0,
					MeasurementUnit: "C",
				},
			},
		},
		"TiresAero|RightRearTire|StartingPressure": &CarSetupItem{
			RawValue: "144.8 kPa",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  144.8,
					NumericalSign:   0,
					MeasurementUnit: "kPa",
				},
			},
		},
		"TiresAero|RightRearTire|TreadRemaining": &CarSetupItem{
			RawValue: "98%, 98%, 99%",
			Parsed: []CarSetupItemParsedValue{
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  98.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
				{
					NumericalValue:  99.0,
					NumericalSign:   0,
					MeasurementUnit: "%",
				},
			},
		},
		"TiresAero|TireCompound|TireCompound": &CarSetupItem{
			RawValue: "Medium",
			Parsed:   make([]CarSetupItemParsedValue, 0),
		},
	},
}

var comparisonCarSetup = &CarSetup{
	Name: "comparison setup",
	Values: map[CarSetupKey]*CarSetupItem{
		"Chassis|Front|HeaveRate": &CarSetupItem{
			"780 N/mm",
			[]CarSetupItemParsedValue{
				{
					NumericalValue:  780.0,
					NumericalSign:   0,
					MeasurementUnit: "N/mm",
				},
			},
		},
		"TiresAero|TireCompound|TireCompound": &CarSetupItem{
			RawValue: "Soft",
		},
	},
}
