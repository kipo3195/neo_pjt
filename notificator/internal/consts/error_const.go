package consts

import "errors"

// token
var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenSignatureInvalid = errors.New("token signature invalid")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenMismatch = errors.New("token mismatch")
var ErrRefreshTokenAuthError = errors.New("refreshToken authentication failed")

// websocket session
var ErrSenderChannelError = errors.New("sender channel error")

// chat room member regist
var ErrChatRoomMemberInvalid = errors.New("chat room member is invalid")

// nats
var ErrPublishToMessageBrokerError = errors.New("publish to message broker failed")
