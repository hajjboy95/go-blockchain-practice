package main

type TXInput struct {
	TXid      []byte
	Vout      int
	ScriptSig string
}

func (in *TXInput) CanUnlockWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
