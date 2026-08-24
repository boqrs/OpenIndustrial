package code

import "github.com/boqrs/zeus/ginx"

var (
	// 通用

	InternalError  = ginx.NewGinError(110000, "Internal Error")
	NotImplemented = ginx.NewGinError(110001, "Not implemented")

	// user
	AccountAlreadyExists        = ginx.NewGinError(120001, "Account already exists")
	InvalidUsername             = ginx.NewGinError(120002, "Invalid username")
	CaptchaExpired              = ginx.NewGinError(120003, "Captcha has expired")
	InvalidCaptcha              = ginx.NewGinError(120004, "Invalid captcha")
	InvalidToken                = ginx.NewGinError(120005, "Invalid token")
	UserNotFound                = ginx.NewGinError(120006, "User not found")
	InvalidCredentialType       = ginx.NewGinError(120007, "Invalid credential type")
	IncorrectUsernameOrPassword = ginx.NewGinError(120008, "Incorrect username or password")
	EmailOrPhoneAlreadyExists   = ginx.NewGinError(120009, "Email or phone already exists")

	// perm
	TokenExpired            = ginx.NewGinError(130001, "Expired token")
	PermissionDenied        = ginx.NewGinError(130002, "No perm found")
	RoleNotFound            = ginx.NewGinError(130003, "Role not found")
	TokenVerificationFailed = ginx.NewGinError(130004, "Token verification failed")
	PermGetFailed           = ginx.NewGinError(130005, "Failed to get user permissions.")
)