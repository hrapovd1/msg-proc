// Модуль dbstorage содержит типы и методы для хранения метрик
// в базе postgresql.
package dbstorage

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/hrapovd1/msg-proc/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBStorage тип для хранения метрик в базу
type DBStorage struct {
	dbConnect *sql.DB
	logger    *log.Logger
	tableName string
}

// NewDBStorage возвращает тип DBStorage по полученному конфигу
func NewDBStorage(dsn string, logger *log.Logger) (*DBStorage, error) {
	db := DBStorage{
		logger:    logger,
		tableName: "",
	}
	dbConnect, err := sql.Open("pgx", dsn)
	db.dbConnect = dbConnect
	return &db, err
}

// Save сохраняет новое значение сообщения
func (ds *DBStorage) Save(ctx context.Context, msgID, message string) {
	update := false
	msg := types.MessageModel{
		ID:      msgID,
		Message: message,
		Status:  "",
	}
	if err := ds.store(ctx, &msg, update); err != nil {
		if ds.logger != nil {
			ds.logger.Println(err)
		}
	}
}

// Update обновляет статус сообщения после обработки
func (ds *DBStorage) Update(ctx context.Context, msgID, status string) {
	update := true
	msg := types.MessageModel{
		ID:     msgID,
		Status: "processed",
	}
	if err := ds.store(ctx, &msg, update); err != nil {
		if ds.logger != nil {
			ds.logger.Println(err)
		}
	}
}

// Close закрывает подключение к БД, необходимо запускать в defer
func (ds *DBStorage) Close() error {
	ds.logger.Print("call DBStorage.Close()")
	return ds.dbConnect.Close()
}

// Ping используется для проверки доступности базы
func (ds *DBStorage) Ping(ctx context.Context) bool {
	if ds.dbConnect == nil {
		return false
	}
	ctxT, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := ds.dbConnect.PingContext(ctxT); err != nil {
		return false
	}
	return true
}

func (ds *DBStorage) GetCount(ctx context.Context, processed bool) int64 {
	var count int64
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: ds.dbConnect}), &gorm.Config{})
	if err != nil {
		ds.logger.Printf("can't connect to db, when sync metric: %v\n", err)
		return count
	}
	tableName := strings.ToLower(types.DBtableName)
	select {
	case <-ctx.Done():
		return count
	default:
		if processed {
			db.Table(tableName).Where("status = ?", "processed").Count(&count)
		} else {
			db.Table(tableName).Count(&count)
		}
		return count
	}
}

// store внутренняя функция сохранения метрики в базу
func (ds *DBStorage) store(ctx context.Context, message *types.MessageModel, isUpdate bool) error {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: ds.dbConnect}), &gorm.Config{})
	if err != nil {
		return err
	}
	tableName := strings.ToLower(types.DBtableName)
	select {
	case <-ctx.Done():
		return nil
	default:
		if ds.tableName == "" {
			if !db.Migrator().HasTable(tableName) {
				if err := db.Table(tableName).Migrator().CreateTable(&types.MessageModel{}); err != nil {
					return err
				}
			}
			ds.tableName = types.DBtableName
		}
		if !isUpdate {
			db.Table(tableName).Create(message)
		} else {
			db.Table(tableName).Model(message).Where("id = ?", message.ID).Update("status", message.Status)
		}
		return nil
	}
}
