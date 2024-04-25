package weather

import (
	"fmt"

	"github.com/johnstarich/go/gowerline/internal/icon"
)

type State int

const (
	stateUnknown State = iota
	stateBlowingDust
	stateBlowingSand
	stateBlowingSnow
	stateDrizzle
	stateFog
	stateFreezingFog
	stateFreezingDrizzle
	stateFreezingRain
	stateFreezingSpray
	stateFrost
	stateHail
	stateHaze
	stateIceCrystals
	stateIceFog
	stateRain
	stateRainShowers
	stateSleet
	stateSmoke
	stateSnow
	stateSnowShowers
	stateThunderstorms
	stateVolcanicAsh
	stateWaterSpouts

	maxStateValue
)

func stateFromEnum(enum string) State {
	for w := stateUnknown; w < maxStateValue; w++ {
		if w.String() == enum {
			return w
		}
	}
	return stateUnknown
}

func (w State) String() string {
	switch w {
	case stateBlowingDust:
		return "blowing_dust"
	case stateBlowingSand:
		return "blowing_sand"
	case stateBlowingSnow:
		return "blowing_snow"
	case stateDrizzle:
		return "drizzle"
	case stateFog:
		return "fog"
	case stateFreezingFog:
		return "freezing_fog"
	case stateFreezingDrizzle:
		return "freezing_drizzle"
	case stateFreezingRain:
		return "freezing_rain"
	case stateFreezingSpray:
		return "freezing_spray"
	case stateFrost:
		return "frost"
	case stateHail:
		return "hail"
	case stateHaze:
		return "haze"
	case stateIceCrystals:
		return "ice_crystals"
	case stateIceFog:
		return "ice_fog"
	case stateRain:
		return "rain"
	case stateRainShowers:
		return "rain_showers"
	case stateSleet:
		return "sleet"
	case stateSmoke:
		return "smoke"
	case stateSnow:
		return "snow"
	case stateSnowShowers:
		return "snow_showers"
	case stateThunderstorms:
		return "thunderstorms"
	case stateVolcanicAsh:
		return "volcanic_ash"
	case stateWaterSpouts:
		return "water_spouts"
	case stateUnknown:
		return ""
	}
	panic(fmt.Sprintf("unexpected state: %d", w))
}

func (w State) Icon() string {
	switch w {
	case
		stateBlowingDust,
		stateBlowingSand,
		stateFog,
		stateHaze,
		stateSmoke:
		return icon.DustCloud
	case
		stateFreezingFog,
		stateIceFog,
		stateFrost,
		stateBlowingSnow,
		stateIceCrystals,
		stateHail,
		stateSleet,
		stateSnow,
		stateSnowShowers:
		return icon.SnowCloud
	case
		stateDrizzle,
		stateRain,
		stateRainShowers,
		stateFreezingDrizzle,
		stateFreezingRain,
		stateFreezingSpray:
		return icon.RainCloud
	case stateThunderstorms:
		return icon.StormCloud
	case stateVolcanicAsh, stateWaterSpouts:
		return icon.Critical
	}
	return icon.Warning
}
