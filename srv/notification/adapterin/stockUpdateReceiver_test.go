package adapterin

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestStockUpdateReceiver(t *testing.T) {
	ts := start(t)

	ts.stockUpdates.EXPECT().RecordStockUpdate(gomock.Any()).Return(nil)

	stockUpdate := stream.StockUpdate{
		ID:          "1",
		WarehouseID: "1",
		Type:        stream.StockUpdateTypeAdd,
		Goods: []stream.StockUpdateGood{
			{GoodID: "1", Quantity: 10, Delta: 10},
		},
		OrderID:       "o",
		TransferID:    "",
		ReservationID: "",
		Timestamp:     time.Now().UnixMilli(),
	}

	payload, err := json.Marshal(stockUpdate)
	require.NoError(t, err)

	_, err = ts.js.Publish(t.Context(), "stock.update.1", payload)
	require.NoError(t, err)

	s, err := ts.js.Stream(t.Context(), stream.StockUpdateStreamConfig.Name)
	require.NoError(t, err)

	i, err := s.Info(t.Context())
	require.NoError(t, err)
	require.Equal(t, uint64(1), i.State.Msgs)
}
