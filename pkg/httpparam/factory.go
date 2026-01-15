package httpparam

type IParamFactory interface {
	GetNewParam() *Param
}

type ParamFactory struct{}

func NewParamFactory() *ParamFactory {
	return &ParamFactory{}
}

func (p *ParamFactory) GetNewParam() *Param {
	return newParam()
}
