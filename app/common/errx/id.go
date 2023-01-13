package errx

func Encode(errorType, sysId, serviceType, serviceId, logicId, oprId, errorId uint32) uint32 {
	// 错误码共32位(下列注释位数是从高位到低位数第几位)
	code := uint32(0)

	// 错误 id (共8位, 在第25-32位上)
	code += errorId

	// 操作 id (共6位, 在第19-24位上)
	code += oprId << 8

	// 业务 id (共6位, 在第13-18位上)
	code += logicId << 14

	// 服务 id (共5位, 在第8-12位上)
	code += serviceId << 20

	// 服务类型 (共1位, 在第7位上)
	code += serviceType << 25

	// 系统id (共5位, 在第2-6位上)
	code += sysId << 26

	// 错误类型 (共1位, 在第1位上)
	code += errorType << 31

	return code
}

func Decode(errCode uint32) (uint32, uint32, uint32, uint32, uint32, uint32, uint32) {
	errorType := errCode >> 31

	errCode -= errorType << 31

	sysId := errCode >> 26

	errCode -= sysId << 26

	serviceType := errCode >> 25

	errCode -= serviceType << 25

	serviceId := errCode >> 20

	errCode -= serviceId << 20

	logicId := errCode >> 14

	errCode -= logicId << 14

	oprId := errCode >> 8

	errCode -= oprId << 8

	errorId := errCode

	return errorType, sysId, serviceType, serviceId, logicId, oprId, errorId
}

func IsSysErr(code uint32) bool {
	errType, _, _, _, _, _, _ := Decode(code)
	return errType == Sys
}

func IsLogicErr(code uint32) bool {
	errType, _, _, _, _, _, _ := Decode(code)
	return errType == Logic
}
