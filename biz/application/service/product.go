package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/bytedance/sonic"

	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req *core_api.CreateProductReq) (resp *core_api.CreateProductResp, err error)
	UpdateProduct(ctx context.Context, req *core_api.UpdateProductReq) (resp *core_api.UpdateProductResp, err error)
	DeleteProduct(ctx context.Context, req *core_api.DeleteProductReq) (resp *core_api.DeleteProductResp, err error)
	GetProducts(ctx context.Context, c *core_api.GetProductsReq) (*core_api.GetProductsResp, error)
	GetProduct(ctx context.Context, c *core_api.GetProductReq) (*core_api.GetProductResp, error)
}

var ProductServiceSet = wire.NewSet(
	wire.Struct(new(ProductService), "*"),
	wire.Bind(new(IProductService), new(*ProductService)),
)

type ProductService struct {
	Config               *config.Config
	CloudMindContent     cloudmind_content.ICloudMindContent
	ProductDomainService service.IProductDomainService
	CloudMindTrade       cloudmind_trade.ICloudMindTrade
	CreateItemKq         *kq.CreateItemKq
	UpdateItemKq         *kq.UpdateItemKq
	DeleteItemKq         *kq.DeleteItemKq
}

func (s *ProductService) CreateProduct(ctx context.Context, req *core_api.CreateProductReq) (resp *core_api.CreateProductResp, err error) {
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	switch req.Type {
	case core_api.Product_Type_Flow_Type:
		updateBalanceResp, err := s.CloudMindTrade.UpdateBalance(ctx, &trade.UpdateBalanceReq{
			UserId: user.UserId,
			Flow:   lo.ToPtr(-req.ProductSize),
		})
		if err != nil {
			return resp, err
		}
		if !updateBalanceResp.Ok {
			return resp, consts.ErrFlowNotEnough
		}
	case core_api.Product_Type_Memory_Type:
		updateBalanceResp, err := s.CloudMindTrade.UpdateBalance(ctx, &trade.UpdateBalanceReq{
			UserId: user.UserId,
			Memory: lo.ToPtr(-req.ProductSize),
		})
		if err != nil {
			return resp, err
		}
		if !updateBalanceResp.Ok {
			return resp, consts.ErrFlowNotEnough
		}
	case core_api.Product_Type_File_Type:
		getFileResp, err := s.CloudMindContent.GetFile(ctx, &content.GetFileReq{
			FileId: req.ObjectId,
		})
		if err != nil {
			return resp, err
		}
		if getFileResp.File.UserId != user.UserId {
			return resp, consts.ErrForbidden
		}
	}

	createProductResp, err := s.CloudMindContent.CreateProduct(ctx, &content.CreateProductReq{
		UserId:      user.UserId,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Urls:        req.Urls,
		Tags:        req.Tags,
		Type:        int64(req.Type),
		Price:       req.Price,
		ProductSize: req.ProductSize,
		ObjectId:    req.ObjectId,
	})
	if err != nil {
		return resp, err
	}
	if _, err = s.CloudMindTrade.AddStock(ctx, &trade.AddStockReq{
		ProductId: createProductResp.ProductId,
		Amount:    1,
	}); err != nil {
		return resp, err
	}

	data, _ := sonic.Marshal(&message.CreateItemMessage{
		ItemId:   createProductResp.ProductId,
		IsHidden: req.Status == int64(core_api.ProductStatus_PrivateProductStatus),
		Labels:   req.Tags,
		Category: core_api.Category_name[int32(core_api.Category_ProductCategory)],
	})
	if err = s.CreateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *core_api.UpdateProductReq) (resp *core_api.UpdateProductResp, err error) {
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	product, err := s.CloudMindContent.GetProduct(ctx, &content.GetProductReq{
		ProductId: req.ProductId,
	})
	if err != nil {
		return resp, err
	}
	if product.UserId != user.UserId {
		return resp, consts.ErrForbidden
	}

	if req.Stock != nil || req.ProductSize != 0 {
		getStock, err := s.CloudMindTrade.GetStock(ctx, &trade.GetStockReq{
			ProductId: req.ProductId,
		})
		if err != nil {
			return resp, err
		}

		// 扣减库存
		if req.Stock != nil {
			if _, err = s.CloudMindTrade.AddStock(ctx, &trade.AddStockReq{
				ProductId: req.ProductId,
				Amount:    req.GetStock() - getStock.Stock,
			}); err != nil {
				return resp, err
			}
		} else {
			req.Stock = lo.ToPtr(getStock.Stock)
		}

		// 扣减流量或内存
		switch product.Type {
		case int64(core_api.Product_Type_Flow_Type):
			updateBalanceResp, err := s.CloudMindTrade.UpdateBalance(ctx, &trade.UpdateBalanceReq{
				UserId: user.UserId,
				Flow:   lo.ToPtr(-(req.ProductSize - product.ProductSize) * (req.GetStock() - getStock.Stock)),
			})
			if err != nil {
				return resp, err
			}
			if !updateBalanceResp.Ok {
				return resp, consts.ErrFlowNotEnough
			}
		case int64(core_api.Product_Type_Memory_Type):
			updateBalanceResp, err := s.CloudMindTrade.UpdateBalance(ctx, &trade.UpdateBalanceReq{
				UserId: user.UserId,
				Memory: lo.ToPtr(-(req.ProductSize - product.ProductSize) * (req.GetStock() - getStock.Stock)),
			})
			if err != nil {
				return resp, err
			}
			if !updateBalanceResp.Ok {
				return resp, consts.ErrFlowNotEnough
			}
		}
	}

	if _, err = s.CloudMindContent.UpdateProduct(ctx, &content.UpdateProductReq{
		ProductId:   req.ProductId,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Urls:        req.Urls,
		Tags:        req.Tags,
		Price:       req.Price,
		ProductSize: req.ProductSize,
	}); err != nil {
		return resp, err
	}

	if req.Status != 0 || req.Tags != nil {
		var isHidden *bool
		if req.Status != 0 {
			isHidden = lo.ToPtr(req.Status == int64(core_api.ProductStatus_PrivateProductStatus))
		}

		data, _ := sonic.Marshal(&message.UpdateItemMessage{
			ItemId:   req.ProductId,
			IsHidden: isHidden,
			Labels:   req.Tags,
		})
		if err = s.UpdateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *core_api.DeleteProductReq) (resp *core_api.DeleteProductResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	// 只能删除自己的帖子
	if err = s.CheckIsMyProduct(ctx, req.ProductId, userData.UserId); err != nil {
		return resp, err
	}

	if _, err = s.CloudMindContent.DeleteProduct(ctx, &content.DeleteProductReq{
		ProductId: req.ProductId,
	}); err != nil {
		return resp, err
	}

	data, _ := sonic.Marshal(&message.DeleteItemMessage{
		ItemId: req.ProductId,
	})
	if err = s.DeleteItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *core_api.GetProductReq) (resp *core_api.GetProductResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	var res *content.GetProductResp
	if res, err = s.CloudMindContent.GetProduct(ctx, &content.GetProductReq{
		ProductId: req.ProductId,
	}); err != nil {
		return resp, err
	}

	// 如果该帖子非公开，并且不是他的，那么他没有权限查看
	if res.Status != int64(core_api.ProductStatus_PublicProductStatus) && res.UserId != userData.GetUserId() {
		return resp, consts.ErrForbidden
	}

	resp = &core_api.GetProductResp{
		Name:        res.Name,
		Description: res.Description,
		Urls:        res.Urls,
		Tags:        res.Tags,
		Type:        res.Type,
		Price:       res.Price,
		ProductSize: res.ProductSize,
		User:        &core_api.User{},
		CreateTime:  res.CreateTime,
	}
	if err = mr.Finish(func() error {
		s.ProductDomainService.LoadAuthor(ctx, resp.User, res.UserId) // 作者
		return nil
	}, func() error {
		s.ProductDomainService.LoadCollectCount(ctx, &resp.CollectCount, req.ProductId) // 收藏量
		return nil
	}, func() error {
		s.ProductDomainService.LoadCollected(ctx, &resp.Collected, userData.GetUserId(), req.ProductId) // 是否收藏
		return nil
	}, func() error {
		s.ProductDomainService.LoadPurchaseCount(ctx, &resp.PurchaseCount, req.ProductId) // 购买量
		return nil
	}, func() error {
		s.ProductDomainService.LoadStock(ctx, &resp.Stock, req.ProductId) // 购买量
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ProductService) GetProducts(ctx context.Context, req *core_api.GetProductsReq) (resp *core_api.GetProductsResp, err error) {
	resp = new(core_api.GetProductsResp)
	userData := adaptor.ExtractUserMeta(ctx)
	var (
		getProductsResp *content.GetProductsResp
		searchOptions   *content.SearchOptions
	)

	if req.AllFieldsKey != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_AllFieldsKey{
				AllFieldsKey: *req.AllFieldsKey,
			},
		}
	}
	if req.Name != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_MultiFieldsKey{
				MultiFieldsKey: &content.SearchField{
					Name: req.Name,
				},
			},
		}
	}

	filter := &content.ProductFilterOptions{
		OnlyUserId:      req.OnlyUserId,
		OnlyTags:        req.OnlyTags,
		OnlySetRelation: req.OnlySetRelation,
	}

	if req.OnlyUserId == nil || req.GetOnlyUserId() != userData.GetUserId() {
		filter.OnlyStatus = lo.ToPtr(int64(core_api.ProductStatus_PublicProductStatus))
	}

	if getProductsResp, err = s.CloudMindContent.GetProducts(ctx, &content.GetProductsReq{
		SearchOptions:        searchOptions,
		ProductFilterOptions: filter,
		PaginationOptions: &basic.PaginationOptions{
			Limit:     req.Limit,
			LastToken: req.LastToken,
			Backward:  req.Backward,
			Offset:    req.Offset,
		},
	}); err != nil {
		return resp, err
	}

	resp.Products = make([]*core_api.Product, len(getProductsResp.Products))
	if err = mr.Finish(lo.Map(getProductsResp.Products, func(item *content.Product, i int) func() error {
		return func() error {
			resp.Products[i] = &core_api.Product{
				ProductId:   item.ProductId,
				Name:        item.Name,
				Description: item.Description,
				Url:         item.Urls[0],
				Tags:        item.Tags,
				Type:        item.Type,
				Price:       item.Price,
				ProductSize: item.ProductSize,
				User:        &core_api.User{},
				CreateTime:  item.CreateTime,
			}
			if err = mr.Finish(func() error {
				s.ProductDomainService.LoadCollectCount(ctx, &resp.Products[i].CollectCount, item.ProductId) // 收藏量
				return nil
			}, func() error {
				s.ProductDomainService.LoadPurchaseCount(ctx, &resp.Products[i].CollectCount, item.ProductId) // 购买量
				return nil
			}, func() error {
				s.ProductDomainService.LoadAuthor(ctx, resp.Products[i].User, item.UserId)
				return nil
			}); err != nil {
				return err
			}
			return nil
		}
	})...); err != nil {
		return resp, err
	}
	resp.Total = getProductsResp.Total
	resp.Token = getProductsResp.Token
	return resp, nil
}

func (s *ProductService) CheckIsMyProduct(ctx context.Context, productId, userId string) (err error) {
	product, err := s.CloudMindContent.GetProduct(ctx, &content.GetProductReq{
		ProductId: productId,
	})
	if err != nil {
		return err
	}
	if product.UserId != userId {
		return consts.ErrForbidden
	}
	return nil
}
