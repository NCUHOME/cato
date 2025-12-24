package param

type Param interface {
	GetCode() string
	GetMessage() string
	SetCode(code string)
	SetMessage(message string)
	GetBody() any
	ResetBody()
}
