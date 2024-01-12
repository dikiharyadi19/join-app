package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/yogabagas/join-app/shared/constant"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrBadRequest          = errors.New("Bad request")
	ErrForbiddenResource   = errors.New("Forbidden resource")
	ErrNotFound            = errors.New("Not Found")
	ErrPreConditionFailed  = errors.New("Precondition failed")
	ErrInternalServerError = errors.New("Internal server error")
	ErrTimeoutError        = errors.New("Timeout error")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrConflict            = errors.New("Conflict")
	ErrMethodNotAllowed    = errors.New("Method not allowed")
)

const (
	StatusCodeGenericSuccess            = "200000"
	StatusCodeAccepted                  = "202000"
	StatusCodeBadRequest                = "400000"
	StatusCodeAlreadyRegistered         = "400001"
	StatusCodeUnauthorized              = "401000"
	StatusCodeForbidden                 = "403000"
	StatusCodeNotFound                  = "404000"
	StatusCodeConflict                  = "409000"
	StatusCodeGenericPreconditionFailed = "412000"
	StatusCodeOTPLimitReached           = "412550"
	StatusCodeNoLinkerExist             = "412553"
	StatusCodeInternalError             = "500000"
	StatusCodeFailedSellBatch           = "500100"
	StatusCodeFailedOTP                 = "503000"
	StatusCodeServiceUnavailable        = "503000"
	StatusCodeTimeoutError              = "504000"
	StatusCodeMethodNotAllowed          = "405000"
)

func GetErrorCode(err error) string {
	err = getErrType(err)

	switch err {
	case ErrBadRequest:
		return StatusCodeBadRequest
	case ErrForbiddenResource:
		return StatusCodeForbidden
	case ErrNotFound:
		return StatusCodeNotFound
	case ErrConflict:
		return StatusCodeConflict
	case ErrUnauthorized:
		return StatusCodeUnauthorized
	case ErrPreConditionFailed:
		return StatusCodeGenericPreconditionFailed
	case ErrInternalServerError:
		return StatusCodeInternalError
	case ErrTimeoutError:
		return StatusCodeTimeoutError
	case ErrMethodNotAllowed:
		return StatusCodeMethodNotAllowed
	case nil:
		return StatusCodeGenericSuccess
	default:
		return StatusCodeInternalError
	}
}

func GetHTTPStatus(code int) string {
	switch code {
	case http.StatusOK:
		return "success"
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusInternalServerError:
		return "internal server error"
	case http.StatusUnauthorized:
		return "unauthorized"
	default:
		return "undefined"
	}
}

func GetHTTPCode(code string) int {
	s := code[0:3]
	i, _ := strconv.Atoi(s)
	return i
}

type JSONResponse struct {
	Data        interface{}            `json:"data,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Code        string                 `json:"code"`
	StatusCode  int                    `json:"status_code"`
	Status      string                 `json:"status"`
	ErrorString string                 `json:"error,omitempty"`
	Error       error                  `json:"-"`
	RealError   string                 `json:"-"`
	Latency     string                 `json:"latency,omitempty"`
	Log         map[string]interface{} `json:"-"`
	HTMLPage    bool                   `json:"-"`
	Result      interface{}            `json:"result,omitempty"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{Code: StatusCodeGenericSuccess, StatusCode: GetHTTPCode(StatusCodeGenericSuccess), Status: GetHTTPStatus(http.StatusOK), Log: map[string]interface{}{}}
}

func (r *JSONResponse) SetData(data interface{}) *JSONResponse {
	r.Data = data
	return r
}

func (r *JSONResponse) SetStatus(status string) *JSONResponse {
	r.Status = status
	return r
}

func (r *JSONResponse) SetCode(code string) *JSONResponse {
	r.Code = code
	return r
}

func (r *JSONResponse) SetStatusCode(statusCode int) *JSONResponse {
	r.StatusCode = statusCode
	return r
}

func (r *JSONResponse) SetHTML() *JSONResponse {
	r.HTMLPage = true
	return r
}

func (r *JSONResponse) SetResult(result interface{}) *JSONResponse {
	r.Result = result
	return r
}

func (r *JSONResponse) SetMessage(msg string) *JSONResponse {
	r.Message = msg
	return r
}

func (r *JSONResponse) SetLatency(latency float64) *JSONResponse {
	r.Latency = fmt.Sprintf("%.2f ms", latency)
	return r
}

