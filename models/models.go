package models

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Car struct {
	Name   string
	States []Carstate
}

func (c *Car) AddState(s Carstate) {
	c.States = append(c.States, s)
}

type Carstate struct {
	Timestamp           string
	IsRaceOn            bool
	TimestampMS         uint32
	EngineMaxRPM        float64
	EngineIdleRPM       float64
	EngineCurrentRPM    float64
	AccelX              float64
	AccelY              float64
	AccelZ              float64
	VelX                float64
	VelY                float64
	VelZ                float64
	AngVelX             float64
	AngVelY             float64
	AngVelZ             float64
	Yaw                 float64
	Pitch               float64
	Roll                float64
	NormSuspFL          float64
	NormSuspFR          float64
	NormSuspRL          float64
	NormSuspRR          float64
	TireSlipFL          float64
	TireSlipFR          float64
	TireSlipRL          float64
	TireSlipRR          float64
	WheelRotFL          float64
	WheelRotFR          float64
	WheelRotRL          float64
	WheelRotRR          float64
	WheelOnRumbleFL     float64
	WheelOnRumbleFR     float64
	WheelOnRumbleRL     float64
	WheelOnRumbleRR     float64
	WheelInPuddleFL     float64
	WheelInPuddleFR     float64
	WheelInPuddleRL     float64
	WheelInPuddleRR     float64
	SurfaceRumbleFL     float64
	SurfaceRumbleFR     float64
	SurfaceRumbleRL     float64
	SurfaceRumbleRR     float64
	TireSlipAngleFL     float64
	TireSlipAngleFR     float64
	TireSlipAngleRL     float64
	TireSlipAngleRR     float64
	TireCombinedSlipFL  float64
	TireCombinedSlipFR  float64
	TireCombinedSlipRL  float64
	TireCombinedSlipRR  float64
	SuspTravelFL        float64
	SuspTravelFR        float64
	SuspTravelRL        float64
	SuspTravelRR        float64
	CarOrdinal          int
	CarClass            int
	CarPerformanceIndex int
	DrivetrainType      int
	NumCylinders        int
	PosX                float64
	PosY                float64
	PosZ                float64
	SpeedMPS            float64
	SpeedKPH            float64
	SpeedMPH            float64
	Power               float64
	Torque              float64
	TireTempFL          float64
	TireTempFR          float64
	TireTempRL          float64
	TireTempRR          float64
	Boost               float64
	Fuel                float64
	Distance            float64
	BestLap             float64
	LastLap             float64
	CurrentLap          float64
	CurrentRaceTime     float64
	LapNumber           int
	RacePosition        int
	Accel               int
	Brake               int
	Clutch              int
	Handbrake           int
	Gear                int
	Steer               int
	NormDrivingLine     int
	NormAIBrakeDiff     int
}

