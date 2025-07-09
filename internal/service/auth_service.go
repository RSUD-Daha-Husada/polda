package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"github.com/RSUD-Daha-Husada/polda-be/internal/repository"
	"github.com/RSUD-Daha-Husada/polda-be/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
	UserRepo  *repository.UserRepository
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		DB: 	  db,
		UserRepo: repository.NewUserRepository(db),
	}
}

func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.UserRepo.FindByUsername(username)
    if err != nil {
        return "", errors.New("username not found")
    }

    if user.Password != nil && !utils.CheckPasswordHash(password, *user.Password) {
        return "", errors.New("invalid password")
    }

    token, err := utils.GenerateJWT(user.UserID)
    if err != nil {
        return "", err
    }

    return token, nil
}

// ✅ Generate kode login dan simpan ke DB, lalu kirim via WhatsApp
func (s *AuthService) GenerateLoginCode(telephone string) error {
	// Cari user via repository
	user, err := s.UserRepo.FindByTelephone(telephone)
	if err != nil {
		return errors.New("user not found")
	}

	now := time.Now()

	// ❗️Nonaktifkan semua kode sebelumnya yang masih aktif & belum dipakai
	s.DB.Model(&model.CodeLoginByWA{}).
		Where("user_id = ? AND used = false AND status = true", user.UserID).
		Update("status", false)

	// Buat kode login baru
	code := utils.RandomCode(4)

	loginCode := model.CodeLoginByWA{
		UserID:     user.UserID,
		Code:       code,
		ValidUntil: now.Add(5 * time.Minute),
		Used:       false,
		Status:     true,
		CreatedAt:  now,
	}

	// Simpan kode baru
	if err := s.DB.Create(&loginCode).Error; err != nil {
		return err
	}

	// Kirim kode ke WA (async)
	go sendWhatsApp(user.Telephone, code)

	return nil
}

// ✅ Verifikasi kode OTP, update status & generate token
func (s *AuthService) LoginWithCode(telephone, code string) (string, error) {
	// Cari user via repository
	user, err := s.UserRepo.FindByTelephone(telephone)
	if err != nil {
		return "", errors.New("user not found")
	}

	var record model.CodeLoginByWA
	if err := s.DB.
		Where("user_id = ? AND code = ? AND used = false AND status = true", user.UserID, code).
		First(&record).Error; err != nil {
		return "", errors.New("code is invalid or already used")
	}

	if time.Now().After(record.ValidUntil) {
		s.DB.Model(&record).Updates(map[string]interface{}{
			"status": false,
		})
		return "", errors.New("code has expired")
	}

	s.DB.Model(&record).Updates(map[string]interface{}{
		"used":   true,
		"status": false,
	})

	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Kirim kode ke WhatsApp melalui WA 
func sendWhatsApp(receiver string, code string) {
	client := &http.Client{Timeout: 5 * time.Second}

	endpoint := "http://192.168.133.20:8101/api/send-message"
	params := url.Values{}
	params.Set("apikey", "wakomplain")
	params.Set("receiver", receiver)
	params.Set("mtype", "text")
	params.Set("text", fmt.Sprintf("Kode login Anda: %s. Berlaku 5 menit.", code))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	resp, err := client.Get(fullURL)
	if err != nil {
		fmt.Println("Gagal kirim WA:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("WA Gateway gagal:", resp.Status)
	}
}
