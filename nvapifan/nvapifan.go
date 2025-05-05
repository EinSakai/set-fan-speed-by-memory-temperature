package nvapifan

/*
#cgo LDFLAGS: -lnvapi64
#include <windows.h>
#include "nvapi.h"

// 全局函数指针
static NVAPI_STATUS (*NvAPI_Initialize)();
static NVAPI_STATUS (*NvAPI_EnumPhysicalGPUs)(NvPhysicalGpuHandle gpuHandles[], NvU32 *count);
static NVAPI_STATUS (*NvAPI_GPU_GetCoolerSettings)(NvPhysicalGpuHandle hGPU, NvU32 coolerIndex, NV_GPU_COOLER_SETTINGS *pCoolerSettings);
static NVAPI_STATUS (*NvAPI_GPU_SetCoolerSettings)(NvPhysicalGpuHandle hGPU, NvU32 coolerIndex, NV_GPU_COOLER_SETTINGS *pCoolerSettings);

static int loadNvapi() {
    HMODULE h = LoadLibraryA("nvapi64.dll");
    if (!h) return 0;
    NvAPI_Initialize = (void*)GetProcAddress(h, "NvAPI_Initialize");
    NvAPI_EnumPhysicalGPUs = (void*)GetProcAddress(h, "NvAPI_EnumPhysicalGPUs");
    NvAPI_GPU_GetCoolerSettings = (void*)GetProcAddress(h, "NvAPI_GPU_GetCoolerSettings");
    NvAPI_GPU_SetCoolerSettings = (void*)GetProcAddress(h, "NvAPI_GPU_SetCoolerSettings");
    return NvAPI_Initialize && NvAPI_EnumPhysicalGPUs && NvAPI_GPU_GetCoolerSettings && NvAPI_GPU_SetCoolerSettings;
}
*/
import "C"
import (
	"errors"
	"log"
)

func init() {
	if C.loadNvapi() == 0 {
		panic("加载 nvapi64.dll 失败")
	}
	if res := C.NvAPI_Initialize(); res != C.NVAPI_OK {
		panic("NvAPI 初始化失败")
	}
}

// SetFanSpeed 将第 index 个冷却器设置为指定百分比 (0–100)
func SetFanSpeed(index uint32, speedPercent uint32) error {
	log.Printf("set gpu fan speed %d/100... \n\n", speedPercent)
	// 枚举 GPU
	var handles [C.NVAPI_MAX_PHYSICAL_GPUS]C.NvPhysicalGpuHandle
	var count C.NvU32
	if res := C.NvAPI_EnumPhysicalGPUs(&handles[0], &count); res != C.NVAPI_OK {
		return errors.New("枚举 GPU 失败")
	}
	// 只拿第一个
	h := handles[0]

	// 填充 COOKER_SETTINGS 结构
	var cool C.NV_GPU_COOLER_SETTINGS
	cool.version = C.NV_GPU_COOLER_SETTINGS_VER_1
	// coolerIndex 通常 0，speedLevelCount 必须至少 1
	if res := C.NvAPI_GPU_GetCoolerSettings(h, C.NvU32(index), &cool); res != C.NVAPI_OK {
		return errors.New("获取冷却器设置失败")
	}

	// 设置目标转速
	cool.coolers[C.NvU32(index)].currentLevel = C.NvU32(speedPercent)
	if res := C.NvAPI_GPU_SetCoolerSettings(h, C.NvU32(index), &cool); res != C.NVAPI_OK {
		return errors.New("设置冷却器转速失败")
	}

	return nil
}
