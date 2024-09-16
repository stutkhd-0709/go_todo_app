package auth

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stutkhd-0709/go_todo_app/clock"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"net/http"
	"time"
)

// go embedを使うことで、指定したファイルを埋め込んでバイナリにしてくれる
// この機能のおかげでgoの特徴であるシングルバイナリを実現することができる

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}
	privkey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	j.PrivateKey = privkey
	j.PublicKey = pubkey
	j.Clocker = c
	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	// jwk -> json web key
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, fmt.Errorf("failed to jwk.ParseKey: %w", err)
	}
	return key, nil
}

const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("github.com/stutkhd-0709/go_todo_app").
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, u.Role).
		Claim(UserNameKey, u.Name).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build token: %w", err)
	}
	if err := j.Store.Save(ctx, tok.JwtID(), u.ID); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}

func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	// auth.JWTer.Clockerフィールドをベースに検証を行うため、WithValidateをfalseにしてるらしい
	// ここで行うのはリクエストからのトークン取得のみ
	token, err := jwt.ParseRequest(r, jwt.WithKey(jwa.RS256, j.PublicKey), jwt.WithValidate(false))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	// トークン検証をこちらで独自におこなう
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	// Redisから削除して手動でexpireさせていることもアリエル
	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}
	return token, nil
}

// FillContext はJWXから取得したデータをcontext.Context型の値に詰める
func (j *JWTer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := j.GetToken(r.Context(), r)
	if err != nil {
		return nil, err
	}
	uid, err := j.Store.Load(r.Context(), token.JwtID())
	if err != nil {
		return nil, err
	}
	ctx := SetUserID(r.Context(), uid)
	ctx = SetRole(ctx, token)
	/*
		http.Requestはimmutableに設計されている。そのためcloneする必要がある
		1. **不変性の維持**:
		   - http.Requestオブジェクトは不変であるべきです。オリジナルのリクエストオブジェクトを変更する代わりに、新しいコンテキストを持つクローンを作成します。
		2. **並行処理の安全性**:
		   - 同じリクエストオブジェクトが複数のゴルーチンで使用される可能性があります。オリジナルのリクエストを変更すると、他のゴルーチンに影響を与える可能性があります。クローンを作成することで、この問題を回避します。
		3. **コードの明確性**:
		   - 新しいコンテキストを持つリクエストオブジェクトを明示的に作成することで、コードの意図が明確になります。これにより、後でコードを読む人が何をしているのかを理解しやすくなります。
	*/
	clone := r.Clone(ctx)
	return clone, nil
}

type userIDKey struct{}
type roleKey struct{}

// キーにインタフェースを入れるのはgo公式が推奨してる
// これは同じキーをセットして衝突を防ぐためである
// ref. https://zenn.dev/hsaki/books/golang-context/viewer/appliedvalue
func SetUserID(ctx context.Context, uid entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func GetUserID(ctx context.Context) (entity.UserID, bool) {
	id, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return id, ok
}

func SetRole(ctx context.Context, tok jwt.Token) context.Context {
	get, ok := tok.Get(RoleKey)
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}
	return context.WithValue(ctx, roleKey{}, get)
}

func GetRole(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(roleKey{}).(string)
	return id, ok
}

func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	if !ok {
		return false
	}
	return role == "admin"
}
