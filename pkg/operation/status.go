package operation

type OperationStatus string

const OpStatusNew = OperationStatus("new")
const OpStatusRunning = OperationStatus("running")
const OpStatusDone = OperationStatus("done")
const OpStatusFailed = OperationStatus("failed")
