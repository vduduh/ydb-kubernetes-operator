package healthcheck

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ydb-platform/ydb-go-genproto/Ydb_Monitoring_V1"
	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Monitoring"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/balancers"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
)

func TestXxx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db, err := ydb.Open(
		ctx,
		os.Getenv("YDB_CONNECTION_STRING"),
		ydb.WithStaticCredentials(os.Getenv("YDB_STATIC_USERNAME"), os.Getenv("YDB_STATIC_PASSWORD")),
		ydb.WithBalancer(balancers.SingleConn()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		// cleanup connection
		if e := db.Close(ctx); e != nil {
			t.Fatalf("close failed: %+v", e)
		}
	}()
	t.Run("monitoring.SelfCheck", func(t *testing.T) {
		if err = retry.Retry(ctx, func(ctx context.Context) (err error) {
			client := Ydb_Monitoring_V1.NewMonitoringServiceClient(ydb.GRPCConn(db))
			response, err := client.SelfCheck(ctx, &Ydb_Monitoring.SelfCheckRequest{
				OperationParams:     nil,
				ReturnVerboseStatus: false,
				MinimumStatus:       0,
				MaximumLevel:        0,
			})
			if err != nil {
				return err
			}
			var result Ydb_Monitoring.SelfCheckResult
			err = response.Operation.Result.UnmarshalTo(&result)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", &result)
			return nil
		}, retry.WithIdempotent(true)); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})
}
