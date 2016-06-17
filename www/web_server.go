package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/trace"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	kDomain    string = "localhost:4000"
	kPackage   string = "github.com/asarcar/go_test.www"
	kDomainT   string = "localhost:4001"
	kMaxMsgLen int    = 300
)

var evlog trace.EventLog

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s+"\n")
	tr := trace.New(kDomain, r.URL.Path)
	defer tr.Finish()
	tr.LazyLog(s, true)
	evlog.Printf("StringEvLog: %s", s)
}

func (s String) String() string {
	return "StringLog: " + string(s)
}

type Struct struct {
	Greeting string
	Punct    string
	Who      string
}

func (s Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s%s %s: %s\n", s.Greeting, s.Punct, s.Who, r.URL.Path[1:])
	tr := trace.New(kDomain, r.URL.Path)
	defer tr.Finish()
	tr.LazyPrintf("Struct called with [Greeting-%s Punct-%s Who-%s]",
		s.Greeting, s.Punct, s.Who)
	evlog.Printf("StructEvLog: [Greeting-%s Punct-%s Who-%s]",
		s.Greeting, s.Punct, s.Who)
}

func StopEventLog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "StopEventLog Called\n")
	evlog.Printf("StopEventLog: Terminating Eventlog")
	evlog.Finish()
}

// Execute: Equivalent sinatra code
// post '/payload' do
//   push = JSON.parse(request.body.read)
//   puts "I got some JSON: #{push.inspect}"
// end
func PayloadFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	evlog.Printf("JSON Payload: " + string(b))
}

// Facebook Messenger Platform
// app.get('/webhook/', function (req, res) {
//   if (req.query['hub.verify_token'] === '<validation_token>') {
//	   res.send(req.query['hub.challenge']);
//   }
//   res.send('Error, wrong validation token');
// })
//
// ForaySys Page Access Token:
// curl -ik -X POST "https://graph.facebook.com/v2.6/me/
// subscribed_apps?access_token=EAAHUFq5NmT4BACvJI8ydiMZ
// AXa81N9yvRCymMTMrum6zNjT7iDMDJuQ7nLkaZAxe9APe3fckiy2Y
// bGsVAmF6MUNkkoEJjNaOdWXZAdYaKO31HLCHd72tzTlWd1Gxj0STu
// SgmhHHDZCZB3JrYLClHKEOSNXhz29KKb8uyJpHicvgZDZD"
func WebHookGet(w http.ResponseWriter, r *http.Request) {
	m, _ := url.ParseQuery(r.URL.RawQuery)
	tokenKey := "hub.verify_token"
	expTokenVal := "foraysys_from_ory_and_vishal"
	tokenVal := m.Get(tokenKey)
	challengeKey := "hub.challenge"
	challengeVal := m.Get(challengeKey)
	evlog.Printf("<" + tokenKey + "," + tokenVal + ">" +
		" : <" + challengeKey + "," + challengeVal + ">")

	if tokenVal != expTokenVal {
		fmt.Fprint(w, "Error, wrong validation token")
		return
	}

	fmt.Fprint(w, challengeVal)
}

type FBId struct {
	Id string `json:"id"`
}

type FBMsg struct {
	Text string `json:"text"`
}

type FBMsgElem struct {
	Sender    FBId  `json:"sender"`
	Recipient FBId  `json:"recipient"`
	Message   FBMsg `json:"message"`
}

type FBEntry struct {
	Messaging []FBMsgElem `json:"messaging"`
}

type FBMessage struct {
	Entry []FBEntry `json:"entry"`
}

// JSON: Decode
// {
//   "object": "page",
//   "entry": [
//     {
//       "id": "1025077210943050",
//       "time": 1465106124974,
//       "messaging": [
//         {
//           "sender": {
//             "id": "10154241354427509"
//           },
//           "recipient": {
//             "id": "1025077210943050"
//           },
//           "timestamp": 1465106124942,
//           "message": {
//             "mid": "mid.1465106124935:76fd485e519583f101",
//             "seq": 29,
//             "text": "Trying 3."
//           }
//         }
//       ]
//     }
//   ]
// }

