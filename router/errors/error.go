package errors

import "net/http"

// General HTTP Error
type Error struct {
	Err    error  `json:"-"`
	Status int    `json:"-"`
	Msg    string `json:"error"`
}

func (e Error) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return e.Err.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Wrap(err error) Error {
	e.Err = err
	return e
}

func (e Error) String() string {
	return e.Msg
}

// BadRequest error with message s
func BadRequest(s string) Error {
	return ErrBadRequest.Message(s)
}

// Unauthorized error with message s
func Unauthorized(s string) Error {
	return ErrUnauthorized.Message(s)
}

// PaymentRequired error with message s
func PaymentRequired(s string) Error {
	return ErrPaymentRequired.Message(s)
}

// Forbidden error with message s
func Forbidden(s string) Error {
	return ErrForbidden.Message(s)
}

// NotFound error with message s
func NotFound(s string) Error {
	return ErrNotFound.Message(s)
}

// MethodNotAllowed error with message s
func MethodNotAllowed(s string) Error {
	return ErrMethodNotAllowed.Message(s)
}

// NotAcceptable error with message s
func NotAcceptable(s string) Error {
	return ErrNotAcceptable.Message(s)
}

// ProxyAuthRequired error with message s
func ProxyAuthRequired(s string) Error {
	return ErrProxyAuthRequired.Message(s)
}

// RequestTimeout error with message s
func RequestTimeout(s string) Error {
	return ErrRequestTimeout.Message(s)
}

// Conflict error with message s
func Conflict(s string) Error {
	return ErrConflict.Message(s)
}

// Gone error with message s
func Gone(s string) Error {
	return ErrGone.Message(s)
}

// LengthRequired error with message s
func LengthRequired(s string) Error {
	return ErrLengthRequired.Message(s)
}

// PreconditionFailed error with message s
func PreconditionFailed(s string) Error {
	return ErrPreconditionFailed.Message(s)
}

// RequestEntityTooLarge error with message s
func RequestEntityTooLarge(s string) Error {
	return ErrRequestEntityTooLarge.Message(s)
}

// RequestURITooLong error with message s
func RequestURITooLong(s string) Error {
	return ErrRequestURITooLong.Message(s)
}

// UnsupportedMediaType error with message s
func UnsupportedMediaType(s string) Error {
	return ErrUnsupportedMediaType.Message(s)
}

// RequestedRangeNotSatisfiable error with message s
func RequestedRangeNotSatisfiable(s string) Error {
	return ErrRequestedRangeNotSatisfiable.Message(s)
}

// ExpectationFailed error with message s
func ExpectationFailed(s string) Error {
	return ErrExpectationFailed.Message(s)
}

// Teapot error with message s
func Teapot(s string) Error {
	return ErrTeapot.Message(s)
}

// MisdirectedRequest error with message s
func MisdirectedRequest(s string) Error {
	return ErrMisdirectedRequest.Message(s)
}

// UnprocessableEntity error with message s
func UnprocessableEntity(s string) Error {
	return ErrUnprocessableEntity.Message(s)
}

// Locked error with message s
func Locked(s string) Error {
	return ErrLocked.Message(s)
}

// FailedDependency error with message s
func FailedDependency(s string) Error {
	return ErrFailedDependency.Message(s)
}

// TooEarly error with message s
func TooEarly(s string) Error {
	return ErrTooEarly.Message(s)
}

// UpgradeRequired error with message s
func UpgradeRequired(s string) Error {
	return ErrUpgradeRequired.Message(s)
}

// PreconditionRequired error with message s
func PreconditionRequired(s string) Error {
	return ErrPreconditionRequired.Message(s)
}

// TooManyRequests error with message s
func TooManyRequests(s string) Error {
	return ErrTooManyRequests.Message(s)
}

// RequestHeaderFieldsTooLarge error with message s
func RequestHeaderFieldsTooLarge(s string) Error {
	return ErrRequestHeaderFieldsTooLarge.Message(s)
}

// UnavailableForLegalReasons error with message s
func UnavailableForLegalReasons(s string) Error {
	return ErrUnavailableForLegalReasons.Message(s)
}

// InternalServererror error with message s
func InternalServerError(s string) Error {
	return ErrInternalServerError.Message(s)
}

// NotImplemented error with message s
func NotImplemented(s string) Error {
	return ErrNotImplemented.Message(s)
}

// BadGateway error with message s
func BadGateway(s string) Error {
	return ErrBadGateway.Message(s)
}

// ServiceUnavailable error with message s
func ServiceUnavailable(s string) Error {
	return ErrServiceUnavailable.Message(s)
}