// func (r *JSONResponse) SetLog(key string, val interface{}) *JSONResponse {
// 	_, file, no, _ := runtime.Caller(1)
// 	log.Errorln(log.Fields{
// 		"code":            r.Code,
// 		"err":             val,
// 		"function_caller": fmt.Sprintf("file %v line no %v", file, no),
// 	}).Errorln("Error API")
// 	r.Log[key] = val
// 	return r
// }

func getErrType(err error) error {
	switch err.(type) {
	case ErrChain:
		errType := err.(ErrChain).Type
		if errType != nil {
			err = errType
		}
	}
	return err
}

func (r *JSONResponse) SetError(err error, a ...string) *JSONResponse {
	r.Code = GetErrorCode(err)
	// r.SetLog("error", err)
	r.RealError = fmt.Sprintf("%+v", err)
	err = getErrType(err)
	r.Error = err
	r.ErrorString = err.Error()
	r.StatusCode = GetHTTPCode(r.Code)
	r.Status = GetHTTPStatus(r.StatusCode)

	if r.StatusCode == http.StatusInternalServerError {
		r.ErrorString = "Internal Server error"
	}
	if len(a) > 0 {
		r.ErrorString = a[0]
	}
	return r
}

func (r *JSONResponse) GetBody() []byte {
	b, _ := json.Marshal(r)
	return b
}

func (r *JSONResponse) Send(w http.ResponseWriter) {
	if r.HTMLPage {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(r.StatusCode)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(r.StatusCode)
		err := json.NewEncoder(w).Encode(r)
		if err != nil {
			log.Println("err", err.Error())
		}
	}
}

// APIStatusSuccess for standard request api status success
func (r *JSONResponse) APIStatusSuccess() *JSONResponse {
	r.Code = constant.StatusCode(constant.StatusSuccess)
	r.Message = constant.StatusText(constant.StatusSuccess)
	return r
}

// APIStatusCreated
func (r *JSONResponse) APIStatusCreated() *JSONResponse {
	r.StatusCode = constant.StatusCreated
	r.Code = constant.StatusCode(constant.StatusCreated)
	r.Message = constant.StatusText(constant.StatusCreated)
	return r
}

// APIStatusAccepted
func (r *JSONResponse) APIStatusAccepted() *JSONResponse {
	r.StatusCode = constant.StatusAccepted
	r.Code = constant.StatusCode(constant.StatusAccepted)
	r.Message = constant.StatusText(constant.StatusAccepted)
	return r
}

// APIStatusNoContent
func (r *JSONResponse) APIStatusNoContent() *JSONResponse {
	r.StatusCode = constant.StatusNoContent
	r.Code = constant.StatusCode(constant.StatusNoContent)
	r.Message = constant.StatusText(constant.StatusNoContent)
	return r
}

// APIStatusErrorUnknown
func (r *JSONResponse) APIStatusErrorUnknown() *JSONResponse {
	r.StatusCode = constant.StatusErrorUnknown
	r.Code = constant.StatusCode(constant.StatusErrorUnknown)
	r.Message = constant.StatusText(constant.StatusErrorUnknown)
	return r
}

// APIStatusInvalidAuthentication
func (r *JSONResponse) APIStatusInvalidAuthentication() *JSONResponse {
	r.StatusCode = constant.StatusInvalidAuthentication
	r.Code = constant.StatusCode(constant.StatusInvalidAuthentication)
	r.Message = constant.StatusText(constant.StatusInvalidAuthentication)
	return r
}

// APIStatusUnauthorized
func (r *JSONResponse) APIStatusUnauthorized() *JSONResponse {
	r.StatusCode = constant.StatusUnauthorized
	r.Code = constant.StatusCode(constant.StatusUnauthorized)
	r.Message = constant.StatusText(constant.StatusUnauthorized)
	return r
}

// APIStatusForbidden
func (r *JSONResponse) APIStatusForbidden() *JSONResponse {
	r.StatusCode = constant.StatusForbidden
	r.Code = constant.StatusCode(constant.StatusForbidden)
	r.Message = constant.StatusText(constant.StatusForbidden)
	return r
}

// APIStatusBadRequest
func (r *JSONResponse) APIStatusBadRequest() *JSONResponse {
	r.StatusCode = constant.StatusErrorForm
	r.Code = constant.StatusCode(constant.StatusErrorForm)
	r.Message = constant.StatusText(constant.StatusErrorForm)
	return r
}

// APIStatusNotFound
func (r *JSONResponse) APIStatusNotFound() *JSONResponse {
	r.StatusCode = constant.StatusNotFound
	r.Code = constant.StatusCode(constant.StatusNotFound)
	r.Message = constant.StatusText(constant.StatusNotFound)
	return r
}
