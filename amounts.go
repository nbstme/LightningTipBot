package main

import (
	"errors"
	"strconv"
	"strings"
)

func getArgumentFromCommand(input string, which int) (output string, err error) {
	if len(strings.Split(input, " ")) < which+1 {
		errmsg := "Message doesn't contain enoough arguments"
		// log.Errorln(errmsg)
		return "None", errors.New(errmsg)
	}
	output = strings.Split(input, " ")[which]
	return output, nil
}

func DecodeAmountFromCommand(input string) (amount int, err error) {
	if len(strings.Split(input, " ")) < 2 {
		errmsg := "Message doesn't contain any amount"
		// log.Errorln(errmsg)
		return 0, errors.New(errmsg)
	}
	amount, err = getAmount(input)
	return amount, nil
}

func getAmount(input string) (amount int, err error) {
	amount, err = strconv.Atoi(strings.Split(input, " ")[1])
	if err != nil {
		return 0, err
	}
	if amount < 1 {
		errmsg := "Error: Amount must be greater than 0"
		// log.Errorln(errmsg)
		return 0, errors.New(errmsg)
	}
	return amount, err
}
