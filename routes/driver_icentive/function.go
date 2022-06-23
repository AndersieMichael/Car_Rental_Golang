package driver_incentive

import "github.com/jmoiron/sqlx"

//GET Incentive
//=============================================================
func getIncentive(tx *sqlx.DB) ([]GetIncentive, error) {
	var (
		data  GetIncentive
		datas []GetIncentive
	)
	query := (`select *
				from "driver_incentive" d`)

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

//GET Incentive BY ID
//=============================================================
func getIncentiveByID(tx *sqlx.DB, id int) (GetIncentive, error) {
	var (
		data GetIncentive
	)

	query := (`select *
	from "driver_incentive" d
	where d."driver_incentive_id" = $1`)

	values := []interface{}{
		id,
	}
	err := tx.QueryRowx(query, values...).StructScan(&data)
	if err != nil {
		return data, err
	}
	return data, err
}

//ADD Incentive
//=============================================================
func AddIncentive(tx *sqlx.DB, input IncentiveForm) (int, error) {

	query := (`insert into "driver_incentive" (booking_id,
			incentive)
			Values ($1,	$2)
			returning "driver_incentive_id"`)
	values := []interface{}{
		input.Booking_ID,
		input.Incentive,
	}
	var id int
	err := tx.QueryRowx(query, values...).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//UPDATE Incentive
//=============================================================
func updateIncentive(tx *sqlx.DB, id int, input IncentiveForm) error {
	query := (`update "driver_incentive"
		set "booking_id" = $1,
		"incentive" = $2
		where driver_incentive_id=$3`)

	values := []interface{}{
		input.Booking_ID,
		input.Incentive,
		id,
	}

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return err
}

//DELETE Incentive
//=============================================================
func deleteIncentive(tx *sqlx.DB, id int) error {
	query := (`delete from "driver_incentive"
	where "driver_incentive_id" =$1`)

	values := []interface{}{
		id,
	}

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return err
}