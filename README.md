# windslammer_exporter

windslammer_exporter is a [Prometheus](https://prometheus.io/) exporter for
[WindSlammer](http://windslammer.net/) weather station data.

The exporter fetches live weather data from WindSlammer (the wireless weather
station atop Ed Levin County Park) and exposes metrics for wind diretion, speed,
temperature, and elevation (for the main "upper" station as well as proxied data
from a station lower in the valley, for estimating lapse rate).

This exporter was created for an extremely niche purpose.  It is very, very
unlikely that you want this.

## Install

Download from [releases](https://github.com/tris/windslammer_exporter/releases)
or run from Docker:

```bash
docker run -d -p 9307:9307 ghcr.io/tris/windslammer_exporter
```

An alternate port may be defined using the `PORT` environment variable.

## Usage

Start the exporter:

```bash
./windslammer_exporter
```

The exporter will be available at:
- Metrics: http://localhost:9307/metrics
- Health check: http://localhost:9307/health

### Port Configuration

You can override the default port (9307) using the `PORT` environment variable:

```bash
PORT=8080 ./windslammer_exporter
# or
export PORT=8080
./windslammer_exporter
```

## Metrics

The exporter fetches fresh data from `http://windslammer.net/cgi-bin/ws.cgi` on each scrape request and exposes the following metrics:

- `windslammer_wind_direction_degrees` - Current wind direction in degrees
- `windslammer_wind_speed_mph` - Current wind speed in miles per hour
- `windslammer_temperature_lower_fahrenheit` - Temperature from lower elevation station in Fahrenheit
- `windslammer_temperature_upper_fahrenheit` - Temperature from upper elevation station in Fahrenheit
- `windslammer_elevation_lower_feet` - Elevation of lower weather station in feet
- `windslammer_elevation_upper_feet` - Elevation of upper weather station in feet

## Example Prometheus config

```yaml
scrape_configs:
  - job_name: 'windslammer'
    static_configs:
      - targets: ['localhost:9307']
    scrape_interval: 30s
```

## Building

```bash
go build -o windslammer_exporter .
```

## Docker

```bash
docker build -t windslammer_exporter .
docker run -d -p 9307:9307 windslammer_exporter
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
