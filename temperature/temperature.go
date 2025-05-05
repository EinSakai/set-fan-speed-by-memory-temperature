package temperature

/*
#cgo LDFLAGS: -lnvapi64
#include <nvapi.h>

// Initialize NVAPI
static int initNVAPI() {
    return (NvAPI_Initialize() == NVAPI_OK);
}

// Return the first GPU handle, or NULL on failure
static NvPhysicalGpuHandle getFirstGpuHandle() {
    NvU32 count = 0;
    static NvPhysicalGpuHandle handles[NVAPI_MAX_PHYSICAL_GPUS];
    if (NvAPI_EnumPhysicalGPUs(handles, &count) != NVAPI_OK || count == 0) {
        return NULL;
    }
    return handles[0];
}

// Populate *temp with memory (junction) temperature; return non-zero on success
static int getMemTemp(NvPhysicalGpuHandle gpu, unsigned int *temp) {
    NV_GPU_THERMAL_SETTINGS ts = {0};
    ts.version = NV_GPU_THERMAL_SETTINGS_VER_1;
    if (NvAPI_GPU_GetThermalSettings(gpu, NVAPI_THERMAL_TARGET_MEMORY, &ts) != NVAPI_OK) {
        return 0;
    }
    if (ts.count == 0) {
        return 0;
    }
    *temp = ts.sensor[0].currentTemp;
    return 1;
}
*/
import "C"
import (
	"errors"
	nvml "github.com/NVIDIA/go-nvml/pkg/nvml"
	"log"
)

type GpuMemoryTemp struct {
	Dev    nvml.Device
	Handle C.NvPhysicalGpuHandle
}

func NewGpuMemoryTemp() *GpuMemoryTemp {
	gpuMemoryTemp := GpuMemoryTemp{}

	//for gpu:
	// NVML init
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Printf("NVML init failed: %v\n", ret.Error())
		panic("NVML init failed!")
	}

	// Get NVML GPU handle
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Printf("NVML DeviceGetCount() failed: %v\n", ret.Error())
		panic("NVML DeviceGetCount failed!")
	}
	if count == 0 {
		panic("No NVIDIA GPUs (NVML)")
	}
	//just get first gpu:
	dev, ret := nvml.DeviceGetHandleByIndex(0)
	if ret != nvml.SUCCESS {
		log.Printf("NVML get GPU handle failed: %v\n", ret.Error())
		panic("NVML get GPU handle failed")
	}
	gpuMemoryTemp.Dev = dev

	//for gpu memory:
	// 1. 初始化 NVAPI
	if C.initNVAPI() == 0 {
		panic("NVAPI 初始化失败") // NvAPI_Initialize must return NVAPI_OK :contentReference[oaicite:4]{index=4}
	}

	// 2. 获取第一个 GPU handle
	gpuHandle := C.getFirstGpuHandle()
	if gpuHandle == nil {
		panic("未检测到任何 NVIDIA GPU") // NvAPI_EnumPhysicalGPUs 返回 non-zero count :contentReference[oaicite:5]{index=5}
	}

	gpuMemoryTemp.Handle = gpuHandle

	return &gpuMemoryTemp
}

func (g *GpuMemoryTemp) GetTemp() (uint64, uint64, error) {
	//读取gpu温度
	coreTemp, ret := g.Dev.GetTemperature(nvml.TEMPERATURE_GPU)
	if ret != nvml.SUCCESS {
		log.Printf("Core Temp error: %v\n", ret.Error())
		return 0, 0, errors.New("get gpu core temperature failed")
	}
	log.Printf("gpu core temp: %d°C \n", coreTemp)

	//读取显存温度
	var memTemp C.uint
	if C.getMemTemp(g.Handle, &memTemp) == 0 {
		log.Println("get gpu memory temperature failed")
		return 0, 0, errors.New("get gpu memory temperature failed")
	}
	log.Printf("gpu memory temp: %d°C \n", memTemp) // 读取成功并返回有效温度 :contentReference[oaicite:6]{index=6}

	return uint64(coreTemp), uint64(memTemp), nil
}

func (g *GpuMemoryTemp) GcGpuMemoryTemp() {
	//卸载 NVAPI
	C.NvAPI_Unload()
	nvml.Shutdown()
}
