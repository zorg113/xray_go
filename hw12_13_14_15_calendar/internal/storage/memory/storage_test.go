package memorystorage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage"
)

//nolint:all
func TestStorage(t *testing.T) {
	t.Run("use storage", func(t *testing.T) {

		data1 := storage.Event{
			ID:          898989,
			Title:       "ToDo: title",
			StartDate:   time.Date(1999, time.April, 05, 20, 00, 00, 00, time.UTC),
			EndDate:     time.Date(1999, time.April, 05, 22, 00, 00, 00, time.UTC),
			Description: "Bla Bla Bla",
			OwnerID:     "CCCC",
			RemindIn:    "02.02.02",
		}
		data2 := storage.Event{
			ID:          898988,
			Title:       "ToDSo: title",
			StartDate:   time.Date(1099, time.April, 05, 20, 00, 00, 00, time.UTC),
			EndDate:     time.Date(1099, time.April, 05, 22, 00, 00, 00, time.UTC),
			Description: "Bla Bla Bl",
			OwnerID:     "CCCXC",
			RemindIn:    "02.03.02",
		}
		data3 := storage.Event{
			ID:          898988,
			Title:       "ToDSo: title",
			StartDate:   time.Date(1999, time.August, 05, 20, 00, 00, 00, time.UTC),
			EndDate:     time.Date(1999, time.August, 05, 22, 00, 00, 00, time.UTC),
			Description: "Bla Bla Bl",
			OwnerID:     "CCCXC",
			RemindIn:    "02.03.02",
		}
		var err error
		stor := New()
		err = stor.CreateEvent(data1)
		require.NoError(t, err)
		err = stor.CreateEvent(data2)
		require.NoError(t, err)
		event, err := stor.GetEvents(
			time.Date(1990, time.April, 05, 20, 00, 00, 00, time.UTC),
			time.Date(2000, time.April, 05, 20, 00, 00, 00, time.UTC))
		require.NoError(t, err)
		require.Equal(t, len(event), 1)
		err = stor.UpdateEvent(data3)
		require.NoError(t, err)
		event, err = stor.GetEvents(
			time.Date(1990, time.April, 05, 20, 00, 00, 00, time.UTC),
			time.Date(2000, time.April, 05, 20, 00, 00, 00, time.UTC))
		require.NoError(t, err)
		require.Equal(t, len(event), 2)
	})
}
