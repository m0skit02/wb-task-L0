package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
	"wb-task-L0/internal/models"
	"wb-task-L0/internal/repository"
)

func newGormMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		return nil, nil, err
	}

	// Говорим sqlmock, что порядок вызовов может быть любым
	mock.MatchExpectationsInOrder(false)

	// Оборачиваем sqlmock.DB в gorm.DB
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}

	return gdb, mock, nil
}

func TestOrderRepo_Create(t *testing.T) {
	db, mock, err := newGormMock()
	require.NoError(t, err)

	repo := repository.NewOrderRepo(db)

	// Заглушка заказа
	order := &models.Order{
		OrderUID:    "order123",
		TrackNumber: "track456",
		Entry:       "entry",
		Locale:      "ru",
		CustomerID:  "cust1",
		DateCreated: time.Now(),
		OofShard:    "oof",
		Delivery: models.Delivery{
			DeliveryID: "del1",
			OrderUID:   "order123",
			Name:       "John",
		},
		Payment: models.Payment{
			PaymentID:   "pay1",
			OrderUID:    "order123",
			Transaction: "tx123",
			Currency:    "RUB",
			Amount:      1000,
		},
		Items: []models.Item{
			{
				ItemID:      "it1",
				OrderUID:    "order123",
				ChrtID:      1,
				TrackNumber: "track456",
				Price:       500,
				Name:        "item1",
			},
		},
	}

	// --- Ожидания SQL ---
	mock.ExpectBegin()

	// Insert в orders
	mock.ExpectExec(`INSERT INTO "orders"`).
		WithArgs(
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.SmID,
			sqlmock.AnyArg(), // DateCreated
			order.OofShard,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Insert в deliveries
	mock.ExpectExec(`INSERT INTO "deliveries"`).
		WithArgs(
			order.Delivery.DeliveryID,
			order.OrderUID,
			order.Delivery.Name,
			order.Delivery.Phone,
			order.Delivery.Zip,
			order.Delivery.City,
			order.Delivery.Address,
			order.Delivery.Region,
			order.Delivery.Email,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Insert в payments
	mock.ExpectExec(`INSERT INTO "payments"`).
		WithArgs(
			order.Payment.PaymentID,
			order.OrderUID,
			order.Payment.Transaction,
			order.Payment.RequestID,
			order.Payment.Currency,
			order.Payment.Provider,
			order.Payment.Amount,
			order.Payment.PaymentDt,
			order.Payment.Bank,
			order.Payment.DeliveryCost,
			order.Payment.GoodsTotal,
			order.Payment.CustomFee,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Insert в items
	mock.ExpectExec(`INSERT INTO "items"`).
		WithArgs(
			order.Items[0].ItemID,
			order.OrderUID,
			order.Items[0].ChrtID,
			order.Items[0].TrackNumber,
			order.Items[0].Price,
			order.Items[0].Rid,
			order.Items[0].Name,
			order.Items[0].Sale,
			order.Items[0].Size,
			order.Items[0].TotalPrice,
			order.Items[0].NmID,
			order.Items[0].Brand,
			order.Items[0].Status,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// --- Вызов ---
	gotUID, err := repo.Create(order)

	// --- Проверки ---
	require.NoError(t, err)
	require.Equal(t, order.OrderUID, gotUID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepo_GetAll(t *testing.T) {
	db, mock, err := newGormMock()
	assert.NoError(t, err)

	repo := repository.NewOrderRepo(db)

	// Заглушка данных
	order := models.Order{
		OrderUID:    "order123",
		TrackNumber: "track456",
		Entry:       "entry",
		Locale:      "ru",
		CustomerID:  "cust1",
		DateCreated: time.Now(),
		OofShard:    "oof",
		Delivery: models.Delivery{
			DeliveryID: "del1",
			OrderUID:   "order123",
			Name:       "John",
		},
		Payment: models.Payment{
			PaymentID:   "pay1",
			OrderUID:    "order123",
			Transaction: "tx123",
			Currency:    "RUB",
			Amount:      1000,
		},
		Items: []models.Item{
			{
				ItemID:      "it1",
				OrderUID:    "order123",
				ChrtID:      1,
				TrackNumber: "track456",
				Price:       500,
				Name:        "item1",
			},
		},
	}

	// --- Моки SQL-запросов GORM ---
	// 1. Основной SELECT orders
	rowsOrders := sqlmock.NewRows([]string{
		"order_uid", "track_number", "entry", "locale", "customer_id", "date_created", "oof_shard",
	}).AddRow(
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.CustomerID, order.DateCreated, order.OofShard,
	)
	mock.ExpectQuery(`SELECT .* FROM "orders"`).WillReturnRows(rowsOrders)

	// 2. Preload deliveries
	rowsDelivery := sqlmock.NewRows([]string{
		"delivery_id", "order_uid", "name",
	}).AddRow(
		order.Delivery.DeliveryID, order.OrderUID, order.Delivery.Name,
	)
	mock.ExpectQuery(`SELECT .* FROM "deliveries"`).
		WithArgs(order.OrderUID).
		WillReturnRows(rowsDelivery)

	// 3. Preload payments
	rowsPayment := sqlmock.NewRows([]string{
		"payment_id", "order_uid", "transaction", "currency", "amount",
	}).AddRow(
		order.Payment.PaymentID, order.OrderUID, order.Payment.Transaction, order.Payment.Currency, order.Payment.Amount,
	)
	mock.ExpectQuery(`SELECT .* FROM "payments"`).
		WithArgs(order.OrderUID).
		WillReturnRows(rowsPayment)

	// 4. Preload items
	rowsItems := sqlmock.NewRows([]string{
		"item_id", "order_uid", "chrt_id", "track_number", "price", "name",
	}).AddRow(
		order.Items[0].ItemID, order.OrderUID, order.Items[0].ChrtID, order.Items[0].TrackNumber, order.Items[0].Price, order.Items[0].Name,
	)
	mock.ExpectQuery(`SELECT .* FROM "items"`).
		WithArgs(order.OrderUID).
		WillReturnRows(rowsItems)

	// --- Вызов ---
	got, err := repo.GetAll()
	assert.NoError(t, err)

	// Проверки
	assert.Len(t, got, 1)
	assert.Equal(t, order.OrderUID, got[0].OrderUID)
	assert.Equal(t, order.Delivery.Name, got[0].Delivery.Name)
	assert.Equal(t, order.Payment.Transaction, got[0].Payment.Transaction)
	assert.Equal(t, order.Items[0].Name, got[0].Items[0].Name)

	// Проверяем что все моки сработали
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepo_GetByID(t *testing.T) {
	db, mock, err := newGormMock()
	assert.NoError(t, err)

	repo := repository.NewOrderRepo(db)

	// Заглушка
	order := models.Order{
		OrderUID:    "order123",
		TrackNumber: "track456",
		Entry:       "entry",
		Locale:      "ru",
		CustomerID:  "cust1",
		DateCreated: time.Now(),
		OofShard:    "oof",
		Delivery: models.Delivery{
			DeliveryID: "del1",
			OrderUID:   "order123",
			Name:       "John",
		},
		Payment: models.Payment{
			PaymentID:   "pay1",
			OrderUID:    "order123",
			Transaction: "tx123",
			Currency:    "RUB",
			Amount:      1000,
		},
		Items: []models.Item{
			{
				ItemID:      "it1",
				OrderUID:    "order123",
				ChrtID:      1,
				TrackNumber: "track456",
				Price:       500,
				Name:        "item1",
			},
		},
	}

	// --- Моки SQL-запросов ---
	// 1. SELECT из orders
	rowsOrders := sqlmock.NewRows([]string{
		"order_uid", "customer_id", "track_number",
	}).AddRow(order.OrderUID, order.CustomerID, order.TrackNumber)

	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE order_uid = \$1 ORDER BY "orders"."order_uid" LIMIT \$2`).
		WithArgs(order.OrderUID, 1).
		WillReturnRows(rowsOrders)

	// 2. Preload Delivery
	rowsDelivery := sqlmock.NewRows([]string{
		"delivery_id", "order_uid", "name",
	}).AddRow(
		order.Delivery.DeliveryID, order.OrderUID, order.Delivery.Name,
	)
	mock.ExpectQuery(`SELECT .* FROM "deliveries" WHERE "deliveries"."order_uid" = \$1`).
		WithArgs(order.OrderUID).
		WillReturnRows(rowsDelivery)

	// 3. Preload Payment
	rowsPayment := sqlmock.NewRows([]string{
		"payment_id", "order_uid", "transaction", "currency", "amount",
	}).AddRow(
		order.Payment.PaymentID, order.OrderUID, order.Payment.Transaction, order.Payment.Currency, order.Payment.Amount,
	)
	mock.ExpectQuery(`SELECT .* FROM "payments" WHERE "payments"."order_uid" = \$1`).
		WithArgs(order.OrderUID).
		WillReturnRows(rowsPayment)

	// 4. Preload Items
	rowsItems := sqlmock.NewRows([]string{
		"item_id", "order_uid", "chrt_id", "track_number", "price", "name",
	}).AddRow(
		order.Items[0].ItemID, order.OrderUID, order.Items[0].ChrtID,
		order.Items[0].TrackNumber, order.Items[0].Price, order.Items[0].Name,
	)

	mock.ExpectQuery(`SELECT .* FROM "items" WHERE "items"."order_uid" = \$1`).
		WithArgs(order.OrderUID).
		WillReturnRows(rowsItems)

	// --- Вызов ---
	got, err := repo.GetByID(order.OrderUID)
	assert.NoError(t, err)

	// Проверки
	assert.Equal(t, order.OrderUID, got.OrderUID)
	assert.Equal(t, order.Delivery.Name, got.Delivery.Name)
	assert.Equal(t, order.Payment.Transaction, got.Payment.Transaction)
	assert.Equal(t, order.Items[0].Name, got.Items[0].Name)

	// Убеждаемся, что все моки сработали
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepo_Delete(t *testing.T) {
	db, mock, err := newGormMock()
	assert.NoError(t, err)

	repo := repository.NewOrderRepo(db)
	orderUID := "order123"

	// --- Моки SQL-запросов ---
	mock.ExpectBegin()

	// 1. Delete items
	mock.ExpectExec(`DELETE FROM "items" WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 2. Delete payments
	mock.ExpectExec(`DELETE FROM "payments" WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 3. Delete deliveries
	mock.ExpectExec(`DELETE FROM "deliveries" WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 4. Delete order
	mock.ExpectExec(`DELETE FROM "orders" WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	// --- Вызов ---
	err = repo.Delete(orderUID)
	assert.NoError(t, err)

	// Проверки
	assert.NoError(t, mock.ExpectationsWereMet())
}
