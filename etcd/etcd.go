package etcd

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"net"
	"strings"
	"time"
)

var(
	Cli *clientv3.Client
	ConfigHasBeenChange =make(chan bool)
)

type LogEntry struct {
	Path string `json:"path"`
	Topic string `json:"topic"`
}

func Init(addr string,time time.Duration)(err error){
     Cli,err=clientv3.New(clientv3.Config{
     	Endpoints: []string{addr},
     	DialTimeout: time,
	 })
/*	if err != nil {
		fmt.Println("connect etcd err,err:",err)
		return err
	}*/
	return nil
}

func  GetConfig(key string)(LogEntrysConfig[]*LogEntry, err error) {
	ctx,_:=context.WithTimeout(context.Background(), time.Second)
	resp, err := Cli.Get(ctx, key)
	if err != nil {
		return
	}
	for _,ev:=range resp.Kvs{
		json.Unmarshal(ev.Value,&LogEntrysConfig)
	}

/*	file,err:=os.Open("./etcd/etcdconfig.json")
	if err != nil {
		fmt.Println("open file failed",err)
	}
	defer file.Close()
	decoder:=json.NewDecoder(file)
	err=decoder.Decode(&LogEntrysConfig)
	if err != nil {
		fmt.Println("decode err:",err)
		return nil, err
	}
	fmt.Println(LogEntrysConfig)*/

	if err != nil {
		return
	}
	return
}

func GetLocalIp()(ip string,err error)  {
	conn,err:=net.Dial("udp","8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localaddr:=conn.LocalAddr().(*net.UDPAddr)
	ip=strings.Split(localaddr.IP.String(),":")[0]
	return
}
