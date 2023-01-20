package article

import (
	"context"
	"errors"
	"fmt"
	pb "go-notes/record/grpc/v1/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pb.UnimplementedArticleServiceServer
	port int
}

func NewService(port int) *Service {
	return &Service{
		port: port,
	}
}

func (s *Service) CreateArticle(ctx context.Context, req *pb.RequestCreateArticle) (*pb.ReplyCreateArticle, error) {
	// 模拟数据 没有使用数据库
	err := s.CheckType(req.Type)
	if err != nil {
		return nil, err
	}
	return &pb.ReplyCreateArticle{ArticleId: 1}, nil
}

func (s *Service) UpdateArticle(ctx context.Context, req *pb.RequestUpdateArticle) (*emptypb.Empty, error) {
	if err := s.checkId(req.ArticleId); err != nil {
		return nil, err
	}
	err := s.CheckType(req.Type)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteArticle(ctx context.Context, req *pb.RequestDeleteArticle) (*emptypb.Empty, error) {
	if err := s.checkId(req.ArticleId); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) QueryArticle(ctx context.Context, req *pb.RequestQueryArticle) (*pb.ReplyQueryArticle, error) {
	if err := s.checkId(req.ArticleId); err != nil {
		return nil, err
	}
	return &pb.ReplyQueryArticle{
		ArticleId: 1,
		Title:     fmt.Sprintf("local:%d", s.port),
		Content:   "字经济是全球未来的发展方向。习主席深刻洞察人类社会发展大势，为我们积极推动数字经济和生产生活深度融合指明了前进方向，也为国际社会共同迈向数字文明新时代贡献了中国方案，必将有力推动构建人类命运共同体。",
		Author:    "张三",
		IsShow:    true,
		Type:      pb.Type_LYRICS,
	}, nil
}

func (s *Service) ArticleList(ctx context.Context, req *empty.Empty) (*pb.ReplyArticleList, error) {
	reply := &pb.ReplyArticleList{}
	reply.Data = append(reply.Data, &pb.ArticleItem{
		ArticleId: 1,
		Title:     "让数字文明造福各国人民",
		Content:   "字经济是全球未来的发展方向。习主席深刻洞察人类社会发展大势，为我们积极推动数字经济和生产生活深度融合指明了前进方向，也为国际社会共同迈向数字文明新时代贡献了中国方案，必将有力推动构建人类命运共同体。",
		Author:    "张三",
		IsShow:    false,
		Type:      pb.Type_LYRICS,
	})
	reply.Data = append(reply.Data, &pb.ArticleItem{
		ArticleId: 2,
		Title:     "生产旺季搞拉闸限电咋回事",
		Content:   "近期，多家上市公司却发布公告称，为配合地区“能耗双控”要求限电停产。正值生产旺季，搞拉闸限电是咋回事？",
		Author:    "李四",
		IsShow:    true,
		Type:      pb.Type_NOVEL,
	})
	return reply, nil
}
func (s *Service) checkId(articleId int64) error {
	if articleId != 1 {
		return errors.New("articleId not exists")
	}
	return nil
}
func (s *Service) CheckType(articleType pb.Type) error {
	switch articleType {
	case pb.Type_PROSE:
	case pb.Type_LYRICS:
	case pb.Type_NOVEL:
	default:
		return errors.New("service type unknown")
	}
	return nil
}
