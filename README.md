## Overview
This is the start of a library for [Tiniyo](https://www.tiniyo.com/). Gotiniyo supports making voice calls and sending text messages.

## Installation
To install gotiniyo, simply run 

`go get github.com/devtiniyo/gotiniyo`.

## SMS Example

	package main

	import (
		"github.com/devtiniyo/gotiniyo"
	)

	func main() {
		authID := "12345454545"
		authToken := "1234545454512345454545"
		tiniyo := gotiniyo.NewTiniyoClient(authID, authToken)

		from := "TINIYO"
		to := "+12345454545"
		message := "Welcome to gotiniyo!"
		tiniyo.SendSMS(from, to, message, "", "")
	}


