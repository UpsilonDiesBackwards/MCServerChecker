package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

type Command string

const (
	CmdOnlinePlayers Command = "getOnlinePlayers"
	CmdServerInfo    Command = "getServerInfo"
	CmdMsgServer     Command = "msgServer"
)

type Packet struct {
	Command Command         `json:"command"`
	Payload json.RawMessage `json:"payload"`
}

type Payload interface{}

type OnlinePlayersPayload struct {
	Payload
	Players []Player `json:"players"`
}

type ServerInfoPayload struct {
	Payload
	ServerInfo ServerInfo `json:"serverInfo"`
}

type MsgServerPayload struct {
	Payload
	Message DiscordMessage `json:"message"`
}

type MsgServerReplyPayload struct {
	Payload
	Success bool `json:"success"`
}

type Player struct {
	Name string `json:"name"`
	Afk  bool   `json:"afk"`
}

type ServerInfo struct {
	Uptime          string  `json:"uptime"`
	Tps             float64 `json:"tps"`
	RamUsage        int64   `json:"ram_usage"`
	ChunksOverworld int64   `json:"chunks_overworld"`
	ChunksNether    int64   `json:"chunks_nether"`
	ChunksEnd       int64   `json:"chunks_end"`
}

type DiscordMessage struct {
	ServerName string `json:"server_name"`
	Channel    string `json:"channel"`
	Username   string `json:"username"`
	Message    string `json:"message"`
}

var conn *net.Conn = nil

func OpenConnection(hostname string, port int) error {
	serverAddr := fmt.Sprintf("%s:%d", hostname, port)
	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}

	conn = &c
	return nil
}

func CloseConnection() {
	if conn != nil {
		_ = (*conn).Close()
	}
	conn = nil
}

func SendAndRecvPacket(command Command, payload Payload) (Payload, error) {
	if conn == nil {
		return nil, errors.New("connection closed")
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	packet := Packet{
		Command: command,
		Payload: payloadData,
	}

	jsonData, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	// Append new line to finish packet
	jsonData = append(jsonData, '\n')

	// Send packet
	writer := bufio.NewWriter(*conn)
	_, err = writer.Write(jsonData)
	if err != nil {
		return nil, err
	}
	_ = writer.Flush()

	// Receive packet
	err = (*conn).SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return nil, err
	}

	var response Packet
	decoder := json.NewDecoder(*conn)

	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	switch response.Command {
	case CmdOnlinePlayers:
		var payload OnlinePlayersPayload
		err = json.Unmarshal(response.Payload, &payload)
		if err != nil {
			return nil, err
		}

		return payload, nil
	case CmdServerInfo:
		var payload ServerInfoPayload
		err = json.Unmarshal(response.Payload, &payload)
		if err != nil {
			return nil, err
		}

		return payload, nil
	case CmdMsgServer:
		var payload MsgServerReplyPayload
		err = json.Unmarshal(response.Payload, &payload)
		if err != nil {
			return nil, err
		}

		return payload, nil
	default:
		return nil, errors.New("unknown command")
	}
}
