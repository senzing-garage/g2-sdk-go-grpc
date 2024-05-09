package szengine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/testfixtures"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szengineapi "github.com/senzing-garage/sz-sdk-go/szengine"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szdiagnosticpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

type GetEntityByRecordIdResponse struct {
	ResolvedEntity struct {
		EntityId int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

var (
	szConfigSingleton        sz.SzConfig
	szConfigManagerSingleton sz.SzConfigManager
	szDiagnosticSingleton    sz.SzDiagnostic
	szEngineSingleton        *Szengine
	grpcAddress              = "localhost:8261"
	grpcConnection           *grpc.ClientConn
	logger                   logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzengine_AddRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := sz.SZ_WITHOUT_INFO
	records := []record.Record{
		truthset.CustomerRecords["1004"],
		truthset.CustomerRecords["1005"],
	}
	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
		testError(test, err)
		printActual(test, actual)
		defer szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
	}
}

func TestSzengine_AddRecord_szBadInput(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "ADD_TEST_ERR_1", "NAME_FULL": "NOBODY NOMATCH"}`)
	testError(test, err)
	record2, err := record.NewRecord(`{"DATA_SOURCE": "BOB", "RECORD_ID": "ADD_TEST_ERR_2", "NAME_FULL": "ERR BAD SOURCE"}`)
	testError(test, err)
	flags := sz.SZ_WITHOUT_INFO

	// This one should succeed.

	actual, err := szEngine.AddRecord(ctx, record1.DataSource, record1.Id, record1.Json, flags)
	testError(test, err)
	defer szEngine.DeleteRecord(ctx, record1.DataSource, record1.Id, sz.SZ_WITHOUT_INFO)
	printActual(test, actual)

	// This one should fail.

	actual, err = szEngine.AddRecord(ctx, "CUSTOMERS", record2.Id, record2.Json, flags)
	assert.True(test, szerror.Is(err, szerror.SzBadInput))
	printActual(test, actual)

	// Clean-up the records we inserted.

	actual, err = szEngine.DeleteRecord(ctx, record1.DataSource, record1.Id, flags)
	testError(test, err)
	defer szEngine.DeleteRecord(ctx, record2.DataSource, record2.Id, sz.SZ_NO_FLAGS)
	printActual(test, actual)
}

func TestSzengine_AddRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := sz.SZ_WITH_INFO
	records := []record.Record{
		truthset.CustomerRecords["1004"],
		truthset.CustomerRecords["1005"],
	}
	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
		testError(test, err)
		printActual(test, actual)
		defer szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
	}
}

func TestSzengine_CloseExport(test *testing.T) {
	// Tested in:
	//  - TestSzengine_ExportCsvEntityReport
	//  - TestSzengine_ExportJsonEntityReport
}

func TestSzengine_CountRedoRecords(test *testing.T) {
	expected := int64(1)
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.CountRedoRecords(ctx)
	testError(test, err)
	printActual(test, actual)
	assert.Equal(test, expected, actual)
}

func TestSzengine_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1009"]
	flags := sz.SZ_WITHOUT_INFO
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
	printActual(test, actual)
	testError(test, err)
	actual, err = szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1010"]
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
	testError(test, err)
	printActual(test, actual)
	actual, err = szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_ExportCsvEntityReport(test *testing.T) {
	expected := []string{
		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
		`1,0,"","","CUSTOMERS","1001"`,
		`1,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
		`1,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
	}
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	csvColumnList := ""
	flags := sz.SZ_EXPORT_INCLUDE_ALL_ENTITIES
	aHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, aHandle)
		testError(test, err)
	}()
	testError(test, err)
	actualCount := 0
	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))
		actualCount += 1
	}
	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportCsvEntityReportIterator(test *testing.T) {
	expected := []string{
		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
		`1,0,"","","CUSTOMERS","1001"`,
		`1,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
		`1,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
	}
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	csvColumnList := ""
	flags := sz.SZ_EXPORT_INCLUDE_ALL_ENTITIES
	actualCount := 0
	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		testError(test, actual.Error)
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))
		actualCount += 1
	}
	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportJsonEntityReport(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.AddRecord(ctx, aRecord.DataSource, aRecord.Id, aRecord.Json, flags)
	testError(test, err)
	printActual(test, actual)
	defer szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.Id, sz.SZ_WITHOUT_INFO)
	// TODO: Figure out correct flags.
	// flags := sz.Flags(sz.SZ_EXPORT_DEFAULT_FLAGS, sz.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS, sz.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS)
	flags = int64(-1)
	aHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, aHandle)
		testError(test, err)
	}()
	testError(test, err)
	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, aHandle)
		testError(test, err)
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}
	testError(test, err)
	assert.True(test, len(jsonEntityReport) > 65536)
}

func TestSzengine_ExportJsonEntityReportIterator(test *testing.T) {
	expected := 1
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := sz.SZ_EXPORT_INCLUDE_ALL_ENTITIES
	actualCount := 0
	for actual := range szEngine.ExportJsonEntityReportIterator(ctx, flags) {
		testError(test, actual.Error)
		printActual(test, actual.Value)
		actualCount += 1
	}
	assert.Equal(test, expected, actualCount)
}

func TestSzengine_FetchNext(test *testing.T) {
	// Tested in:
	//  - TestSzengine_ExportJsonEntityReport
}

func TestSzengine_FindInterestingEntitiesByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	flags := int64(0)
	actual, err := szEngine.FindInterestingEntitiesByEntityId(ctx, entityID, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	actual, err := szEngine.FindInterestingEntitiesByRecordId(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := sz.SZ_FIND_NETWORK_DEFAULT_FLAGS
	actual, err := szEngine.FindNetworkByEntityId(ctx, entityList, maxDegrees, buildOutDegree, maxEntities, flags)
	testErrorNoFail(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.Id + `"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := sz.SZ_FIND_NETWORK_DEFAULT_FLAGS
	actual, err := szEngine.FindNetworkByRecordId(ctx, recordList, maxDegrees, buildOutDegree, maxEntities, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startEntityId := getEntityId(truthset.CustomerRecords["1001"])
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := sz.SZ_NO_EXCLUSIONS
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityId_excluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(startRecord)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(startRecord) + `}]}`
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityId_excludingAndIncluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(startRecord)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(startRecord) + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityId_including(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(startRecord)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := sz.SZ_NO_EXCLUSIONS
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := sz.SZ_NO_EXCLUSIONS
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordId_excluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordId_excludingAndIncluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordId_including(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetActiveConfigId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetActiveConfigId(ctx)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetEntityByEntityId(ctx, entityId, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetEntityByRecordId(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetRedoRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetRedoRecord(ctx)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetStats(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetStats(ctx)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetVirtualEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetVirtualEntityByRecordId(ctx, recordList, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_HowEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.HowEntityByEntityId(ctx, entityId, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	err := szEngine.PrimeEngine(ctx)
	testError(test, err)
}

func TestSzengine_ProcessRedoRecord(test *testing.T) {
	// TODO: Implement TestSzengine_ProcessRedoRecord
	// ctx := context.TODO()
	// szEngine := getTestObject(ctx, test)
	// flags := sz.SZ_WITHOUT_INFO
	// actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	// testError(test, err)
	// printActual(test, actual)
}

func TestSzengine_ProcessRedoRecord_withInfo(test *testing.T) {
	// TODO: Implement TestSzengine_ProcessRedoRecord_withInfo
	// ctx := context.TODO()
	// szEngine := getTestObject(ctx, test)
	// flags := sz.SZ_WITH_INFO
	// actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	// testError(test, err)
	// printActual(test, actual)
}

func TestSzengine_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_WITHOUT_INFO
	actual, err := szEngine.ReevaluateEntity(ctx, entityId, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.ReevaluateEntity(ctx, entityId, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_WITHOUT_INFO
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := sz.SZ_NO_SEARCH_PROFILE
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_StreamExportCsvEntityReport(test *testing.T) {
	// TODO: Write TestSzengine_StreamExportCsvEntityReport
}

func TestSzengine_StreamExportJsonEntityReport(test *testing.T) {
	// TODO: Write TestSzengine_StreamExportJsonEntityReport
}

func TestSzengine_SearchByAttributes_searchProfile(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := sz.SZ_NO_SEARCH_PROFILE // TODO: Figure out the search profile
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.WhyEntities(ctx, entityId1, entityId2, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.WhyRecordInEntity(ctx, record.DataSource, record.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.WhyRecords(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, flags)
	testError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzengine_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
}

func TestSzengine_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	actual := szEngine.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzengine_AsInterface(test *testing.T) {
	expected := int64(0)
	ctx := context.TODO()
	szEngine := getSzEngineAsInterface(ctx)
	actual, err := szEngine.CountRedoRecords(ctx)
	testError(test, err)
	printActual(test, actual)
	assert.Equal(test, expected, actual)
}

func TestSzengine_Initialize(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	instanceName := "Test name"
	settings := "{}"
	verboseLogging := sz.SZ_NO_LOGGING
	configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
	err := szEngine.Initialize(ctx, instanceName, settings, configId, verboseLogging)
	testError(test, err)
}

func TestSzengine_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	configId, err := szEngine.GetActiveConfigId(ctx)
	testError(test, err)
	err = szEngine.Reinitialize(ctx, configId)
	testError(test, err)
	printActual(test, configId)
}

func TestSzengine_Destroy(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	err := szEngine.Destroy(ctx)
	testError(test, err)
	szEngineSingleton = nil
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
}

func getEntityId(record record.Record) int64 {
	return getEntityIdForRecord(record.DataSource, record.Id)
}

func getEntityIdForRecord(datasource string, id string) int64 {
	ctx := context.TODO()
	var result int64 = 0
	szEngine := getSzEngine(ctx)
	response, err := szEngine.GetEntityByRecordId(ctx, datasource, id, sz.SZ_WITHOUT_INFO)
	if err != nil {
		return result
	}
	getEntityByRecordIdResponse := &GetEntityByRecordIdResponse{}
	err = json.Unmarshal([]byte(response), &getEntityByRecordIdResponse)
	if err != nil {
		return result
	}
	return getEntityByRecordIdResponse.ResolvedEntity.EntityId
}

func getEntityIdString(record record.Record) string {
	entityId := getEntityId(record)
	return strconv.FormatInt(entityId, 10)
}

func getEntityIdStringForRecord(datasource string, id string) string {
	entityId := getEntityIdForRecord(datasource, id)
	return strconv.FormatInt(entityId, 10)
}

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getSzConfig(ctx context.Context) sz.SzConfig {
	_ = ctx
	if szConfigSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &szconfig.Szconfig{
			GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
		}
	}
	return szConfigSingleton
}

func getSzConfigManager(ctx context.Context) sz.SzConfigManager {
	_ = ctx
	if szConfigManagerSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{
			GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
		}
	}
	return szConfigManagerSingleton
}

func getSzDiagnostic(ctx context.Context) sz.SzDiagnostic {
	_ = ctx
	if szDiagnosticSingleton == nil {
		grpcConnection := getGrpcConnection()
		szDiagnosticSingleton = &szdiagnostic.Szdiagnostic{
			GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(grpcConnection),
		}
	}
	return szDiagnosticSingleton
}

func getSzEngine(ctx context.Context) *Szengine {
	_ = ctx
	if szEngineSingleton == nil {
		grpcConnection := getGrpcConnection()
		szEngineSingleton = &Szengine{
			GrpcClient: szpb.NewSzEngineClient(grpcConnection),
		}
	}
	return szEngineSingleton
}

func getSzEngineAsInterface(ctx context.Context) sz.SzEngine {
	return getSzEngine(ctx)
}

func getTestObject(ctx context.Context, test *testing.T) *Szengine {
	_ = test
	return getSzEngine(ctx)
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func testError(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if szerror.Is(err, szerror.SzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzBadInput) {
			fmt.Printf("\nBad user input error detected. \n\n")
		}
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	options := []interface{}{
		&logging.OptionCallerSkip{Value: 4},
	}
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szengineapi.IdMessages, options...)
	if err != nil {
		return createError(5901, err)
	}
	err = setupSenzingConfiguration(ctx)
	if err != nil {
		return createError(5920, err)
	}
	err = setupPurgeRepository(ctx)
	if err != nil {
		return createError(5921, err)
	}
	err = setupAddRecords()
	if err != nil {
		return createError(5922, err)
	}
	return err
}

func setupAddRecords() error {
	var err error = nil
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	flags := sz.SZ_WITHOUT_INFO
	for _, record := range records {
		_, err = szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
		if err != nil {
			return err
		}
	}
	return err
}

func setupSenzingConfiguration(ctx context.Context) error {
	now := time.Now()

	// Create an in memory Senzing configuration.

	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	// Create a string representation of the in-memory configuration.

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	// Close szConfig in-memory object.

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5906, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	szConfigManager := getSzConfigManager(ctx)
	configComment := fmt.Sprintf("Created by szengine_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}

	return err
}

func setupPurgeRepository(ctx context.Context) error {
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.PurgeRepository(ctx)
	return err
}

func teardown() error {
	var err error = nil
	return err
}
