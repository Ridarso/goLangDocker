package main

import (
	"context"
	"errors"
	"fmt"

	pb "docker/vessel-service/proto/vessel"

	"github.com/micro/go-micro"
)

type Repository interface {
	FindAvailable(*pb.Spesification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

type service struct {
	repo Repository
}

func (repo *VesselRepository) FindAvailable(spec *pb.Spesification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel found by the spec")
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Spesification, res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{
			Id:        "1000",
			Name:      "Bombb",
			Capacity:  200,
			MaxWeight: 5000,
		},
	}
	repo := &VesselRepository{vessels}

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
