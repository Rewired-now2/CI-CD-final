package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
	ParcelStatusDelivered  = "delivered"
)

type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

type ParcelService struct {
	store ParcelStore
}

func NewParcelService(store ParcelStore) ParcelService {
	return ParcelService{store: store}
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	);
	`)
	return err
}

func (s ParcelService) Register(client int, address string) (Parcel, error) {
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := s.store.Add(parcel)
	if err != nil {
		return parcel, err
	}

	parcel.Number = id

	fmt.Printf("Новая посылка № %d на адрес %s\n", parcel.Number, parcel.Address)

	return parcel, nil
}

func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}

	for _, p := range parcels {
		fmt.Printf("Parcel #%d | %s | %s\n", p.Number, p.Status, p.Address)
	}
	return nil
}

func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}

	var next string

	switch parcel.Status {
	case ParcelStatusRegistered:
		next = ParcelStatusSent
	case ParcelStatusSent:
		next = ParcelStatusDelivered
	default:
		return nil
	}

	fmt.Printf("Status updated: %s -> %s\n", parcel.Status, next)

	return s.store.SetStatus(number, next)
}

func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address)
}

func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}

func main() {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		fmt.Println("DB init error:", err)
		return
	}

	store := NewParcelStore(db)
	service := NewParcelService(store)

	client := 1
	address := "Test address"

	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	service.NextStatus(p.Number)
	service.PrintClientParcels(client)

	service.Delete(p.Number)
}
