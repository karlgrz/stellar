package application

import (
	"fmt"

	"github.com/containerd/containerd"
	"github.com/ehazlett/element"
	"github.com/ehazlett/stellar"
	api "github.com/ehazlett/stellar/api/services/application/v1"
	"github.com/ehazlett/stellar/client"
	ptypes "github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
)

const (
	serviceID = "stellar.services.application.v1"
)

var (
	empty = &ptypes.Empty{}
)

type service struct {
	containerdAddr string
	namespace      string
	dataDir        string
	agent          *element.Agent
}

func New(containerdAddr, namespace, dataDir string, agent *element.Agent) (*service, error) {
	return &service{
		containerdAddr: containerdAddr,
		namespace:      namespace,
		dataDir:        dataDir,
		agent:          agent,
	}, nil
}

func (s *service) Register(server *grpc.Server) error {
	api.RegisterApplicationServer(server, s)
	return nil
}

func (s *service) ID() string {
	return serviceID
}

func (s *service) Start() error {
	return nil
}

func (s *service) containerd() (*containerd.Client, error) {
	return stellar.DefaultContainerd(s.containerdAddr, s.namespace)
}

func (s *service) client() (*client.Client, error) {
	peer, err := s.agent.LocalNode()
	if err != nil {
		return nil, err
	}
	return client.NewClient(peer.Addr)
}

func (s *service) nodeName() string {
	return s.agent.Config().NodeName
}

func (s *service) nodeClient(id string) (*client.Client, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	nodes, err := c.Cluster().Nodes()
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		if node.Name == id {
			return client.NewClient(node.Addr)
		}
	}

	return nil, fmt.Errorf("node %s not found in cluster", id)
}
