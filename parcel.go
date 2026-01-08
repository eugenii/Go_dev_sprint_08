package main

import (
	"database/sql"
)

// ParcelStore предоставляет доступ к данным посылок в БД
type ParcelStore struct {
	db *sql.DB
}

// NewParcelStore создает новый объект ParcelStore для работы с БД
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add добавляет новую посылку в БД и возвращает её идентификатор
func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	// верните идентификатор последней добавленной записи
	return 0, nil
}

// Get возвращает посылку по её номеру
func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	return p, nil
}

// GetByClient возвращает список посылок клиента по его идентификатору
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	return res, nil
}

// SetStatus обновляет статус посылки по её номеру
func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	return nil
}

// SetAddress обновляет адрес посылки по её номеру, если её статус "registered"
func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	return nil
}

// Delete удаляет посылку по её номеру, если её статус "registered"
func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	return nil
}
