package main

import (
	"bytes"
	"encoding/json"
	"go-advanced/internal/auth"
	"go-advanced/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initBD() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "metlov.nm1@yandex.ru",
		Password: "$2a$10$/5CETCJz1UIZ2g0DcIszX.S77vv5G421gpp4A0qGwLHWm.kaxZYBy",
		Name:     "nikmet",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "metlov.nm1@yandex.ru").
		Delete(&user.User{})
}

func TestLoginSucces(t *testing.T) {
	// Prepare
	db := initBD()
	initData(db)
	// Test
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "metlov.nm1@yandex.ru",
		Password: "12345",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse

	err = json.Unmarshal(body, &resData)

	if err != nil {
		t.Fatal(err)
	}

	if resData.Token == "" {
		t.Fatal("Token is empty")
	}
	removeData(db)
}

func TestLoginFail(t *testing.T) {
	// Prepare
	db := initBD()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "metlov.nm1@yx.ru",
		Password: "$2a$10$/5CETCJz1UIZ2g0DcIszX.S77vv5G421gpp4A0qGwLHWm.kaxZYBy",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 400 {
		t.Fatalf("expected %d got %d", 400, res.StatusCode)
	}
	removeData(db)
}
