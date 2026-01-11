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
	// реализуем добавление строки в таблицу parcel, используем данные из переменной p
	res, err := s.db.Exec(
		"INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client), sql.Named("status", p.Status), sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	identifier, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// возвращаем идентификатор последней добавленной записи
	return int(identifier), nil
}

// Get возвращает посылку по её номеру
func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуем чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow(
		"SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	// заполняем объект Parcel данными из таблицы
	p := Parcel{}
	// Извлекаем данные из строки и сканируем их в переменные
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		// Проверяем, если запись не найдена
		if err == sql.ErrNoRows {
			return Parcel{}, sql.ErrNoRows
		}
		// Возвращаем ошибку, если произошла другая проблема
		return Parcel{}, err
	}

	return p, nil
}

// GetByClient возвращает список посылок клиента по его идентификатору
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуем чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// заполняем срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

// SetStatus обновляет статус посылки по её номеру
func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуем обновление статуса в таблице parcel
	_, err := s.db.Exec(
		"UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

// SetAddress обновляет адрес посылки по её номеру, если её статус "registered"
func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуем обновление адреса в таблице parcel
	_, err := s.db.Exec(
		"UPDATE parcel SET address = :address WHERE number = :number AND status = 'registered'",
		sql.Named("address", address), sql.Named("number", number))
	// менять адрес можно только если значение статуса registered
	if err != nil {
		return err
	}

	return nil
}

// Delete удаляет посылку по её номеру, если её статус "registered"
func (s ParcelStore) Delete(number int) error {
	// реализуем удаление строки из таблицы parcel
	_, err := s.db.Exec(
		"DELETE FROM parcel WHERE number = :number AND status = 'registered'",
		sql.Named("number", number))
	// удалять строку можно только если значение статуса registered
	if err != nil {
		return err
	}

	return nil
}