// ExportCSV writes a car's telemetry to a CSV file named <CarName>.csv
func (c *Car) ExportCSV() error {
	filename := fmt.Sprintf("%s.csv", c.Name)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Header row
	header := []string{
		"timestamp",
		"isRaceOn",
		"timestampMS",
		"engine_max_rpm",
		"engine_idle_rpm",
		"engine_current_rpm",
		"accel_x",
		"accel_y",
		"accel_z",
		"vel_x",
		"vel_y",
		"vel_z",
		"ang_vel_x",
		"ang_vel_y",
		"ang_vel_z",
		"yaw",
		"pitch",
		"roll",
		"norm_susp_fl",
		"norm_susp_fr",
		"norm_susp_rl",
		"norm_susp_rr",
		"tire_slip_fl",
		"tire_slip_fr",
		"tire_slip_rl",
		"tire_slip_rr",
		"wheel_rot_fl",
		"wheel_rot_fr",
		"wheel_rot_rl",
		"wheel_rot_rr",
		"wheel_on_rumble_fl",
		"wheel_on_rumble_fr",
		"wheel_on_rumble_rl",
		"wheel_on_rumble_rr",
		"wheel_in_puddle_fl",
		"wheel_in_puddle_fr",
		"wheel_in_puddle_rl",
		"wheel_in_puddle_rr",
		"surface_rumble_fl",
		"surface_rumble_fr",
		"surface_rumble_rl",
		"surface_rumble_rr",
		"tire_slip_angle_fl",
		"tire_slip_angle_fr",
		"tire_slip_angle_rl",
		"tire_slip_angle_rr",
		"tire_combined_slip_fl",
		"tire_combined_slip_fr",
		"tire_combined_slip_rl",
		"tire_combined_slip_rr",
		"susp_travel_fl",
		"susp_travel_fr",
		"susp_travel_rl",
		"susp_travel_rr",
		"car_ordinal",
		"car_class",
		"car_performance_index",
		"drivetrain_type",
		"num_cylinders",
		"pos_x",
		"pos_y",
		"pos_z",
		"speed_mps",
		"speed_kph",
		"speed_mph",
		"power",
		"torque",
		"tire_temp_fl",
		"tire_temp_fr",
		"tire_temp_rl",
		"tire_temp_rr",
		"boost",
		"fuel",
		"distance",
		"best_lap",
		"last_lap",
		"current_lap",
		"current_race_time",
		"lap_number",
		"race_position",
		"accel",
		"brake",
		"clutch",
		"handbrake",
		"gear",
		"steer",
		"norm_driving_line",
		"norm_ai_brake_diff",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	// Data rows
	for _, s := range c.States {
		record := []string{
			s.Timestamp,
			boolToStr(s.IsRaceOn),
			u32(s.TimestampMS),
			f64(s.EngineMaxRPM),
			f64(s.EngineIdleRPM),
			f64(s.EngineCurrentRPM),
			f64(s.AccelX),
			f64(s.AccelY),
			f64(s.AccelZ),
			f64(s.VelX),
			f64(s.VelY),
			f64(s.VelZ),
			f64(s.AngVelX),
			f64(s.AngVelY),
			f64(s.AngVelZ),
			f64(s.Yaw),
			f64(s.Pitch),
			f64(s.Roll),
			f64(s.NormSuspFL),
			f64(s.NormSuspFR),
			f64(s.NormSuspRL),
			f64(s.NormSuspRR),
			f64(s.TireSlipFL),
			f64(s.TireSlipFR),
			f64(s.TireSlipRL),
			f64(s.TireSlipRR),
			f64(s.WheelRotFL),
			f64(s.WheelRotFR),
			f64(s.WheelRotRL),
			f64(s.WheelRotRR),
			f64(s.WheelOnRumbleFL),
			f64(s.WheelOnRumbleFR),
			f64(s.WheelOnRumbleRL),
			f64(s.WheelOnRumbleRR),
			f64(s.WheelInPuddleFL),
			f64(s.WheelInPuddleFR),
			f64(s.WheelInPuddleRL),
			f64(s.WheelInPuddleRR),
			f64(s.SurfaceRumbleFL),
			f64(s.SurfaceRumbleFR),
			f64(s.SurfaceRumbleRL),
			f64(s.SurfaceRumbleRR),
			f64(s.TireSlipAngleFL),
			f64(s.TireSlipAngleFR),
			f64(s.TireSlipAngleRL),
			f64(s.TireSlipAngleRR),
			f64(s.TireCombinedSlipFL),
			f64(s.TireCombinedSlipFR),
			f64(s.TireCombinedSlipRL),
			f64(s.TireCombinedSlipRR),
			f64(s.SuspTravelFL),
			f64(s.SuspTravelFR),
			f64(s.SuspTravelRL),
			f64(s.SuspTravelRR),
			i(s.CarOrdinal),
			i(s.CarClass),
			i(s.CarPerformanceIndex),
			i(s.DrivetrainType),
			i(s.NumCylinders),
			f64(s.PosX),
			f64(s.PosY),
			f64(s.PosZ),
			f64(s.SpeedMPS),
			f64(s.SpeedKPH),
			f64(s.SpeedMPH),
			f64(s.Power),
			f64(s.Torque),
			f64(s.TireTempFL),
			f64(s.TireTempFR),
			f64(s.TireTempRL),
			f64(s.TireTempRR),
			f64(s.Boost),
			f64(s.Fuel),
			f64(s.Distance),
			f64(s.BestLap),
			f64(s.LastLap),
			f64(s.CurrentLap),
			f64(s.CurrentRaceTime),
			i(s.LapNumber),
			i(s.RacePosition),
			i(s.Accel),
			i(s.Brake),
			i(s.Clutch),
			i(s.Handbrake),
			i(s.Gear),
			i(s.Steer),
			i(s.NormDrivingLine),
			i(s.NormAIBrakeDiff),
		}

		if err := w.Write(record); err != nil {
			return err
		}
	}

	fmt.Println("Wrote CSV:", filename)
	return nil
}

// Helpers
func f64(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
func u32(v uint32) string  { return strconv.FormatUint(uint64(v), 10) }
func i(v int) string       { return strconv.Itoa(v) }
func boolToStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
