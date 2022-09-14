package nascan

// Package nascan https://github.com/dean2021/go-masscan

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type (
	Address struct {
		Addr     string `xml:"addr,attr"`
		AddrType string `xml:"addrtype,attr"`
	}
	State struct {
		State     string `xml:"state,attr"`
		Reason    string `xml:"reason,attr"`
		ReasonTTL string `xml:"reason_ttl,attr"`
	}
	Host struct {
		XMLName xml.Name `xml:"host"`
		Endtime string   `xml:"endtime,attr"`
		Address Address  `xml:"address"`
		Ports   Ports    `xml:"ports>port"`
	}
	Ports []struct {
		Protocol string  `xml:"protocol,attr"`
		Portid   string  `xml:"portid,attr"`
		State    State   `xml:"state"`
		Service  Service `xml:"service"`
	}
	Service struct {
		Name   string `xml:"name,attr"`
		Banner string `xml:"banner,attr"`
	}
	Masscaner struct {
		SystemPath string
		Args       []string
		Ports      string
		Ranges     string
		Rate       string
		Exclude    string
		Result     []byte
	}
)

func (p Ports) ToPortList() (portList []string) {
	for _, port := range p {
		portList = append(portList, port.Portid)
	}
	return
}

func (m *Masscaner) SetSystemPath(systemPath string) {
	if systemPath != "" {
		m.SystemPath = systemPath
	}
}
func (m *Masscaner) SetArgs(arg ...string) {
	m.Args = arg
}
func (m *Masscaner) SetPorts(ports string) {
	m.Ports = ports
}
func (m *Masscaner) SetRanges(ranges string) {
	m.Ranges = ranges
}
func (m *Masscaner) SetRate(rate string) {
	m.Rate = rate
}
func (m *Masscaner) SetExclude(exclude string) {
	m.Exclude = exclude
}

// Start scanning
func (m *Masscaner) Run() error {
	var (
		cmd        *exec.Cmd
		outb, errs bytes.Buffer
	)
	if m.Rate != "" {
		m.Args = append(m.Args, "--rate")
		m.Args = append(m.Args, m.Rate)
	}
	if m.Ranges != "" {
		m.Args = append(m.Args, "--range")
		m.Args = append(m.Args, m.Ranges)
	}
	if m.Ports != "" {
		m.Args = append(m.Args, "-p")
		m.Args = append(m.Args, m.Ports)
	}
	if m.Exclude != "" {
		m.Args = append(m.Args, "--exclude")
		m.Args = append(m.Args, m.Exclude)
	}
	m.Args = append(m.Args, "-oX")
	m.Args = append(m.Args, "-")
	cmd = exec.Command(m.SystemPath, m.Args...)
	fmt.Println(cmd.Args)
	cmd.Stdout = &outb
	cmd.Stderr = &errs
	err := cmd.Run()
	if err != nil {
		if errs.Len() > 0 {
			return errors.New(errs.String())
		}
		return err
	}
	m.Result = outb.Bytes()
	return nil
}

// Parse scans result.
func (m *Masscaner) Parse() ([]Host, error) {
	var hosts []Host
	decoder := xml.NewDecoder(bytes.NewReader(m.Result))
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "host" {
				var host Host
				err := decoder.DecodeElement(&host, &se)
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}
				hosts = append(hosts, host)
			}
		default:
		}
	}
	return hosts, nil
}
func NewMasscaner() *Masscaner {
	return &Masscaner{
		SystemPath: "masscan",
	}
}
