package auth

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/open-policy-agent/opa/rego"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewOPAInterceptor(ctx context.Context, authProvider AuthenticationProvider, policy func(*rego.Rego)) (Interceptor, error) {
	p, err := rego.New(
		rego.Query("res = data.piro.allow"),
		policy,
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot compile policy: %w", err)
	}
	return &opaInterceptor{
		Policy: p,
		Auth:   authProvider,
	}, nil
}

type opaInterceptor struct {
	Policy rego.PreparedEvalQuery
	Auth   AuthenticationProvider
}

type policyInput struct {
	Method   string        `json:"method"`
	Metadata metadata.MD   `json:"metadata"`
	Message  interface{}   `json:"message"`
	Auth     *AuthResponse `json:"auth,omitempty"`
}

func (i *opaInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		auth, err := i.getAuth(ctx)
		if err != nil {
			return nil, err
		}

		md, _ := metadata.FromIncomingContext(ctx)
		input := policyInput{
			Method:   info.FullMethod,
			Metadata: md,
			Message:  req,
			Auth:     auth,
		}
		err = i.eval(ctx, input)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (i *opaInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		auth, err := i.getAuth(ctx)
		if err != nil {
			return err
		}

		md, _ := metadata.FromIncomingContext(ctx)

		return handler(srv, &interceptingStream{
			ServerStream: ss,
			eval: func(msg interface{}) error {
				input := policyInput{
					Method:   info.FullMethod,
					Metadata: md,
					Message:  msg,
					Auth:     auth,
				}
				err = i.eval(ctx, input)
				if err != nil {
					return err
				}
				return nil
			},
		})
	}
}

type interceptingStream struct {
	grpc.ServerStream
	eval func(msg interface{}) error
	done bool
	mu   sync.Mutex
}

func (s *interceptingStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.done {
		return nil
	}
	s.done = true

	err = s.eval(m)
	if err != nil {
		return err
	}
	return nil
}

func (i *opaInterceptor) eval(ctx context.Context, input policyInput) error {
	result, err := i.Policy.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		logrus.WithError(err).Error("cannot evaluate policy")
		return status.Error(codes.Internal, "cannot evaluate policy")
	}
	if len(result) == 0 {
		logrus.WithError(err).Error("policy does not define data.piro.allow query")
		return status.Error(codes.Internal, "invalid policy")
	}

	if _, ok := input.Metadata["x-auth-token"]; ok {
		input.Metadata["x-auth-token"] = []string{"some-value"}
	}
	dmp, _ := json.Marshal(input)
	allowed, ok := result[0].Bindings["res"].(bool)
	logrus.WithField("input", string(dmp)).WithField("allowed", allowed).Debug("evaluating request")

	if !allowed || !ok {
		return status.Error(codes.Unauthenticated, "not allowed")
	}
	return nil
}

func (i *opaInterceptor) getAuth(ctx context.Context) (*AuthResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil
	}
	tkn := md.Get("x-auth-token")
	if len(tkn) == 0 {
		return nil, nil
	}

	aresp, err := i.Auth.Authenticate(ctx, tkn[0])
	if err != nil {
		log.WithError(err).Warn("authentication failure")
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	return aresp, nil
}
