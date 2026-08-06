package iec104

// CauseOfTransmission constants represent the reason for data transmission.
const (
	COT_PERIODIC          CauseOfTransmission = 1
	COT_BACKGROUND_SCAN   CauseOfTransmission = 2
	COT_SPONTANEOUS       CauseOfTransmission = 3
	COT_INITIALIZED       CauseOfTransmission = 4
	COT_REQUEST           CauseOfTransmission = 5
	COT_ACTIVATION        CauseOfTransmission = 6
	COT_ACTIVATION_CON    CauseOfTransmission = 7
	COT_DEACTIVATION      CauseOfTransmission = 8
	COT_DEACTIVATION_CON  CauseOfTransmission = 9
	COT_ACTIVATION_TERM   CauseOfTransmission = 10
	COT_ACTIVATION_TERMINATION   CauseOfTransmission = 10 // 激活终止
	COT_INTERROGATED_STATION CauseOfTransmission = 20
)
