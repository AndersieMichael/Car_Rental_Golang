package membership

import "github.com/jmoiron/sqlx"

//GET Membership
//=============================================================
func getMembership(tx *sqlx.DB) ([]GetMembership, error) {
	var (
		data  GetMembership
		datas []GetMembership
	)
	query := (`select *from "membership" m`)

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


//GET Membership BY ID
//=============================================================
func GetMembershipByID(tx *sqlx.DB, id int) (GetMembership, error) {
	var (
		data  GetMembership
	)
	query := (`select *from "membership" m
	where "membership_id"=$1 `)
	values:= []interface{}{
		id,
	}

	err := tx.QueryRowx(query,values...).StructScan(&data)
	if err != nil {
		return data, err
	}
	return data, err
}


//ADD MEMBERSHIP
//=============================================================
func addMembership(tx *sqlx.DB, input MembershipFormat) (int, error) {
	query := (`insert into membership (name,
		daily_discount)
		Values ($1,	$2)
		returning "membership_id"`)
	values:= []interface{}{
		input.Name,
		input.Daily_discount,
	}
	var id int
	err := tx.QueryRowx(query,values...).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//UPDATE MEMBERSHIP
//=============================================================
func updateMembership(tx *sqlx.DB, id int, input MembershipFormat) (error) {
	query := (`update "membership"
	set "name" = $1,
	"daily_discount"=$2
	where membership_id=$3`)
	values:= []interface{}{
		input.Name,
		input.Daily_discount,
		id,
	}
	_, err := tx.Exec(query,values...)
	if err != nil {
		return err
	}
	return err
}

//DELETE MEMBERSHIP
//=============================================================
func deleteMembership(tx *sqlx.DB, id int) (error) {
	query := (`delete from "membership"
	where "membership_id" =$1`)
	values:= []interface{}{
		id,
	}
	_, err := tx.Exec(query,values...)
	if err != nil {
		return err
	}
	return err
}