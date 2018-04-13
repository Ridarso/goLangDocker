package main

import (
	/*
		"context"
	*/

	"fmt"
	"log"

	// micro "consignment-service/go-micro"
	pb "docker/consignment-service/proto/consignment"
	vesselProto "docker/vessel-service/proto/vessel"

	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
	/*
		pb "consignment-service/proto/consignment"
		"golang.org/x/net/context"
		"google.golang.org/grpc"
		"google.golang.org/grpc/reflection"
	*/)

// const (
// 	port = ":50051"
// )

// type IRespository interface {
// 	Create(*pb.Consignment) (*pb.Consignment, error)
// 	GetAll() []*pb.Consignment
// }

type Respository interface {
	// consignments []*pb.Consignment
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	// repo IRespository
	repo         Respository
	vesselClient vesselProto.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Spesification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vesse; : %s", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}
	req.VesselId = vesselResponse.Vessel.Id

	//save consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// return &pb.Response{Created: true, Consignment: consignment}, nil
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}

	/*lis, err := net.Listen("tcp", port)
	if err != nil{
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &service{repo})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil{
		log.Fatalf("Failed to serve: %v", err)
	}
	*/

	//Create New Serivce
	srv := micro.NewService(
		//this name must match with package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("lastest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	//INIT will parse the command line flags
	srv.Init()

	//REGISTER SERVER
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	//Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
