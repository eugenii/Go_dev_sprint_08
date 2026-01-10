package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") 
	if err != nil {
		return err
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
    id, err := store.Add(parcel)
    if err != nil {
        t.Errorf("Ошибка при добавлении посылки: %v", err)
    }
    if id == 0 {
        t.Errorf("Идентификатор посылки равен 0")
    }
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
    retrievedParcel, err := store.Get(id)
    if err != nil {
        t.Errorf("Ошибка при получении посылки: %v", err)
    }

    // Проверяем совпадение полей
    if retrievedParcel.Number != id {
        t.Errorf("Неверный номер посылки: ожидалось %d, получено %d", id, retrievedParcel.Number)
    }
    if retrievedParcel.Client != parcel.Client {
        t.Errorf("Неверный клиент: ожидалось %d, получено %d", parcel.Client, retrievedParcel.Client)
    }
    if retrievedParcel.Status != parcel.Status {
        t.Errorf("Неверный статус: ожидалось %s, получено %s", parcel.Status, retrievedParcel.Status)
    }
    if retrievedParcel.Address != parcel.Address {
        t.Errorf("Неверный адрес: ожидалось %s, получено %s", parcel.Address, retrievedParcel.Address)
    }
    if retrievedParcel.CreatedAt != parcel.CreatedAt {
        t.Errorf("Неверная дата создания: ожидалось %s, получено %s", parcel.CreatedAt, retrievedParcel.CreatedAt)
    }
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
    err = store.Delete(id)
    if err != nil {
        t.Errorf("Ошибка при удалении посылки: %v", err)
    }

    // Проверяем, что посылка больше не существует
    _, err = store.Get(id)
    if err == nil {
        t.Errorf("Посылка всё ещё существует после удаления")
    }
    if err != sql.ErrNoRows {
        t.Errorf("Неверная ошибка при попытке получения удаленной посылки: %v", err)
    }
	
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")  // настройте подключение к БД

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")  // настройте подключение к БД

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// set status
	// обновите статус, убедитесь в отсутствии ошибки

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")  // настройте подключение к БД

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
