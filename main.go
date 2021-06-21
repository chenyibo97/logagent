package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	config2 "studygo2/logagent/config"
	"studygo2/logagent/etcd"
	"studygo2/logagent/kafka"
	"studygo2/logagent/tail_log"
	"time"
)





func main(){
	//1.初始化kafka连接
	config,err:=config2.ReadConfigFile()
	if err!=nil{
		fmt.Println("read config file fail,err:",err)
	}
	addr:=[]string{config.KafkaAddr}
     err=kafka.Init(addr)
     if err!=nil{
     	fmt.Println("init kafka producer failed,err:",err)
	 }

	//2.打开日志文件准备收集
	//failname:="./my.log"
	/*err=tail_log.Init(config.LogFileName)
	if err != nil {
		fmt.Println("init taillog  failed,err:",err)
	}*/

	//2.0初始化etcd
     err=etcd.Init(config.EtcdAddr,time.Duration(config.EtcdTimeout)*time.Second)
	if err != nil {
		fmt.Println("connect etcd failed,err:",err)
	}

    //2.1 从ETCD获取日志收集项的配置信息
    var logEntryConfig []*etcd.LogEntry
     ipstr,err:=etcd.GetLocalIp()
	if err != nil {
		panic(err)
	}
	etcdConfkey:=fmt.Sprintf(config.EtcdAgentKey,ipstr)
	fmt.Println(etcdConfkey)
     ReadConfigFromEtcd(etcdConfkey,&logEntryConfig)



    //3收集日志发往kafka
     //3.1循环每一个日志配置项，创建tailobj
	//ctx, cancel := context.WithCancel(context.Background())
    CreateTailObj(&logEntryConfig)
	//2.2派哨兵见识日志收集项目的配置信息的变更
	go WatchConfigChange(etcdConfkey)
     //3.2发往kafka


	//读取日志发送到kafka
   // run()
   //fmt.Println("a")


    for{
    /*    select{
        case <-etcd.ConfigHasBeenChange:{
        	fmt.Println("a new config has receive")
             cancel()
			ReadConfigFromEtcd(&logEntryConfig)
			CreateTailObj(ctx,cancel,&logEntryConfig)
		}
		default:
			time.Sleep(5*time.Microsecond)
		}*/
	}
}

func ReadConfigFromEtcd(key string,logEntryConfig *[]*etcd.LogEntry){
	var err error
	*logEntryConfig,err=etcd.GetConfig(key)
	if err != nil {
		fmt.Println("etcd config failed,err： ",err)
	}else{
		fmt.Println("get config from etcd sucess:",err)
	}
}


func CreateTailObj(logEntryConfig *[]*etcd.LogEntry){
	ctx, cancel := context.WithCancel(context.Background())
	for _,v:=range *logEntryConfig{
		tailStruct,err:=tail_log.Init(v.Path)
		tailStruct.TailTaskCancel=cancel
		tailStruct.Ctx=ctx
		tail_log.TailTaskMap[v.Topic]=tailStruct
		go tail_log.Run(v.Topic)
		if err != nil {
			fmt.Println("init taillog  failed,err:",err)
		}
	}
}

func WatchConfigChange(key string){
	ch:=etcd.Cli.Watch(context.Background(),key)
	for wresp:=range ch{
		for  _,ev:=range wresp.Events{
			var newConfig[]*etcd.LogEntry
			err:=json.Unmarshal(ev.Kv.Value,&newConfig)
			if err != nil {
				fmt.Println("json marshal failed,err:",err)
			}
			switch ev.Type {
			case clientv3.EventTypePut:
				for _,v:=range newConfig{
					if v2,ok:=tail_log.TailTaskMap[v.Topic];ok{
						if v2.TailObj.Filename==v.Path{
							continue
						}else{
							v2.TailTaskCancel()
							tail_log.CreatNewtTask(v)
						}

					}else{
						tail_log.CreatNewtTask(v)
					}
				}
			//case clientv3.EventTypeDelete:
			default:
				time.Sleep(50*time.Microsecond)

			}
			//tail_log.NewConfigChan<-
		}
	}
}