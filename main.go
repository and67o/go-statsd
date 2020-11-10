package main

func main() {
	a := App{}
	a.Initialize()

	a.Run()
}

//func main() {
//	c, err := statsd.New() // Connect to the UDP port 8125 by default.
//	if err != nil {
//		log.Print(err)
//	}
//	defer c.Close()
//
//	total, err := 2,nil
//
//	fmt.Println(total)
//	c.Gauge("num_goroutine", total)
//
//}
