package _default

import (
	"encoding/binary"
	"fmt"
	"go-tcp/filet/requests"
	"io"
	"net"
)

type WorkRequest struct {
	info *requests.RequestInfo
	step byte
	text string
}

func init() {
	requests.RegisterRequestType(3, func(reqInfo *requests.RequestInfo) requests.Request { return &WorkRequest{info: reqInfo} })
}

func MakeWorkRequest(step byte, text string) *WorkRequest {
	return &WorkRequest{
		info: &requests.RequestInfo{
			Id:            3,
			WantsResponse: !(step == 1 && text == "no"),
		},
		text: text,
		step: step,
	}
}

func (wr *WorkRequest) Name() string                { return "Work" }
func (wr *WorkRequest) Info() *requests.RequestInfo { return wr.info }
func (wr *WorkRequest) DataSize() uint32            { return uint32(len([]byte(wr.text))) }

func (wr *WorkRequest) SerializeTo(conn *net.Conn) error {

	err := binary.Write(*conn, binary.BigEndian, wr.step)
	if err != nil {
		return fmt.Errorf("error while sending step : %v\n", err)
	}

	data := []byte(wr.text)
	n, err := (*conn).Write(data)
	if n != len(data) {
		return fmt.Errorf("didn't send as much text as I had : %v\n", err)
	}

	return nil
}
func (wr *WorkRequest) DeserializeFrom(conn *net.Conn) (requests.Request, error) {

	text_length := make([]byte, 32/8)
	_, err := (*conn).Read(text_length)
	if err != nil {
		return wr, err
	}

	step_buff := make([]byte, 1)
	_, err = (*conn).Read(step_buff)
	if err != nil {
		return wr, err
	}
	step := step_buff[0]

	data := make([]byte, binary.BigEndian.Uint32(text_length))
	_, err = io.ReadFull(*conn, data)

	if err != nil {
		return wr, err
	}

	wr.text = string(data)
	wr.step = step

	return wr, err
}

func (wr *WorkRequest) GetResult() requests.Request {

	return nil
}

func (wr *WorkRequest) GetText() string {
	return wr.text
}

func (wr *WorkRequest) GetStep() byte {
	return wr.step
}

func (wr *WorkRequest) Answer(text string) *WorkRequest {

	return MakeWorkRequest(wr.step, text)
}
