package test

import (
	"context"
	pb "go-notes/record/grpc/v1/proto"
	"go-notes/record/grpc/v1/server"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// var articleClient pb.ArticleServiceClient
func NewClient() (pb.ArticleServiceClient, error) {
	ctx := context.Background()
	// 指定自定义认证
	var opts []grpc.DialOption
	// 开启安全的选项
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// token验证
	opts = append(opts, grpc.WithPerRPCCredentials(new(customerTokenAuth)))
	// 连接rpc
	articleCoon, err := grpc.DialContext(ctx, "", opts...)
	if err != nil {
		return nil, err
	}
	// TODO 正常来说rpcClient是一个全局变量，此时为了方便省略
	//defer articleCoon.Close()
	articleClient := pb.NewArticleServiceClient(articleCoon)
	return articleClient, nil
}

type customerTokenAuth struct{}

// 实现interface接口 是否开启传输安全 TLS
func (c customerTokenAuth) RequireTransportSecurity() bool {
	return false
}

// 获取元数据
func (c customerTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  server.AppId,
		"appKey": server.AppKey,
	}, nil
}

func TestCreateArticle(t *testing.T) {
	c, err := NewClient()
	assert.Equal(t, nil, err)

	reply, err := c.CreateArticle(context.Background(), &pb.RequestCreateArticle{
		Title:   "aa",
		Content: "elkjgldekngh",
		Author:  "王五",
		IsShow:  false,
		Type:    pb.Type_LYRICS,
	})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, int64(2), reply.ArticleId)
	assert.Equal(t, int64(1), reply.ArticleId)

}

func TestUpdateArticle(t *testing.T) {
	c, err := NewClient()
	assert.Equal(t, nil, err)
	// failed
	{
		_, err := c.UpdateArticle(context.Background(), &pb.RequestUpdateArticle{
			ArticleId: 2,
			Title:     "ww",
			Content:   "123456",
			Author:    "王五",
			IsShow:    false,
			Type:      pb.Type_NOVEL,
		})
		assert.NotEqual(t, nil, err)
	}
	// pass
	{
		_, err := c.UpdateArticle(context.Background(), &pb.RequestUpdateArticle{
			ArticleId: 1,
			Title:     "ww",
			Content:   "123456",
			Author:    "王五",
			IsShow:    false,
			Type:      pb.Type_NOVEL,
		})
		assert.Equal(t, nil, err)
	}
}

func TestDeleteArticle(t *testing.T) {
	c, err := NewClient()
	assert.Equal(t, nil, err)

	// failed
	{
		_, err := c.DeleteArticle(context.Background(), &pb.RequestDeleteArticle{
			ArticleId: 3,
		})
		assert.NotEqual(t, nil, err)
	}

	// pass
	{
		_, err := c.DeleteArticle(context.Background(), &pb.RequestDeleteArticle{
			ArticleId: 1,
		})
		assert.Equal(t, nil, err)
	}
}

func TestQueryArticle(t *testing.T) {
	c, err := NewClient()
	assert.Equal(t, nil, err)

	// failed
	{
		reply, err := c.QueryArticle(context.Background(), &pb.RequestQueryArticle{
			ArticleId: 3,
		})
		assert.NotEqual(t, nil, err)
		assert.Equal(t, (*pb.ReplyQueryArticle)(nil), reply)
	}
	// pass
	{
		reply, err := c.QueryArticle(context.Background(), &pb.RequestQueryArticle{
			ArticleId: 1,
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, int64(1), reply.ArticleId)
		assert.Equal(t, "张三", reply.Author)
		assert.Equal(t, "让数字文明造福各国人民", reply.Title)
		assert.Equal(t, "字经济是全球未来的发展方向。习主席深刻洞察人类社会发展大势，为我们积极推动数字经济和生产生活深度融合指明了前进方向，也为国际社会共同迈向数字文明新时代贡献了中国方案，必将有力推动构建人类命运共同体。", reply.Content)
		assert.Equal(t, true, reply.IsShow)
		assert.Equal(t, pb.Type_LYRICS, reply.Type)
	}
}

func TestQueryArticleList(t *testing.T) {
	c, err := NewClient()
	assert.Equal(t, nil, err)
	{
		reply, err := c.ArticleList(context.Background(), &empty.Empty{})
		assert.Equal(t, nil, err)
		data := reply.Data
		assert.Equal(t, int64(1), data[0].ArticleId)
		assert.Equal(t, "张三", data[0].Author)
		assert.Equal(t, "让数字文明造福各国人民", data[0].Title)
		assert.Equal(t, "字经济是全球未来的发展方向。习主席深刻洞察人类社会发展大势，为我们积极推动数字经济和生产生活深度融合指明了前进方向，也为国际社会共同迈向数字文明新时代贡献了中国方案，必将有力推动构建人类命运共同体。", data[0].Content)
		assert.Equal(t, false, data[0].IsShow)
		assert.Equal(t, pb.Type_LYRICS, data[0].Type)

		assert.Equal(t, int64(2), data[1].ArticleId)
		assert.Equal(t, "李四", data[1].Author)
		assert.Equal(t, "生产旺季搞拉闸限电咋回事", data[1].Title)
		assert.Equal(t, "近期，多家上市公司却发布公告称，为配合地区“能耗双控”要求限电停产。正值生产旺季，搞拉闸限电是咋回事？", data[1].Content)
		assert.Equal(t, true, data[1].IsShow)
		assert.Equal(t, pb.Type_LYRICS, data[1].Type)
	}
}
