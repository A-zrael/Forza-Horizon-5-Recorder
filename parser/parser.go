package parser

import (
	"encoding/binary"
	"fmt"
	"forza/models"
	"math"
	"time"
)

func RawtoCarstate(data []byte) (models.Carstate, error) {
	const (
		requiredBytes    = 323 // need up through byte 322 (indexing from 0)
		cutStart         = 232 // bytes we keep before the gap
		cutEnd           = 244 // start of the chunk to skip
		patchedTotalSize = 311
	)

	if len(data) < requiredBytes {
		return models.Carstate{}, fmt.Errorf("packet too small: %d bytes", len(data))
	}

	// Mirror the Python parser: drop bytes 232–243 (12 bytes) before unpacking.
	patched := make([]byte, 0, patchedTotalSize)
	patched = append(patched, data[:cutStart]...)
	patched = append(patched, data[cutEnd:requiredBytes]...)

	f32 := func(off int) float64 {
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(patched[off : off+4])))
	}
	i32 := func(off int) int {
		return int(int32(binary.LittleEndian.Uint32(patched[off : off+4])))
	}
	u16 := func(off int) uint16 {
		return binary.LittleEndian.Uint16(patched[off : off+2])
	}
	i8 := func(off int) int8 { return int8(patched[off]) }
	u8 := func(off int) uint8 { return patched[off] }

	speedFromVec := math.Sqrt(f32(32)*f32(32) + f32(36)*f32(36) + f32(40)*f32(40))
	speedMps := f32(244)
	if speedMps == 0 {
		speedMps = speedFromVec
	}

	return models.Carstate{
		Timestamp:           time.Now().Format(time.RFC3339Nano),
		IsRaceOn:            i32(0) != 0,
		TimestampMS:         binary.LittleEndian.Uint32(patched[4:8]),
		EngineMaxRPM:        f32(8),
		EngineIdleRPM:       f32(12),
		EngineCurrentRPM:    f32(16),
		AccelX:              f32(20),
		AccelY:              f32(24),
		AccelZ:              f32(28),
		VelX:                f32(32),
		VelY:                f32(36),
		VelZ:                f32(40),
		AngVelX:             f32(44),
		AngVelY:             f32(48),
		AngVelZ:             f32(52),
		Yaw:                 f32(56),
		Pitch:               f32(60),
		Roll:                f32(64),
		NormSuspFL:          f32(68),
		NormSuspFR:          f32(72),
		NormSuspRL:          f32(76),
		NormSuspRR:          f32(80),
		TireSlipFL:          f32(84),
		TireSlipFR:          f32(88),
		TireSlipRL:          f32(92),
		TireSlipRR:          f32(96),
		WheelRotFL:          f32(100),
		WheelRotFR:          f32(104),
		WheelRotRL:          f32(108),
		WheelRotRR:          f32(112),
		WheelOnRumbleFL:     f32(116),
		WheelOnRumbleFR:     f32(120),
		WheelOnRumbleRL:     f32(124),
		WheelOnRumbleRR:     f32(128),
		WheelInPuddleFL:     f32(132),
		WheelInPuddleFR:     f32(136),
		WheelInPuddleRL:     f32(140),
		WheelInPuddleRR:     f32(144),
		SurfaceRumbleFL:     f32(148),
		SurfaceRumbleFR:     f32(152),
		SurfaceRumbleRL:     f32(156),
		SurfaceRumbleRR:     f32(160),
		TireSlipAngleFL:     f32(164),
		TireSlipAngleFR:     f32(168),
		TireSlipAngleRL:     f32(172),
		TireSlipAngleRR:     f32(176),
		TireCombinedSlipFL:  f32(180),
		TireCombinedSlipFR:  f32(184),
		TireCombinedSlipRL:  f32(188),
		TireCombinedSlipRR:  f32(192),
		SuspTravelFL:        f32(196),
		SuspTravelFR:        f32(200),
		SuspTravelRL:        f32(204),
		SuspTravelRR:        f32(208),
		CarOrdinal:          i32(212),
		CarClass:            i32(216),
		CarPerformanceIndex: i32(220),
		DrivetrainType:      i32(224),
		NumCylinders:        i32(228),
		PosX:                f32(232),
		PosY:                f32(236),
		PosZ:                f32(240),
		SpeedMPS:            speedMps,
		SpeedKPH:            speedMps * 3.6,
		SpeedMPH:            speedMps * 2.23694,
		Power:               f32(248),
		Torque:              f32(252),
		TireTempFL:          f32(256),
		TireTempFR:          f32(260),
		TireTempRL:          f32(264),
		TireTempRR:          f32(268),
		Boost:               f32(272),
		Fuel:                f32(276),
		Distance:            f32(280),
		BestLap:             f32(284),
		LastLap:             f32(288),
		CurrentLap:          f32(292),
		CurrentRaceTime:     f32(296),
		LapNumber:           int(u16(300)),
		RacePosition:        int(u8(302)),
		Accel:               int(u8(303)),
		Brake:               int(u8(304)),
		Clutch:              int(u8(305)),
		Handbrake:           int(u8(306)),
		Gear:                int(u8(307)),
		Steer:               int(i8(308)),
		NormDrivingLine:     int(i8(309)),
		NormAIBrakeDiff:     int(i8(310)),
	}, nil
}
