package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	pb "service/service"

	"log"
	"strconv"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/yaml.v2"
)

type server struct {
	pb.UnimplementedServiceServer
}

func newReverseServer() *server {
	return &server{}
}

type conf struct {
	Rest_port int64 `yaml:"rest_port"`
	Grpc_port int64 `yaml:"grpc_port"`
}

type test struct {
	name           string
	firstquestion  string
	secondquestion string
	thirdquestion  string
}

type answer struct {
	firstanswer  string
	secondanswer string
	thirdanswer  string
}

var psychotest = test{name: "Тест по gRPC", firstquestion: "Какой протокол использует gRPC?", secondquestion: "Формат обмена данными в gRPC?", thirdquestion: "Какие типы API есть в gRPC?"}
var prisontest = test{name: "Тюремный тест", firstquestion: "Вы заходите в хату, на стене нарисованы ворота, а на полу мяч, вас просят забить гол. Что будете делать?", secondquestion: "Вам дают веник и говорят сыграть на гитаре. Ваши действия?", thirdquestion: "Просят сыграть на батарее, как на баяне. Что будешь делать?"}

var prisonanswer1 = answer{firstanswer: "Попросить дать пас", secondanswer: "Ударить по нарисованному мячу", thirdanswer: "Промолчать"}
var prisonanswer2 = answer{firstanswer: "Промолчать", secondanswer: "Отдать веник и попросить настроить", thirdanswer: "Начать играть на импровизированной гитаре"}
var prisonanswer3 = answer{firstanswer: "Попросить раздуть меха.", secondanswer: "Попытаться сыграть на импровизированном баяне", thirdanswer: "Промолчать"}

var psychoanswer1 = answer{firstanswer: "HTTP/2", secondanswer: "HTTP", thirdanswer: "TCP"}
var psychoanswer2 = answer{firstanswer: "JSON", secondanswer: "Protobuf", thirdanswer: "XML"}
var psychoanswer3 = answer{firstanswer: "Унарный", secondanswer: "Бинарный", thirdanswer: "Ни одного"}

func getFileName() string {
	dir, err := os.Open("./configs")
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	filename := files[0].Name()

	return filename
}

func (c *conf) getConf(filename string) *conf {
	// dir, err := os.Open("./configs")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer dir.Close()

	yamlFile, err := ioutil.ReadFile("./configs/" + filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func (s *server) Do(c context.Context, request *pb.Request) (response *pb.Response, err error) {
	tests := map[int]test{1: psychotest, 2: prisontest}
	response = &pb.Response{
		Message: "Please choose a number of a test : " + tests[1].name + " (1) | " + tests[2].name + " (2)",
	}
	return response, nil
}

func (s *server) Get(c context.Context, request *pb.TestRequest) (response *pb.TestResponse, err error) {
	tests := map[int]test{1: psychotest, 2: prisontest}
	prisonanswers := map[int]answer{1: prisonanswer1, 2: prisonanswer2, 3: prisonanswer3}
	psychoanswers := map[int]answer{1: psychoanswer1, 2: psychoanswer2, 3: psychoanswer3}
	numtest := request.Message
	num, err := strconv.Atoi(numtest)
	if err != nil {
		log.Fatal(err)
	}
	if num == 1 {
		response = &pb.TestResponse{
			Message: &pb.Test{Message: "Ответьте на данные вопросы (ответ - одно трехзначное число с ответами) :",
				Firstquestion:  tests[num].firstquestion,
				Answersfirst:   psychoanswers[1].firstanswer + " (1) | " + psychoanswers[1].secondanswer + " (2) | " + psychoanswers[1].thirdanswer + " (3)",
				Secondquestion: tests[num].secondquestion,
				Answerssecond:  psychoanswers[2].firstanswer + " (1) | " + psychoanswers[2].secondanswer + " (2) | " + psychoanswers[2].thirdanswer + " (3)",
				Thirdquestion:  tests[num].thirdquestion,
				Answerthird:    psychoanswers[3].firstanswer + " (1) | " + psychoanswers[3].secondanswer + " (2) | " + psychoanswers[3].thirdanswer + " (3)",
			},
		}
	} else {
		response = &pb.TestResponse{
			Message: &pb.Test{Message: "Ответьте на данные вопросы (ответ - одно трехзначное число с ответами) :",
				Firstquestion:  tests[num].firstquestion,
				Answersfirst:   prisonanswers[1].firstanswer + " (1) | " + prisonanswers[1].secondanswer + " (2) | " + prisonanswers[1].thirdanswer + " (3)",
				Secondquestion: tests[num].secondquestion,
				Answerssecond:  prisonanswers[2].firstanswer + " (1) | " + prisonanswers[2].secondanswer + " (2) | " + prisonanswers[2].thirdanswer + " (3)",
				Thirdquestion:  tests[num].thirdquestion,
				Answerthird:    prisonanswers[3].firstanswer + " (1) | " + prisonanswers[3].secondanswer + " (2) | " + prisonanswers[3].thirdanswer + " (3)",
			},
		}
	}
	return response, nil
}

func (s *server) Answer(c context.Context, request *pb.AnswerRequest) (response *pb.AnswerResponse, err error) {
	number := request.Message
	num, err := strconv.Atoi(number)
	if err != nil {
		log.Fatal(err)
	}
	firstnum := num / 100
	secondnum := (num / 10) % 10
	thirdnum := num % 10
	trueanswers := 0
	if firstnum == 1 {
		trueanswers++
	}
	if secondnum == 2 {
		trueanswers++
	}
	if thirdnum == 2 {
		trueanswers++
	}
	if trueanswers == 3 {
		response = &pb.AnswerResponse{
			Message: "Вы успешно прошли тест.",
		}
	} else if trueanswers == 2 {
		response = &pb.AnswerResponse{
			Message: "Вы прошли тест с одной ошибкой.",
		}
	} else {
		response = &pb.AnswerResponse{
			Message: "Вы не прошли тест!",
		}
	}

	return response, nil
}

func readconf() {

}

func main() {
	filename := getFileName()

	var c conf
	c.getConf(filename)

	numgrpc := strconv.Itoa(int(c.Grpc_port))
	numhttp := strconv.Itoa(int(c.Rest_port))

	lis, err := net.Listen("tcp", ":"+numgrpc)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	defer lis.Close()

	mux := http.NewServeMux()
	mux.Handle("/metric", promhttp.Handler())

	// http.Handle("/metrics", promhttp.Handler())

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	grpc_prometheus.Register(grpcServer)
	reverseServer := newReverseServer()
	pb.RegisterServiceServer(grpcServer, reverseServer)

	log.Println("server started on ", []grpc.ServerOption{}, "port is : ", numgrpc)
	log.Println("server http started on ", numhttp)

	go func() {
		if err := http.ListenAndServe(":8082", promhttp.Handler()); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	log.Fatal(grpcServer.Serve(lis))

}
