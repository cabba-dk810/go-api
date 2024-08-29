// packageとは、コードを構造化するための機能
// 別のファイルで定義した型や関数をこのファイルで使用できるようにするための機能
// 例）import model-> UserやUserResponseが別ファイルで使用可能
package model

// Goの標準パッケージ
// 正確な日付を表すtime.Time型を提供し、日付の操作が簡単に可能
// GORMでtime.Time型を自動的にデータベースの日時型にマッピングする
import "time"


// struct(構造体)はデータをひとまとめにしたもの
// Go には明示的なクラス定義がない
// クラスと似た役割を果たす

// テーブルと1対1になる構造体を定義する
// json, gormは構造体タグ
// json: Json変換時のフィールド名
// （例） json:"id"-> JSONでこのフィールドは "id" として表現される
// gorm: テーブルとのマッピングやカラムの振る舞いの制御
// （例） gorm:"primary_key"-> このフィールドはデータベースのプライマリーキー
// マッピングはデフォルトで実施される
type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


// APIのレスポンスとして使用する構造体も定義しておく
type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
}