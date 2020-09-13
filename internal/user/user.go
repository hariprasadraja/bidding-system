package user

import (
	context "context"

	log "github.com/micro/go-micro/v2/logger"
)

type Handler struct{}

// Call is a single request handler called via client.Call or the generated client code
func (h *Handler) Call(ctx context.Context, req *Request, rsp *Response) error {
	log.Info("Received User.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (h *Handler) Stream(ctx context.Context, req *StreamingRequest, stream User_StreamStream) error {
	log.Infof("Received User.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (h *Handler) PingPong(ctx context.Context, stream User_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (h *Handler) Create(ctx context.Context, req *CreateRequest, resp *CreateResponse) error {
	log.Info(req.String())
	return nil
}

func (h *Handler) Update(ctx context.Context, req *UpdateRequest, resp *UpdateResponse) error {
	log.Info(req.String())
	return nil
}

func (h *Handler) Get(ctx context.Context, req *GetRequest, resp *GetResponse) error {
	log.Info(req.String())
	return nil
}

func (h *Handler) Delete(ctx context.Context, req *DeleteRequest, resp *DeleteResponse) error {
	log.Info(req.String())
	return nil
}

func (h *Handler) Login(ctx context.Context, req *LoginRequest, resp *LoginResponse) error {
	log.Info(req.String())
	return nil
}
