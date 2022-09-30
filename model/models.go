package model

import "time"

type Sleep struct {
	Records []struct {
		ID             int       `json:"id"`
		UserID         int       `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Start          time.Time `json:"start"`
		End            time.Time `json:"end"`
		TimezoneOffset string    `json:"timezone_offset"`
		Nap            bool      `json:"nap"`
		ScoreState     string    `json:"score_state"`
		Score          struct {
			StageSummary struct {
				TotalInBedTimeMilli         int `json:"total_in_bed_time_milli"`
				TotalAwakeTimeMilli         int `json:"total_awake_time_milli"`
				TotalNoDataTimeMilli        int `json:"total_no_data_time_milli"`
				TotalLightSleepTimeMilli    int `json:"total_light_sleep_time_milli"`
				TotalSlowWaveSleepTimeMilli int `json:"total_slow_wave_sleep_time_milli"`
				TotalRemSleepTimeMilli      int `json:"total_rem_sleep_time_milli"`
				SleepCycleCount             int `json:"sleep_cycle_count"`
				DisturbanceCount            int `json:"disturbance_count"`
			} `json:"stage_summary"`
			SleepNeeded struct {
				BaselineMilli             int `json:"baseline_milli"`
				NeedFromSleepDebtMilli    int `json:"need_from_sleep_debt_milli"`
				NeedFromRecentStrainMilli int `json:"need_from_recent_strain_milli"`
				NeedFromRecentNapMilli    int `json:"need_from_recent_nap_milli"`
			} `json:"sleep_needed"`
			RespiratoryRate            float64 `json:"respiratory_rate"`
			SleepPerformancePercentage float64 `json:"sleep_performance_percentage"`
			SleepConsistencyPercentage float64 `json:"sleep_consistency_percentage"`
			SleepEfficiencyPercentage  float64 `json:"sleep_efficiency_percentage"`
		} `json:"score"`
	} `json:"records"`
	NextToken string `json:"next_token"`
}

type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

type Recovery struct {
	Records []struct {
		CycleID    int       `json:"cycle_id"`
		SleepID    int       `json:"sleep_id"`
		UserID     int       `json:"user_id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		ScoreState string    `json:"score_state"`
		Score      struct {
			UserCalibrating  bool    `json:"user_calibrating"`
			RecoveryScore    float64 `json:"recovery_score"`
			RestingHeartRate float64 `json:"resting_heart_rate"`
			HrvRmssdMilli    float64 `json:"hrv_rmssd_milli"`
			Spo2Percentage   float64 `json:"spo2_percentage"`
			SkinTempCelsius  float64 `json:"skin_temp_celsius"`
		} `json:"score"`
	} `json:"records"`
	NextToken string `json:"next_token"`
}

type Cycle struct {
	Records []struct {
		ID             int         `json:"id"`
		UserID         int         `json:"user_id"`
		CreatedAt      time.Time   `json:"created_at"`
		UpdatedAt      time.Time   `json:"updated_at"`
		Start          time.Time   `json:"start"`
		End            interface{} `json:"end"`
		TimezoneOffset string      `json:"timezone_offset"`
		ScoreState     string      `json:"score_state"`
		Score          struct {
			Strain           float64 `json:"strain"`
			Kilojoule        float64 `json:"kilojoule"`
			AverageHeartRate int     `json:"average_heart_rate"`
			MaxHeartRate     int     `json:"max_heart_rate"`
		} `json:"score"`
	} `json:"records"`
	NextToken string `json:"next_token"`
}

