package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/RSUD-Daha-Husada/polda-be/helpers"
	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"github.com/RSUD-Daha-Husada/polda-be/internal/repository"
	"github.com/RSUD-Daha-Husada/polda-be/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	DB       *gorm.DB
	UserRepo *repository.UserRepository
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		DB:       db,
		UserRepo: repository.NewUserRepository(db),
	}
}

func (s *AuthService) Login(username, password, ipAddress, userAgent string) (string, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		_ = helpers.SaveLoginLog(s.DB, nil, "login", userAgent, ipAddress, "failed", "username not found")
		return "", errors.New("username tidak ditemukan!")
	}

	if user.Password != nil && !utils.CheckPasswordHash(password, *user.Password) {
		_ = helpers.SaveLoginLog(s.DB, &user.UserID, "login", userAgent, ipAddress, "failed", "invalid password")
		return "", errors.New("password anda salah")
	}

	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		_ = helpers.SaveLoginLog(s.DB, &user.UserID, "login", userAgent, ipAddress, "failed", "failed to generate token")
		return "", err
	}

	expiredAt := time.Now().Add(24 * time.Hour) 
	if err := utils.SaveAccessToken(s.DB, user.UserID, token, userAgent, ipAddress, expiredAt); err != nil {
		fmt.Println("Failed to save access token:", err)
	}

	_ = helpers.SaveLoginLog(s.DB, &user.UserID, "login", userAgent, ipAddress, "success", "login successful")

	return token, nil
}


func (s *AuthService) GenerateLoginCode(telephone string) error {
	user, err := s.UserRepo.FindByTelephone(telephone)
	if err != nil {
		return errors.New("user not found")
	}

	now := time.Now()

	s.DB.Model(&model.CodeLoginByWA{}).
		Where("user_id = ? AND used = false AND status = true", user.UserID).
		Update("status", false)

	code := utils.RandomCode(4)

	loginCode := model.CodeLoginByWA{
		CodeLoginByWAID: uuid.New(), 
		UserID:          user.UserID,
		Code:            code,
		ValidUntil:      now.Add(5 * time.Minute),
		Used:            false,
		Status:          true,
		CreatedAt:       now,
	}

	if err := s.DB.Create(&loginCode).Error; err != nil {
		return err
	}

	go func() {
		err := sendWhatsApp(user.Telephone, code)
		if err != nil {
			_ = helpers.SaveLoginLog(s.DB, &user.UserID, "send_whatsapp", "", "", "failed", err.Error())
		} else {
			_ = helpers.SaveLoginLog(s.DB, &user.UserID, "send_whatsapp", "", "", "success", "WA code sent successfully")
		}
	}()

	return nil
}

func (s *AuthService) LoginWithCode(username, code, userAgent, ipAddress string) (string, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		_ = helpers.SaveLoginLog(s.DB, nil, "login_with_code", userAgent, ipAddress, "failed", "user not found")
		return "", errors.New("user not found")
	}

	var record model.CodeLoginByWA
	if err := s.DB.
		Where("user_id = ? AND code = ? AND used = false AND status = true", user.UserID, code).
		First(&record).Error; err != nil {
		_ = helpers.SaveLoginLog(s.DB, &user.UserID, "login_with_code", userAgent, ipAddress, "failed", "code invalid or already used")
		return "", errors.New("code is invalid or already used")
	}

	if time.Now().After(record.ValidUntil) {
		s.DB.Model(&record).Update("status", false)
		_ = helpers.SaveLoginLog(s.DB, &user.UserID, "login_with_code", userAgent, ipAddress, "failed", "code has expired")
		return "", errors.New("code has expired")
	}

	s.DB.Model(&record).Updates(map[string]interface{}{
		"used":   true,
		"status": false,
	})

	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		_ = helpers.SaveLoginLog(s.DB, &user.UserID, "login_with_code", userAgent, ipAddress, "failed", "failed to generate token")
		return "", err
	}

	_ = helpers.SaveLoginLog(s.DB, &user.UserID, "login_with_code", userAgent, ipAddress, "success", "login successful")

	return token, nil
}

func sendWhatsApp(receiver string, code string) error {
	client := &http.Client{Timeout: 5 * time.Second}

	endpoint := "http://192.168.133.20:8101/api/send-message"
	params := url.Values{}
	params.Set("apikey", "warm")
	params.Set("receiver", receiver)
	params.Set("mtype", "text")
	params.Set("text", fmt.Sprintf("Kode login Anda: %s. Berlaku 5 menit.", code))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	resp, err := client.Get(fullURL)
	if err != nil {
		return fmt.Errorf("gagal kirim WA: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("wa gateway gagal: %s", resp.Status)
	}

	return nil
}
