from johnstarich.weather.parse import raw_weather_info
from johnstarich.interval import Interval
from johnstarich.segment import segment, segment_default


weather_status_icons = {
    'clear': '☀️',
    'clear night': '🌙',
    'few clouds': '⛅️',
    'few clouds night': '☁️',
    'clouds': '☁️',
    'rain': '🌧',
    'shower rain': '🌦',
    'shower rain night': '🌧',
    'thunderstorm': '⛈',
    'snow': '🌨',
    'mist': '💨',
    'disaster': '🌪',
    'invalid': '🚫',
}


# Yahoo! Weather codes:
# https://developer.yahoo.com/weather/documentation.html#codes
weather_status_mappings = {
    0: ['tornado', 'disaster'],
    1: ['tropical storm', 'disaster'],
    2: ['hurricane', 'disaster'],
    3: ['severe thunderstorms', 'thunderstorm'],
    4: ['thunderstorms', 'thunderstorm'],
    5: ['mixed rain and snow', 'snow'],
    6: ['mixed rain and sleet', 'snow'],
    7: ['mixed snow and sleet', 'snow'],
    8: ['freezing drizzle', 'rain'],
    9: ['drizzle', 'shower rain'],
    10: ['freezing rain', 'rain'],
    11: ['showers', 'shower rain'],
    12: ['showers', 'shower rain'],
    13: ['snow flurries', 'snow'],
    14: ['light snow showers', 'snow'],
    15: ['blowing snow', 'snow'],
    16: ['snow', 'snow'],
    17: ['hail', 'snow'],
    18: ['sleet', 'snow'],
    19: ['dust', 'mist'],
    20: ['foggy', 'mist'],
    21: ['haze', 'mist'],
    22: ['smoky', 'mist'],
    23: ['blustery', 'mist'],
    24: ['windy', 'mist'],
    25: ['cold', 'clear'],
    26: ['cloudy', 'clouds'],
    27: ['mostly cloudy (night)', 'clouds'],
    28: ['mostly cloudy (day)', 'clouds'],
    29: ['partly cloudy (night)', 'few clouds'],
    30: ['partly cloudy (day)', 'few clouds'],
    31: ['clear (night)', 'clear night'],
    32: ['sunny', 'clear'],
    33: ['fair (night)', 'clear night'],
    34: ['fair (day)', 'clear'],
    35: ['mixed rain and hail', 'snow'],
    36: ['hot', 'clear'],
    37: ['isolated thunderstorms', 'thunderstorm'],
    38: ['scattered thunderstorms', 'thunderstorm'],
    39: ['scattered thunderstorms', 'thunderstorm'],
    40: ['scattered showers', 'shower rain'],
    41: ['heavy snow', 'snow'],
    42: ['scattered snow showers', 'snow'],
    43: ['heavy snow', 'snow'],
    44: ['partly cloudy', 'few clouds'],
    45: ['thundershowers', 'thunderstorm'],
    46: ['snow showers', 'snow'],
    47: ['isolated thundershowers', 'thunderstorm'],
    3200: ['not available', 'invalid'],
}


segment_kwargs = {
    'highlight_groups': [
        'weather_temp_gradient',
        'weather_temp',
        'weather'
    ],
}


update_interval = Interval(5 * 60)
last_status = ''
last_gradient = 0


