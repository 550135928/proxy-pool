package api

import (
	"context"
	"proxy-pool/config"
	"proxy-pool/databases"
	"proxy-pool/model"
	"time"

	log "github.com/sirupsen/logrus"
)

// Service http handle
type Service struct {
	db   *databases.ORM
	conf *config.Config
}

// NewService service
func NewService(db *databases.ORM, conf *config.Config) *Service {
	return &Service{
		db:   db,
		conf: conf,
	}
}

type (
	// ProxyRsp 代理信息结构体
	ProxyRsp struct {
		ID        uint64    `gorm:"column:id" json:"id"`
		IP        string    `gorm:"column:ip" json:"ip"`
		Port      int       `gorm:"column:port" json:"port"`
		Schema    string    `gorm:"column:schema" json:"schema"`
		CheckTime time.Time `gorm:"column:check_time" json:"last_check_time"`
	}
)

// CopyFromProxy proxy to rsp
func (p *ProxyRsp) CopyFromProxy(proxy *model.Proxy) {
	p.ID = proxy.ID
	p.IP = proxy.IP
	p.Port = proxy.Port
	p.Schema = proxy.Schema
	p.CheckTime = proxy.CheckTime
}

// GetOneProxy 获取一个代理
func (s *Service) GetOneProxy(c context.Context) (*ProxyRsp, error) {
	rsp := &ProxyRsp{}
	if err := s.db.DB.Raw(`SELECT r1.id, r1.schema, r1.ip, r1.port,
	r1.check_time FROM proxy r1 order by RANDOM() LIMIT 1`).
		Row().Scan(&rsp.ID, &rsp.Schema, &rsp.IP, &rsp.Port, &rsp.CheckTime); err != nil {
		return nil, NoFound
	}
	// rsp.CheckTime =
	return rsp, nil
}

// GetAllProxy 获取所有代理
func (s *Service) GetAllProxy(c context.Context) ([]*ProxyRsp, error) {
	var rsps = make([]*ProxyRsp, 0)
	if err := s.db.DB.Table("proxy").
		Select([]string{"id", "ip", "`port`", "`schema`", "check_time"}).
		Where("is_deleted=?", 0).
		Find(&rsps).Error; err != nil {
		return nil, ServerError
	}
	return rsps, nil
}

// DeleteOneProxy 删除一个代理
func (s *Service) DeleteOneProxy(c context.Context, id int) error {
	if err := s.db.DB.Table("proxy").
		Where("id=?", id).
		Updates(map[string]interface{}{"is_deleted": 1}).
		Error; err != nil {
		log.Errorf("DeleteOneProxy id:%d, err:%#v", id, err)
		return ServerError
	}
	return nil
}
