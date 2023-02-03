package trace

import (
	"Tiktok/config"
	"Tiktok/pkg/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// initTracer 初始化一个全局的TracerProvider，traceID将在log中打印出来，添加jaeger作为可观测性后端
func initTracer() *sdktrace.TracerProvider {
	// Create the Jaeger exporter
	endpoint := fmt.Sprintf("http://%v:%v/%v",
		config.JaegerSetting.Host,
		config.JaegerSetting.Port,
		config.JaegerSetting.Path)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		log.Error("jaeger init failed", zap.Error(err))
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

// Set 封装了对InitTracer的调用，同时使用了官方中间件，r参数是gin的Engine
func Set(r *gin.Engine) {
	// 获取一个全局TracerProvider
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("Error shutting down tracer provider", zap.Error(err))
		}
	}()
	// 使用open-telemetry官方中间件
	r.Use(otelgin.Middleware(config.ServerSetting.Name))
}
