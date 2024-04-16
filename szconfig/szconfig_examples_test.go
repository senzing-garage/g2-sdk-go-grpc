//go:build linux

package szconfig

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzConfig_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzConfig_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	result := szConfig.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzConfig_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	dataSourceCode := "GO_TEST"
	result, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DSRC_ID":1001}
}

func ExampleSzConfig_CloseConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfig_CreateConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfig_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	dataSourceCode := "TEST"
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfig_GetDataSources() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	result, err := szConfig.GetDataSources(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleSzConfig_ImportConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	mockConfigHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configDefinition, err := szConfig.ExportConfig(ctx, mockConfigHandle)
	if err != nil {
		fmt.Println(err)
	}
	configHandle, err := szConfig.ImportConfig(ctx, configDefinition)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfig_ExportConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
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
	fmt.Println(truncate(configDefinition, 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"INTERNAL":"No"},{"ATTR_ID":1003,"...
}

func ExampleSzConfig_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	err := szConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfig_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	instanceName := "Test module name"
	settings := "{}"
	verboseLogging := int64(0)
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60114002" error.
	}
	// Output:
}

func ExampleSzConfig_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	err := szConfig.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60114001" error.
	}
	// Output:
}
