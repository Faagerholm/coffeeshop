package main

import (
	"context"
	"fmt"

	"github.com/faagerholm/coffee/email/proto"
	"github.com/fatih/color"
)

type server struct {
	proto.UnimplementedEmailServiceServer
}

func (s server) SendEmail(ctx context.Context, req *proto.SendEmailRequest) (*proto.Empty, error) {
	// log email to stdout
	c := color.New(color.FgYellow)
	fmt.Printf(`------------------------------------------
Email sent to %s with subject %s and body:
%s
`, c.Sprintf(req.GetEmail()), c.Sprintf(req.GetSubject()), c.Sprintf(req.GetBody()))

	return &proto.Empty{}, nil
}
