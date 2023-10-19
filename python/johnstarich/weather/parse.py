#!/usr/bin/env python3

from datetime import datetime, timezone
import os
import subprocess

from johnstarich.selectors import uname
import dns.resolver
import maxminddb
import requests


def location() -> (float, float):
    kernel = uname.get_kernel()
    if kernel == 'Darwin':
        try:
            process = subprocess.run(
                ['CoreLocationCLI', '-once', 'YES',
                 '-format', '%latitude\n%longitude'],
                stdout=subprocess.PIPE,
                encoding='utf8',
                timeout=2,
            )
        except FileNotFoundError:
            print("CoreLocationCLI not installed")
            return (None, None)
        lines = str(process.stdout).splitlines()
        if len(lines) < 2:
            return (None, None)
        latitude = lines[0]
        longitude = lines[1]
        return (latitude, longitude)
    elif kernel == 'Linux':
        try:
            answer = dns.resolver.resolve_at("resolver1.opendns.com", "myip.opendns.com", "A")
            ip = str(answer[0].to_text())
            return ip2location(ip)
        except Exception as e:
            return (None, None)
    else:
        raise Exception("Location services not currently supported on this platform: {kernel}".format(kernel=kernel))


def ip2location(ip: str) -> (float, float):
    """IP Geolocation by DB-IP https://db-ip.com"""
    db = os.path.expanduser('~/.local/lib/johnstarich-powerline/dbip-city.mmdb')
    with maxminddb.open_database(db) as reader:
        location = reader.get(ip)['location']
        return location['latitude'], location['longitude']


def parse_iso_8601(date: str) -> datetime:
    date = date.rsplit(sep='/', maxsplit=1)[0]
    date = datetime.fromisoformat(date)
    return date

def recent_value(d: dict, now: datetime):
    values = map(
        lambda v: {
            'value': v['value'],
            'validTime': parse_iso_8601(v['validTime']),
        },
        d['values']
    )
    current = None
    for entry in values:
        if current is None:
            current = entry
        if entry['validTime'] < now and entry['validTime'] > current['validTime']:
            current = entry
    return current['value']


def weather(latitude: float, longitude: float) -> dict:
    response = requests.get(
        'https://api.weather.gov/points/{lat},{lon}'
        .format(lat=latitude, lon=longitude),
        timeout=5,
    )
    forecast_url = response.json()['properties']['forecastGridData']
    response = requests.get(forecast_url, timeout=5)
    json = response.json()
    actual_lat, actual_lon = json['geometry']['coordinates'][0][0]
    props = json['properties']

    now = datetime.now(tz=timezone.utc)
    return {
        'code': 0,
        'coordinates': {
            'latitude': actual_lat,
            'longitude': actual_lon,
        },
        'humidity': recent_value(props['relativeHumidity'], now) / 100.0,
        'apparentTemperature': recent_value(props['apparentTemperature'], now),
        'temperature': recent_value(props['temperature'], now),
        'wind': {
            'chill': recent_value(props['windChill'], now),
            'direction': recent_value(props['windDirection'], now),
            'speed': recent_value(props['windSpeed'], now),
        },
        'units': {
            'temperature': 'C',
            'apparentTemperature': 'C',
        },
    }


last_location = (None, None)


def raw_weather_info():
    global last_location
    latitude, longitude = location()
    error = None
    if latitude is None or longitude is None:
        error = 'Error getting current location.'
        if last_location[0] is None or last_location[1] is None:
            return {'error': error}
        latitude = last_location[0]
        longitude = last_location[1]
    else:
        last_location = (latitude, longitude)
    try:
        info = weather(latitude, longitude)
        if error is not None:
            info['error'] = error
        return info
    except Exception as e:
        return {'error': 'Error parsing weather information: ' + str(e)}


if __name__ == '__main__':
    print(raw_weather_info())
