package main

import (
	"context"
	"fmt"
	pb "github.com/rainset/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewShortenerClient(conn)

	// функция, в которой будем отправлять сообщения
	TestRequests(c)

}

func TestRequests(c pb.ShortenerClient) {

	// add link
	resp1, err := c.AddURL(context.Background(), &pb.AddURLRequest{
		Url: "https://yandex.ru/test/",
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp1.Error != "" {
		fmt.Println(resp1.Error)
	}

	log.Println("")
	log.Println("Результат AddURL:", resp1.Result)

	strSlice := strings.Split(resp1.Result, "/")
	hash := strSlice[len(strSlice)-1]

	resp2, err := c.GetURL(context.Background(), &pb.GetURLRequest{
		Hash: hash,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp2.Error != "" {
		fmt.Println(resp2.Error)
	}

	log.Println("")
	log.Println("Результат GetURL:", resp2.Result)

	resp3, err := c.Stats(context.Background(), &pb.StatsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	if resp3.Error != "" {
		fmt.Println(resp3.Error)
	}

	log.Println("")
	log.Println("Результат Stats:", resp3)

	resp4, err := c.AddBatchURL(context.Background(), &pb.AddBatchURLRequest{
		Urls: []*pb.BatchUrlRequest{
			{
				Correlation_ID: "1",
				OriginalUrl:    "http://example.com/1/",
			},
			{
				Correlation_ID: "2",
				OriginalUrl:    "http://example.com/2/",
			},
			{
				Correlation_ID: "3",
				OriginalUrl:    "http://example.com/3/",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp4.Error != "" {
		fmt.Println(resp4.Error)
	}

	log.Println("")
	log.Println("Результат AddBatchURL:", resp4.Urls)
}
