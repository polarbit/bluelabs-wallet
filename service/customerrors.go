package service

type DbError struct {
	Err error
}

type ServiceError struct {
	Msg string
}

var (
	ErrWalletNotFound                        = &ServiceError{Msg: "wallet not found"}
	ErrWalletAlreadyExists                   = &ServiceError{Msg: "wallet already exists"}
	ErrTransactionConsistency                = &ServiceError{Msg: "transaction failed due to consistency but retriable"}
	ErrTransactionAlreadyExistsByRefNo       = &ServiceError{Msg: "a transaction already exists with same refno"}
	ErrTransactionAlreadyExistsByFingerprint = &ServiceError{Msg: "a transaction already exists with same fingerprint"}
	ErrTransactionNotFound                   = &ServiceError{Msg: "transaction not found"}
)

func (e *ServiceError) Error() string {
	return e.Msg
}

func (e *DbError) Error() string {
	return e.Err.Error()
}

func NewDbError(err error) *DbError {
	return &DbError{Err: err}
}
