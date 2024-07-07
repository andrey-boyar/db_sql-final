package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	/*реализуйте добавление строки в таблицу parcel, используйте данные из переменной p, получите идентификатор новой вставленной записи*/
	res, err := s.db.Exec("INSERT INTO parcel (client, address, status, created_at) VALUES (:client, :address, :status, :created_at)",
		//sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("address", p.Address),
		sql.Named("status", p.Status),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить идентификатор последней вставки: %w", err)
	}
	// верните идентификатор последней добавленной записи
	return int(lastInsertID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	rows := s.db.QueryRow("SELECT number, client, address, status, created_at FROM parcel WHERE number = ?", sql.Named("number", p.Number))
	// Сканируйте значения из строки в поля объекта Parcel
	err := rows.Scan(&p.Number, &p.Client, &p.Address, &p.Status, &p.CreatedAt)
	if err != nil {
		return Parcel{}, fmt.Errorf("не удалось отсканировать строку: %w", err)
	}
	// заполните объект Parcel данными из таблицы
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// Создайте фрагмент для хранения полученных посылок
	var parcell []Parcel
	rows, err := s.db.Query("SELECT number, client, address, status, created_at FROM parcel WHERE client = ?")
	if err != nil {
		return parcell, fmt.Errorf("не удалось подготовить оператор select: %w", err)
	}
	defer rows.Close()
	// Итерация по строкам, возвращенным запросом
	for rows.Next() {
		// Создаем объект Parcel для хранения данных из каждого ряда
		var p Parcel
		// Сканируем значения из строки в поля объекта Parcel
		err = rows.Scan(&p.Number, &p.Client, &p.Address, &p.Status, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("не удалось отсканировать строку: %w", err)
		}
		// Добавить объект Parcel к фрагменту parcels
		parcell = append(parcell, p)
	}
	// Проверка на наличие ошибок во время итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации строк: %w", err)
	}
	return parcell, nil
	// заполните срез Parcel данными из таблицы
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?",
		sql.Named("status", ParcelStatusRegistered),
		sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("не удалось выполнить оператор обновления: %w", err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	/* реализуйте обновление адреса в таблице parcel
	менять адрес можно только если значение статуса registered*/
	_, err := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ? AND status = ?",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		//return
		return fmt.Errorf("не удалось выполнить оператор обновления: %w", err)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	/* реализуйте удаление строки из таблицы parcel
	удалять строку можно только если значение статуса registered*/
	stmt, err := s.db.Prepare("DELETE FROM parcel WHERE number = ? AND status = ?")
	if err != nil {
		return fmt.Errorf("не удалось подготовить оператор удаления: %w", err)
	}
	defer stmt.Close()
	// Выполните подготовленный оператор с номером и статусом участка
	_, err = stmt.Exec(number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("не удалось выполнить оператор delete: %w", err)
	}
	return nil
}
