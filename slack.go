package main

import (
   "time"
   "log"
   "os"
   "io/ioutil"
   "encoding/xml"
   "encoding/json"
   "bytes"
   "strings"
   "strconv"
   "regexp"
   "os/exec"
   "net/http"
   "github.com/tidwall/gjson"
)

func auth(p string) bool { /* only serve request from slackbot */
    r, _ := regexp.Compile("Slackbot")
    if r.MatchString(p) == false {return false}
    return true
}
func check(p string) bool { /* security measure for harmful/unwanted command */
    r1, _ := regexp.Compile("(?i)shutdown"); r2, _ := regexp.Compile("(?i)reboot"); r3, _ := regexp.Compile("(?i)stop"); r4, _ := regexp.Compile("(?i)kill"); r5, _ := regexp.Compile("(?i)rm"); r6, _ := regexp.Compile("(?i)mv"); r7, _ := regexp.Compile("(?i)yum"); r8, _ := regexp.Compile("(?i)apt")
    if r1.MatchString(p) == true || r2.MatchString(p) == true || r3.MatchString(p) == true || r4.MatchString(p) == true || r5.MatchString(p) == true || r6.MatchString(p) == true || r7.MatchString(p) == true || r8.MatchString(p) == true {return false}
    return true
}
func check2(p1 string, p2 string) bool { /* ensure date/time abides with RFC3339 */
    format := "2006/01/02 03:04"
    if _, err1 := time.Parse(format, p1); err1 != nil {panic(err1)}; if _, err2 := time.Parse(format, p2); err2 != nil {panic(err2)}
    return true
}
func get(url string, p1 string, p2 string, p3 []byte) []byte { /* retrieve the target data from target platform intended */
    var req *http.Request; var body []byte /* handle differently for each method API call */
    if p2 == "GET" {log.Println("URL:>"+url)
      req, _ = http.NewRequest(p2, url, nil)
      req.Header.Set("Authorization", p1)
      if len(p3) != 0 {req.Header.Set("Accept", "application/*+json;version=30.0")}}
    if p2 == "POST" {log.Println("URL:>"+url)
      req, _ = http.NewRequest(p2, url, nil)
      req.Header.Set("Authorization", p1)
      req.Header.Set("Accept", "application/*+xml;version=30.0")
      req.Header.Set("Content-Type", "application/json")}
    if p2 == "PROPFIND" {log.Println("URL:>"+url)
      req, _ = http.NewRequest(p2, url, bytes.NewBuffer(p3))
      req.Header.Set("Authorization", p1)
      req.Header.Set("Content-Type", "application/xml")}

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {panic(err);}

    if p2 == "POST" {
      body = []byte(resp.Header.Get("X-Vmware-Vcloud-Access-Token"))
    } else {
      body, _ = ioutil.ReadAll(resp.Body)
      defer resp.Body.Close()}
    return body
}
func post(body string) { /* post parsed data to slack platform as collaboration platform */
    var pay = []byte(`{
    "channel": "G01BUA4K03G",
    "username": "السلام عليكم",
    "icon_emoji": ":exclamation:",
    "attachments": [
	{
        "text": "
`)
    var load = []byte(`",
	"color": "#0000FF"
	}]
}`)
    var payload = append(pay, append([]byte(body), load...)...)
    req, _ := http.NewRequest("POST", "https://hooks.slack.com/services/$TENANT/$CHANNEL/$WEBHOOK", bytes.NewBuffer(payload))
    req.Header.Set("Content-type", "application/json")

    client := &http.Client{}
    _, err := client.Do(req)
    if err != nil {panic(err);}
}
func dev(w http.ResponseWriter, r *http.Request) { /* remote command shell, excluding continouos stdout & excessive stdout */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if auth(r.Header.Get("User-Agent")) && check(r.FormValue("text")) { /* parsing */
        result1, err1 := exec.Command("bash", "-c", r.FormValue("text")).Output(); if err1 != nil {panic(err1)}
        w.WriteHeader(http.StatusOK); log.Println("CMD:>"+r.FormValue("text"))
        w.Write(append([]byte("[root@DEV-SVR /root/go/src/test]# "), append([]byte("bash -c "+r.FormValue("text")+"\n"), []byte(result1)...)...))
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}
func rmt1(w http.ResponseWriter, r *http.Request) { /* remote command shell, excluding continouos stdout & excessive stdout */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if auth(r.Header.Get("User-Agent")) && check(r.FormValue("text")) { /* parsing */
        result1, err1 := exec.Command("bash", "-c", r.FormValue("text")).Output(); if err1 != nil {panic(err1)}
        w.WriteHeader(http.StatusOK); log.Println(r.FormValue("user_id")+":>"+r.FormValue("text"))
        w.Write(append([]byte("[root@JUMPHOST /root/go/src/slack]# "), append([]byte("bash -c "+r.FormValue("text")+"\n"), []byte(result1)...)...))
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}
func rmt2(w http.ResponseWriter, r *http.Request) { /* remote command shell, excluding continouos stdout & excessive stdout */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if auth(r.Header.Get("User-Agent")) && check(r.FormValue("text")) { /* parsing */
        result1, err1 := exec.Command("expect", "-f", "/root/go/src/test/ssh.sh", r.FormValue("text")).Output(); if err1 != nil {panic(err1)}
        w.WriteHeader(http.StatusOK); log.Println(r.FormValue("user_id")+":>"+r.FormValue("text"))
        w.Write(append([]byte("[root@JUMPHOST /root/go/src/slack]# "), append([]byte("expect -f /root/go/src/test/ssh.sh "+r.FormValue("text")+"\n"), []byte(result1)...)...))
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}
func rmt3(w http.ResponseWriter, r *http.Request) { /* query zabbix graph based on item id parameter */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if _, err := strconv.Atoi(r.FormValue("text")); auth(r.Header.Get("User-Agent")) && err == nil { /* parsing */
        result1, err1 := exec.Command("bash", "/root/go/src/test/graph.sh", "-i", r.FormValue("text"), "-u", r.FormValue("user_id")).Output(); if err1 != nil {panic(err1)}
        w.WriteHeader(http.StatusOK); log.Println(r.FormValue("user_id")+":>"+r.FormValue("text"))
        w.Write(append([]byte("[root@JUMPHOST /root/go/src/slack]# "), append([]byte("bash /root/go/src/test/graph.sh -i "+r.FormValue("text")+"\n"), []byte(result1)...)...))
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}
func jira(w http.ResponseWriter, r *http.Request) { /* query jira ticket based on date/time parameter */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if p := strings.Split(r.FormValue("text"), "-"); auth(r.Header.Get("User-Agent")) && check2(p[1], p[2]) { /* parsing */
        start := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(p[1], ":", "%3A"), ":", "%3A"), "/", "%2F"), " ", "%20");end := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(p[2], ":", "%3A"), ":", "%3A"), "/", "%2F"), " ", "%20")
        r2, _ := regexp.Compile("(?i)ID"); r3, _ := regexp.Compile("(?i)IF"); r4, _ := regexp.Compile("(?i)EU"); r5, _ := regexp.Compile("(?i)SG"); var pid, url1 string
        if r2.MatchString(p[0]) == true {pid = "IDMS"}
        if r3.MatchString(p[0]) == true {pid = "IFMS"}
        if r4.MatchString(p[0]) == true {pid = "EUGIO"}
        if r5.MatchString(p[0]) == true {pid = "SGGIO"}
        url1 = "https://$JIRADOMAIN.atlassian.net/rest/api/2/search?jql=project%20%3D%20"+pid+"%20AND%20createdDate%20%3E%3D%20%22"+start+"%22%20AND%20createdDate%20%3C%3D%20%22"+end+"%22"
        body := get(url1, "Basic bS5naGlmZmFyaUBasdasdasdzxczxczxczxc", "GET", []byte(``))
        summ := "issues.#.fields.summary";summary := gjson.GetBytes(body, summ)
        var result1 = pid+" tickets "+p[1]+" - "+p[2]+":\n"; for _, v := range summary.Array() {result1 += v.Str+"\n"}
