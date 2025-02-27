package main

import (
  "errors"
  "fmt"
  "log"
  "net"
  "net/http"
)

func main() {
  http.HandleFunc("/", handler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
  ip, err := externalIP()
  if err != nil {
    fmt.Println(err)
  }
  fmt.Fprintf(w, "%s", ip)
}

func externalIP() (string, error) {
  ifaces, err := net.Interfaces()
  if err != nil {
    return "", err
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
      return "", err
    }
    for _, addr := range addrs {
      var ip net.IP
      switch v := addr.(type) {
      case *net.IPNet:
        ip = v.IP
      case *net.IPAddr:
        ip = v.IP
      }
      if ip == nil || ip.IsLoopback() {
        continue
      }
      ip = ip.To4()
      if ip == nil {
        continue // not an ipv4 address
      }
      return ip.String(), nil
    }
  }
  return "", errors.New("are you connected to the network?")
}
