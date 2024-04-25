package weather

import (
	"fmt"

	"github.com/johnstarich/go/gowerline/internal/icon"
)

type weatherEnum int

const (
	weatherUnknown weatherEnum = iota
	weatherBlowingDust
	weatherBlowingSand
	weatherBlowingSnow
	weatherDrizzle
	weatherFog
	weatherFreezingFog
	weatherFreezingDrizzle
	weatherFreezingRain
	weatherFreezingSpray
	weatherFrost
	weatherHail
	weatherHaze
	weatherIceCrystals
	weatherIceFog
	weatherRain
	weatherRainShowers
	weatherSleet
	weatherSmoke
	weatherSnow
	weatherSnowShowers
	weatherThunderstorms
	weatherVolcanicAsh
	weatherWaterSpouts

	maxWeatherValue
)

func weatherValueFromEnum(enum string) weatherEnum {
	for w := weatherUnknown; w < maxWeatherValue; w++ {
		if w.String() == enum {
			return w
		}
	}
	return weatherUnknown
}

func (w weatherEnum) String() string {
	switch w {
	case weatherBlowingDust:
		return "blowing_dust"
	case weatherBlowingSand:
		return "blowing_sand"
	case weatherBlowingSnow:
		return "blowing_snow"
	case weatherDrizzle:
		return "drizzle"
	case weatherFog:
		return "fog"
	case weatherFreezingFog:
		return "freezing_fog"
	case weatherFreezingDrizzle:
		return "freezing_drizzle"
	case weatherFreezingRain:
		return "freezing_rain"
	case weatherFreezingSpray:
		return "freezing_spray"
	case weatherFrost:
		return "frost"
	case weatherHail:
		return "hail"
	case weatherHaze:
		return "haze"
	case weatherIceCrystals:
		return "ice_crystals"
	case weatherIceFog:
		return "ice_fog"
	case weatherRain:
		return "rain"
	case weatherRainShowers:
		return "rain_showers"
	case weatherSleet:
		return "sleet"
	case weatherSmoke:
		return "smoke"
	case weatherSnow:
		return "snow"
	case weatherSnowShowers:
		return "snow_showers"
	case weatherThunderstorms:
		return "thunderstorms"
	case weatherVolcanicAsh:
		return "volcanic_ash"
	case weatherWaterSpouts:
		return "water_spouts"
	case weatherUnknown:
		return ""
	}
	panic(fmt.Sprintf("unexpected weatherValue: %d", w))
}

func (w weatherEnum) Icon() string {
	return icon.StormCloud
}
