package server

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type streamTransport struct {
}

func (s *streamTransport) handlePacket(t ServerTransport, packet *Packet, trInfo *traceInfo) {

	entityID := packet.EntityID()
	_ = entityID

	sm := packet.Method()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		if trInfo != nil {
			trInfo.tr.LazyLog(&fmtStringer{"Malformed method name %q", []interface{}{sm}}, true)
			trInfo.tr.SetError()
		}
		errDesc := fmt.Sprintf("malformed method name: %q", packet.Method())
		if err := t.WriteStatus(packet, status.New(codes.Unimplemented, errDesc)); err != nil {
			if trInfo != nil {
				trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
				trInfo.tr.SetError()
			}
			channelz.Warningf(logger, s.channelzID, "grpc: Server.handlePacket failed to write status: %v", err)
		}
		if trInfo != nil {
			trInfo.tr.Finish()
		}
		return
	}
	service := sm[:pos]
	method := sm[pos+1:]

	srv, knownService := s.services[service]
	if knownService {
		if md, ok := srv.methods[method]; ok {
			s.processUnaryRPC(t, packet, srv, md, trInfo)
			return
		}
		if sd, ok := srv.packets[method]; ok {
			s.processPacketingRPC(t, packet, srv, sd, trInfo)
			return
		}
	}
	// Unknown service, or known server unknown method.
	if unknownDesc := s.opts.unknownPacketDesc; unknownDesc != nil {
		s.processPacketingRPC(t, packet, nil, unknownDesc, trInfo)
		return
	}
	var errDesc string
	if !knownService {
		errDesc = fmt.Sprintf("unknown service %v", service)
	} else {
		errDesc = fmt.Sprintf("unknown method %v for service %v", method, service)
	}
	if trInfo != nil {
		trInfo.tr.LazyPrintf("%s", errDesc)
		trInfo.tr.SetError()
	}
	if err := t.WriteStatus(packet, status.New(codes.Unimplemented, errDesc)); err != nil {
		if trInfo != nil {
			trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
			trInfo.tr.SetError()
		}
		channelz.Warningf(logger, s.channelzID, "grpc: Server.handlePacket failed to write status: %v", err)
	}
	if trInfo != nil {
		trInfo.tr.Finish()
	}
}
