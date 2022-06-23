package booking

import (
	"time"

	"github.com/jmoiron/sqlx"
)

//GET BOOKING
//=============================================================
func getBooking(tx *sqlx.DB) ([]GetBooking, error) {
	var (
		data  GetBooking
		datas []GetBooking
	)
	query := (`select booking_id,
				customer_id,
				cars_id,
				extract(epoch from "start_time")::bigint as "start_time",
				extract(epoch from "end_time")::bigint as "end_time",
				total_cost,
				finished,
				discount,
				booktype_id,
				driver_id,
				total_driver_cost  
				from "booking" b`)

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

//GET BOOKING BY ID
//=============================================================
func getBookingByID(tx *sqlx.DB, id int) (GetBooking, error) {
	var (
		data GetBooking
	)

	query := (`select booking_id,
				customer_id,
				cars_id,
				extract(epoch from "start_time")::bigint as "start_time",
				extract(epoch from "end_time")::bigint as "end_time",
				total_cost,
				finished,
				discount,
				booktype_id,
				driver_id,
				total_driver_cost 
				from "booking" b
	where b."booking_id" = $1`)

	values := []interface{}{
		id,
	}
	err := tx.QueryRowx(query, values...).StructScan(&data)
	if err != nil {
		return data, err
	}
	return data, err
}

//ADD BOOKING
//=============================================================
func addBooking(tx *sqlx.DB, input BookingForm) (int, error) {
	var startT = time.Unix(input.Start_time, 0).Format("2006-01-02")
	var endT = time.Unix(input.End_time, 0).Format("2006-01-02")

	query := (`insert into booking (customer_id,
			cars_id,
			start_time,
			end_time,
			total_cost,
			finished,
			discount,
			booktype_id,
			driver_id,
			total_driver_cost)
			Values ($1,	$2,	$3,	$4,	$5,	$6,	$7,	$8,	$9,	$10)
			returning "booking_id"`)
	values := []interface{}{
		input.Customer_ID,
		input.Cars_ID,
		startT,
		endT,
		input.Total_cost,
		input.Finished,
		input.Discount,
		input.Booktype_ID,
		input.Driver_ID,
		input.Total_driver_cost,
	}
	var id int
	err := tx.QueryRowx(query, values...).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//ADD BOOKING
//=============================================================
func addBookingV2(tx *sqlx.DB, input BookingFormV2,
	startT string,
	endT string,
	total_cost int,
	discount int,
	total_driver_cost int) (int, error) {

	query := (`insert into booking (customer_id,
			cars_id,
			start_time,
			end_time,
			total_cost,
			finished,
			discount,
			booktype_id,
			driver_id,
			total_driver_cost)
			Values ($1,	$2,	$3,	$4,	$5,	$6,	$7,	$8,	$9,	$10)
			returning "booking_id"`)
	values := []interface{}{
		input.Customer_ID ,
		input.Cars_ID ,
		startT ,
		endT ,
		total_cost ,
		input.Finished ,
		discount ,
		input.Booktype_ID ,
		input.Driver_ID ,
		total_driver_cost,
	}
	var id int
	err := tx.QueryRowx(query, values...).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//UPDATE BOOKING
//=============================================================
func updateBooking(tx *sqlx.DB, id int, input BookingForm) error {
	var startT = time.Unix(input.Start_time, 0).Format("2006-01-02")
	var endT = time.Unix(input.End_time, 0).Format("2006-01-02")
	query := (`update "booking"
		set "customer_id" = $1,
		"cars_id" = $2,         
		"start_time"=$3,
		"end_time" = $4,
		"total_cost" = $5,         
		"finished"=$6,
		"discount"=$7,
		"booktype_id" = $8,
		"driver_id" = $9,         
		"total_driver_cost"=$10
		where booking_id=$11`)

	values := []interface{}{
		input.Customer_ID,
		input.Cars_ID,
		startT,
		endT,
		input.Total_cost,
		input.Finished,
		input.Discount,
		input.Booktype_ID,
		input.Driver_ID,
		input.Total_driver_cost,
		id,
	}

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return err
}

//DELETE BOOKING
//=============================================================
func deleteBooking(tx *sqlx.DB, id int) error {
	query := (`delete from "booking"b
	where "booking_id" =$1`)

	values := []interface{}{
		id,
	}

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return err
}
