package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"ticketapp/internal/repositories"
	"ticketapp/internal/services"
	"ticketapp/internal/utils"
)

type AuthHandler struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.RefreshTokenRepository
	jwt       *services.JWTService
	otp       *services.OTPService
}

func NewAuthHandler(
	userRepo repositories.UserRepository,
	tokenRepo repositories.RefreshTokenRepository,
	jwt *services.JWTService,
	otp *services.OTPService,
) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwt:       jwt,
		otp:       otp,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/* =========================
   LOGIN
========================= */

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil || utils.ComparePassword(user.PasswordHash, req.Password) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	h.issueTokens(w, user.ID.String(), user.Role)
}

/* =========================
   REFRESH (FIXED)
========================= */

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	hash := services.HashToken(c.Value)

	token, err := h.tokenRepo.GetValid(hash)
	if err != nil {

		switch err {

		// 🚨 REAL ATTACK
		case repositories.ErrTokenReuse:
			_ = h.tokenRepo.RevokeAll(token.UserID)
			http.Error(w, "token reuse detected", http.StatusUnauthorized)
			return

		// ❌ Normal auth failure
		case repositories.ErrTokenExpired,
			repositories.ErrTokenRevoked,
			repositories.ErrTokenNotFound:
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return

		default:
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}

	// 🔄 ROTATE refresh token (ONLY correct way)
	newRefresh := uuid.NewString()
	newHash := services.HashToken(newRefresh)

	if err := h.tokenRepo.Rotate(
		token,
		newHash,
		time.Now().Add(7*24*time.Hour),
	); err != nil {
		http.Error(w, "token rotation failed", http.StatusInternalServerError)
		return
	}

	user, err := h.userRepo.GetByID(token.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.jwt.GenerateAccessToken(
		user.ID.String(),
		user.Role,
	)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	// ✅ Set NEW refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		HttpOnly: true,
		Secure:   false, // true in prod HTTPS
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
		"role":         user.Role,
	})
}

/* =========================
   VERIFY OTP (UNCHANGED)
========================= */

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
		Code   string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	secret, err := h.userRepo.GetOTPSecret(uid)
	if err != nil || !h.otp.Verify(secret, req.Code) {
		http.Error(w, "invalid otp", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(uid)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	h.issueTokens(w, user.ID.String(), user.Role)
}

/* =========================
   FORGOT / RESET (UNCHANGED)
========================= */

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email string }
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		return
	}

	token := uuid.NewString()
	hash := services.HashToken(token)

	h.userRepo.StoreResetToken(user.ID, hash, time.Now().Add(15*time.Minute))
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string
		Password string
	}
	json.NewDecoder(r.Body).Decode(&req)

	hash := services.HashToken(req.Token)
	userID, err := h.userRepo.ValidateResetToken(hash)
	if err != nil {
		http.Error(w, "invalid or expired", http.StatusBadRequest)
		return
	}

	pw, _ := utils.HashPassword(req.Password)
	h.userRepo.UpdatePassword(userID, pw)
	h.tokenRepo.RevokeAll(userID)

	w.WriteHeader(http.StatusOK)
}
