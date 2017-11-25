package encoder

//Encoder Encoder
type Encoder interface {

	//Encode Encode
	Encode(data []byte) []byte
	//Decode Decode
	Decode(data []byte) []byte
}
