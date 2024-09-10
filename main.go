package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	log "log/slog"
	"os"
	"strings"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/fatih/color"
)

func main() {
	uid := os.Getuid()
	if uid != 0 {
		fmt.Println("System>", color.YellowString(fmt.Sprintf("id [%d]>", os.Getuid())), color.RedString("Emperium system must be started with as root, please use sudo"))
		os.Exit(1)
	}

	err := dumpFiles()
	if err != nil {
		fmt.Println("System>", color.RedString(err.Error()))
	}

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Error("rlimit")
		panic(err)
	}

	fmt.Print(tie) // Print the main logo

	s := InitSecurity(4)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter User> ")
	text, _ := reader.ReadString('\n')
	user := strings.Trim(text, " \n")
	fmt.Println("Security Status>", color.RedString("User ["), color.WhiteString(user), color.RedString("] does not exist, system will be disabled in 2 hours!"))
	start := time.Now()

	go func() {
		time.Sleep(time.Hour * 2)
		fmt.Println("Security Status>", color.RedString("Exiting system, Imperial guard are alerted to your presence"))
		os.Exit(1)
	}()

	// Create our templatae map spec
	mapSpec := ebpf.MapSpec{
		Type:       ebpf.Hash,
		KeySize:    1, // 4 bytes for u32
		ValueSize:  20,
		MaxEntries: 1, // We'll have 5 maps inside this map
	}
	var watchedMaps [4]*ebpf.Map // Hold the four maps that we care about

	name, contents := create_maps(1000) // Generate unique names/contents for the maps
	// Lets create many maps
	for i := range name {
		mapSpec.Name = name[i]
		mapSpec.Contents = []ebpf.MapKV{
			{Key: uint8(1), Value: []byte(contents[i])},
		}
		m, err := ebpf.NewMap(&mapSpec)
		if err != nil {
			log.Error("map create fail ")
			panic(err)
		}
		defer m.Close()
		if i < 4 { // Just grab the first four maps for now // TODO:
			watchedMaps[i] = m
		}
	}

	fmt.Println("Data system>", color.GreenString("Ready"))
	// Start the Tie Fighter security systems
	s.Status()

	s.keyWatch(watchedMaps)
	end := time.Now()
	diff := end.Sub(start)
	data, err := EncryptMessage([]byte("James Earl Jones"), fmt.Sprintf("%s %s", user, diff.String()))
	if err != nil {
		panic(err)
	}
	fmt.Println("Security>", color.BlueString("Root Key>"), color.GreenString(data))
	fmt.Println("Security>", color.BlueString("Root Key>"), "Please make sure you copy your \"Root key\" into the CTF Slack channel!")

	fmt.Println("System>", color.GreenString("Shutting Down!"))
	time.Sleep(time.Second * 3)
	fmt.Printf("\n\n\n")
	fmt.Printf("Outro:\n========\n")
	outro := `You call Blue Hex on her holocomm. You tell her: the chocolates are so good! Oh, and also you shut down the TIE fighter product lines, all of them. It will take a couple of years for the engineers to repair the damage. This is more time than you need to clear up the skies around you, and safely move the base to another system. While hacking into the Imperial mainframe, you also discovered where Bajeroff Lake resides. Now, if only your ship was able to fly. Beeping sounds? IP-V6, what do you say? The _Yellow Stripe_ is ready at last? And there you go, Jephen'Tsa, rushing to your next adventure.`
	fmt.Println(outro)
	fmt.Printf("\n\nMay the Force accompany you.\n")
}

// rot13(alphabets) + rot5(numeric)
// func encode(input string) string {

// 	var result []rune
// 	rot5map := map[rune]rune{'0': '5', '1': '6', '2': '7', '3': '8', '4': '9', '5': '0', '6': '1', '7': '2', '8': '3', '9': '4'}

// 	for _, i := range input {
// 		switch {
// 		case !unicode.IsLetter(i) && !unicode.IsNumber(i):
// 			result = append(result, i)
// 		case i >= 'A' && i <= 'Z':
// 			result = append(result, 'A'+(i-'A'+13)%26)
// 		case i >= 'a' && i <= 'z':
// 			result = append(result, 'a'+(i-'a'+13)%26)
// 		case i >= '0' && i <= '9':
// 			result = append(result, rot5map[i])
// 		case unicode.IsSpace(i):
// 			result = append(result, ' ')
// 		}
// 	}
// 	return string(result[:])
// }

func EncryptMessage(key []byte, message string) (string, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}
