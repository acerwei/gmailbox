package central

import (
	"acerwei/gmailbox/encoder"
	"acerwei/gmailbox/gmailbox"
	"acerwei/gmailbox/util"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	gmail "google.golang.org/api/gmail/v1"
)

var (
	gStartDate = "#"
	gEndDate   = "#"
	gPath      string
	service    *gmail.Service
	gEncoder   encoder.Encoder
)

//MessageOption MessageOption
type MessageOption struct {
}

//Get Get
func (f *MessageOption) Get() (key, val string) {
	key = "format"
	val = "full"
	return
}

func getMailsByLabel(wg *sync.WaitGroup, label string) int {
	defer wg.Done()
	pageToken := ""
	seqNo := 0
	for {
		req := service.Users.Messages.List("me").Q(fmt.Sprintf("label:%s after:%s before:%s", label, gStartDate, gEndDate))
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			fmt.Printf("[Warning] Unable to retrieve messages: %v\n", err)
		}

		fmt.Printf("Processing %v messages...\n", len(r.Messages))
		for _, m := range r.Messages {
			msg, err := service.Users.Messages.Get("me", m.Id).Do(&MessageOption{})
			//fmt.Println(msg)
			if err != nil {
				fmt.Printf("Unable to retrieve message %v: %v\n", m.Id, err)
			}
			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("[Warning] unable to marshal message %v", err)
			}
			data = gEncoder.Encode(data)
			fileName := fmt.Sprintf("%s/%s/%d.mail", gPath, label, seqNo)
			err = ioutil.WriteFile(fileName, data, os.ModePerm)
			seqNo++

		}
		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
	return seqNo
}

func listLabels() []*gmail.Label {
	user := "me"
	r, err := service.Users.Labels.List(user).Do()
	if err != nil {
		fmt.Printf("[Warning] Unable to retrieve labels. %v\n", err)
	}
	if len(r.Labels) > 0 {
		fmt.Print("Labels:\n")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n", l.Name)
			folder := fmt.Sprintf("%s/%s", gPath, l.Name)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				fmt.Printf("[Warning] create folder %s error:%v\n", folder, err)
			}
		}
	} else {
		fmt.Print("No labels found.")
	}
	return r.Labels
}

//Initialize Initialize
func Initialize(path, encodeAlgo, encKey string) error {
	gPath = path
	if encodeAlgo == "simple" {
		gEncoder = encoder.NewDeaultSimpleEncoder()
	} else if encodeAlgo == "blowfish" {
		var err error
		gEncoder, err = encoder.NewBlowFishEncoder(encKey)
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported encryption algorithms")
	}
	return nil
}

//Retrieve Retrieve
func Retrieve(startDate, endDate string) {
	//fmt.Printf("startDate=%s, endDate=%s\n", startDate, endDate)
	gStartDate = startDate
	gEndDate = endDate
	service = gmailbox.Authorize()
	os.RemoveAll(gPath)
	labels := listLabels()
	wg := &sync.WaitGroup{}
	wg.Add(len(labels))
	for _, l := range labels {
		go func(label string) {
			count := getMailsByLabel(wg, label)
			fmt.Printf("label=%s, messages=%d\n", label, count)
		}(l.Name)
	}
	wg.Wait()
	fmt.Println("Done.")
}

//DecodeMail decode an email
func DecodeMail(fileToDecode string) {
	data, err := ioutil.ReadFile(fileToDecode)
	if err != nil {
		fmt.Printf("[ERROR] fail to read file %s\n", fileToDecode)
		return
	}
	data = gEncoder.Decode(data)
	msg := &gmail.Message{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		fmt.Printf("[WARNING] fail to decode message %s", fileToDecode)
		return
	}
	fmt.Println(msg)
	for _, part := range msg.Payload.Parts {
		//fmt.Println(part.Body.Data)
		base64Content, err := base64.URLEncoding.DecodeString(part.Body.Data)
		if err != nil {
			fmt.Printf("[WARNING] base64 convert error=%v", err)
		} else {
			fmt.Printf("[Body] %s\n", string(base64Content))
			plainBody := util.Text(base64Content)
			fmt.Printf("[Plain Body] %s\n", plainBody)
		}
	}
}