// GatewayTimeout error with message s
func GatewayTimeout(s string) Error {
	return ErrGatewayTimeout.Message(s)
}

// HTTPVersionNotSupported error with message s
func HTTPVersionNotSupported(s string) Error {
	return ErrHTTPVersionNotSupported.Message(s)
}

// VariantAlsoNegotiates error with message s
func VariantAlsoNegotiates(s string) Error {
	return ErrVariantAlsoNegotiates.Message(s)
}

// InsufficientStorage error with message s
func InsufficientStorage(s string) Error {
	return ErrInsufficientStorage.Message(s)
}

// LoopDetected error with message s
func LoopDetected(s string) Error {
	return ErrLoopDetected.Message(s)
}

// NotExtended error with message s
func NotExtended(s string) Error {
	return ErrNotExtended.Message(s)
}

// NetworkAuthenticationRequired error with message s
func NetworkAuthenticationRequired(s string) Error {
	return ErrNetworkAuthenticationRequired.Message(s)
}

func New(s int) Error {
	return Error{
		Status: s,
		Msg:    http.StatusText(s),
	}
}

func (e Error) Message(s string) Error {
	e.Msg = s
	return e
}

// Errors
var (
	ErrBadRequest                   = New(http.StatusBadRequest)                   // 400
	ErrUnauthorized                 = New(http.StatusUnauthorized)                 // 401
	ErrPaymentRequired              = New(http.StatusPaymentRequired)              // 402
	ErrForbidden                    = New(http.StatusForbidden)                    // 403
	ErrNotFound                     = New(http.StatusNotFound)                     // 404
	ErrMethodNotAllowed             = New(http.StatusMethodNotAllowed)             // 405
	ErrNotAcceptable                = New(http.StatusNotAcceptable)                // 406
	ErrProxyAuthRequired            = New(http.StatusProxyAuthRequired)            // 407
	ErrRequestTimeout               = New(http.StatusRequestTimeout)               // 408
	ErrConflict                     = New(http.StatusConflict)                     // 409
	ErrGone                         = New(http.StatusGone)                         // 410
	ErrLengthRequired               = New(http.StatusLengthRequired)               // 411
	ErrPreconditionFailed           = New(http.StatusPreconditionFailed)           // 412
	ErrRequestEntityTooLarge        = New(http.StatusRequestEntityTooLarge)        // 413
	ErrRequestURITooLong            = New(http.StatusRequestURITooLong)            // 414
	ErrUnsupportedMediaType         = New(http.StatusUnsupportedMediaType)         // 415
	ErrRequestedRangeNotSatisfiable = New(http.StatusRequestedRangeNotSatisfiable) // 416
	ErrExpectationFailed            = New(http.StatusExpectationFailed)            // 417
	ErrTeapot                       = New(http.StatusTeapot)                       // 418
	ErrMisdirectedRequest           = New(http.StatusMisdirectedRequest)           // 421
	ErrUnprocessableEntity          = New(http.StatusUnprocessableEntity)          // 422
	ErrLocked                       = New(http.StatusLocked)                       // 423
	ErrFailedDependency             = New(http.StatusFailedDependency)             // 424
	ErrTooEarly                     = New(http.StatusTooEarly)                     // 425
	ErrUpgradeRequired              = New(http.StatusUpgradeRequired)              // 426
	ErrPreconditionRequired         = New(http.StatusPreconditionRequired)         // 428
	ErrTooManyRequests              = New(http.StatusTooManyRequests)              // 429
	ErrRequestHeaderFieldsTooLarge  = New(http.StatusRequestHeaderFieldsTooLarge)  // 431
	ErrUnavailableForLegalReasons   = New(http.StatusUnavailableForLegalReasons)   // 451

	ErrInternalServerError           = New(http.StatusInternalServerError)           // 500
	ErrNotImplemented                = New(http.StatusNotImplemented)                // 501
	ErrBadGateway                    = New(http.StatusBadGateway)                    // 502
	ErrServiceUnavailable            = New(http.StatusServiceUnavailable)            // 503
	ErrGatewayTimeout                = New(http.StatusGatewayTimeout)                // 504
	ErrHTTPVersionNotSupported       = New(http.StatusHTTPVersionNotSupported)       // 505
	ErrVariantAlsoNegotiates         = New(http.StatusVariantAlsoNegotiates)         // 506
	ErrInsufficientStorage           = New(http.StatusInsufficientStorage)           // 507
	ErrLoopDetected                  = New(http.StatusLoopDetected)                  // 508
	ErrNotExtended                   = New(http.StatusNotExtended)                   // 510
	ErrNetworkAuthenticationRequired = New(http.StatusNetworkAuthenticationRequired) // 511
)
