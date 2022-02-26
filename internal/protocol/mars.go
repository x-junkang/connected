package protocol

type MarsHeader struct {
	HeaderLength  int32
	ClientVersion int32
	Cmd           int32
	Sequence      int32
	BodyLength    int32
}

const MarsHeaderLength = 4 * 5