type Workout struct {
	Records []struct {
		ID             int       `json:"id"`
		UserID         int       `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Start          time.Time `json:"start"`
		End            time.Time `json:"end"`
		TimezoneOffset string    `json:"timezone_offset"`
		SportID        int       `json:"sport_id"`
		ScoreState     string    `json:"score_state"`
		Score          struct {
			Strain              float64 `json:"strain"`
			AverageHeartRate    int     `json:"average_heart_rate"`
			MaxHeartRate        int     `json:"max_heart_rate"`
			Kilojoule           float64 `json:"kilojoule"`
			PercentRecorded     float64 `json:"percent_recorded"`
			DistanceMeter       float64 `json:"distance_meter"`
			AltitudeGainMeter   float64 `json:"altitude_gain_meter"`
			AltitudeChangeMeter float64 `json:"altitude_change_meter"`
			ZoneDuration        struct {
				ZoneZeroMilli  int `json:"zone_zero_milli"`
				ZoneOneMilli   int `json:"zone_one_milli"`
				ZoneTwoMilli   int `json:"zone_two_milli"`
				ZoneThreeMilli int `json:"zone_three_milli"`
				ZoneFourMilli  int `json:"zone_four_milli"`
				ZoneFiveMilli  int `json:"zone_five_milli"`
			} `json:"zone_duration"`
		} `json:"score"`
	} `json:"records"`
	NextToken interface{} `json:"next_token"`
}

type All struct {
	Records []struct {
		ID             int       `json:"id"`
		UserID         int       `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Start          time.Time `json:"start"`
		End            time.Time `json:"end"`
		TimezoneOffset string    `json:"timezone_offset"`
		Nap            bool      `json:"nap"`

		SportID    int    `json:"sport_id"`
		ScoreState string `json:"score_state"`
		Score      struct {
			UserCalibrating  bool    `json:"user_calibrating"`
			RecoveryScore    float64 `json:"recovery_score"`
			RestingHeartRate float64 `json:"resting_heart_rate"`
			HrvRmssdMilli    float64 `json:"hrv_rmssd_milli"`
			Spo2Percentage   float64 `json:"spo2_percentage"`
			SkinTempCelsius  float64 `json:"skin_temp_celsius"`
			StageSummary     struct {
				TotalInBedTimeMilli         int `json:"total_in_bed_time_milli"`
				TotalAwakeTimeMilli         int `json:"total_awake_time_milli"`
				TotalNoDataTimeMilli        int `json:"total_no_data_time_milli"`
				TotalLightSleepTimeMilli    int `json:"total_light_sleep_time_milli"`
				TotalSlowWaveSleepTimeMilli int `json:"total_slow_wave_sleep_time_milli"`
				TotalRemSleepTimeMilli      int `json:"total_rem_sleep_time_milli"`
				SleepCycleCount             int `json:"sleep_cycle_count"`
				DisturbanceCount            int `json:"disturbance_count"`
			} `json:"stage_summary"`
			SleepNeeded struct {
				BaselineMilli             int `json:"baseline_milli"`
				NeedFromSleepDebtMilli    int `json:"need_from_sleep_debt_milli"`
				NeedFromRecentStrainMilli int `json:"need_from_recent_strain_milli"`
				NeedFromRecentNapMilli    int `json:"need_from_recent_nap_milli"`
			} `json:"sleep_needed"`
			RespiratoryRate            float64 `json:"respiratory_rate"`
			SleepPerformancePercentage float64 `json:"sleep_performance_percentage"`
			SleepConsistencyPercentage float64 `json:"sleep_consistency_percentage"`
			SleepEfficiencyPercentage  float64 `json:"sleep_efficiency_percentage"`
			Strain                     float64 `json:"strain"`
			AverageHeartRate           int     `json:"average_heart_rate"`
			MaxHeartRate               int     `json:"max_heart_rate"`
			Kilojoule                  float64 `json:"kilojoule"`
			PercentRecorded            float64 `json:"percent_recorded"`
			DistanceMeter              float64 `json:"distance_meter"`
			AltitudeGainMeter          float64 `json:"altitude_gain_meter"`
			AltitudeChangeMeter        float64 `json:"altitude_change_meter"`
			ZoneDuration               struct {
				ZoneZeroMilli  int `json:"zone_zero_milli"`
				ZoneOneMilli   int `json:"zone_one_milli"`
				ZoneTwoMilli   int `json:"zone_two_milli"`
				ZoneThreeMilli int `json:"zone_three_milli"`
				ZoneFourMilli  int `json:"zone_four_milli"`
				ZoneFiveMilli  int `json:"zone_five_milli"`
			} `json:"zone_duration"`
		} `json:"score"`
	} `json:"records"`
	NextToken string `json:"next_token"`
}
