package hygge_srv

import (
    "github.com/sk3wlabs/go-serial"
    "github.com/sk3wlabs/go-serial/enumerator"
)

const (
    BridgeVID = "239a"
    BridgePID = "800c"
)

type Server struct {
    sd serial.Port
}

func NewServer() (*Server, error) {
    srv := &Server{}
    err := srv.Init()
    if err != nil {
        return nil, err
    }
    return srv, nil
}

func (s *Server) Close() error {
    return s.sd.Close()
}

func (srv *Server) Init() error {
    dpl, err := enumerator.GetDetailedPortsList()
    if err != nil {
        return err
    }

    mode := serial.Mode{
        BaudRate:          9600,
        DataBits:          8,
        Parity:            serial.NoParity,
        StopBits:          serial.OneStopBit,
        InitialStatusBits: nil,
    }

    for _, d := range dpl {
        if d.IsUSB && d.VID == BridgeVID && d.PID == BridgePID {
            srv.sd, err = serial.Open(d.Name, &serial.Mode{
                BaudRate:          0,
                DataBits:          0,
                Parity:            0,
                StopBits:          0,
                InitialStatusBits: nil,
            })
            if err != nil {
                return err
            }
            break
        }
    }

    return nil
}
