package main

import (
	"bufio"
	"flag"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"os"
	"strconv"
	"time"
)

func main() {
	var address,
		username,
		password,
		database,
		hostname,
		namespace string
	interval := 60 * time.Second

	flag.StringVar(&address, "a", "http://localhost:8086", "InfluxDB HTTP endpoint. Default: http://localhost:8086")
	flag.StringVar(&username, "u", "", "InfluxDB username")
	flag.StringVar(&password, "p", "", "InfluxDB password")
	flag.StringVar(&database, "d", "monitoring", "InfluxDB database. Default: monitoring")
	flag.StringVar(&namespace, "n", "gpuz", "InfluxDB measurement title. Default: gpuz")
	flag.StringVar(&hostname, "h", "", "Hostname for current working machine. By default OS Hostname will be used")
	flag.Func("i", "Interval in seconds between measurements. Default: 60s", func(s string) error {
		intervalFromFlag, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if interval < 1 {
			return fmt.Errorf("interval can not be below 0")
		}
		interval = time.Duration(intervalFromFlag) * time.Second
		return nil
	})
	flag.Parse()

	if hostname == "" {
		h, err := os.Hostname()
		if err != nil {
			h = "undefined"
		}
		hostname = h
	}
	hostname = simplifyString(hostname)

	// Create client
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     address,
		Username: username,
		Password: password,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		in := bufio.NewScanner(os.Stdin)
		in.Scan()
	}
	collector := NewCollector()

	fmt.Printf("Collector has been started with parameters: endpoint %s, username %s, namespace %s, database %s, hostname %s\n",
		address, username, namespace, database, hostname)

	for {
		tags, values, err := collector.GetInfluxRow(hostname)
		if err != nil {
			fmt.Println("ERROR: ", err)
			time.Sleep(interval)
			continue
		}
		fmt.Println(tags, values)

		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Precision: "s",
			Database:  database,
		})
		if err != nil {
			fmt.Println("ERROR: ", err)
			time.Sleep(interval)
			continue
		}

		point, err := client.NewPoint(namespace, tags, values)
		if err != nil {
			fmt.Println("ERROR: ", err)
			time.Sleep(interval)
			continue
		}
		bp.AddPoint(point)

		err = influxClient.Write(bp)
		if err != nil {
			fmt.Println("ERROR: ", err)
			time.Sleep(interval)
			continue
		}

		fmt.Println("Success: metrics were sent")
		time.Sleep(interval)
	}
}
