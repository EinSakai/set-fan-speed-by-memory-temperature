package main

import (
	"github.com/set-fan-speed-by-memory-temperature/nvapifan"
	"github.com/set-fan-speed-by-memory-temperature/temperature"
	"log"
	"time"
)

// To set *all* coolers (if supported) to 50%:
const NVAPI_COOLER_TARGET_ALL = 0xFFFFFFFF

func main() {
	gpuMemoryTemp := temperature.NewGpuMemoryTemp()
	defer gpuMemoryTemp.GcGpuMemoryTemp()

	var temp uint64
	for {
		gpu, memory, err := gpuMemoryTemp.GetTemp()

		if err != nil {
			log.Println("get gpu core or gpu memory temperature failed! %s", err.Error())
		}

		if gpu > memory {
			temp = gpu
		} else {
			temp = memory
		}

		log.Printf("temperature is %d°C \n", temp)
		if temp <= 30 {
			//set fan speed to 0
			err := nvapifan.SetFanSpeed(NVAPI_COOLER_TARGET_ALL, 0)
			if err != nil {
				log.Println(err)
				panic(err)
			}
		} else {
			//set fan speed to temp
			temp = temp + 10
			if temp > 100 {
				temp = 100
			}
			err := nvapifan.SetFanSpeed(NVAPI_COOLER_TARGET_ALL, uint32(temp))
			if err != nil {
				log.Println(err)
				panic(err)
			}
		}

		time.Sleep(time.Second)
	}
}
