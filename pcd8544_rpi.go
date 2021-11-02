package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

func GetUpTime() int32 {
	sysi := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(sysi)
	if err != nil {
		fmt.Printf("Get sysinfo err ->[%s]", err.Error())
	}
	uptime := sysi.Uptime
	return uptime

}

func GetCPULoads() [3]uint32 {
	sysi := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(sysi)
	if err != nil {
		fmt.Printf("Get sysinfo err ->[%s]", err.Error())
	}
	avgCPULoads := sysi.Loads
	return avgCPULoads
}
func GetRamInfo() (totalram uint32, freeram uint32) {
	sysi := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(sysi)
	if err != nil {
		fmt.Printf("Get sysinfo err ->[%s]", err.Error())
	}
	totalram = sysi.Totalram
	freeram = sysi.Freeram
	return
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetCPUTemp() string {
	TEMP_FILE := "/sys/class/thermal/thermal_zone0/temp"
	f, _ := os.Open(TEMP_FILE)

	var b []byte = make([]byte, 32)
	n, err := f.Read(b)
	if err != nil {
		fmt.Printf("Open file->[%s] Failed", TEMP_FILE)
	}
	f.Close()               //close the file
	data := string(b[:n-1]) // remove '\n' from the origin line
	//fmt.Printf("temp = %s\n", data)

	tempFloat64, e := strconv.ParseFloat(data, 64)
	if e != nil {
		fmt.Printf("parseFloat error->[%s]", e.Error())
	}

	tempFloat64 = tempFloat64 / 1000.0
	//fmt.Printf("tempFloat64 = %f\n", tempFloat64)

	return fmt.Sprintf("%2.2f", tempFloat64)

}

func main() {

	//define the gpio pin for pcd8544
	//pin setup
	var (
		SCLK uint8 = 17
		DIN  uint8 = 18
		DC   uint8 = 27
		CS   uint8 = 22
		RST  uint8 = 23
		BL   uint8 = 4
	)

	var contrast uint8 = 45

	fmt.Printf("Raspberry Pi Nokia5110 sysinfo display\n")

	//Init LCD
	pin := LCDInit(SCLK, DIN, DC, CS, RST, BL, contrast)

	LCDClear()

	pin.LCDShowRpiLogo()

	for {
		LCDClear()

		//timeinfo
		timeObj := time.Now()
		month := timeObj.Month()
		day := timeObj.Day()
		hour := timeObj.Hour()
		minute := timeObj.Minute()
		second := timeObj.Second()

		timeInfo := fmt.Sprintf("%2d/%2d %d:%d:%d", int(month), day, hour, minute, second)
		timeInfoBytes := []byte(timeInfo)

		//ipinfo
		ipInfo := []byte("No Connect")
		ip, err := externalIP()
		if err == nil {
			ipInfo = []byte(ip.String())
		}

		//cputemp info
		cpuTemp := GetCPUTemp()
		//fmt.Printf("cputemp = %s\n", cpuTemp)
		cpuTempInfo := fmt.Sprintf("TEM %sC", cpuTemp)
		//fmt.Println(ip.String())

		//cpuinfo
		avgCpuLoad := GetCPULoads()[0] / 1000
		cpuinfo := fmt.Sprintf("CPU %d%s", avgCpuLoad, "%")

		//system uptime
		uptime := GetUpTime()
		uptimeDays := uptime / 86400
		uptimeHours := (uptime / 3600) - (uptimeDays * 24)
		uptimeMinus := (uptime / 60) - (uptimeDays * 1440) - (uptimeHours * 60)
		uptimeInfo := fmt.Sprintf("Up %dD%dH%dM", uptimeDays, uptimeHours, uptimeMinus)

		//ram info
		totalRam, freeRam := GetRamInfo()
		totalRam = totalRam / 1024 / 1024 // -> MB
		freeRam = freeRam / 1024 / 1024   // -> MB
		usedRam := totalRam - freeRam
		ram_load := (usedRam * 100) / totalRam
		ramInfo := fmt.Sprintf("RAM %.3dM %.2d%s", usedRam, ram_load, "%")

		LCDDrawString(0, 0, []byte(uptimeInfo))  //line0
		LCDDrawString(0, 1, timeInfoBytes)       //line1
		LCDDrawString(0, 2, []byte(cpuinfo))     //line2
		LCDDrawString(0, 3, []byte(ramInfo))     //line3
		LCDDrawString(0, 4, []byte(cpuTempInfo)) //line4
		LCDDrawString(0, 5, ipInfo)              //line5

		pin.LCDDisplay()

		time.Sleep(1000)

	}

}
