package consts

import "errors"

// db
var ErrDBresultNotFound = errors.New("db result not found")
var ErrDB = errors.New("db error")

// token
var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenSignatureInvalid = errors.New("token signature invalid")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenMismatch = errors.New("token mismatch")
var ErrRefreshTokenAuthError = errors.New("refreshToken authentication failed")

// nats
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
var ErrRoomUpdateDateTypeError = errors.New("room update date type error")

// chat room title
var ErrChatRoomTypeMismatch = errors.New("chat room type mismatch")                        // 요청한 타입에 해당하는 방이 없음.
var ErrChatRoomKeyCheck = errors.New("chat room key check")                                // 내가 참여하지 않은 방이거나, 룸이 존재하지 않음.
var ErrChatRoomCreateMemberIsNotExist = errors.New("chat room create member is not exist") // 채팅방 생성시 참여자 정보가 누락됨.
