package repository

import (
	"avito/internal/domain"
	"database/sql"
	"fmt"
	"log"
)

type AvitoRepository interface {
	Create(user domain.User) error
	GetByEmail(email string) (domain.User, error)

	GetHouseByAddress(address string) (domain.House, error)
	CreateHouse(house domain.House) (domain.House, error)

	CreateFlat(flat domain.Flat) (domain.Flat, error)
	UpdateForMod(id int, status string) (domain.Flat, error)
	GetFlatsByHouseId(houseId int, status string) ([]domain.Flat, error)
}

type avitoRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) AvitoRepository {
	return &avitoRepository{db: db}
}

func (r *avitoRepository) Create(user domain.User) error {
	resp, err := r.db.Exec("INSERT INTO users (email, password, user_type) VALUES ($1, $2, $3)", user.Email, user.Password, user.UserType)
	if err != nil {
		log.Printf("Error occured in database: %s", err.Error())
		return err
	}

	rowsAffected, err := resp.RowsAffected()
	if err != nil {
		log.Printf("Error in reading database's response: %s", err.Error())
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No rows was affected")
		return err
	}

	return nil
}

func (r *avitoRepository) GetByEmail(email string) (domain.User, error) {
	row := r.db.QueryRow("SELECT id, email, password, user_type FROM users WHERE email=$1", email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("no user")
		}
		return user, err
	}

	return user, nil
}

func (r *avitoRepository) CreateHouse(house domain.House) (domain.House, error) {
	var developer interface{}
	if house.Developer == nil {
		developer = nil
	} else {
		developer = house.Developer
	}

	var houseResp domain.House
	err := r.db.QueryRow("INSERT INTO houses (address, year, developer) VALUES ($1, $2, $3) RETURNING id, address, year, developer, created_at, updated_at",
		house.Address, house.Year, developer).Scan(&houseResp.ID, &houseResp.Address, &houseResp.Year,
		&houseResp.Developer, &houseResp.CreatedAt, &houseResp.UpdatedAt, &houseResp.LastFlatAddedAt)

	if err != nil {
		return houseResp, fmt.Errorf("error inserting new house: %s", err.Error())
	}

	return houseResp, nil
}

func (r *avitoRepository) GetHouseByAddress(address string) (domain.House, error) {
	row := r.db.QueryRow("SELECT id, address, year, developer FROM houses WHERE address=$1", address)
	var house domain.House
	err := row.Scan(house)
	if err != nil {
		if err == sql.ErrNoRows {
			return house, fmt.Errorf("no house at this address")
		}
		return house, err
	}

	return house, nil
}

func (r *avitoRepository) CreateFlat(flat domain.Flat) (domain.Flat, error) {
	var flatResp domain.Flat
	trx, err := r.db.Begin()
	if err != nil {
		return flatResp, err
	}
	defer trx.Rollback()
	err = trx.QueryRow("INSERT INTO flats (house_id, price, rooms, flat_number) VALUES ($1,$2, $3, $4) RETURNING id, house_id, price, rooms, status, created_at, updated_at, flat_number", flat.HouseID, flat.Price, flat.Rooms, flat.FlatNumber).Scan(
		&flatResp.ID,
		&flatResp.HouseID,
		&flatResp.Price,
		&flatResp.Rooms,
		&flatResp.Status,
		&flatResp.CreatedAt,
		&flatResp.UpdatedAt,
		&flatResp.FlatNumber,
	)

	if err != nil {
		return flatResp, fmt.Errorf("error inserting new flat: %s", err.Error())
	}

	_, err = r.db.Exec("UPDATE houses SET updated_at=NOW() WHERE id=$1", flatResp.HouseID)
	if err != nil {
		return domain.Flat{}, err
	}

	err = trx.Commit()
	if err != nil {
		return domain.Flat{}, err
	}

	return flatResp, nil
}

func (r *avitoRepository) UpdateForMod(id int, status string) (domain.Flat, error) {
	var flatResp domain.Flat
	err := r.db.QueryRow("UPDATE flats SET status=$1, updated_at=NOW() WHERE id=$2 RETURNING id, house_id, price, rooms, status, created_at, updated_at, flat_number", status, id).Scan(
		&flatResp.ID,
		&flatResp.HouseID,
		&flatResp.Price,
		&flatResp.Rooms,
		&flatResp.Status,
		&flatResp.CreatedAt,
		&flatResp.UpdatedAt,
		&flatResp.FlatNumber,
	)
	if err != nil {
		return domain.Flat{}, fmt.Errorf("error updating flats status %s", err.Error())
	}

	return flatResp, nil
}

func (r *avitoRepository) GetFlatsByHouseId(houseId int, status string) ([]domain.Flat, error) {
	flatsArray := []domain.Flat{}
	var flatsTemp domain.Flat
	var rows *sql.Rows
	var err error
	if status == "moderator" {
		rows, err = r.db.Query("SELECT id, house_id, price, rooms, status, created_at, updated_at, flat_number FROM flats WHERE house_id=$1", houseId)
		if err != nil {
			if err == sql.ErrNoRows {
				return []domain.Flat{}, fmt.Errorf("no house with id %d", houseId)
			}
			return []domain.Flat{}, fmt.Errorf("error reading rows")
		}
	}
	if status == "client" {
		rows, err = r.db.Query("SELECT id, house_id, price, rooms, status, created_at, updated_at, flat_number FROM flats WHERE house_id=$1 AND status='approved'", houseId)
		if err != nil {
			if err == sql.ErrNoRows {
				return []domain.Flat{}, fmt.Errorf("no house with id %d", houseId)
			}
			return []domain.Flat{}, fmt.Errorf("error reading rows")
		}
	}

	for rows.Next() {
		err := rows.Scan(
			&flatsTemp.ID,
			&flatsTemp.HouseID,
			&flatsTemp.Price,
			&flatsTemp.Rooms,
			&flatsTemp.Status,
			&flatsTemp.CreatedAt,
			&flatsTemp.UpdatedAt,
			&flatsTemp.FlatNumber,
		)
		if err != nil {
			return []domain.Flat{}, fmt.Errorf("error reading rows")
		}
		flatsArray = append(flatsArray, flatsTemp)
	}
	defer rows.Close()

	return flatsArray, nil
}
