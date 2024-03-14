package internalgrpc

import (
	"context"
	"time"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func getLoggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		responseStatus := status.New(codes.OK, "OK")
		if err != nil {
			responseStatus = status.Convert(err)
		}

		latency := time.Since(start)
		currentTime := time.Now().Format("02/Jan/2006:15:04:05")
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Error("gRPC Error: logging interceptor error, can't get metadata")
		}
		userAgent := getInfoFromMetadata(md, "user-agent")
		host := getClientHost(ctx)
		method := info.FullMethod
		reqInfo := method + " " + responseStatus.Code().String()
		logStr := host + " " + currentTime + " " + reqInfo + " " + latency.String() + " " + userAgent
		logger.Info(logStr)

		return resp, err
	}
}

func getInfoFromMetadata(md metadata.MD, key string) string {
	userAgents, ok := md[key]
	if !ok {
		return ""
	}
	return userAgents[0]
}

func getClientHost(ctx context.Context) string {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	return peerInfo.Addr.String()
}
