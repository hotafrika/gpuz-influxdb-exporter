package main

import (
	"fmt"
	"github.com/hotafrika/gpuz-reader"
	"github.com/pkg/errors"
	"math"
)

var RefSensors = map[string][]string{
	"gpu_clock":       []string{"GPU Clock"},
	"memory_clock":    []string{"Memory Clock"},
	"gpu_temperature": []string{"GPU Temperature"},
	"hot_spot":        []string{"Hot Spot"},
	"gpu_power":       []string{"Board Power Draw", "GPU Power"},
	"gpu_load":        []string{"GPU Load"},
	"gpu_voltage":     []string{"GPU Voltage"},
	"cpu_temperature": []string{"CPU Temperature"},
	"memory_used":     []string{"Memory Used", "Memory Used (Dedicated)"},
}

var RefRecords = map[string]string{
	"card_name":    "CardName",
	"vendor_id":    "VendorID",
	"device_id":    "DeviceID",
	"subvendor_id": "SubvendorID",
	"sybsys_id":    "SubsysID",
}

type Collector struct {
	metrics []string
	sm      *gpuz.SharedMemory
}

func NewCollector(metrics ...string) Collector {
	var necMetrics []string
	if len(metrics) != 0 {
		for _, m := range metrics {
			_, ok := RefSensors[m]
			if ok {
				necMetrics = append(necMetrics, m)
			}
		}
	}
	if len(necMetrics) == 0 {
		for k := range RefSensors {
			necMetrics = append(necMetrics, k)
		}
	}
	return Collector{
		metrics: necMetrics,
		sm:      gpuz.DefaultSharedMemory(),
	}
}

func (c Collector) GetInfluxRow(hostname string) (map[string]string, map[string]interface{}, error) {
	tags := make(map[string]string)
	values := make(map[string]interface{})
	stat, err := c.sm.GetStat()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get data: ")
	}
	for _, m := range c.metrics {
		if len(RefSensors[m]) < 1 {
			continue
		}
		value, ok := stat.GetSensorValue(RefSensors[m][0])
		if !ok {
			if len(RefSensors[m]) > 1 {
				value, ok = stat.GetSensorValue(RefSensors[m][1])
				if !ok {
					continue
				}
			}
			continue
		}
		if math.IsNaN(value) {
			value = 0
		}
		values[m] = roundFloat(value)
	}

	if len(values) == 0 {
		return nil, nil, fmt.Errorf("no data found for these types of metrics")
	}

	name, ok := stat.GetRecord(RefRecords["card_name"])
	if !ok {
		name = "undefined"
	}
	vendorID, ok := stat.GetRecord(RefRecords["vendor_id"])
	if !ok {
		vendorID = "undefined"
	}
	deviceID, ok := stat.GetRecord(RefRecords["device_id"])
	if !ok {
		deviceID = "undefined"
	}
	subvendorID, ok := stat.GetRecord(RefRecords["subvendor_id"])
	if !ok {
		subvendorID = "undefined"
	}
	subsysID, ok := stat.GetRecord(RefRecords["sybsys_id"])
	if !ok {
		subsysID = "undefined"
	}

	hostname = simplifyString(hostname)
	if hostname == "" {
		hostname = "undefined"
	}
	tags["host"] = hostname
	tags["card_name"] = simplifyString(name)
	tags["device_id"] = simplifyString(vendorID + deviceID + subvendorID + subsysID)

	return tags, values, nil
}
