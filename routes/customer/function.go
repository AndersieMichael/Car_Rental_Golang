package customer

import (
	"github.com/jmoiron/sqlx"
)

//GET CUSTOMER
//
func getCustomer(tx *sqlx.DB) ([]GetCustomer, error) {
	var (
		data  GetCustomer
		datas []GetCustomer
	)

	query := (`select c.customer_id ,c."name" ,c.nik ,c.phone_number 
			from "customer" c `)
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

//GET CUSTOMER BY ID
//
func getCustomerByID(tx *sqlx.DB, id int) (GetCustomer, error) {
	var (
		data  GetCustomer
		// datas []GetCustomer
	)

	query := (`select c.customer_id ,c."name" ,c.nik ,c.phone_number 
			from "customer" c 
			where c."customer_id" = $1`)

	values := []interface{}{
		id,
	}
	err := tx.QueryRowx(query, values...).StructScan(&data)
	if err != nil {
		return data, err
	}
	return data, err
}

//ADD CUSTOMER
//
func addCustomer(tx *sqlx.DB, input CustomerForm) (int, error) {
	query := `insert into customer (name,nik,phone_number,membership_id,password)
				values ($1,$2,$3,$4,$5)
				returning "customer_id"`

	Values := []interface{}{
		input.Name,
		input.Nik,
		input.Phone_number,
		input.Membership_ID,
		input.Password,
	}

	var id int
	err := tx.QueryRowx(query,
		Values...).Scan(&id)

	if err != nil {
		return id, err
	}
	return id, err
}

//UPDATE CUSTOMER
//
func updateCustomer(tx *sqlx.DB, id int, input CustomerUpdateForm) (error) {
	query := `update customer
				set "name" = $1,
				"nik" = $2,         
				"phone_number"=$3,
				"membership_id"=$4
				where customer_id=$5`

	Values := []interface{}{
		input.Name,
		input.Nik,
		input.Phone_number,
		input.Membership_ID,
		id,
	}

	_, err := tx.Exec(query,
		Values...)

	if err != nil {
		return  err
	}
	return err
}


//DELETE CUSTOMER
//
func deleteCustomer(tx *sqlx.DB, id int) (error) {
	query := `delete from customer 
	where customer_id =$1`

	Values := []interface{}{
		id,
	}

	_, err := tx.Exec(query,
		Values...)

	if err != nil {
		return  err
	}
	return err
}
