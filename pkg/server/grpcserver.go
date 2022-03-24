package server

import (
	"challenge/pkg/api"
	"challenge/util"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GRPCServer struct{}

func (s *GRPCServer) MakeShortLink(ctx context.Context, link *api.Link) (*api.Link, error) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	client := &http.Client{}

	var data = fmt.Sprintf("{ \"long_url\": \"%s\", \"domain\": \"bit.ly\" }", link.GetData())
	var body = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://api-ssl.bitly.com/v4/shorten", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.BitlyOauthToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(bodyText, &jsonBody)
	if err != nil {
		log.Fatal(err)
	}

	if shortLink, linkExist := jsonBody["link"].(string); linkExist {
		return &api.Link{Data: shortLink}, nil
	} else {
		err = fmt.Errorf("%v", jsonBody["message"])
		return nil, err
	}
}

func (s *GRPCServer) ReadMetadata(ctx context.Context, placeholder *api.Placeholder) (*api.Placeholder, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	return &api.Placeholder{Data: fmt.Sprintf("%v", md.Get("i-am-random-key"))}, nil
}

type TimersServer struct {
	//map: timer name -> streams of subscribers to this timer
	Data map[string][]api.ChallengeService_StartTimerServer
}

func (s *TimersServer) GetSubscribers(timerName string) []api.ChallengeService_StartTimerServer {
	subs, ok := s.Data[timerName]
	if !ok {
		log.Fatal(fmt.Sprintf("GetSubscriber: timer \"%s\" not running", timerName))
	}
	return subs
}

func (s *TimersServer) AddSubscriber(timerName string, newSub api.ChallengeService_StartTimerServer, client *http.Client) {
	subs, ok := s.Data[timerName]
	if !ok {
		log.Fatal(fmt.Sprintf("AddSubscriber: timer \"%s\" not running", timerName))
	}
	s.Data[timerName] = append(subs, newSub)
	lifeTimeTimer := time.NewTimer(time.Second * time.Duration(TS.CheckTimer(timerName, client)))
	<-lifeTimeTimer.C
}

func (s *TimersServer) CheckTimer(timerName string, client *http.Client) int {
	_, exist := s.Data[timerName]
	if !exist {
		return 0
	}

	req, err := http.NewRequest("GET", "https://timercheck.io"+"/"+timerName, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	var jsonBody map[string]interface{}
	err = json.Unmarshal(bodyText, &jsonBody)
	if err != nil {
		log.Fatal(err)
	}

	if secondsRemaining, timerExist := jsonBody["seconds_remaining"].(float64); timerExist {
		return int(secondsRemaining)
	} else {
		return 0
	}
}

func (s *TimersServer) RunTimer(timer *api.Timer, client *http.Client) {
	for {
		secondsRemaining := s.CheckTimer(timer.GetName(), client)
		if secondsRemaining == 0 {
			break
		}
		timer.Seconds = int64(secondsRemaining)
		for _, sub := range s.GetSubscribers(timer.GetName()) {
			err := sub.Send(timer)
			if err != nil {
				log.Fatal(err)
			}
		}
		frequencyTimer := time.NewTimer(time.Second * time.Duration(timer.GetFrequency()))
		<-frequencyTimer.C
	}

	delete(s.Data, timer.GetName())
}

var TS TimersServer

func (s *GRPCServer) StartTimer(timer *api.Timer, server api.ChallengeService_StartTimerServer) error {
	client := &http.Client{}
	if _, exist := TS.Data[timer.GetName()]; !exist {
		if TS.Data == nil {
			TS.Data = make(map[string][]api.ChallengeService_StartTimerServer)
		}

		TS.Data[timer.GetName()] = make([]api.ChallengeService_StartTimerServer, 0, 1)

		req, err := http.NewRequest("GET", "https://timercheck.io"+"/"+timer.GetName()+"/"+strconv.Itoa(int(timer.GetSeconds())), nil)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		go TS.RunTimer(timer, client)
	}
	TS.AddSubscriber(timer.GetName(), server, client)

	return nil
}