//        post("<@"+r.FormValue("user_id")+">\n"+result1)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(result1)); log.Println("CMD:>"+r.FormValue("text"))
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}
func tv(w http.ResponseWriter, r *http.Request) { /* list teamvault files */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if auth(r.Header.Get("User-Agent")) && r.FormValue("text") != "" { /* parsing */
        url1 := "https://$NEXTCLOUDDOMAIN.iiji.id/remote.php/dav/files/m.ghiffari/"+r.FormValue("text")
        body := get(url1, "Basic bS5naGlmZmasdasdasdzxczxczxc", "PROPFIND", []byte(`<?xml version="1.0" encoding="UTF-8"?><d:propfind xmlns:d="DAV:"><d:prop xmlns:oc="http://owncloud.org/ns"><d:getlastmodified/><oc:size/><d:getcontenttype/></d:prop></d:propfind>`))
       type Prop struct { /* define teamvault xml to parse the received data */
       XMLName     xml.Name `xml:"prop"`
       Lastmod     string   `xml:"getlastmodified"`
       Length      string    `xml:"size"`}
       type Props struct {
       XMLName     xml.Name `xml:"propstat"`
       Prop     Prop   `xml:"prop"`}
       type Response struct {
       XMLName     xml.Name `xml:"response"`
       Href     string   `xml:"href"`
       Props     Props  `xml:"propstat"`}
       type Xml struct {
       XMLName     xml.Name `xml:"multistatus"`
       Resp     []Response  `xml:"response"`}
        result0 := Xml{}; if err := xml.Unmarshal(body, &result0); err != nil {panic(err)}
        var result1 string; for k, v := range result0.Resp {
        if k == 0 {result1 += r.FormValue("text")+":\n"} else
        {result1 += strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(v.Href,"/remote.php/dav/files/m.ghiffari/",""),"%20"," "),"%23","#"),"%26","&")+" | "+v.Props.Prop.Length+" | "+v.Props.Prop.Lastmod+"\n"}}
