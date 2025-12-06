package consts

import "errors"

// db
var ErrDBresultNotFound = errors.New("db result not found")

// token
var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenSignatureInvalid = errors.New("token signature invalid")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenMismatch = errors.New("token mismatch")
var ErrRefreshTokenAuthError = errors.New("refreshToken authentication failed")

// chat
var ErrPublishToMessageBrokerError = errors.New("publish to message broker failed")

// otp
var ErrFailedToDecodePEMBlock = errors.New("failed to decode PEM block")
var ErrFailedToParsePublicKey = errors.New("failed to parse public key")
var ErrFailedToEncryptOtpKey = errors.New("failed to encrypt OTP key")
var ErrOtpNotFound = errors.New("failed to otp not found")
var ErrInvalidOtpType = errors.New("failed to otp type")

// chat room
var ErrInvalidChatRoomMember = errors.New("invalid chat room member")
var ErrRoomKeyAlreadyExist = errors.New("room key already exist")
var ErrRoomTypeCheckError = errors.New("room type check error")
var ErrRoomSecretFlagCheckError = errors.New("room secret flag check error")
var ErrRoomSecretCheckError = errors.New("room secret check error")
