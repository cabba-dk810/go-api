package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"
)

// サーバーの起動に必要な依存関係を初期化している
func main() {
	// ⑥ データベースと接続
	// return *gorm.DB
	db := db.NewDB()

	// ⑤ リポジトリ、バリデーターの作成
	// IUserRepositoryとIUserValidatorを引数に、NewUserRepositoryを呼び出す
	// return *IUserRepository
	userRepository := repository.NewUserRepository(db)
	// IUserValidatorを引数に、NewUserValidatorを呼び出す
	// return *IUserValidator
	userValidator := validator.NewUserValidator()

	// ④ ユースケースの作成
	// IUserRepositoryとIUserValidatorを引数に、NewUserUsecaseを呼び出す
	// return *IUserUsecase
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)

	// ③ コントローラーの作成
	// IUserUsecaseを引数に、NewUserControllerを呼び出す
	// return *IUserController
	userController := controller.NewUserController(userUsecase)

	// ② ルーターの作成
	// ECHOの構造体を作成すると同時に、ルーティングを設定するため、routerパッケージを使用
	// return *echo.Echo
	e := router.NewRouter(userController)

	// ① サーバー起動
	// 　eは、*echo.Echo型でEcho構造体へのポインタ
	// HTTPサーバーの設定、ルーティング、ミドルウェアの管理が可能
	// （echoを使わない場合は、"net/http"等のパッケージを使ってHTTPサーバーを実装する）
	//
	// ex)
	// Start(): サーバーを起動するメソッド
	// Use(): ミドルウェアを設定するメソッド
	// GET(), POST(), PUT(), DELETE(): ルーティングを設定するメソッド
	// Logger: ログを出力するミドルウェア
	e.Logger.Fatal(e.Start(":8080"))
}