//        post("<@"+r.FormValue("user_id")+">\n"+result1)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(result1)); log.Println("CMD:>"+r.FormValue("text"))
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}
func gio(w http.ResponseWriter, r *http.Request) { /* list GIOv2 VM */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if p := strings.Split(r.FormValue("text"), "-"); auth(r.Header.Get("User-Agent")) && p[0] != "" && p[1] != "" { /* parsing */
        url1 := "https://vcloud.biznetgiocloud.com/api/sessions"; var uid string
        r2, _ := regexp.Compile("(?i)AF[0-9]*04|mon"); r3, _ := regexp.Compile("(?i)AF[0-9]*24|iij"); r4, _ := regexp.Compile("(?i)AF[0-9]*56|client1"); r5, _ := regexp.Compile("(?i)AF[0-9]*65|client2"); r6, _ := regexp.Compile("(?i)AF[0-9]*96|client3")
        if r2.MatchString(p[0]) == true {uid = "bS5naGlmZmFyaUBBRasdasdasdaszxczxczxc"}
        if r3.MatchString(p[0]) == true {uid = "bS5naGlmZmFyaUBBRzxczxczxczxcasdasdasd"}
        if r4.MatchString(p[0]) == true {uid = "b3BlcmF0aW9uLXVzZXasdasdassdazxczxczxczx"}
        if r5.MatchString(p[0]) == true {uid = "YnBkcC1jbG91ZEBpazxczxczxcasdasdasdasd"}
        if r6.MatchString(p[0]) == true {uid = "dGVtcG9yYXJ5LXVzZasdasdasdzxczxczxcxzc"}
        body1 := get(url1, "Basic "+uid, "POST", []byte(``))
        url2 := "https://vcloud.biznetgiocloud.com/api/vms/query?type=vm&page="+p[1]+"&pageSize=25&format=records"
        body2 := get(url2, "Bearer "+string(body1), "GET", []byte(`gio2`))
       type VM struct { /* define GIOv2 JSON to parse the received data */
       Name     string   `json:"name"`
       Mem     int   `json:"memoryMB"`
       CPU     int   `json:"numberOfCpus"`
       IP     string   `json:"ipAddress"`
       Net     string   `json:"networkname"`
       Stat     string   `json:"status"`
       NG     string  `json:"catalogName"`}
       type Json struct {
       VM     []VM  `json:"record"`}
        var result0 Json; if err := json.Unmarshal(body2, &result0); err != nil {panic(err)}
        var result1 = p[0]+" page "+p[1]+":\n"; for _, v := range result0.VM {
        if len(v.NG) < 5 { result1 += v.Name+" | "+strconv.Itoa(v.Mem)+" | "+strconv.Itoa(v.CPU)+" | "+v.IP+" | "+v.Net+" | "+v.Stat+"\n"}}
//        post("<@"+r.FormValue("user_id")+">\n"+result1)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(result1)); log.Println("CMD:>"+r.FormValue("text"))
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}
func test(w http.ResponseWriter, r *http.Request) { /* man page */
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
      r.ParseForm()
      if auth(r.Header.Get("User-Agent")) && r.FormValue("text") != "" { /* parsing */
        result1, _ := ioutil.ReadFile("/root/go/src/test/usage.txt")
//        post("<@"+r.FormValue("user_id")+">\n"+result1)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(result1)); log.Println("CMD:>usage")
      } else { w.WriteHeader(http.StatusBadRequest); }
    default:
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte(`{"text": "not found"}`))
    }
}

func main() { /* listen each port & endpoint and serve the corresponding content */
    f, err := os.OpenFile("/var/log/a/access.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {panic(err)}; defer f.Close()
    log.SetOutput(f)
    http.HandleFunc("/远注", dev)
    http.HandleFunc("/远门", rmt1)
    http.HandleFunc("/远代", rmt2)
    http.HandleFunc("/杂比", rmt3)
    http.HandleFunc("/询问", jira)
    http.HandleFunc("/文件", tv)
    http.HandleFunc("/虚拟", gio)
    http.HandleFunc("/试试", test)
    http.ListenAndServe(":10485", nil)
}
