package headers

import (
	"os"
	"reflect"
	"testing"
)

func TestSessionInfoHeader(t *testing.T) {
	t.Run("valid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		output, err := ReadSessionInfo(f, expectedTelemetryHeader.SessionInfoOffset, expectedTelemetryHeader.SessionInfoLength)
		if err != nil {
			t.Errorf("failed to parse session header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(*output, expectedSessionInfo) {
			t.Errorf("expected session header does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedSessionInfo, *output)
		}
	})

	t.Run("invalid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/invalid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ReadSessionInfo(f, expectedTelemetryHeader.SessionInfoOffset, expectedTelemetryHeader.SessionInfoLength)
		if err == nil {
			t.Error("expected session parsing of invalid file to return an error")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		f, err := os.Open("../.testing/empty_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ReadSessionInfo(f, expectedTelemetryHeader.SessionInfoOffset, expectedTelemetryHeader.SessionInfoLength)
		if err == nil {
			t.Error("expected session parsing of empty file to return an error")
		}
	})

	t.Run("test GetDriver valid", func(t *testing.T) {
		driver := expectedSessionInfo.GetDriver()

		if !reflect.DeepEqual(*driver, expectedDriver) {
			t.Errorf("expected GetDriver output does not match actual. \nexpected: %+v\n \nactual: %+v\n", driver, expectedDriver)
		}
	})

	t.Run("test GetDriver not found", func(t *testing.T) {
		newSession := expectedSessionInfo
		newSession.DriverInfo.DriverCarIdx = 9999

		driver := newSession.GetDriver()

		if driver != nil {
			t.Errorf("expected GetDriver result to be nil. received: %+v\n", driver)
		}
	})
}

var expectedDriver = Drivers{
	AbbrevName:              nil,
	CarClassColor:           16777215,
	CarClassDryTireSetLimit: "0 %",
	CarClassEstLapTime:      69.3118,
	CarClassID:              0,
	CarClassLicenseLevel:    0,
	CarClassMaxFuelPct:      "1.000 %",
	CarClassPowerAdjust:     "0.000 %",
	CarClassRelSpeed:        0,
	CarClassShortName:       nil,
	CarClassWeightPenalty:   "0.000 kg",
	CarDesignStr:            "10,ff00bf,ff00bf,00d1ff",
	CarID:                   161,
	CarIdx:                  0,
	CarIsAI:                 0,
	CarIsElectric:           0,
	CarIsPaceCar:            0,
	CarNumber:               "64",
	CarNumberDesignStr:      "0,0,ffffff,777777,000000",
	CarNumberRaw:            64,
	CarPath:                 "mercedesw13",
	CarScreenName:           "Mercedes-AMG W13 E Performance",
	CarScreenNameShort:      "Mercedes W13",
	CarSponsor1:             0,
	CarSponsor2:             0,
	CurDriverIncidentCount:  12,
	HelmetDesignStr:         "53,00d1ff,ff00bf,00d1ff",
	IRating:                 1,
	Initials:                nil,
	IsSpectator:             0,
	LicColor:                "0xundefined",
	LicLevel:                1,
	LicString:               "R 0.01",
	LicSubLevel:             1,
	SuitDesignStr:           "17,00d1ff,00d1ff,00d1ff",
	TeamID:                  0,
	TeamIncidentCount:       12,
	TeamName:                "George v Rensburg",
	UserID:                  450313,
	UserName:                "George v Rensburg",
}

var expectedSessionInfo = Session{
	CameraInfo: CameraInfo{
		Groups: []Groups{
			{
				Cameras: []Cameras{
					{
						CameraName: "CamNose",
						CameraNum:  1,
					},
				},
				GroupName: "Nose",
				GroupNum:  1,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamGearbox",
						CameraNum:  1,
					},
				},
				GroupName: "Gearbox",
				GroupNum:  2,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamRoll Bar",
						CameraNum:  1,
					},
				},
				GroupName: "Roll Bar",
				GroupNum:  3,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamLF Susp",
						CameraNum:  1,
					},
				},
				GroupName: "LF Susp",
				GroupNum:  4,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamLR Susp",
						CameraNum:  1,
					},
				},
				GroupName: "LR Susp",
				GroupNum:  5,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamGyro",
						CameraNum:  1,
					},
				},
				GroupName: "Gyro",
				GroupNum:  6,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamRF Susp",
						CameraNum:  1,
					},
				},
				GroupName: "RF Susp",
				GroupNum:  7,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamRR Susp",
						CameraNum:  1,
					},
				},
				GroupName: "RR Susp",
				GroupNum:  8,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamCockpit",
						CameraNum:  1,
					},
				},
				GroupName: "Cockpit",
				GroupNum:  9,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "Scenic_01",
						CameraNum:  1,
					},
					{
						CameraName: "Scenic_02",
						CameraNum:  2,
					},
					{
						CameraName: "Scenic_03",
						CameraNum:  3,
					},
					{
						CameraName: "Scenic_04",
						CameraNum:  4,
					},
					{
						CameraName: "Scenic_05",
						CameraNum:  5,
					},
					{
						CameraName: "Scenic_06",
						CameraNum:  6,
					},
					{
						CameraName: "Scenic_07",
						CameraNum:  7,
					},
					{
						CameraName: "Scenic_08",
						CameraNum:  8,
					},
					{
						CameraName: "Scenic_09",
						CameraNum:  9,
					},
					{
						CameraName: "Scenic_10",
						CameraNum:  10,
					},
				},
				GroupName: "Scenic",
				GroupNum:  10,
				IsScenic:  true,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamTV1_00",
						CameraNum:  1,
					},
					{
						CameraName: "CamTV1_02",
						CameraNum:  2,
					},
					{
						CameraName: "CamTV1_03",
						CameraNum:  3,
					},
					{
						CameraName: "CamTV1_04",
						CameraNum:  4,
					},
					{
						CameraName: "CamTV1_05",
						CameraNum:  5,
					},
					{
						CameraName: "CamTV1_06",
						CameraNum:  6,
					},
					{
						CameraName: "CamTV1_07",
						CameraNum:  7,
					},
					{
						CameraName: "CamTV1_08",
						CameraNum:  8,
					},
					{
						CameraName: "CamTV1_01",
						CameraNum:  9,
					},
				},
				GroupName: "TV1",
				GroupNum:  11,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamTV2_10",
						CameraNum:  1,
					},
					{
						CameraName: "CamTV2_03",
						CameraNum:  2,
					},
					{
						CameraName: "CamTV2_04",
						CameraNum:  3,
					},
					{
						CameraName: "CamTV2_05",
						CameraNum:  4,
					},
					{
						CameraName: "CamTV2_06",
						CameraNum:  5,
					},
					{
						CameraName: "CamTV2_07",
						CameraNum:  6,
					},
					{
						CameraName: "CamTV2_08",
						CameraNum:  7,
					},
					{
						CameraName: "CamTV2_01",
						CameraNum:  8,
					},
					{
						CameraName: "CamTV2_00",
						CameraNum:  9,
					},
					{
						CameraName: "CamTV2_02",
						CameraNum:  10,
					},
					{
						CameraName: "CamTV2_09",
						CameraNum:  11,
					},
				},
				GroupName: "TV2",
				GroupNum:  12,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamTV3_00",
						CameraNum:  1,
					},
					{
						CameraName: "CamTV3_02",
						CameraNum:  2,
					},
					{
						CameraName: "CamTV3_03",
						CameraNum:  3,
					},
					{
						CameraName: "CamTV3_04",
						CameraNum:  4,
					},
					{
						CameraName: "CamTV3_06",
						CameraNum:  5,
					},
					{
						CameraName: "CamTV3_05",
						CameraNum:  6,
					},
					{
						CameraName: "CamTV3_07",
						CameraNum:  7,
					},
					{
						CameraName: "CamTV3_08",
						CameraNum:  8,
					},
					{
						CameraName: "CamTV3_01",
						CameraNum:  9,
					},
				},
				GroupName: "TV3",
				GroupNum:  13,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamTV4_00",
						CameraNum:  1,
					},
					{
						CameraName: "CamTV4_14",
						CameraNum:  2,
					},
					{
						CameraName: "CamTV4_01",
						CameraNum:  3,
					},
					{
						CameraName: "CamTV4_02",
						CameraNum:  4,
					},
					{
						CameraName: "CamTV4_03",
						CameraNum:  5,
					},
					{
						CameraName: "CamTV4_05",
						CameraNum:  6,
					},
					{
						CameraName: "CamTV4_06",
						CameraNum:  7,
					},
					{
						CameraName: "CamTV4_08",
						CameraNum:  8,
					},
					{
						CameraName: "CamTV4_07",
						CameraNum:  9,
					},
					{
						CameraName: "CamTV4_09",
						CameraNum:  10,
					},
					{
						CameraName: "CamTV4_10",
						CameraNum:  11,
					},
					{
						CameraName: "CamTV4_11",
						CameraNum:  12,
					},
					{
						CameraName: "CamTV4_12",
						CameraNum:  13,
					},
					{
						CameraName: "CamTV4_13",
						CameraNum:  14,
					},
				},
				GroupName: "TV Static",
				GroupNum:  14,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamTV1_00",
						CameraNum:  1,
					},
					{
						CameraName: "CamTV1_01",
						CameraNum:  2,
					},
					{
						CameraName: "CamTV1_02",
						CameraNum:  3,
					},
					{
						CameraName: "CamTV1_03",
						CameraNum:  4,
					},
					{
						CameraName: "CamTV1_04",
						CameraNum:  5,
					},
					{
						CameraName: "CamTV1_05",
						CameraNum:  6,
					},
					{
						CameraName: "CamTV1_06",
						CameraNum:  7,
					},
					{
						CameraName: "CamTV1_07",
						CameraNum:  8,
					},
					{
						CameraName: "CamTV1_08",
						CameraNum:  9,
					},
					{
						CameraName: "CamTV2_00",
						CameraNum:  10,
					},
					{
						CameraName: "CamTV2_01",
						CameraNum:  11,
					},
					{
						CameraName: "CamTV2_02",
						CameraNum:  12,
					},
					{
						CameraName: "CamTV2_03",
						CameraNum:  13,
					},
					{
						CameraName: "CamTV2_04",
						CameraNum:  14,
					},
					{
						CameraName: "CamTV2_05",
						CameraNum:  15,
					},
					{
						CameraName: "CamTV2_06",
						CameraNum:  16,
					},
					{
						CameraName: "CamTV2_07",
						CameraNum:  17,
					},
					{
						CameraName: "CamTV2_08",
						CameraNum:  18,
					},
					{
						CameraName: "CamTV2_09",
						CameraNum:  19,
					},
					{
						CameraName: "CamTV2_10",
						CameraNum:  20,
					},
					{
						CameraName: "CamTV3_00",
						CameraNum:  21,
					},
					{
						CameraName: "CamTV3_01",
						CameraNum:  22,
					},
					{
						CameraName: "CamTV3_02",
						CameraNum:  23,
					},
					{
						CameraName: "CamTV3_03",
						CameraNum:  24,
					},
					{
						CameraName: "CamTV3_04",
						CameraNum:  25,
					},
					{
						CameraName: "CamTV3_05",
						CameraNum:  26,
					},
					{
						CameraName: "CamTV3_06",
						CameraNum:  27,
					},
					{
						CameraName: "CamTV3_07",
						CameraNum:  28,
					},
					{
						CameraName: "CamTV3_08",
						CameraNum:  29,
					},
					{
						CameraName: "CamTV4_07",
						CameraNum:  30,
					},
					{
						CameraName: "CamTV4_08",
						CameraNum:  31,
					},
					{
						CameraName: "CamTV4_09",
						CameraNum:  32,
					},
					{
						CameraName: "CamTV4_11",
						CameraNum:  33,
					},
					{
						CameraName: "CamTV4_13",
						CameraNum:  34,
					},
					{
						CameraName: "CamTV4_14",
						CameraNum:  35,
					},
				},
				GroupName: "TV Mixed",
				GroupNum:  15,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamPit Lane",
						CameraNum:  1,
					},
				},
				GroupName: "Pit Lane",
				GroupNum:  16,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamPit Lane 2",
						CameraNum:  1,
					},
				},
				GroupName: "Pit Lane 2",
				GroupNum:  17,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamBlimp",
						CameraNum:  1,
					},
				},
				GroupName: "Blimp",
				GroupNum:  18,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamChopper",
						CameraNum:  1,
					},
				},
				GroupName: "Chopper",
				GroupNum:  19,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamChase",
						CameraNum:  1,
					},
				},
				GroupName: "Chase",
				GroupNum:  20,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamFar Chase",
						CameraNum:  1,
					},
				},
				GroupName: "Far Chase",
				GroupNum:  21,
				IsScenic:  false,
			},
			{
				Cameras: []Cameras{
					{
						CameraName: "CamRear Chase",
						CameraNum:  1,
					},
				},
				GroupName: "Rear Chase",
				GroupNum:  22,
				IsScenic:  false,
			},
		},
	},
	CarSetup: map[string]interface{}{
		"Chassis": map[string]interface{}{
			"Front": map[string]interface{}{
				"HeaveRate":  "750 N/mm",
				"RideHeight": "25.0 mm",
				"RollRate":   "400 N/mm",
			},
			"LeftFront": map[string]interface{}{
				"Camber":       "-3.15 deg",
				"CornerWeight": "1902 N",
				"ToeIn":        "-0.05 deg",
			},
			"LeftRear": map[string]interface{}{
				"Camber":       "-1.57 deg",
				"CornerWeight": "2231 N",
				"ToeIn":        "+0.05 deg",
			},
			"Rear": map[string]interface{}{
				"FuelLevel":  "45 Kg",
				"HeaveRate":  "180 N/mm",
				"RideHeight": "60.0 mm",
				"RollRate":   "300 N/mm",
			},
			"RightFront": map[string]interface{}{
				"Camber":       "-3.15 deg",
				"CornerWeight": "1902 N",
				"ToeIn":        "-0.05 deg",
			},
			"RightRear": map[string]interface{}{
				"Camber":       "-1.57 deg",
				"CornerWeight": "2231 N",
				"ToeIn":        "+0.05 deg",
			},
		},
		"DriveBrake": map[string]interface{}{
			"BrakeSystemConfig": map[string]interface{}{
				"BaseBrakeBias":  "52.5% (BBAL)",
				"BrakeMigration": "6 (BMIG)",
				"DynamicRamping": "20% pedal",
				"FineBrakeBias":  "-1.0 (BB+/BB-)",
				"TotalBrakeBias": "56.5% front",
			},
			"Differential": map[string]interface{}{
				"Entry":     "1 (ENTRY)",
				"HighSpeed": "6 (HISPD)",
				"Middle":    "4 (MID)",
				"Preload":   "0 Nm",
			},
			"PowerUnitConfig": map[string]interface{}{
				"EngineBraking":  "5 (EB)",
				"MguKDeployMode": "Balanced",
			},
		},
		"TiresAero": map[string]interface{}{
			"AeroCalculator": map[string]interface{}{
				"AeroBalance":     "48.14%",
				"DownforceToDrag": "4.501:1",
				"FrontRhAtSpeed":  "18.0 mm",
				"RearRhAtSpeed":   "33.0 mm",
			},
			"AeroPackage": map[string]interface{}{
				"DownforceTrim":   "High",
				"FrontFlapOffset": "0.50 deg",
				"RearWingGurney":  "0 mm",
			},
			"LeftFrontTire": map[string]interface{}{
				"LastHotPressure":  "24.1 psi",
				"LastTempsOMI":     "89C, 96C, 101C",
				"StartingPressure": "165.5 kPa",
				"TreadRemaining":   "99%, 98%, 98%",
			},
			"LeftRearTire": map[string]interface{}{
				"LastHotPressure":  "21.2 psi",
				"LastTempsOMI":     "88C, 93C, 95C",
				"StartingPressure": "144.8 kPa",
				"TreadRemaining":   "99%, 98%, 98%",
			},
			"RightFrontTire": map[string]interface{}{
				"LastHotPressure":  "23.9 psi",
				"LastTempsIMO":     "98C, 93C, 86C",
				"StartingPressure": "165.5 kPa",
				"TreadRemaining":   "97%, 97%, 98%",
			},
			"RightRearTire": map[string]interface{}{
				"LastHotPressure":  "21.0 psi",
				"LastTempsIMO":     "93C, 90C, 84C",
				"StartingPressure": "144.8 kPa",
				"TreadRemaining":   "98%, 98%, 99%",
			},
			"TireCompound": map[string]interface{}{
				"TireCompound": "Medium",
			},
		},
		"UpdateCount": 11,
	},
	DriverInfo: DriverInfo{
		DriverCarEngCylinderCount: 6,
		DriverCarEstLapTime:       69.3118,
		DriverCarFuelKgPerLtr:     0.75,
		DriverCarFuelMaxLtr:       146,
		DriverCarGearNeutral:      1,
		DriverCarGearNumForward:   8,
		DriverCarGearReverse:      1,
		DriverCarIdleRPM:          4000,
		DriverCarIdx:              0,
		DriverCarIsElectric:       0,
		DriverCarMaxFuelPct:       1,
		DriverCarRedLine:          13000,
		DriverCarSLBlinkRPM:       11629,
		DriverCarSLFirstRPM:       10560,
		DriverCarSLLastRPM:        11473,
		DriverCarSLShiftRPM:       11057,
		DriverCarVersion:          "2024.05.28.02",
		DriverHeadPosX:            0.386,
		DriverHeadPosY:            0.0,
		DriverHeadPosZ:            0.341,
		DriverIncidentCount:       12,
		DriverPitTrkPct:           0.055219,
		DriverSetupIsModified:     0,
		DriverSetupLoadTypeName:   "user",
		DriverSetupName:           "ARA_23S1_W13_RBR_R_2.sto",
		DriverSetupPassedTech:     1,
		DriverUserID:              450313,
		Drivers: []Drivers{
			{
				AbbrevName:              nil,
				CarClassColor:           16777215,
				CarClassDryTireSetLimit: "0 %",
				CarClassEstLapTime:      69.3118,
				CarClassID:              0,
				CarClassLicenseLevel:    0,
				CarClassMaxFuelPct:      "1.000 %",
				CarClassPowerAdjust:     "0.000 %",
				CarClassRelSpeed:        0,
				CarClassShortName:       nil,
				CarClassWeightPenalty:   "0.000 kg",
				CarDesignStr:            "10,ff00bf,ff00bf,00d1ff",
				CarID:                   161,
				CarIdx:                  0,
				CarIsAI:                 0,
				CarIsElectric:           0,
				CarIsPaceCar:            0,
				CarNumber:               "64",
				CarNumberDesignStr:      "0,0,ffffff,777777,000000",
				CarNumberRaw:            64,
				CarPath:                 "mercedesw13",
				CarScreenName:           "Mercedes-AMG W13 E Performance",
				CarScreenNameShort:      "Mercedes W13",
				CarSponsor1:             0,
				CarSponsor2:             0,
				CurDriverIncidentCount:  12,
				HelmetDesignStr:         "53,00d1ff,ff00bf,00d1ff",
				IRating:                 1,
				Initials:                nil,
				IsSpectator:             0,
				LicColor:                "0xundefined",
				LicLevel:                1,
				LicString:               "R 0.01",
				LicSubLevel:             1,
				SuitDesignStr:           "17,00d1ff,00d1ff,00d1ff",
				TeamID:                  0,
				TeamIncidentCount:       12,
				TeamName:                "George v Rensburg",
				UserID:                  450313,
				UserName:                "George v Rensburg",
			},
		},
		PaceCarIdx: -1,
	},
	RadioInfo: RadioInfo{
		Radios: []Radios{
			{
				Frequencies: []Frequencies{
					{
						CanScan:       1,
						CanSquawk:     1,
						CarIdx:        -1,
						ClubID:        0,
						EntryIdx:      -1,
						FrequencyName: "@ALLTEAMS",
						FrequencyNum:  0,
						IsDeletable:   0,
						IsMutable:     1,
						Muted:         0,
						Priority:      12,
					},
					{
						CanScan:       1,
						CanSquawk:     1,
						CarIdx:        -1,
						ClubID:        0,
						EntryIdx:      -1,
						FrequencyName: "@DRIVERS",
						FrequencyNum:  1,
						IsDeletable:   0,
						IsMutable:     1,
						Muted:         0,
						Priority:      15,
					},
					{
						CanScan:       1,
						CanSquawk:     1,
						CarIdx:        0,
						ClubID:        0,
						EntryIdx:      -1,
						FrequencyName: "@TEAM",
						FrequencyNum:  2,
						IsDeletable:   0,
						IsMutable:     0,
						Muted:         0,
						Priority:      60,
					},
					{
						CanScan:       1,
						CanSquawk:     1,
						CarIdx:        -1,
						ClubID:        50,
						EntryIdx:      -1,
						FrequencyName: "@CLUB",
						FrequencyNum:  3,
						IsDeletable:   0,
						IsMutable:     1,
						Muted:         0,
						Priority:      20,
					},
					{
						CanScan:       1,
						CanSquawk:     1,
						CarIdx:        -1,
						ClubID:        0,
						EntryIdx:      -1,
						FrequencyName: "@ADMIN",
						FrequencyNum:  4,
						IsDeletable:   0,
						IsMutable:     0,
						Muted:         0,
						Priority:      90,
					},
					{
						CanScan:       1,
						CanSquawk:     1,
						CarIdx:        -1,
						ClubID:        0,
						EntryIdx:      -1,
						FrequencyName: "@RACECONTROL",
						FrequencyNum:  5,
						IsDeletable:   0,
						IsMutable:     0,
						Muted:         0,
						Priority:      80,
					},
					{
						CanScan:       1,
						CanSquawk:     1,
						CarIdx:        -1,
						ClubID:        0,
						EntryIdx:      0,
						FrequencyName: "@PRIVATE",
						FrequencyNum:  6,
						IsDeletable:   0,
						IsMutable:     0,
						Muted:         0,
						Priority:      70,
					},
				},
				HopCount:            2,
				NumFrequencies:      7,
				RadioNum:            0,
				ScanningIsOn:        1,
				TunedToFrequencyNum: 0,
			},
		},
		SelectedRadioNum: 0,
	},
	SessionInfo: SessionInfo{
		Sessions: []Sessions{
			{
				ResultsAverageLapTime: -1,
				ResultsFastestLap: []ResultsFastestLap{
					{
						CarIdx:      0,
						FastestLap:  4,
						FastestTime: 68,
					},
				},
				ResultsLapsComplete:    -1,
				ResultsNumCautionFlags: 0,
				ResultsNumCautionLaps:  0,
				ResultsNumLeadChanges:  0,
				ResultsOfficial:        0,
				ResultsPositions: []interface{}{
					map[string]interface{}{
						"CarIdx":            0,
						"ClassPosition":     0,
						"FastestLap":        4,
						"FastestTime":       68.6711,
						"Incidents":         0,
						"JokerLapsComplete": 0,
						"Lap":               4,
						"LapsComplete":      8,
						"LapsDriven":        0.0,
						"LapsLed":           0,
						"LastTime":          -1.0,
						"Position":          1,
						"ReasonOutId":       0,
						"ReasonOutStr":      "Running",
						"Time":              68.6711,
					},
				},
				SessionEnforceTireCompoundChange: 0,
				SessionLaps:                      "unlimited",
				SessionName:                      "TESTING",
				SessionNum:                       0,
				SessionNumLapsToAvg:              0,
				SessionRunGroupsUsed:             0,
				SessionSkipped:                   0,
				SessionSubType:                   nil,
				SessionTime:                      "unlimited",
				SessionTrackRubberState:          "moderate usage",
				SessionType:                      "Offline Testing",
			},
		},
	},
	SplitTimeInfo: SplitTimeInfo{
		Sectors: []Sectors{
			{
				SectorNum:      0,
				SectorStartPct: 0.0,
			},
			{
				SectorNum:      1,
				SectorStartPct: 0.271918,
			},
			{
				SectorNum:      2,
				SectorStartPct: 0.668198,
			},
		},
	},
	WeekendInfo: WeekendInfo{
		BuildTarget:            "Members",
		BuildType:              "Release",
		BuildVersion:           "2024.06.10.01",
		Category:               "Road",
		DCRuleSet:              "None",
		EventType:              "Test",
		HeatRacing:             0,
		LeagueID:               0,
		MaxDrivers:             0,
		MinDrivers:             0,
		NumCarClasses:          1,
		NumCarTypes:            1,
		Official:               0,
		QualifierMustStartRace: 0,
		RaceWeek:               0,
		SeasonID:               0,
		SeriesID:               0,
		SessionID:              0,
		SimMode:                "full",
		SubSessionID:           0,
		TeamRacing:             0,
		TelemetryOptions: TelemetryOptions{
			TelemetryDiskFile: "",
		},
		TrackAirPressure:      "27.69 Hg",
		TrackAirTemp:          "23.89 C",
		TrackAltitude:         "677.30 m",
		TrackCity:             "Spielberg",
		TrackCleanup:          1,
		TrackConfigName:       "Grand Prix",
		TrackCountry:          "Austria",
		TrackDirection:        "neutral",
		TrackDisplayName:      "Red Bull Ring",
		TrackDisplayShortName: "Spielberg",
		TrackDynamicTrack:     1,
		TrackFogLevel:         "0 %",
		TrackID:               403,
		TrackLatitude:         "47.220305 m",
		TrackLength:           "4.28 km",
		TrackLengthOfficial:   "4.32 km",
		TrackLongitude:        "14.766722 m",
		TrackName:             "spielberg gp",
		TrackNorthOffset:      "1.5876 rad",
		TrackNumTurns:         9,
		TrackPitSpeedLimit:    "80.00 kph",
		TrackRelativeHumidity: "55 %",
		TrackSkies:            "Clear",
		TrackSurfaceTemp:      "38.89 C",
		TrackType:             "road course",
		TrackVersion:          "2024.05.22.01",
		TrackWeatherType:      "Static",
		TrackWindDir:          "2.36 rad",
		TrackWindVel:          "4.02 m/s",
		WeekendOptions: WeekendOptions{
			CommercialMode:             "consumer",
			CourseCautions:             "off",
			Date:                       "2024-04-01",
			EarthRotationSpeedupFactor: 1,
			FastRepairsLimit:           "unlimited",
			FogLevel:                   "0 %",
			GreenWhiteCheckeredLimit:   0,
			HardcoreLevel:              1,
			HasOpenRegistration:        0,
			IncidentLimit:              "unlimited",
			IsFixedSetup:               0,
			NightMode:                  "variable",
			NumJokerLaps:               0,
			NumStarters:                0,
			QualifyScoring:             "best lap",
			RelativeHumidity:           "55 %",
			Restarts:                   "single file",
			ShortParadeLap:             0,
			Skies:                      "Clear",
			StandingStart:              0,
			StartingGrid:               "single file",
			StrictLapsChecking:         "default",
			TimeOfDay:                  "12:00 pm",
			Unofficial:                 1,
			WeatherTemp:                "23.89 C",
			WeatherType:                "Static",
			WindDirection:              "SE",
			WindSpeed:                  "14.48 km/h",
		},
	},
}
