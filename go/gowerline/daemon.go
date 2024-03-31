package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	ipc "github.com/james-barrow/golang-ipc"
	"github.com/pkg/errors"
)

const (
	socketName = "gowerline"
	daemonFlag = "--daemon"
)

func startClient(_ context.Context, args []string) (*ipc.Client, error) {
	fileName := filepath.Join(os.TempDir(), socketName+".lock")
	err := os.RemoveAll(fileName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to remove pre-existing lock")
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		// TODO if failed due to EXCL, then skip booting server
		return nil, errors.WithMessage(err, "failed to acquire startup lock")
	}

	pathToThisBinary := os.Args[0]
	args = append([]string{daemonFlag}, args...)
	cmd := exec.CommandContext(context.Background(), pathToThisBinary, args...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout // TODO nil these out
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{file}
	err = cmd.Start()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to start server")
	}
	client, err := ipc.StartClient(socketName, nil)
	return client, errors.WithMessage(err, "failed to start client")
}

func runServer(ctx context.Context) error {
	server, err := ipc.StartServer(socketName, nil)
	if err != nil {
		return err
	}
	log.Println("Starting server...", server.Status())
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			server.Close()
			return nil
		default:
			log.Println("Server reading message...")
		}
		var responseType int
		var response []byte
		message, err := server.Read()
		if err == nil {
			responseType, response, err = handleMessage(message)
			log.Println("Server handled message:", responseType, string(response), err)
		}
		if err != nil {
			log.Println(err)
			responseType = messageErr
			response = []byte(err.Error())
		}
		err = server.Write(responseType, response)
		if err != nil {
			log.Println(err)
		}
		if server.StatusCode() == ipc.Closed {
			log.Println("IPC socket closed. Stopping server...")
			return nil
		}
	}
}

const (
	messageErr = 1 + iota
	messageStatus
	messageStatusResponse
)

func handleMessage(message *ipc.Message) (int, []byte, error) {
	switch message.MsgType {
	case messageErr:
		return 0, nil, errors.Errorf("failed processing request: %s", string(message.Data))
	case messageStatus:
		return generateStatus()
	default:
		return 0, nil, errors.Errorf("unexpected message type: %d", message.MsgType)
	}
}

func handleClientMessage(message *ipc.Message) (data []byte, recognizedType bool, err error) {
	switch {
	case message.MsgType < messageErr:
		return nil, false, nil
	case message.MsgType == messageErr:
		return nil, false, errors.Errorf("server failed to process request: %s", string(message.Data))
	default:
		return message.Data, true, nil
	}
}

func generateStatus() (int, []byte, error) {
	var buf bytes.Buffer
	err := status(&buf)
	return messageStatusResponse, buf.Bytes(), err
}
