package main

import (
	"flag"
	"fmt"
	"forza/models"
	"forza/parser"
	"forza/recorder"
	"log"
	"strconv"
	"strings"
)

func main() {

	fmt.Println("Starting multi-car recorder")

	portFlag := flag.String("ports", "5030-5040", "Comma-separated list or inclusive range (start-end) of UDP ports to listen on")
	flag.Parse()

	ports, err := parsePorts(*portFlag)
	if err != nil {
		log.Fatalf("invalid -ports value %q: %v", *portFlag, err)
	}

	// Shared packet channel
	packetStream := make(chan recorder.RawPacket, 1000)
	recordingStarted := false

	for _, port := range ports {
		err := recorder.Listen(port, packetStream)
		if err != nil {
			panic(err)
		}
	}

	// Car registry: port → car object
	cars := make(map[string]*models.Car)
	readyCars := make(map[string]bool)              // carID → seen IsRaceOn == true
	firstRaceOn := make(map[string]models.Carstate) // carID → first IsRaceOn packet (buffer until start)

	for pkt := range packetStream {

		// 1. If car doesn't exist yet, create one
		if _, exists := cars[pkt.CarID]; !exists {
			cars[pkt.CarID] = &models.Car{
				Name:   fmt.Sprintf("Car-%s", pkt.CarID),
				States: []models.Carstate{},
			}
			fmt.Println("Detected new car on port:", pkt.CarID)
		}

		car := cars[pkt.CarID]

		// 2. Parse packet
		state, err := parser.RawtoCarstate(pkt.Data)
		if err != nil {
			fmt.Println("Parse error:", err)
			continue
		}

		if !recordingStarted {
			if state.IsRaceOn {
				// Capture the first race-on packet for this car.
				if !readyCars[pkt.CarID] {
					firstRaceOn[pkt.CarID] = state
					readyCars[pkt.CarID] = true
				}
			} else {
				// Mark the car as seen but not ready yet.
				if _, seen := readyCars[pkt.CarID]; !seen {
					readyCars[pkt.CarID] = false
				}
			}

			if allCarsReady(readyCars) {
				recordingStarted = true
				fmt.Printf("All %d players reporting race on. Starting recording.\n", len(readyCars))

				// Seed each car with its first race-on packet.
				for id, s := range firstRaceOn {
					cars[id].AddState(s)
				}
			} else {
				readyCount, total := countReady(readyCars)
				fmt.Printf("Waiting for race start: %d/%d players reported race on\n", readyCount, total)
				continue
			}
		} else {
			car.AddState(state)
		}

		fmt.Printf("[%s] Speed: %.1f MPH | Gear %d\n",
			car.Name, state.SpeedMPH, state.Gear)

		// 3. Race end detection
		if recordingStarted && allCarsFinished(cars) {
			fmt.Println("All cars finished the race!")
			break
		}
	}

	fmt.Println("Total cars:", len(cars))
	for id, c := range cars {
		fmt.Printf("-- %s has %d states\n", id, len(c.States))
	}

	fmt.Println("Exporting CSV files...")
	for id, car := range cars {
		err := car.ExportCSV()
		if err != nil {
			fmt.Printf("Failed to write CSV for %s: %v\n", id, err)
		}
	}

}

func allCarsFinished(cars map[string]*models.Car) bool {
	if len(cars) == 0 {
		return false
	}

	for _, c := range cars {
		if len(c.States) == 0 {
			return false
		}
		last := c.States[len(c.States)-1]
		if last.IsRaceOn {
			return false
		}
	}
	return true
}

func allCarsReady(ready map[string]bool) bool {
	if len(ready) == 0 {
		return false
	}
	for _, isOn := range ready {
		if !isOn {
			return false
		}
	}
	return true
}

func countReady(ready map[string]bool) (int, int) {
	total := len(ready)
	count := 0
	for _, isOn := range ready {
		if isOn {
			count++
		}
	}
	return count, total
}

func parsePorts(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("ports string cannot be empty")
	}

	// Range form: "start-end"
	if strings.Contains(input, "-") && !strings.Contains(input, ",") {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("range must be start-end")
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start port: %w", err)
		}
		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end port: %w", err)
		}

		step := 1
		if start > end {
			step = -1
		}

		var ports []string
		for p := start; ; p += step {
			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("port %d out of range", p)
			}
			ports = append(ports, strconv.Itoa(p))
			if p == end {
				break
			}
		}
		return ports, nil
	}

	// Comma-separated list
	var ports []string
	for _, part := range strings.Split(input, ",") {
		val := strings.TrimSpace(part)
		if val == "" {
			return nil, fmt.Errorf("empty port entry in list")
		}
		p, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid port %q: %w", val, err)
		}
		if p < 1 || p > 65535 {
			return nil, fmt.Errorf("port %d out of range", p)
		}
		ports = append(ports, strconv.Itoa(p))
	}
	return ports, nil
}
