package passwordsim

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestSend(t *testing.T) {
	message := make(chan string, bufferSize)
	passwords := filepath.Join("test", "passwords.txt")
	var wg sync.WaitGroup
	go send(message, passwords)

	receivedPasswords := []string{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for password := range message {
			receivedPasswords = append(receivedPasswords, password)
		}
	}()

	wg.Wait()
	expectedPasswords := []string{
		"// This file is for testing",
		"password",
		"123456",
		"123456789",
		"guest",
		"qwerty",
		"12345678",
		"111111",
		"12345",
		"col123456",
		"123123",
		"1234567",
		"1234",
		"1234567890",
		"000000",
		"555555",
		"666666",
		"123321",
		"654321",
		"7777777",
		"123",
		"correct horse battery staple",
		"incorrect horse battery staple",
		"incorrect horse battery st@ple",
	}
	if !reflect.DeepEqual(receivedPasswords, expectedPasswords) {
		t.Errorf("Passwords received do not match contents of expectedPasswords")
	}
}

type checkPasswordsTest struct {
	expectedContent string
	passwordToCheck string
	threshold       float64
}

func TestCheckPasswords(t *testing.T) {
	passwords := filepath.Join("test", "passwords.txt")

	tests := []checkPasswordsTest{
		{"1234567 0.125\n12345678 0.125\n123456789 0.2222222222222222\n123456 0.25\n1234567890 0.3\n", "12345677", 0.3},
		{"1234567 0.125\n12345678 0.125\n123456789 0.2222222222222222\n123456 0.25\n1234567890 0.3\n12345 0.375\n1234 0.5\n" +
			"col123456 0.5555555555555556\n123 0.625\n123123 0.625\n123321 0.625\n7777777 0.75\n111111 0.875\n555555 0.875\n" +
			"654321 0.875\n666666 0.875\n// This file is for testing 1\n000000 1\ncorrect horse battery staple 1\nguest 1\n" +
			"incorrect horse battery st@ple 1\nincorrect horse battery staple 1\npassword 1\nqwerty 1\n", "12345677", 1.1},
		{"\n", "12345677", -1},
		{"correct horse battery staple 0\nincorrect horse battery staple 0.06666666666666667\n" +
			"incorrect horse battery st@ple 0.1\n", "correct horse battery staple", 0.3},
		{"password 0\n", "password", 0},
	}

	for idx, testCase := range tests {
		tempfile, err := os.CreateTemp("", "passwordsim-")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tempfile.Name())
		output := tempfile.Name()
		passwordToCheck := testCase.passwordToCheck
		threshold := testCase.threshold
		CheckPasswords(passwords, output, passwordToCheck, threshold)

		expectedContent := testCase.expectedContent
		tempfile.Seek(0, 0)
		reader := bufio.NewScanner(tempfile)
		lines := []string{}
		for reader.Scan() {
			lines = append(lines, reader.Text())
		}
		content := strings.Join(lines, "\n") + "\n"
		if !reflect.DeepEqual(content, expectedContent) {
			t.Errorf("Test case #%d failed. Expected %q, got %q", idx+1, expectedContent, content)
		}
	}
}
