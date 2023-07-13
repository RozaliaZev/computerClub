package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Event struct {
	Time      string
	EventType int
	Client    *Client
	Table     int
}

type Client struct {
	NameClient     string
	PresenceClient bool
}

type Table struct {
	Number     int
	Revenue    int
	Occupied   int
	TotalTime  int
	LastClient string
	StartTime  []string
	StopTime   []string
}

var hourRate int

func main() {
	// Чтение входного файла
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Считывание числа столов
	scanner.Scan()
	numTables := 0
	fmt.Sscanf(scanner.Text(), "%d", &numTables)

	// Считывание времени работы
	scanner.Scan()
	var startTime, endTime string
	fmt.Sscanf(scanner.Text(), "%s %s", &startTime, &endTime)

	// Считывание стоимости часа
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &hourRate)

	// Создание таблиц
	tables := make([]Table, numTables)
	for i := 0; i < numTables; i++ {
		tables[i].Number = i + 1
	}

	// Обработка событий
	events := make([]Event, 0)
	for scanner.Scan() {
		line := scanner.Text()
		event := parseEvent(line, &tables)
		events = append(events, event)
	}

	// Обработка времени начала работы
	fmt.Println(startTime)

	presenceClients := make([]string, 0)
	// Выполнение событий
	for _, event := range events {
		if event.Table != 0 {
			fmt.Printf("%s %v %v %d\n", event.Time, event.EventType, event.Client.NameClient, event.Table)
		} else {
			fmt.Printf("%s %v %v\n", event.Time, event.EventType, event.Client.NameClient)
		}

		err := performEvent(&event, &tables, startTime, endTime)
		if err != nil {
			fmt.Printf("%s 13 %s\n", event.Time, err.Error())
		}

		if event.Client.PresenceClient {
			if !contains(presenceClients, event.Client.NameClient) {
				presenceClients = append(presenceClients, event.Client.NameClient)
			}
		} else {
			presenceClients = remove(presenceClients, event.Client.NameClient)
		}

	}

	sort.Strings(presenceClients)
	for _, presenceClient := range presenceClients {
		fmt.Printf("%s 11 %s\n", endTime, fmt.Errorf(presenceClient))
		addEndTimeForTable(presenceClient, endTime, &tables)
	}

	// Обработка времени окончания работы
	fmt.Println(endTime)

	revenue := 0
	// Вывод информации о столах
	for _, table := range tables {
		for i := 0; i < len(table.StartTime); i++ {
			table.TotalTime += getTimeDifference(table.StartTime[i], table.StopTime[i])
		}

		if table.TotalTime%60 == 0 {
			table.Revenue = hourRate * (table.TotalTime/60)
		} else {
			table.Revenue = hourRate * (table.TotalTime/60 + 1)
		}
		revenue += table.Revenue
		fmt.Printf("%d %d %s\n", table.Number, table.Revenue, formatTime(table.TotalTime))
	}
	fmt.Println("revenue for 1 day:", revenue)
}

func parseEvent(line string, tables *[]Table) Event {
	parts := strings.Split(line, " ")
	event := Event{
		Time:      parts[0],
		EventType: 0,
		Client:    &Client{NameClient: "", PresenceClient: true},
		Table:     0,
	}

	switch parts[1] {
	case "1":
		event.EventType = 1
		event.Client.NameClient = parts[2]
	case "2":
		event.EventType = 2
		event.Client.NameClient = parts[2]
		event.Table = parseInt(parts[3])
	case "3":
		event.EventType = 3
		event.Client.NameClient = parts[2]
	case "4":
		event.EventType = 4
		event.Client.NameClient = parts[2]
	}

	return event
}

func parseInt(s string) int {
	var num int
	fmt.Sscanf(s, "%d", &num)
	return num
}

func performEvent(event *Event, tables *[]Table, startTime, endTime string) error {
	switch event.EventType {
	case 1:
		return handleClientArrival(event, tables, startTime, endTime)
	case 2:
		return handleClientSeat(event, tables)
	case 3:
		return handleClientWait(event, tables)
	case 4:
		return handleClientLeave(event, tables)
	}

	return nil
}

func handleClientArrival(event *Event, tables *[]Table, startTime, endTime string) error {

	if startTime > endTime {
		if !(event.Time >= startTime && event.Time <= endTime) {
			return fmt.Errorf("NotOpenYet")
		}
	} else {
		if !(event.Time >= startTime && event.Time <= endTime) {
			return fmt.Errorf("NotOpenYet")
		}
	}

	return nil
}

func handleClientSeat(event *Event, tables *[]Table) error {
	for _, table := range *tables {
		if table.LastClient == event.Client.NameClient {
			return fmt.Errorf("PlaceIsBusy")
		}
	}

	for i, table := range *tables {
		if table.Number == event.Table {
			if table.Occupied == 1 {
				return fmt.Errorf("PlaceIsBusy")
			}
			(*tables)[i].StartTime = append((*tables)[i].StartTime, event.Time)
			(*tables)[i].Occupied = 1
			(*tables)[i].LastClient = event.Client.NameClient
			return nil
		}
	}

	return nil
}

func handleClientWait(event *Event, tables *[]Table) error {
	for i, table := range *tables {
		if table.Occupied == 0 {
			return fmt.Errorf("ICanWaitNoLonger!")
		} else if event.Table == table.Number {
			(*tables)[i].StopTime = append((*tables)[i].StopTime, event.Time)
			break
		}
	}

	return nil
}

func handleClientLeave(event *Event, tables *[]Table) error {
	if !event.Client.PresenceClient {
		return fmt.Errorf("ClientUnknown")
	}

	event.Client.PresenceClient = false

	for i, table := range *tables {
		if table.LastClient == event.Client.NameClient {
			(*tables)[i].Occupied = 0
			(*tables)[i].StopTime = append((*tables)[i].StopTime, event.Time)
			(*tables)[i].LastClient = ""
			break
		}
	}

	return nil
}

func getTimeDifference(startTime string, endTime string) int {
	var startHour, startMinute, endHour, endMinute int
	fmt.Sscanf(startTime, "%d:%d", &startHour, &startMinute)
	fmt.Sscanf(endTime, "%d:%d", &endHour, &endMinute)

	startTimeMins := startHour*60 + startMinute
	endTimeMins := endHour*60 + endMinute

	return endTimeMins - startTimeMins
}

func remove(slice []string, value string) []string {
	for i, element := range slice {
		if element == value {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}

func contains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}

func addEndTimeForTable(presenceClient string, endTime string, tables *[]Table) {
	for i, table := range *tables {
		if table.LastClient == presenceClient {
			(*tables)[i].StopTime = append((*tables)[i].StopTime, endTime)
		}
	}
}

func formatTime(minutes int) string {
    hours := minutes / 60
    minutes = minutes % 60

    return fmt.Sprintf("%02d:%02d", hours, minutes)
}
