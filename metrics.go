package main

import (
	"strconv"
	"time"

	"github.com/mindprince/gonvml"
	"github.com/pkg/errors"
)

var (
	averageDuration = 10 * time.Second
)

type Metrics struct {
	Version string
	Devices []*Device
}

type Device struct {
	Index                 string
	MinorNumber           string
	Name                  string
	UUID                  string
	Temperature           float64
	PowerUsage            float64
	PowerUsageAverage     float64
	MemoryTotal           float64
	MemoryUsed            float64
	UtilizationMemory     float64
	UtilizationGPU        float64
	UtilizationGPUAverage float64
}

func collectMetrics() (*Metrics, error) {
	if err := gonvml.Initialize(); err != nil {
		return nil, errors.Wrap(err, "Initialize is failed")
	}
	defer gonvml.Shutdown()

	version, err := gonvml.SystemDriverVersion()
	if err != nil {
		return nil, errors.Wrap(err, "SystemDriverVersion is failed")
	}

	metrics := &Metrics{
		Version: version,
	}

	numDevices, err := gonvml.DeviceCount()
	if err != nil {
		return nil, errors.Wrap(err, "DeviceCount is failed")
	}

	for index := 0; index < int(numDevices); index++ {
		device, err := gonvml.DeviceHandleByIndex(uint(index))
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d DeviceHandleByIndex is failed", index)
		}

		uuid, err := device.UUID()
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d UUID is failed", index)
		}

		name, err := device.Name()
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d Name is failed", index)
		}

		minorNumber, err := device.MinorNumber()
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d MinorNumber is failed", index)
		}

		temperature, err := device.Temperature()
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d Temperature is failed", index)
		}

		powerUsage, err := device.PowerUsage()
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d PowerUsage is failed", index)
		}

		powerUsageAverage, err := device.AveragePowerUsage(averageDuration)
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d AveragePowerUsage is failed", index)
		}

		memoryTotal, memoryUsed, err := device.MemoryInfo()
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d MemoryInfo is failed", index)
		}

		utilizationGPU, utilizationMemory, err := device.UtilizationRates()
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d UtilizationRates is failed", index)
		}

		utilizationGPUAverage, err := device.AverageGPUUtilization(averageDuration)
		if err != nil {
			return nil, errors.Wrapf(err, "index:%d AverageGPUUtilization is failed", index)
		}

		metrics.Devices = append(metrics.Devices,
			&Device{
				Index:                 strconv.Itoa(index),
				MinorNumber:           strconv.Itoa(int(minorNumber)),
				Name:                  name,
				UUID:                  uuid,
				Temperature:           float64(temperature),
				PowerUsage:            float64(powerUsage),
				PowerUsageAverage:     float64(powerUsageAverage),
				MemoryTotal:           float64(memoryTotal),
				MemoryUsed:            float64(memoryUsed),
				UtilizationMemory:     float64(utilizationMemory),
				UtilizationGPU:        float64(utilizationGPU),
				UtilizationGPUAverage: float64(utilizationGPUAverage),
			})
	}

	return metrics, nil
}
