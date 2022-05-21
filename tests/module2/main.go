package main

import (
	"github.com/pkg/errors"
	"database/sql"
)

func getValueFromDB() (int,error){
	return 0, errors.Wrap(sql.ErrNoRows, "no such row")
}

func main() {
	row,err :=getValueFromDB()
	if err!=nil {
		if sql.ErrNoRows== errors.Cause(err){
			// handle specifically
			println("not server error, logic to handle no row: ",err.Error())
		} else {
			// unknown error
			println("unknown error: ", err.Error())
		}
		return
	}
	println("row value: ",row)
}