package headers

import (
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

// ReadSessionInfo extracts and parses the SessionInfo YAML from the given ibt file.
//
// Additional cleaning (mostly trimming some trailing bytes and spaces) is needed to ensure the YAML is correctly parsed.
func ReadSessionInfo(reader Reader, offset, size int) (*Session, error) {
	sessionBuf := make([]byte, size)

	_, err := reader.ReadAt(sessionBuf, int64(offset))
	if err != nil {
		return nil, err
	}

	dec := charmap.Windows1252.NewDecoder()

	sessionBuf, err = dec.Bytes(sessionBuf)
	if err != nil {
		return nil, err
	}

	rawYaml := strings.TrimRight(string(sessionBuf[:]), "\x00")
	rawYaml = strings.TrimRight(rawYaml, ".")
	rawYaml = strings.TrimSpace(rawYaml)

	output := Session{}

	if err := yaml.Unmarshal([]byte(rawYaml), &output); err != nil {
		return nil, err
	}

	return &output, nil
}

// GetDriver attempts to find the current driver's information from the SessionInfo object.
//
// If the driver is not found (possibly spectating), it will return nil.
func (s *Session) GetDriver() *Drivers {
	idx := s.DriverInfo.DriverCarIdx

	for _, driver := range s.DriverInfo.Drivers {
		if driver.CarIdx == idx {
			return &driver
		}
	}

	return nil
}

// Session in which the ibt file occurred. This represents an actual iRacing session and not just information for a single ibt file.
type Session struct {
	CameraInfo    CameraInfo             `yaml:"CameraInfo"`
	CarSetup      map[string]interface{} `yaml:"CarSetup"`
	DriverInfo    DriverInfo             `yaml:"DriverInfo"`
	RadioInfo     RadioInfo              `yaml:"RadioInfo"`
	SessionInfo   SessionInfo            `yaml:"SessionInfo"`
	SplitTimeInfo SplitTimeInfo          `yaml:"SplitTimeInfo"`
	WeekendInfo   WeekendInfo            `yaml:"WeekendInfo"`
}

// Cameras available for the a given camera group.
type Cameras struct {
	CameraName string `yaml:"CameraName"`
	CameraNum  int    `yaml:"CameraNum"`
}

// Groups of cameras available for spectating.
//
// For example: TV1 and it's related cameras.
type Groups struct {
	Cameras   []Cameras `yaml:"Cameras"`
	GroupName string    `yaml:"GroupName"`
	GroupNum  int       `yaml:"GroupNum"`
	IsScenic  bool      `yaml:"IsScenic,omitempty"`
}

// CameraInfo contains all available.
type CameraInfo struct {
	Groups []Groups `yaml:"Groups"`
}

// Drivers provides information on a driver present in the session.
type Drivers struct {
	AbbrevName              interface{} `yaml:"AbbrevName"`
	CarClassColor           int         `yaml:"CarClassColor"`
	CarClassDryTireSetLimit string      `yaml:"CarClassDryTireSetLimit"`
	CarClassEstLapTime      float64     `yaml:"CarClassEstLapTime"`
	CarClassID              int         `yaml:"CarClassID"`
	CarClassLicenseLevel    int         `yaml:"CarClassLicenseLevel"`
	CarClassMaxFuelPct      string      `yaml:"CarClassMaxFuelPct"`
	CarClassPowerAdjust     string      `yaml:"CarClassPowerAdjust"`
	CarClassRelSpeed        int         `yaml:"CarClassRelSpeed"`
	CarClassShortName       interface{} `yaml:"CarClassShortName"`
	CarClassWeightPenalty   string      `yaml:"CarClassWeightPenalty"`
	CarDesignStr            string      `yaml:"CarDesignStr"`
	CarID                   int         `yaml:"CarID"`
	CarIdx                  int         `yaml:"CarIdx"`
	CarIsAI                 int         `yaml:"CarIsAI"`
	CarIsElectric           int         `yaml:"CarIsElectric"`
	CarIsPaceCar            int         `yaml:"CarIsPaceCar"`
	CarNumber               string      `yaml:"CarNumber"`
	CarNumberDesignStr      string      `yaml:"CarNumberDesignStr"`
	CarNumberRaw            int         `yaml:"CarNumberRaw"`
	CarPath                 string      `yaml:"CarPath"`
	CarScreenName           string      `yaml:"CarScreenName"`
	CarScreenNameShort      string      `yaml:"CarScreenNameShort"`
	CarSponsor1             int         `yaml:"CarSponsor_1"`
	CarSponsor2             int         `yaml:"CarSponsor_2"`
	CurDriverIncidentCount  int         `yaml:"CurDriverIncidentCount"`
	HelmetDesignStr         string      `yaml:"HelmetDesignStr"`
	IRating                 int         `yaml:"IRating"`
	Initials                interface{} `yaml:"Initials"`
	IsSpectator             int         `yaml:"IsSpectator"`
	LicColor                string      `yaml:"LicColor"`
	LicLevel                int         `yaml:"LicLevel"`
	LicString               string      `yaml:"LicString"`
	LicSubLevel             int         `yaml:"LicSubLevel"`
	SuitDesignStr           string      `yaml:"SuitDesignStr"`
	TeamID                  int         `yaml:"TeamID"`
	TeamIncidentCount       int         `yaml:"TeamIncidentCount"`
	TeamName                string      `yaml:"TeamName"`
	UserID                  int         `yaml:"UserID"`
	UserName                string      `yaml:"UserName"`
}

// DriverInfo provides information regarding the driver in the car.
type DriverInfo struct {
	DriverCarEngCylinderCount int       `yaml:"DriverCarEngCylinderCount"`
	DriverCarEstLapTime       float64   `yaml:"DriverCarEstLapTime"`
	DriverCarFuelKgPerLtr     float64   `yaml:"DriverCarFuelKgPerLtr"`
	DriverCarFuelMaxLtr       int       `yaml:"DriverCarFuelMaxLtr"`
	DriverCarGearNeutral      int       `yaml:"DriverCarGearNeutral"`
	DriverCarGearNumForward   int       `yaml:"DriverCarGearNumForward"`
	DriverCarGearReverse      int       `yaml:"DriverCarGearReverse"`
	DriverCarIdleRPM          int       `yaml:"DriverCarIdleRPM"`
	DriverCarIdx              int       `yaml:"DriverCarIdx"`
	DriverCarIsElectric       int       `yaml:"DriverCarIsElectric"`
	DriverCarMaxFuelPct       int       `yaml:"DriverCarMaxFuelPct"`
	DriverCarRedLine          int       `yaml:"DriverCarRedLine"`
	DriverCarSLBlinkRPM       int       `yaml:"DriverCarSLBlinkRPM"`
	DriverCarSLFirstRPM       int       `yaml:"DriverCarSLFirstRPM"`
	DriverCarSLLastRPM        int       `yaml:"DriverCarSLLastRPM"`
	DriverCarSLShiftRPM       int       `yaml:"DriverCarSLShiftRPM"`
	DriverCarVersion          string    `yaml:"DriverCarVersion"`
	DriverHeadPosX            float64   `yaml:"DriverHeadPosX"`
	DriverHeadPosY            float64   `yaml:"DriverHeadPosY"`
	DriverHeadPosZ            float64   `yaml:"DriverHeadPosZ"`
	DriverIncidentCount       int       `yaml:"DriverIncidentCount"`
	DriverPitTrkPct           float64   `yaml:"DriverPitTrkPct"`
	DriverSetupIsModified     int       `yaml:"DriverSetupIsModified"`
	DriverSetupLoadTypeName   string    `yaml:"DriverSetupLoadTypeName"`
	DriverSetupName           string    `yaml:"DriverSetupName"`
	DriverSetupPassedTech     int       `yaml:"DriverSetupPassedTech"`
	DriverUserID              int       `yaml:"DriverUserID"`
	Drivers                   []Drivers `yaml:"Drivers"`
	PaceCarIdx                int       `yaml:"PaceCarIdx"`
}

// Frequencies provide information of of the given radio channel.
type Frequencies struct {
	CanScan       int    `yaml:"CanScan"`
	CanSquawk     int    `yaml:"CanSquawk"`
	CarIdx        int    `yaml:"CarIdx"`
	ClubID        int    `yaml:"ClubID"`
	EntryIdx      int    `yaml:"EntryIdx"`
	FrequencyName string `yaml:"FrequencyName"`
	FrequencyNum  int    `yaml:"FrequencyNum"`
	IsDeletable   int    `yaml:"IsDeletable"`
	IsMutable     int    `yaml:"IsMutable"`
	Muted         int    `yaml:"Muted"`
	Priority      int    `yaml:"Priority"`
}

// Radios is a single available radio channel and its related frequencies.
type Radios struct {
	Frequencies         []Frequencies `yaml:"Frequencies"`
	HopCount            int           `yaml:"HopCount"`
	NumFrequencies      int           `yaml:"NumFrequencies"`
	RadioNum            int           `yaml:"RadioNum"`
	ScanningIsOn        int           `yaml:"ScanningIsOn"`
	TunedToFrequencyNum int           `yaml:"TunedToFrequencyNum"`
}

// RadioInfo contains all available radios and the currently selected radio.
type RadioInfo struct {
	Radios           []Radios `yaml:"Radios"`
	SelectedRadioNum int      `yaml:"SelectedRadioNum"`
}

// ResultsFastestLap provides information regarding the session's fastest lap.
type ResultsFastestLap struct {
	CarIdx      int `yaml:"CarIdx"`
	FastestLap  int `yaml:"FastestLap"`
	FastestTime int `yaml:"FastestTime"`
}

// Sessions provides information for a sub-session, such as Practice or Qualifying or Race.
type Sessions struct {
	ResultsAverageLapTime            int                 `yaml:"ResultsAverageLapTime"`
	ResultsFastestLap                []ResultsFastestLap `yaml:"ResultsFastestLap"`
	ResultsLapsComplete              int                 `yaml:"ResultsLapsComplete"`
	ResultsNumCautionFlags           int                 `yaml:"ResultsNumCautionFlags"`
	ResultsNumCautionLaps            int                 `yaml:"ResultsNumCautionLaps"`
	ResultsNumLeadChanges            int                 `yaml:"ResultsNumLeadChanges"`
	ResultsOfficial                  int                 `yaml:"ResultsOfficial"`
	ResultsPositions                 interface{}         `yaml:"ResultsPositions"`
	SessionEnforceTireCompoundChange int                 `yaml:"SessionEnforceTireCompoundChange"`
	SessionLaps                      string              `yaml:"SessionLaps"`
	SessionName                      string              `yaml:"SessionName"`
	SessionNum                       int                 `yaml:"SessionNum"`
	SessionNumLapsToAvg              int                 `yaml:"SessionNumLapsToAvg"`
	SessionRunGroupsUsed             int                 `yaml:"SessionRunGroupsUsed"`
	SessionSkipped                   int                 `yaml:"SessionSkipped"`
	SessionSubType                   interface{}         `yaml:"SessionSubType"`
	SessionTime                      string              `yaml:"SessionTime"`
	SessionTrackRubberState          string              `yaml:"SessionTrackRubberState"`
	SessionType                      string              `yaml:"SessionType"`
}

// SessionInfo is simply a collection of sub-sessions.
type SessionInfo struct {
	Sessions []Sessions `yaml:"Sessions"`
}

// Sectors provide details for a single track sector.
type Sectors struct {
	SectorNum      int     `yaml:"SectorNum"`
	SectorStartPct float64 `yaml:"SectorStartPct"`
}

// SplitTimeInfo contains all track sector information.
type SplitTimeInfo struct {
	Sectors []Sectors `yaml:"Sectors"`
}
type TelemetryOptions struct {
	TelemetryDiskFile string `yaml:"TelemetryDiskFile"`
}

// WeekendOptions contains session rules, physics, environment, and weather information.
type WeekendOptions struct {
	CommercialMode             string `yaml:"CommercialMode"`
	CourseCautions             string `yaml:"CourseCautions"`
	Date                       string `yaml:"Date"`
	EarthRotationSpeedupFactor int    `yaml:"EarthRotationSpeedupFactor"`
	FastRepairsLimit           string `yaml:"FastRepairsLimit"`
	FogLevel                   string `yaml:"FogLevel"`
	GreenWhiteCheckeredLimit   int    `yaml:"GreenWhiteCheckeredLimit"`
	HardcoreLevel              int    `yaml:"HardcoreLevel"`
	HasOpenRegistration        int    `yaml:"HasOpenRegistration"`
	IncidentLimit              string `yaml:"IncidentLimit"`
	IsFixedSetup               int    `yaml:"IsFixedSetup"`
	NightMode                  string `yaml:"NightMode"`
	NumJokerLaps               int    `yaml:"NumJokerLaps"`
	NumStarters                int    `yaml:"NumStarters"`
	QualifyScoring             string `yaml:"QualifyScoring"`
	RelativeHumidity           string `yaml:"RelativeHumidity"`
	Restarts                   string `yaml:"Restarts"`
	ShortParadeLap             int    `yaml:"ShortParadeLap"`
	Skies                      string `yaml:"Skies"`
	StandingStart              int    `yaml:"StandingStart"`
	StartingGrid               string `yaml:"StartingGrid"`
	StrictLapsChecking         string `yaml:"StrictLapsChecking"`
	TimeOfDay                  string `yaml:"TimeOfDay"`
	Unofficial                 int    `yaml:"Unofficial"`
	WeatherTemp                string `yaml:"WeatherTemp"`
	WeatherType                string `yaml:"WeatherType"`
	WindDirection              string `yaml:"WindDirection"`
	WindSpeed                  string `yaml:"WindSpeed"`
}

// WeekendInfo contains all session metadata.
type WeekendInfo struct {
	BuildTarget            string           `yaml:"BuildTarget"`
	BuildType              string           `yaml:"BuildType"`
	BuildVersion           string           `yaml:"BuildVersion"`
	Category               string           `yaml:"Category"`
	DCRuleSet              string           `yaml:"DCRuleSet"`
	EventType              string           `yaml:"EventType"`
	HeatRacing             int              `yaml:"HeatRacing"`
	LeagueID               int              `yaml:"LeagueID"`
	MaxDrivers             int              `yaml:"MaxDrivers"`
	MinDrivers             int              `yaml:"MinDrivers"`
	NumCarClasses          int              `yaml:"NumCarClasses"`
	NumCarTypes            int              `yaml:"NumCarTypes"`
	Official               int              `yaml:"Official"`
	QualifierMustStartRace int              `yaml:"QualifierMustStartRace"`
	RaceWeek               int              `yaml:"RaceWeek"`
	SeasonID               int              `yaml:"SeasonID"`
	SeriesID               int              `yaml:"SeriesID"`
	SessionID              int              `yaml:"SessionID"`
	SimMode                string           `yaml:"SimMode"`
	SubSessionID           int              `yaml:"SubSessionID"`
	TeamRacing             int              `yaml:"TeamRacing"`
	TelemetryOptions       TelemetryOptions `yaml:"TelemetryOptions"`
	TrackAirPressure       string           `yaml:"TrackAirPressure"`
	TrackAirTemp           string           `yaml:"TrackAirTemp"`
	TrackAltitude          string           `yaml:"TrackAltitude"`
	TrackCity              string           `yaml:"TrackCity"`
	TrackCleanup           int              `yaml:"TrackCleanup"`
	TrackConfigName        string           `yaml:"TrackConfigName"`
	TrackCountry           string           `yaml:"TrackCountry"`
	TrackDirection         string           `yaml:"TrackDirection"`
	TrackDisplayName       string           `yaml:"TrackDisplayName"`
	TrackDisplayShortName  string           `yaml:"TrackDisplayShortName"`
	TrackDynamicTrack      int              `yaml:"TrackDynamicTrack"`
	TrackFogLevel          string           `yaml:"TrackFogLevel"`
	TrackID                int              `yaml:"TrackID"`
	TrackLatitude          string           `yaml:"TrackLatitude"`
	TrackLength            string           `yaml:"TrackLength"`
	TrackLengthOfficial    string           `yaml:"TrackLengthOfficial"`
	TrackLongitude         string           `yaml:"TrackLongitude"`
	TrackName              string           `yaml:"TrackName"`
	TrackNorthOffset       string           `yaml:"TrackNorthOffset"`
	TrackNumTurns          int              `yaml:"TrackNumTurns"`
	TrackPitSpeedLimit     string           `yaml:"TrackPitSpeedLimit"`
	TrackRelativeHumidity  string           `yaml:"TrackRelativeHumidity"`
	TrackSkies             string           `yaml:"TrackSkies"`
	TrackSurfaceTemp       string           `yaml:"TrackSurfaceTemp"`
	TrackType              string           `yaml:"TrackType"`
	TrackVersion           string           `yaml:"TrackVersion"`
	TrackWeatherType       string           `yaml:"TrackWeatherType"`
	TrackWindDir           string           `yaml:"TrackWindDir"`
	TrackWindVel           string           `yaml:"TrackWindVel"`
	WeekendOptions         WeekendOptions   `yaml:"WeekendOptions"`
}
