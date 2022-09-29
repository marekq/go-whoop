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
