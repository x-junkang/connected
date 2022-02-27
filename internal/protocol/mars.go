package protocol

type MarsHeader struct {
	HeaderLength  uint32
	ClientVersion uint32
	Cmd           uint32
	Sequence      uint32
	BodyLength    uint32
}

const MarsHeaderLength = 4 * 5
