package main

import (
	"bufio"
	"fmt"
	log "log/slog"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/fatih/color"
)

func main() {
	dumpFiles()
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Error("rlimit")
		panic(err)
	}

	fmt.Print(tie) // Print the main logo

	s := InitSecurity()
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
	//var validMaps = map[int]bool
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
	data := encode(fmt.Sprintf("%s %s", user, diff.String()))
	fmt.Println(data)

}

// rot13(alphabets) + rot5(numeric)
func encode(input string) string {

	var result []rune
	rot5map := map[rune]rune{'0': '5', '1': '6', '2': '7', '3': '8', '4': '9', '5': '0', '6': '1', '7': '2', '8': '3', '9': '4'}

	for _, i := range input {
		switch {
		case !unicode.IsLetter(i) && !unicode.IsNumber(i):
			result = append(result, i)
		case i >= 'A' && i <= 'Z':
			result = append(result, 'A'+(i-'A'+13)%26)
		case i >= 'a' && i <= 'z':
			result = append(result, 'a'+(i-'a'+13)%26)
		case i >= '0' && i <= '9':
			result = append(result, rot5map[i])
		case unicode.IsSpace(i):
			result = append(result, ' ')
		}
	}
	return string(result[:])
}
