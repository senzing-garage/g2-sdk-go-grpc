//go:build linux

package szconfigmanager

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_AddConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	szConfigManager := getSzConfigManager(ctx)
	configComment := "Example configuration"
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_GetConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configDefinition, err := szConfigManager.GetConfig(ctx, configID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configDefinition, defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleSzconfigmanager_GetConfigs() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configList, err := szConfigManager.GetConfigs(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configList, 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleSzconfigmanager_GetDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	szConfigManager := getSzConfigManager(ctx)
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configComment := "Example configuration"
	newDefaultConfigID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		fmt.Println(err)
	}
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzconfigmanager_SetDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configID, err := szConfigManager.GetDefaultConfigID(ctx) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		fmt.Println(err)
	}
	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	err := szConfigManager.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzconfigmanager_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfigmanager_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmananger_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	result := szConfigManager.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_Initialize() {
	// TODO: Write ExampleSzconfigmanager_Initialize
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	// ctx := context.TODO()
	// grpcConnection := getGrpcConnection()
	// szConfigManager := &SzConfigManager{
	// 	GrpcClient: szpb.NewG2ConfigMgrClient(grpcConnection),
	// }
	// moduleName := "Test module name"
	// iniParams := "{}"
	// verboseLogging := int64(0)
	// err := szConfigManager.Init(ctx, moduleName, iniParams, verboseLogging)
	// if err != nil {
	// 	// This should produce a "senzing-60124002" error.
	// }
	// // Output:
}

func ExampleSzconfigmanager_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	err := szConfigManager.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
