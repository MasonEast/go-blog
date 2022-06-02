package consumer

import "fmt"

//传入消费者数量
func Get(consumerNum int) {



	//下面试试多个partition
	for i := 0; i < consumerNum; i++ {
		consumer := new(Consumer)
		err := consumer.InitConsumer()
		if err != nil {
			fmt.Errorf("fail to init consumer, err:%v", err)
			return
		}
		go consumer.GetMessageToAll(1)
	}

}
