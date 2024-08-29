// プログラムのエントリーポイントとして使用する場合は、main関数が必要
// main関数を使用するため、パッケージ名をmainにする
package main

import (
	"fmt"
	"go-rest-api/db"
	"go-rest-api/model"
)

// AutoMigrate: モデル構造体のポインタを引数に取り、テーブルを更新する
func main() {
	dbConnect := db.NewDB()
	defer fmt.Println("Successfuly Migrated")
	defer db.Close(dbConnect)
	dbConnect.AutoMigrate(model.User{})
}
