package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// インターフェースでuserControllerが持つべき振る舞い（メソッド）を定義
// 引数：echo.Context
// 返り値：error
type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

// userController構造体を返しているが、IUserControllerで抽象化している
func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

// userControllerの構造体のSignUpメソッドを定義
// 引数はecho.Context
// メソッド例）
// c.Bind(&someStruct): リクエストボディを構造体にバインド
// c.Param("id"): URLパラメータの取得
// c.QueryParam("key"): クエリパラメータの取得
// c.FormValue("field"): フォームデータの取得
// c.JSON(statusCode, data): JSONレスポンスの送信
// c.String(statusCode, "message"): プレーンテキストレスポンスの送信
//
// 返り値 error
// エラーがあればerrorにデータを詰める
func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}

	// リクエストボディをuser構造体に変換
	if err := c.Bind(&user); err != nil {
		// 変換に失敗した場合は、４００とエラー内容をクライアントに返す
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Signupのビジネスロジックである、userUsecase.SignUpを実行
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		// エラーがあれば、５００とエラー内容をクライアントに返す
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// ユーザー登録が成功した場合は、２０１とユーザー情報をクライアントに返す
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	// リクエストボディをuser構造体に変換
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// LogInのビジネスロジックである、userUsecase.LogInを実行
	// 成功時は、トークン文字列を返す
	tokenString, err := uc.uu.LogIn(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// クッキーに関する処理
	// 新しいHTTPクッキーを作成
	cookie := new(http.Cookie)
	cookie.Name = "token"
	// クッキーの値にトークン文字列をセット
	cookie.Value = tokenString
	// クッキーの有効期限を24時間後に設定
	cookie.Expires = time.Now().Add(24 * time.Hour)
	// クッキーのパスを"/"に設定
	cookie.Path = "/"
	// クッキーのドメインをAPI_DOMAINに設定。このドメインに対してのみクッキーが有効
	cookie.Domain = os.Getenv("API_DOMAIN")
	// HTTPS接続でのみクッキーを送信
	cookie.Secure = true
	// jsでクッキーにアクセスできない設定
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	// 200と空のレスポンスをクライアントに返す
	return c.NoContent(http.StatusOK)
}

// クッキーを削除する
func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
