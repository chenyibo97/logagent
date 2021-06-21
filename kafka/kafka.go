package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

//this mod is working to writting log to kafa

var (
	client  sarama.SyncProducer   //声明一个全局的kafka生产者

)
func Init(addr []string)(err error){
	config:=sarama.NewConfig()
	config.Producer.RequiredAcks=sarama.WaitForAll
	config.Producer.Partitioner=sarama.NewRandomPartitioner
	config.Producer.Return.Successes=true
	/*msg:=&sarama.ProducerMessage{}
	msg.Topic="web_log"
	msg.Value=sarama.StringEncoder("this is a test log")*/
	client,err=sarama.NewSyncProducer(addr,config)
	if err!=nil{
		fmt.Println("produce close,err:",err)
		return err
	}
	return nil
}


func SendToKafka(topic string,data string){

	msg:=&sarama.ProducerMessage{}
	msg.Topic=topic
	msg.Value=sarama.StringEncoder(data)

	pid,offset,err:=client.SendMessage(msg)
	if err!=nil{
		fmt.Println("send mes failed,err:",err)
		return
	}
	fmt.Println("pid ,offset:",pid,offset)
}