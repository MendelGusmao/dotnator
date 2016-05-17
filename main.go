package main

import (
	"fmt"
	"hash/crc64"
	"net/mail"
	"os"
	"strconv"
	"strings"
)

/*

From https://support.google.com/mail/answer/10313
	"Gmail doesn't recognize dots as characters within usernames,
	you can add or remove the dots from a Gmail address without
	changing the actual destination address; they'll all go to
	your inbox, and only yours."

	This piece of code provides a deterministic way of generating different e-mail addresses
	from an e-mail address by putting dots between its characters. So, after receiving an spam
	sent to one of your "new" addresses, you can track and have the information of what service
	forwarded your address.

	It's just a proof of concept and should not be used seriously. Some smart spammers might just
	remove the dots (and/or the plus sign and what comes after) from the addresses they receive.
	Also, the services you use can simply consider your address as not valid.

	Usage: dotnator <email address> <service name or address> [salt]
*/

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:\n\tdotnator <email address> <service name or address> [salt]")
		os.Exit(1)
	}

	if _, err := mail.ParseAddress(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "invalid email address")
		os.Exit(1)
	}

	email := strings.Split(os.Args[1], "@")
	username, host := email[0], email[1]
	service := os.Args[2]
	salt := ""

	if len(os.Args) == 4 {
		salt = os.Args[3]
	}

	plus := ""

	if i := strings.Index(username, "+"); i > -1 {
		plus = username[i:]
		username = username[:i]
	}

	key := append([]byte(salt), []byte(service)...)
	crc := crc64.Checksum(key, crc64.MakeTable(crc64.ECMA))
	crcp := fmt.Sprintf("%063s", strconv.FormatInt(int64(crc), 2))
	index := int(service[0]) % (63 - len(username))
	name := ""

	for i := 0; i < len(username); i++ {
		name += string(username[i])

		if crcp[i+index] == '1' && i < len(username)-1 {
			name += "."
		}
	}

	fmt.Printf("%s%s@%s\n", name, plus, host)
}
