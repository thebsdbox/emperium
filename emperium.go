package main

import (
	"fmt"
	log "log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cilium/ebpf"
	"github.com/fatih/color"
)

// pass map pointers once all the maps have been created as finding existing maps seems impossible
func (s *securityLevel) keyWatch(watchedMaps [4]*ebpf.Map) {

	_, first := os.LookupEnv("SKIPFIRST")
	_, second := os.LookupEnv("SKIPSECOND")
	_, third := os.LookupEnv("SKIPTHIRD")
	_, fourth := os.LookupEnv("SKIPFOURTH")

	var wg, wg2 sync.WaitGroup
	wg2.Add(1)
	// lockWait := make(chan struct{})
	var once sync.Once

	if !first {
		s.firstLock(watchedMaps[0])
	}
	if !second {
		go func() {
			s.secondLock(&wg, &wg2, &once)
		}()
	}
	wg2.Wait()
	if !third {
		s.thirdLock(watchedMaps[2])
	}
	if !fourth {
		s.fourthLock(watchedMaps[3])
	}
	wg.Wait()

}

func (s *securityLevel) firstLock(m *ebpf.Map) error {
	var value [20]byte
	err := m.Lookup(uint8(1), &value)
	if err != nil {
		log.Error(fmt.Sprintf("This shouldn't occur [%v]", err))
	}
	strValue := string(value[:])
	correctValue := Reverse(strValue)
	// Horrific string manipulation as we can't get JUST the name of the map
	a := strings.Split(m.String(), "(")
	b := strings.Split(a[1], ")")
	fmt.Println("Data system>", color.RedString("Map ["), color.WhiteString(string(b[0])), color.RedString("] is corrupt!"))

	for {
		err := m.Lookup(uint8(1), &value)
		if err != nil {
			log.Error(fmt.Sprintf("This shouldn't occur [%v]", err))
		}
		time.Sleep(time.Second * 5) // We check this map
		strValue := string(value[:])
		if strValue == correctValue {
			s.Unlock(0)
			break // pop out the loop
		}
	}
	return nil
}
func (s *securityLevel) secondLock(wg *sync.WaitGroup, wg2 *sync.WaitGroup, once *sync.Once) error {
	var value [20]byte
	var m *ebpf.Map
	wg.Add(1)
	waitCounter := 1
	mapName := fmt.Sprintf("empire_%s", RandStringBytesMaskImprSrcSB(3))
	for {
		found, correct := false, false
		for id := ebpf.MapID(0); ; {
			var err error
			id, err = ebpf.MapGetNextID(ebpf.MapID(id))
			if err != nil {
				break
			}
			m, err = ebpf.NewMapFromID(id)
			if err != nil {
				panic(err)
			}
			info, err := m.Info()
			if err != nil {
				panic(err)
			}

			if info.Name == mapName {
				found = true
				err := m.Lookup(uint8(1), &value)

				if err != nil {
					log.Info(fmt.Sprintf("%v %d %d %d,", err, m.ValueSize(), m.FD(), id))
					fmt.Println("Data system>", color.YellowString("Map ["), color.WhiteString(mapName), color.YellowString("] has no data!"))
					continue
				}
				strValue := string(value[:])
				if strValue == "brRz3HVSVzC6RXrBC2Y7" {
					correct = true
					m.Close()
					break // pop out the loop
				} else {
					continue
				}
			}
			m.Close()
		}

		// After checking all maps, see if the map was found
		if found {
			if correct {

				once.Do(func() { wg2.Done() })
				if waitCounter == 1 {
					wg.Done()
					waitCounter = 0
				}
				s.Unlock(1)
			} else {
				fmt.Println("Data system>", color.YellowString("Map ["), color.WhiteString(mapName), color.YellowString("] has incorrect data!"))
			}
		} else {
			fmt.Println("Data system>", color.RedString("Map ["), color.WhiteString(mapName), color.RedString("] is missing!"))
			s.Lock(1)
			if waitCounter == 0 {
				wg.Add(1)
				waitCounter = 1
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func (s *securityLevel) thirdLock(m *ebpf.Map) error {
	fmt.Println("Authentication>", color.YellowString("Waiting for Auth on port 2000"))

	// Start a listener on localhost (port 2000)
	go func() {
		l, err := net.Listen("tcp", "127.0.0.1:2000")
		if err != nil {
			panic(err)
		}
		defer l.Close()
		for {
			// Wait for a connection.
			conn, _ := l.Accept()
			// We will ignore a lot of the errors (as we expect things to be a little broken)
			// if err != nil {
			// 	//log.Fatal(err)
			// }
			go func(conn net.Conn) {
				buf := make([]byte, 1024)
				_, err := conn.Read(buf)
				if err != nil {
					//fmt.Printf("Error reading: %#v\n", err)
					//return
				}
				fmt.Println("Authentication>", color.GreenString("Authorisation recieved"))

				// fmt.Printf("Message received: %s\n", string(buf[:len]))

				conn.Write([]byte("Message received.\n"))
				conn.Close()
			}(conn)
		}
	}()

	// Connect to e2e endpoint with a second timeout
	for {
		time.Sleep(time.Second)

		conn, err := net.DialTimeout("tcp", "127.0.0.1:2001", time.Second)
		if err != nil {
			//log.Error(fmt.Sprintf("Dial failed: %v", err.Error()))
			continue
		}
		_, err = conn.Write([]byte("The Grid, a digital frontier"))
		if err != nil {
			//log.Error("Write data failed: %v ", err.Error())
		}

		// buffer to get data
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			//log.Error("Read data failed:", err.Error())
		} else {
			//println("Received message: %s", string(received))
			conn.Close()
			break
		}
		// Wait for a second and connect again
		time.Sleep(time.Second)
	}
	return nil
}

func (s *securityLevel) fourthLock(m *ebpf.Map) error {
	fmt.Println("Connect>", color.YellowString("Empire local Mainframe"))
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		Port: 9000,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	response := make([]byte, 3) // Get map identifier

	// Attempt to send data (3 bytes) to the remote address
	for {

		_, err = conn.Write([]byte("SYN\n"))
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
		}
		conn.SetReadDeadline(time.Now().Add(time.Second))

		_, _, err = conn.ReadFromUDP(response)
		if err != nil {
			fmt.Println("Connect>", color.RedString("Mainframe connection failure"))
			// log.Error(fmt.Sprintf("%v", err))
		} else {
			if string(response) != "ACK" {
				fmt.Println("Connect>", color.RedString("Mainframe sent incorrect response [%s]", response))
			} else {
				fmt.Println("Connect>", color.GreenString("Mainframe acknowledged response"))
				s.Unlock(2)
				break // pop out the loop
			}
		}
		time.Sleep(time.Second * 10)
	}
	return nil
}
