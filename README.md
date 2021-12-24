### GPU-Z InfluxDB exporter
Collects GPU metrics from launched GPU-Z utility via shared memory and sends these metrics to InfluxDB endpoint.

Currently, it collects following metrics:
* GPU Clock (MHz)
* Memory Clock (MHz)
* GPU Temperature (Celsius degrees)
* Hot Spot (Celsius degrees)
* GPU Power (Watts)
* GPU Load (%)
* GPU Voltage (Volts)
* CPU Temperature (Celsius degrees)
* Memory Used (MB)

To use this exporter you need to specify some flags:
> *-a* - InfluxDB HTTP endpoint. Default: http://localhost:8086
> 
> *-u* - InfluxDB username. Default: empty
> 
> *-p* - InfluxDB password. Default: empty
> 
> *-n* - namespace for metrics ("measurement" in InfluxDB terminology). Default: gpuz
>
> *-d* - database for metrics. Default: monitoring
> 
> *-i* - interval in seconds between metrics collecting. Default: 60s
> 
> *-h* - hostname of working machine (if you need to change hostname of current machine). Default: OS hostname

###Example
>*gpuz-influxdb-exporter.exe -a http://localhost:8086 -u someuser -p somepassword -n gpuz -d monitoring -i 60*
