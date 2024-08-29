package db

// 外部のパッケージをインポートした場合は、次のコマンドでインストールする
// # go mod tidy
//
// 外部パッケージ：　"github.com/joho/godotenv", "gorm.io/gorm", "gorm.io/driver/postgres"
//
// - os-> ファイル操作、環境変数の取得等
// - log-> ログの出力先変更、ログ出力後の処理変更
// - godotenv-> .envファイルの読み込み
// - gorm-> GORMのデータベース接続
// - gorm.io/driver/postgres-> GORMとPostgreSQLの接続
import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// *gorm.DB型（GORMのデータベース接続へのポインタ）を返す関数
// （ポインタ型の説明は省略）
// ->プログラムの起動時に呼び出され、データベース接続を確立する
func NewDB() *gorm.DB {
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err) // ログの出力後、プログラムを終了する
		}
	}
	fmt.Println(os.Getenv("POSTGRES_USER"))

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	// gorm.Configは設定の構造体のため、ポインタを指定
	// ポインタを指定することで、コピーのオーバーヘッドを避け、元の構造体を変更することができる
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("DB接続成功")
	return db
}

// データベース接続を閉じる関数
// - アプリケーションの終了時
// - データベース接続の再設定が必要な場合
// - リソースの即時解放が必要な場合
func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	sqlDB.Close()
}
