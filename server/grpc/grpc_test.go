package grpc

import (
	"testing"

	"github.com/micro/go-micro/registry/mock"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/micro/examples/greeter/srv/proto/hello"
)

// server is used to implement helloworld.GreeterServer.
type sayServer struct{}

// SayHello implements helloworld.GreeterServer
func (s *sayServer) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func TestGRPCServer(t *testing.T) {
	r := mock.NewRegistry()
	s := NewServer(
		server.Name("foo"),
		server.Registry(r),
	)

	pb.RegisterSayHandler(s, &sayServer{})

	if err := s.Start(); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	if err := s.Register(); err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	// check registration
	services, err := r.GetService("foo")
	if err != nil || len(services) == 0 {
		t.Fatal("failed to get service: %v # %d", err, len(services))
	}

	defer func() {
		if err := s.Deregister(); err != nil {
			t.Fatalf("failed to deregister: %v", err)
		}

		if err := s.Stop(); err != nil {
			t.Fatalf("failed to stop: %v", err)
		}
	}()

	cc, err := grpc.Dial(s.Options().Address, grpc.WithInsecure())
	if err != nil {
		t.Fatal("failed to dial server: %v", err)
	}

	testMethods := []string{"Say.Hello", "/helloworld.Say/Hello", "/greeter.helloworld.Say/Hello"}

	for _, method := range testMethods {
		rsp := pb.Response{}

		if err := grpc.Invoke(context.Background(), method, &pb.Request{Name: "John"}, &rsp, cc); err != nil {
			t.Fatal("error calling server: %v", err)
		}

		if rsp.Msg != "Hello John" {
			t.Fatalf("Got unexpected response %v", rsp.Msg)
		}
	}
}
