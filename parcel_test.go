package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"
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
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при подключению к БД: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE parcel (
            number INTEGER PRIMARY KEY AUTOINCREMENT,
            client INTEGER NOT NULL,
            status TEXT NOT NULL,
            address TEXT NOT NULL,
            created_at TEXT NOT NULL
        )
    `)
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	if err != nil {
		t.Errorf("Ошибка при добавлении посылки: %v", err)
	}
	if id == 0 {
		t.Errorf("Идентификатор посылки равен 0")
	}

	retrievedParcel, err := store.Get(id)
	if err != nil {
		t.Errorf("Ошибка при получении посылки: %v", err)
	}

	if retrievedParcel.Number != id {
		t.Errorf("Неверный номер посылки: ожидалось %d, получено %d", id, retrievedParcel.Number)
	}

	err = store.Delete(id)
	if err != nil {
		t.Errorf("Ошибка при удалении посылки: %v", err)
	}

	_, err = store.Get(id)
	if err == nil {
		t.Errorf("Посылка всё ещё существует после удаления")
	} else if err != sql.ErrNoRows {
		t.Errorf("Неверная ошибка при попытке получения удаленной посылки: ожидалось sql.ErrNoRows, получено '%v'", err)
	}
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при подключении к БД: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE parcel (
            number INTEGER PRIMARY KEY AUTOINCREMENT,
            client INTEGER NOT NULL,
            status TEXT NOT NULL,
            address TEXT NOT NULL,
            created_at TEXT NOT NULL
        )
    `)
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	if err != nil {
		t.Errorf("Ошибка при добавлении посылки: %v", err)
	}
	if id == 0 {
		t.Errorf("Идентификатор посылки равен 0")
	}

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	if err != nil {
		t.Errorf("Ошибка при обновлении адреса: %v", err)
	}

	updatedParcel, err := store.Get(id)
	if err != nil {
		t.Errorf("Ошибка при получении посылки: %v", err)
	} else {
		if updatedParcel.Address != newAddress {
			t.Errorf("Адрес не обновился. Ожидалось %s, получено %s", newAddress, updatedParcel.Address)
		}
	}
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		t.Errorf("Ошибка при подключении к БД: %v", err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	service := NewParcelService(store)
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

	// Переход к следующему статусу
	err = service.NextStatus(id)
	if err != nil {
		t.Errorf("Ошибка при обновлении статуса: %v", err)
	}

	// check
	updatedParcel, err := store.Get(id)
	if err != nil {
		t.Errorf("Ошибка при получении посылки: %v", err)
	} else {
		expectedStatus := ParcelStatusSent // Ожидаемый статус после первого вызова NextStatus
		if updatedParcel.Status != expectedStatus {
			t.Errorf("Статус не обновился. Ожидалось %s, получено %s", expectedStatus, updatedParcel.Status)
		}
	}

	// Проверка второго перехода статуса
	err = service.NextStatus(id)
	if err != nil {
		t.Errorf("Ошибка при обновлении статуса: %v", err)
	}

	updatedParcel, err = store.Get(id)
	if err != nil {
		t.Errorf("Ошибка при получении посылки: %v", err)
	} else {
		expectedStatus := ParcelStatusDelivered // Ожидаемый статус после второго вызова NextStatus
		if updatedParcel.Status != expectedStatus {
			t.Errorf("Статус не обновился. Ожидалось %s, получено %s", expectedStatus, updatedParcel.Status)
		}
	}

	// Проверка третьего вызова (статус уже delivered)
	err = service.NextStatus(id)
	if err != nil {
		t.Errorf("Ошибка при обновлении статуса: %v", err)
	}

	updatedParcel, err = store.Get(id)
	if err != nil {
		t.Errorf("Ошибка при получении посылки: %v", err)
	} else {
		expectedStatus := ParcelStatusDelivered // Статус не должен измениться
		if updatedParcel.Status != expectedStatus {
			t.Errorf("Статус неожиданно изменился. Ожидалось %s, получено %s", expectedStatus, updatedParcel.Status)
		}
	}
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при подключению к БД: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE parcel (
            number INTEGER PRIMARY KEY AUTOINCREMENT,
            client INTEGER NOT NULL,
            status TEXT NOT NULL,
            address TEXT NOT NULL,
            created_at TEXT NOT NULL
        )
    `)
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := make(map[int]Parcel)

	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	for i := range parcels {
		id, err := store.Add(parcels[i])
		if err != nil {
			t.Errorf("Ошибка при добавлении посылки: %v", err)
		}
		if id == 0 {
			t.Errorf("Идентификатор посылки равен 0")
		}

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	if err != nil {
		t.Errorf("Ошибка при получении посылок: %v", err)
	}

	if len(storedParcels) != len(parcels) {
		t.Errorf("Количество полученных посылок (%d) не совпадает с количеством добавленных (%d)", len(storedParcels), len(parcels))
	}

	for _, storedParcel := range storedParcels {
		expectedParcel, exists := parcelMap[storedParcel.Number]
		if !exists {
			t.Errorf("Посылка с номером %d не найдена среди добавленных", storedParcel.Number)
			continue
		}

		if storedParcel.Client != expectedParcel.Client {
			t.Errorf("Неверный клиент: ожидалось %d, получено %d", expectedParcel.Client, storedParcel.Client)
		}
		if storedParcel.Status != expectedParcel.Status {
			t.Errorf("Неверный статус: ожидалось %s, получено %s", expectedParcel.Status, storedParcel.Status)
		}
		if storedParcel.Address != expectedParcel.Address {
			t.Errorf("Неверный адрес: ожидалось %s, получено %s", expectedParcel.Address, storedParcel.Address)
		}
		if storedParcel.CreatedAt != expectedParcel.CreatedAt {
			t.Errorf("Неверная дата создания: ожидалось %s, получено %s", expectedParcel.CreatedAt, storedParcel.CreatedAt)
		}
	}
}
