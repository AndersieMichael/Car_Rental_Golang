package driver

import "github.com/jmoiron/sqlx"

//GET Driver
//=============================================================
func getDriver(tx *sqlx.DB) ([]GetDriver, error) {
	var (
		data  GetDriver
		datas []GetDriver
	)
	query := (`select *
				from "driver" d`)

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

//GET Driver BY ID
//=============================================================
func GetDriverByID(tx *sqlx.DB, id int) (GetDriver, error) {
	var (
		data GetDriver
	)

	query := (`select *
	from "driver" d
	where d."driver_id" = $1`)

	values := []interface{}{
		id,
	}
	err := tx.QueryRowx(query, values...).StructScan(&data)
	if err != nil {
		return data, err
	}
	return data, err
}

//ADD Driver
//=============================================================
func addDriver(tx *sqlx.DB, input DriverForm) (int, error) {

	query := (`insert into driver (name,
			nik,
			phone_number,
			daily_cost)
			Values ($1,	$2,	$3,	$4)
			returning "driver_id"`)
	values := []interface{}{
		input.Name,
		input.Nik,
		input.Phone_number,
		input.Daily_cost,
	}
	var id int
	err := tx.QueryRowx(query, values...).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//UPDATE Driver
//=============================================================
func updateDriver(tx *sqlx.DB, id int, input DriverForm) error {
	query := (`update "driver"
		set "name" = $1,
		"nik" = $2,         
		"phone_number"=$3,
		"daily_cost" = $4
		where driver_id=$5`)

	values := []interface{}{
		input.Name,
		input.Nik,
		input.Phone_number,
		input.Daily_cost,
		id,
	}

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return err
}


//DELETE Driver
//=============================================================
func deleteDriver(tx *sqlx.DB, id int) error {
	query := (`delete from "driver"b
	where "driver_id" =$1`)

	values := []interface{}{
		id,
	}

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return err
}