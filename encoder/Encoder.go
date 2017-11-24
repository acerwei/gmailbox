package encoder

//Encoder Encoder
type Encoder interface {

	//Encode Encode
	Encode(data []byte) []byte
	//Decode Decode
	Decode(data []byte) []byte
}

//SimpleEncoder simple masking
type SimpleEncoder struct {
	Mask byte
	Step byte
}

//NewDeaultSimpleEncoder NewDeaultSimpleEncoder
func NewDeaultSimpleEncoder() *SimpleEncoder {
	return &SimpleEncoder{
		Mask: 90,
		Step: 15,
	}
}

//Encode Encode
func (e *SimpleEncoder) Encode(data []byte) []byte {
	n := len(data)
	for i := 0; i < n; i++ {
		data[i] ^= (e.Mask + (byte(i) & e.Step))
	}
	return data
}

//Decode Decode
func (e *SimpleEncoder) Decode(data []byte) []byte {
	return e.Encode(data)
}
