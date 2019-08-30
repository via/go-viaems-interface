package viaems

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Received an inappropriate response from the target EMS
var MalformedTargetResponse = errors.New("Malformed response from EMS")

type WireTarget struct {
	*bufio.ReadWriter

	requestChannel  chan request
	updates         []chan Status
	debug           chan string
	update_requests chan chan Status
	feed_fields     []string
}

func (target *WireTarget) GetStatusUpdates() chan Status {
	c := make(chan Status, 100)
	target.update_requests <- c
	return c
}

func (target *WireTarget) GetTable(name string) (TableConfig, error) {
	ret := TableConfig{}

	res, err := target.Command("get config.tables." + name)
	if err != nil {
		return ret, err
	}

	attrs := strings.Split(res, " ")
	for _, attr := range attrs {
		kvs := strings.Split(attr, "=")
		if len(kvs) != 2 {
			return ret, MalformedTargetResponse
		}
		switch kvs[0] {
		case "name":
			ret.Name = kvs[1]
		case "naxis":
			ret.AxisCount, err = strconv.Atoi(kvs[1])
		case "rows":
			ret.RowCount, err = strconv.Atoi(kvs[1])
		case "cols":
			ret.ColumnCount, err = strconv.Atoi(kvs[1])
		case "rowname":
			ret.RowName = kvs[1]
		case "colname":
			ret.ColumnName = kvs[1]
		case "rowlabels":
			ret.RowLabels = strings.Split(strings.Trim(kvs[1], "[]"), ",")
		case "collabels":
			ret.ColumnLabels = strings.Split(strings.Trim(kvs[1], "[]"), ",")
		default:
		}
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

//* name=ve naxis=2 rows=16 rowname=RPM cols=16 colname=MAP rowlabels=[250.0,500.0,900.0,1200.0,1600.0,2000.0,2400.0,3000.0,3600.0,4000.0,4400.0,5200.0,5800.0,6400.0,6800.0,7200.0] collabels=[20.0,30.0,40.0,50.0,60.0,70.0,80.0,90.0,100.0,120.0,140.0,160.0,180.0,200.0,220.0,240.0]

func (target *WireTarget) ListTables() ([]string, error) {
	res, err := target.Command("list config.tables.")
	if err != nil {
		return []string{}, err
	}

	tablenames := strings.Split(res, " ")
	/* Remove prefixes */
	for i, name := range tablenames {
		tablenames[i] = strings.TrimPrefix(name, "config.tables.")
	}
	return tablenames, nil
}

func (target *WireTarget) Command(node string) (string, error) {
	c := make(chan string)
	req := request{
		request_str: node + "\n",
		notify:      c,
	}
	target.requestChannel <- req
	res := <-c
	res = strings.TrimRight(res, "\n\r ")
	if strings.HasPrefix(res, "- ") {
		return "", errors.New(strings.TrimPrefix(res, "- "))
	}

	return strings.TrimPrefix(res, "* "), nil
}

func (target *WireTarget) parseStatusUpdate(line string) (Status, error) {
	line = strings.TrimRight(line, "\r\n ")
	vals := strings.Split(line, ",")
	status := Status{
    WallTime: time.Now(),
    Sensors: make(map[string]*SensorStatus),
    Fueling: make(map[string]float64),
    Decoder: make(map[string]string),
    Ignition: make(map[string]float64),
  }

	if len(vals) != len(target.feed_fields) {
		return Status{}, errors.New("Feed field length mismatch")
	}

	var err error
	for i, v := range vals {
		switch fieldname := target.feed_fields[i]; {
		case strings.HasPrefix(fieldname, "status.sensors."):
			sensor := strings.Split(fieldname, ".")[2]
			if status.Sensors[sensor] == nil {
				status.Sensors[sensor] = &SensorStatus{}
			}

			if strings.HasSuffix(fieldname, ".fault") {
			} else {
				status.Sensors[sensor].Value, err = strconv.ParseFloat(v, 64)
			}

		case fieldname == "status.current_time":
			status.CpuTime, err = strconv.ParseFloat(v, 64)

    case strings.HasPrefix(fieldname, "status.ignition"):
      status.Ignition[strings.Split(fieldname, ".")[2]], err = strconv.ParseFloat(v, 64)

    case strings.HasPrefix(fieldname, "status.fueling"):
      status.Fueling[strings.Split(fieldname, ".")[2]], err = strconv.ParseFloat(v, 64)

    case strings.HasPrefix(fieldname, "status.decoder"):
      status.Decoder[strings.Split(fieldname, ".")[2]] = v

		default:
		}
		if err != nil {
			return status, MalformedTargetResponse
		}
	}
	return status, nil
}

func wireInputLoop(buf *WireTarget, result chan string) {
	for {
		line, err := buf.ReadString('\n')
		if err == nil {
			result <- line
		}
	}
}

type request struct {
	request_str string
	notify      chan string
}

func (target *WireTarget) process() {
	inputChannel := make(chan string, 100)
	go wireInputLoop(target, inputChannel)

	var current_request request
	var request_pending bool = false

	for {
		var line string
		var reqchan chan request
		var updatereq chan Status
		/* If we're not actively handling a request, we're open to receiving a new
		 * one */
		if !request_pending {
			reqchan = target.requestChannel
		}

		select {
		case line = <-inputChannel:
			if strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "- ") {
				current_request.notify <- line
				request_pending = false
			} else if strings.HasPrefix(line, "# ") {
				select {
				case target.debug <- line:
				default:
				}
			} else {
				status, err := target.parseStatusUpdate(line)
        if err != nil {
          fmt.Println(err)
        }
				for _, client := range target.updates {
					// Iterate though registered update clients, send nonblocking updates
					select {
					case client <- status:
					default:
					}
				}
			}
		case current_request = <-reqchan:
			request_pending = true
			target.WriteString(current_request.request_str)
			target.Flush()
		case updatereq = <-target.update_requests:
			target.updates = append(target.updates, updatereq)
		}
	}
}

func (target *WireTarget) InitializeTarget() {
  /*Get all possible statuses*/
  res, _ := target.Command("list status.")
  statuses := strings.Split(res, " ")

	/* Set up feed line */
	res, _ = target.Command(fmt.Sprintf("set config.feed %v", strings.Join(statuses, ",")))
	target.feed_fields = statuses
}

func OpenTCPInterface(addr string) (*WireTarget, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	wt := &WireTarget{
		ReadWriter:      rw,
		requestChannel:  make(chan request),
		updates:         make([]chan Status, 0),
		debug:           make(chan string),
		update_requests: make(chan chan Status),
		feed_fields:     []string{},
	}

	go wt.process()
	go wt.InitializeTarget()

	return wt, nil
}
