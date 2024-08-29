package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	LogIn(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		// バリデーションエラーがあれば、空のUserResponseとエラーをクライアントに返す
		return model.UserResponse{}, err
	}
	// パスワードをハッシュ化(10はコスト因子)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	// 新しいUserオブジェクトを作成
	newUser := model.User{Email: user.Email, Password: string(hash)}

	// ユーザーの保存。永続化はリポジトリで行う
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) LogIn(user model.User) (string, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}

	// 空のUser構造体を作成
	storedUser := model.User{}
	// User構造体を使って、リポジトリのGetUserByEmailメソッドを呼び出す
	// 取得されたユーザー情報をstoredUserに格納
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}

	//　パスワード検証
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}

	// jwt.NewWithClaimsを使って、新しいJWTトークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})

	//　トークンに署名をして、トークン文字列を生成
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
