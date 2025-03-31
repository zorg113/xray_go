package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/app"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage/memory"
	cnd_pb "github.com/zorg113/xray_go/hw12_13_14_15_calendar/pkg/calendar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_ServerGRPC(t *testing.T) {
	// Server initialization
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})
	log, err := logger.New("INFO", "")
	require.NoError(t, err)
	stor := memorystorage.New()
	app := app.New(log, stor)
	svc := Server{
		app: app,
	}
	cnd_pb.RegisterCalendarServer(srv, &svc)

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Error(err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	resolver.SetDefaultScheme("passthrough")
	conn, err := grpc.NewClient("", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	t.Cleanup(func() {
		conn.Close()
	})
	require.NoError(t, err)
	ctx := context.Background()
	client := cnd_pb.NewCalendarClient(conn)
	res, err := client.CreateEvent(ctx, &cnd_pb.CreateEventRequest{
		Event: &cnd_pb.Event{
			Id:          1,
			Title:       "Test event",
			Description: "This is a test event",
			StartDate: &timestamppb.Timestamp{
				Seconds: int64(time.Date(1979, time.April, 7, 20, 00, 00, 00, time.UTC).Unix()), //nolint:gofumpt,unconvert
			},
			EndDate: &timestamppb.Timestamp{
				Seconds: int64(time.Date(1979, time.April, 8, 20, 00, 00, 00, time.UTC).Unix()), //nolint:gofumpt,unconvert
			},
			RemindIn: "15m",
			OwnerId:  1,
		},
	})
	if err != nil {
		t.Error(err)
	}
	require.NotNil(t, res)
	events, err := stor.GetEvents(time.Date(1979, time.April, 5, 20, 00, 00, 00, time.UTC), time.Date(1979, time.April, 9, 20, 00, 00, 00, time.UTC)) //nolint:lll,gofumpt
	require.NoError(t, err)
	require.Equal(t, len(events), 1)
	require.Equal(t, events[0].Title, "Test event")
	require.Equal(t, events[0].Description, "This is a test event")
	require.Equal(t, events[0].StartDate, time.Date(1979, time.April, 7, 20, 00, 00, 00, time.UTC)) //nolint:gofumpt
	require.Equal(t, events[0].EndDate, time.Date(1979, time.April, 8, 20, 00, 00, 00, time.UTC))
	require.Equal(t, events[0].RemindIn, "15m")
	require.Equal(t, events[0].OwnerID, "1")

	resupd, err := client.UpdateEvent(ctx,
		&cnd_pb.UpdateEventRequest{
			Event: &cnd_pb.Event{
				Id:          1,
				Title:       "Test event",
				Description: "This is a test event",
				StartDate: &timestamppb.Timestamp{
					Seconds: int64(time.Date(1979, time.April, 7, 20, 00, 00, 00, time.UTC).Unix()), //nolint:gofumpt,unconvert
				},
				EndDate: &timestamppb.Timestamp{
					Seconds: int64(time.Date(1979, time.April, 8, 22, 00, 00, 00, time.UTC).Unix()), //nolint:gofumpt,unconvert
				},
				RemindIn: "15m",
				OwnerId:  1,
			},
		})
	require.NoError(t, err)
	require.NotNil(t, resupd)
	events, err = stor.GetEvents(time.Date(1979, time.April, 5, 20, 00, 00, 00, time.UTC), time.Date(1979, time.April, 9, 20, 00, 00, 00, time.UTC)) //nolint:lll,gofumpt
	require.NoError(t, err)
	require.Equal(t, events[0].EndDate, time.Date(1979, time.April, 8, 22, 00, 00, 00, time.UTC)) //nolint:gofumpt
	require.NoError(t, err)
}
