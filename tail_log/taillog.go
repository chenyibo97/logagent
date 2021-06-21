package tail_log

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"studygo2/logagent/etcd"
	"studygo2/logagent/kafka"
	"time"
)
type TailTask struct {
	TailObj *tail.Tail
	Ctx context.Context
	TailTaskCancel func()
}


var (
	 TailTaskMap=make(map[string]*TailTask,16)
     NewConfigChan=make(chan []*etcd.LogEntry,0)
)


func Init(filename string)(tailObj *TailTask,err error){
	tailObj=new(TailTask)
	config:=tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		MustExist: false,
		Poll: true,
	}
	tailObj.TailObj,err=tail.TailFile(filename,config)
	if err!=nil{
		fmt.Println("tail failed:",err)
		return
	}else{
		fmt.Println("tail sucess")
	}
	return
}
func RestartTask(topic string){

}


func CreatNewtTask(v *etcd.LogEntry){
	    ctx, cancel := context.WithCancel(context.Background())
		tailStruct,err:=Init(v.Path)
		tailStruct.TailTaskCancel=cancel
		tailStruct.Ctx=ctx
		TailTaskMap[v.Topic]=tailStruct
		if err != nil {
			fmt.Println("init taillog  failed,err:",err)
		}
		go Run(v.Topic)
}



func Run(topic string){
	fmt.Println("gorutine has start for watching topic:",topic)
	tailStruct:=TailTaskMap[topic]
	//读取日志发送到kafka
	for{
		select {
		case <-tailStruct.Ctx.Done():{
			fmt.Println("topic exit:",topic)
			return
		}
		case line:=<-tailStruct.TailObj.Lines:
			fmt.Println(topic,line.Text)
			kafka.SendToKafka(topic,line.Text)
		default:
			time.Sleep(100*time.Microsecond)
		}
	}
}

/*func run (){
	for{
		select {
		case newConf:=<-newConfigChan:
			fmt.Println("a new config has receive")
		default:
			time.Sleep(100*time.Microsecond)
		}
	}
}*/

/*func ReadChan() <-chan *tail.Line{
	return tailObj.Lines
}*/