package service

import (
	"context"
	"github.com/carloseduribeiro/working-with-grpc-go/internal/database"
	"github.com/carloseduribeiro/working-with-grpc-go/internal/pb"
	"io"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(db database.Category) CategoryService {
	return CategoryService{CategoryDB: db}
}

func (c CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryListResponse, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		panic(err)
	}
	categoriesResponse := make([]*pb.Category, 0, len(categories))
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
	return &pb.CategoryListResponse{
		Categories: categoriesResponse,
	}, nil
}

func (c CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.FindByID(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	listResponse := &pb.CategoryListResponse{}
	for {
		receivedCategory, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(listResponse)
		}
		if err != nil {
			return err
		}
		createdCategory, err := c.CategoryDB.Create(receivedCategory.Name, receivedCategory.Description)
		if err != nil {
			return err
		}
		listResponse.Categories = append(listResponse.Categories, &pb.Category{
			Id:          createdCategory.ID,
			Name:        createdCategory.Name,
			Description: createdCategory.Description,
		})
	}
}