// var token = "<page_access_token>";
//
// function sendTextMessage(sender, text) {
//   messageData = {
//	   text:text
//   }
//   request({
//		 url: 'https://graph.facebook.com/v2.6/me/messages',
//		 qs: {access_token:token},
//		 method: 'POST',
//	   json: {
//  		 recipient: {id:sender},
//  		 message: messageData,
//  	 }
//   }, function(error, response, body) {
//        if (error) {
//          console.log('Error sending message: ', error);
//        } else if (response.body.error) {
//          console.log('Error: ', response.body.error);
//        }
//   });
// }
func SendTextMessage(sender string, text string) {
	evlog.Printf("SendTextMessage: sender %s: text %s", sender, text)
	token := "EAAHUFq5NmT4BAMvk5t8UQkvtL1aalrOjbqFygMas9scw6mFPPbnfHjiqoHR71O5CKmREHDrnpqxPvtJoYXn7kLGhm6Memn31dSjqlHQaA1JVMWRg591Ls6ZCBFBw74JrgoVqDNnLgMzfssSmQVwTe1SJqQR0gdHdZA4gbtWwZDZD"
	urlp, errp := url.Parse("https://graph.facebook.com/v2.6/me/messages")
	// urlp, errp := url.Parse("http://" + kDomainT + "/webhookT")
	if errp != nil {
		panic(errp)
	}

	q := urlp.Query()
	q.Add("access_token", token)
	urlp.RawQuery = q.Encode()
	urls := urlp.String()
	evlog.Printf("URI: %s", urls)

	fbMsgElem := FBMsgElem{
		Recipient: FBId{Id: sender},
		Message:   FBMsg{Text: text},
	}
	jsonBytes, _ := json.Marshal(fbMsgElem)
	evlog.Printf(string(jsonBytes))
	jsonByteReader := bytes.NewReader(jsonBytes)
	req, err := http.NewRequest("POST", urls, jsonByteReader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	evlog.Printf("response Status:", resp.Status)
	evlog.Printf("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	evlog.Printf("response Body:", string(body))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// app.post('/webhook/', function (req, res) {
//   messaging_events = req.body.entry[0].messaging;
//   for (i = 0; i < messaging_events.length; i++) {
//  	 event = req.body.entry[0].messaging[i];
//  	 sender = event.sender.id;
// 	   if (event.message && event.message.text) {
// 		   text = event.message.text;
//       // Handle a text message from this sender
//	   }
//   }
//   res.sendStatus(200);
// });
func WebHookPostHelper(w http.ResponseWriter, r *http.Request, sendmsg bool) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	evlog.Printf("Body: %s", string(b))

	var fbMsg FBMessage
	err = json.Unmarshal(b, &fbMsg)
	if err != nil {
		panic(err)
	}
	defer w.WriteHeader(http.StatusOK)

	numentries := len(fbMsg.Entry)
	evlog.Printf("fbMsg: #Entries %d", numentries)
	if numentries == 0 {
		return
	}
	msgelems := fbMsg.Entry[0].Messaging
	nummsgelems := len(msgelems)
	evlog.Printf("MsgElems: #msgs %d", nummsgelems)
	if nummsgelems == 0 {
		return
	}
	for i := 0; i < nummsgelems; i++ {
		msg := msgelems[i]
		text := msg.Message.Text
		sender := msg.Sender.Id
		evlog.Printf("Msg: Sender=%s, Recipient=%s, Message=\"%s\"",
			sender, msg.Recipient.Id, text)
		if sendmsg {
			SendTextMessage(sender,
				"Text received, echo: "+text[0:min(kMaxMsgLen, len(text))])
		}
	}
}

func WebHookPost(w http.ResponseWriter, r *http.Request) {
	WebHookPostHelper(w, r, true)
}

func WebHookTPost(w http.ResponseWriter, r *http.Request) {
	WebHookPostHelper(w, r, false)
}

func WebHook(w http.ResponseWriter, r *http.Request) {
	evlog.Printf("Method " + r.Method + ": URL " + r.URL.String())
	str := "Headers {"
	for field, attrs := range r.Header {
		str += " <" + field + ":["
		for _, attr := range attrs {
			str += " " + attr
		}
		str += " ]>"
	}
	evlog.Printf(str + "}")
	if r.Method == "GET" {
		WebHookGet(w, r)
		return
	}
	if r.Method == "POST" {
		WebHookPost(w, r)
		return
	}
}

func main() {
	evlog = trace.NewEventLog(kPackage, kDomain)
	http.Handle("/string", String("I'm a frayed knot."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.HandleFunc("/payload", PayloadFunc)
	http.HandleFunc("/webhook", WebHook)
	http.HandleFunc("/stopeventlog", StopEventLog)

	// go func() {
	// 	http.HandleFunc("/webhookT", WebHookTPost)
	// 	log.Fatal(http.ListenAndServe(kDomainT, nil))
	// }()

	log.Fatal(http.ListenAndServe(kDomain, nil))
}
