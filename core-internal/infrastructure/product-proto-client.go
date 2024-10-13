package infrastructure

import (
	"context"
	product_proto "pyre-promotion/core-internal/proto"
	"pyre-promotion/core-internal/utils"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductProtoClientInfra struct {
	ProductProtoClient product_proto.ProductProtoClient
}

func NewProductProtoClientInfra() *ProductProtoClientInfra {
	client, err := grpc.NewClient(utils.GlobalEnv.GRPCProductHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	productProtoClient := product_proto.NewProductProtoClient(client)

	return &ProductProtoClientInfra{
		ProductProtoClient: productProtoClient,
	}
}

func (p *ProductProtoClientInfra) GetProduct(productIds []string) (res *product_proto.GetProductResponse, err error) {
	getProductRequest := &product_proto.GetProductRequest{
		ProductIds: productIds,
	}

	res, err = p.ProductProtoClient.GetProduct(context.Background(), getProductRequest)
	if err != nil {
		return res, err
	}
	return res, nil
}