def weather(pl, unit: str='C', temp_low: float=0, temp_high: float=100,
            **kwargs) -> list:
    global last_status, last_gradient
    if not update_interval.should_run():
        return segment(last_status, gradient_level=last_gradient,
                       **segment_kwargs)

    if temp_low >= temp_high:
        raise ValueError('temp_low cannot be higher then or '
                         'the same as temp_high')
    weather = raw_weather_info()
    warning_str = None
    if 'error' in weather and 'location' in weather['error']:
        warning_str = ' ⚠️ 🌎 '
        if len(weather.keys()) == 1:
            update_interval.start()
            if not last_status.endswith(warning_str):
                last_status += warning_str
            return segment(last_status, gradient_level=last_gradient,
                           **segment_kwargs)
    elif 'error' in weather:
        update_interval.start()
        warning_str = ' ⚠️ '
        if '⚠️' not in last_status:
            last_status += warning_str
        return segment(last_status, gradient_level=last_gradient,
                       **segment_kwargs)
    print(weather)
    temperature = weather['temperature']
    input_unit = weather['units']['temperature']
    humidity = weather['humidity']
    additional_variance = 0

    temp_in_fahrenheit = convert_temperature(temperature, input_unit, 'F')
    if temp_in_fahrenheit >= 80 and humidity >= 40:
        display_temperature = heat_index(temperature, humidity,
                                         input_unit, unit)
        if display_temperature != temp_in_fahrenheit:
            additional_variance = convert_temperature(1.3, 'F', unit) - \
                convert_temperature(0, 'F', unit)
    elif temp_in_fahrenheit <= 50:
        display_temperature = convert_temperature(weather['wind']['chill'],
                                                  input_unit, unit)
    else:
        display_temperature = convert_temperature(temperature,
                                                  input_unit, unit)

    gradient = 100 * (display_temperature - temp_low) / (temp_high - temp_low)
    if display_temperature > temp_high:
        gradient = 100
    elif display_temperature < temp_low:
        gradient = 0

    variance = ''
    if additional_variance != 0:
        display_temperature = round(display_temperature)
        variance = '±' + str(round(abs(additional_variance), 1))
    else:
        display_temperature = round(display_temperature, 1)

    contents = '{icon}  {temp}{var}°{unit}{warning}'.format(
        icon=extract_icon(weather),
        temp=display_temperature,
        unit=unit,
        var=variance,
        warning=warning_str if warning_str is not None else '',
    )
    update_interval.start()
    last_status = contents
    last_gradient = gradient
    return segment(contents, gradient_level=gradient, **segment_kwargs)


def extract_icon(weather: dict) -> str:
    weather_code = weather['code']
    if weather_code not in weather_status_mappings:
        return weather_status_icons['invalid']
    return weather_status_icons[weather_status_mappings[weather_code][1]]


heat_index_constants = [
    -42.379, 2.04901523, 10.14333127, -0.22475541, -6.83783e-3,
    -5.481717e-2, 1.22874e-3, 8.5282e-4, -1.99e-6
]


def heat_index(temperature: float, relative_humidity: float,
               input_unit: str, output_unit: str) -> float:
    if input_unit not in 'KCF':
        raise ValueError('Invalid input unit: ' + input_unit)
    if output_unit not in 'KCF':
        raise ValueError('Invalid output unit: ' + output_unit)
    global heat_index_constants
    HI_C = heat_index_constants
    T = temperature_conversions[input_unit + 'F'](temperature)
    R = relative_humidity
    HI = HI_C[0] + HI_C[1]*T + HI_C[2]*R + HI_C[3]*T*R + HI_C[4]*T*T + \
        HI_C[5]*R*R + HI_C[6]*T*T*R + HI_C[7]*T*R*R + HI_C[8]*T*T*R*R
    return temperature_conversions['F' + output_unit](HI)


temperature_conversions = {
    'CC': lambda t: t,
    'CF': lambda t: t * 9/5 + 32,
    'CK': lambda t: t + 273.15,
    'FC': lambda t: (t - 32) * 5/9,
    'FF': lambda t: t,
    'FK': lambda t: (t + 459.67) * 5/9,
    'KC': lambda t: t - 273.15,
    'KF': lambda t: t * 9/5 - 459.67,
    'KK': lambda t: t,
}


def convert_temperature(temperature: float, input_unit: str, output_unit: str):
    if input_unit not in 'KCF':
        raise ValueError('Input unit is not valid: ' + input_unit)
    if output_unit not in 'KCF':
        raise ValueError('Output unit is not valid: ' + output_unit)
    return temperature_conversions[input_unit + output_unit](temperature)
