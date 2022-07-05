package p2p

import (
	"sync"

	"github.com/cloudslit/cloudslit/fullnode/pkg/util/json"

	"github.com/cloudslit/cloudslit/fullnode/pkg/schema"
)

var server = &Server{
	Server: make([]*schema.ServerInfo, 0),
}

type Server struct {
	Server []*schema.ServerInfo
	sync.RWMutex
}

//
//func CheckServerExist(peerid string) bool {
//	return server.checkServerExist(peerid)
//}
//
//func AddServer(ser *schema.ServerInfo) {
//	server.addServer(ser)
//}
//
//func DelServer(peerid string) {
//	server.delServer(peerid)
//}
//
//func GetServer() []*schema.ServerInfo {
//	return server.getServer()
//}

//func (s *Server) checkServerExist(peerid string) bool {
//	if len(s.Server) == 0 {
//		return false
//	}
//	s.Lock()
//	defer s.Unlock()
//	for key, value := range s.Server {
//		if value.PeerId == peerid {
//			s.Server[key].WithTime()
//		}
//	}
//	return false
//}
//
//func (s *Server) addServer(server *schema.ServerInfo) {
//	if s.checkServerExist(server.PeerId) {
//		return
//	}
//	s.Lock()
//	defer s.Unlock()
//	s.Server = append(s.Server, server)
//}
//
//func (s *Server) delServer(peerid string) {
//	if !s.checkServerExist(peerid) {
//		return
//	}
//	s.Lock()
//	defer s.Unlock()
//	for i := 0; i < len(s.Server); i++ {
//		if s.Server[i].PeerId == peerid {
//			switch {
//			case i == 0:
//				s.Server = s.Server[1:]
//			case i == len(s.Server)-1:
//				s.Server = s.Server[:i]
//			default:
//				s.Server = append(s.Server[:i], s.Server[i+1:]...)
//			}
//			i-- // 如果索引i被去掉后，原来索引i+1的会移动到索引i
//		}
//	}
//}
//
//func (s *Server) getServer() []*schema.ServerInfo {
//	s.RLock()
//	defer s.RUnlock()
//	return s.Server
//}
//
func Generate(message string) (server *schema.ServerInfo) {
	_ = json.Unmarshal([]byte(message), &server)
	return
}
