//go:build linux

package szdiagnostic

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_CheckDatastorePerformance() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	secondsToRun := 1
	result, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 25))
	// Output: {"numRecordsInserted":...
}

func ExampleSzdiagnostic_GetDatastoreInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	result, err := szDiagnostic.GetDatastoreInfo(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"dataStores":[{"id":"CORE", "type":"sqlite3","location":"/tmp/sqlite/G2C.db"}]}
}

func ExampleSzdiagnostic_GetFeature() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	featureID := int64(1)
	result, err := szDiagnostic.GetFeature(ctx, featureID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"LIB_FEAT_ID":1,"FTYPE_CODE":"NAME","ELEMENTS":[{"FELEM_CODE":"TOKENIZED_NM","FELEM_VALUE":"ROBERT|SMITH"},{"FELEM_CODE":"CATEGORY","FELEM_VALUE":"PERSON"},{"FELEM_CODE":"CULTURE","FELEM_VALUE":"ANGLO"},{"FELEM_CODE":"GIVEN_NAME","FELEM_VALUE":"Robert"},{"FELEM_CODE":"SUR_NAME","FELEM_VALUE":"Smith"},{"FELEM_CODE":"FULL_NAME","FELEM_VALUE":"Robert Smith"}]}
}

func ExampleSzdiagnostic_PurgeRepository() {
	// For more information, visit https://github.com/Senzing/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzdiagnostic_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzdiagnostic_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	result := szDiagnostic.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
