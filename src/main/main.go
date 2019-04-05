package main

const subsidy = 50
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func main() {
	bc := NewBlockchain("0xismail111111")
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
