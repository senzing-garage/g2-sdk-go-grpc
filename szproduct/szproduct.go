/*
Package szproduct implements a client for the service.
*/
package szproduct

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	szproductapi "github.com/senzing-garage/sz-sdk-go/szproduct"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

type Szproduct struct {
	GrpcClient     szpb.SzProductClient
	isTrace        bool // Performance optimization
	logger         logging.Logging
	observerOrigin string
	observers      subject.Subject
}

// ----------------------------------------------------------------------------
// sz-sdk-go.SzProduct interface methods
// ----------------------------------------------------------------------------

/*
The Destroy method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szproduct) Destroy(ctx context.Context) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(3)
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}
	return err
}

/*
The GetLicense method retrieves information about the currently used license by the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing license metadata.
    See the example output.
*/
func (client *Szproduct) GetLicense(ctx context.Context) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetLicenseRequest{}
	response, err := client.GrpcClient.GetLicense(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}
	return result, err
}

/*
The GetVersion method returns the version of the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing metadata about the Senzing Engine version being used.
    See the example output.
*/
func (client *Szproduct) GetVersion(ctx context.Context) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19)
		defer func() { client.traceExit(20, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetVersionRequest{}
	response, err := client.GrpcClient.GetVersion(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8006, err, details)
		}()
	}
	return result, err
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
The GetObserverOrigin method returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szproduct) GetObserverOrigin(ctx context.Context) string {
	_ = ctx
	return client.observerOrigin
}

/*
The Initialize method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szproduct) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(10, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szproduct) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, observer.GetObserverID(ctx))
		defer func() { client.traceExit(22, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers == nil {
		client.observers = &subject.SimpleSubject{}
	}
	err = client.observers.RegisterObserver(ctx, observer)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverID(ctx),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szproduct) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(13, logLevelName)
		defer func() { client.traceExit(14, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	err = client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8009, err, details)
		}()
	}
	return err
}

/*
The SetObserverOrigin method sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szproduct) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szproduct) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(23, observer.GetObserverID(ctx))
		defer func() { client.traceExit(24, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8010, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szproduct) getLogger() logging.Logging {
	var err error
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingLogger(ComponentID, szproductapi.IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *Szproduct) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szproduct) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}
