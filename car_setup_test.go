package ibt

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/teamjorge/ibt/headers"
	"golang.org/x/exp/maps"
)

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

	t.Run("full setup test", func(t *testing.T) {

		setup := ParseCarSetup(header.SessionInfo)

		if setup.Name != expectedFullCarSetup.Name {
			t.Errorf("expected setup name to be %s. received: %s", expectedFullCarSetup.Name, setup.Name)
		}

		if setup.Update != expectedFullCarSetup.Update {
			t.Errorf("expected setup update to be %d. received: %d", expectedFullCarSetup.Update, setup.Update)
		}

		for itemName, itemValue := range expectedFullCarSetup.Values {
			if !reflect.DeepEqual(setup.Values[itemName], itemValue) {
				t.Errorf("item %s has unexpected values. \nexpected: %+v\n \nactual: %+v\n",
					itemName, itemValue, setup.Values[itemName])
			}
		}
	})
}

func TestCarSetupDifference(t *testing.T) {
	files, err := filepath.Glob(".testing/mercedesw13*.ibt")
	if err != nil {
		t.Error(err)
	}

	stubs, err := ParseStubs(files...)
	if err != nil {
		t.Error(err)
	}

	setups := make([]*CarSetup, 0)

	for _, stub := range stubs {
		setups = append(setups, stub.CarSetup())
	}

	for i := 0; i < len(setups); i++ {
		if i+1 >= len(setups) {
			break
		}

		comparison := CompareSetups(setups[i], setups[i+1]).Differences()
		if len(comparison) == 0 {
			continue
		}

		comparisonKeys := CarSetupKeys(maps.Keys(comparison))
		sort.Sort(comparisonKeys)

		for _, key := range comparisonKeys {
			diff := comparison[key]

			fmt.Printf("%s - %s - %s - %s\n", key.Category(), key.SubCategory(), key.ItemName(), diff.RawDifference)
		}

		fmt.Println()
	}

	t.Fail()
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
