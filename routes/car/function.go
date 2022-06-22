package car

import "github.com/jmoiron/sqlx"

//GET CAR
//
func getCar(tx *sqlx.DB) ([]GetCar, error) {
	var (
		data  GetCar
		datas []GetCar
	)
	query := (`select * from "cars" c`)

	rows, err := tx.Queryx(query)
	if err != nil {
		return datas, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.StructScan(&data)
		if err != nil {
			return datas, err
		}
		datas = append(datas, data)
	}
	return datas, err
}

//GET CAR BY ID
//
func getCarByID(tx *sqlx.DB, id int) (GetCar, error) {
	var (
		data GetCar
	)

	query := (`select * from "cars" c
			where c."cars_id" = $1`)

	values := []interface{}{
		id,
	}
	err := tx.QueryRowx(query, values...).StructScan(&data)
	if err != nil {
		return data, err
	}
	return data, err
}

//ADD CAR
//
func addCar(tx *sqlx.DB, input CarForm) (int, error) {
	query := (`insert into "cars" (name,rent_price_daily,stock)
			values($1,$2,$3)
			returning "cars_id"`)
	values := []interface{}{
		input.Name,
		input.Rent_price_daily,
		input.Stock,
	}

	var id int
	err := tx.QueryRowx(query, values...).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//UPDATE CAR
//
func updateCar(tx *sqlx.DB, id int, input CarForm) error {
	query := (`update "cars"
		set "name" = $1,
		"rent_price_daily" = $2,         
		"stock"=$3
		where cars_id=$4`)
	values := []interface{}{
		input.Name,
		input.Rent_price_daily,
		input.Stock,
		id,
	}
	_, err := tx.Exec(query, values...)

	if err != nil {
		return err
	}

	return err
}

//DELETE CAR
//
func deleteCar(tx *sqlx.DB, id int) error {
	query := (`delete from "cars"c
		where "cars_id" =$1`)

	values := []interface{}{
		id,
	}

	_, err := tx.Exec(query, values...)

	if err != nil {
		return err
	}

	return err
}
