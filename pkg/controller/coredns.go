package controller

import (
	"fmt"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/log"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/models"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/service"
	"github.com/gin-gonic/gin"
)

type CorednsController interface {
	RegisterNode(c *gin.Context)
	GetNode(c *gin.Context)
	AddDomain(c *gin.Context)
}

func NewCorednsController() CorednsController {
	return &corednsController{
		corednsService: service.NewCorednsService(),
	}
}

type corednsController struct {
	corednsService service.CorednsService
}

// @Summary RegisterNode
// @Description RegisterNode
// @Param   body   body   models.NodeInfoList true   "Node info list you want to register."
// @Success 201 {object} models.Response OK
// @Failure 400 {object} models.Response Bad Request
// @Failure 401 {object} models.Response Unauthorized
// @Failure 403 {object} models.Response Forbidden
// @Failure 500 {object} models.Response Internal Server Error
// @router /node [post]
func (this *corednsController) RegisterNode(c *gin.Context) {
	nodeInfoList := &models.NodeInfoList{}
	err := c.ShouldBindJSON(nodeInfoList)
	if err != nil {
		message := fmt.Sprintf("Json Unmarshal error: %v", err.Error())
		log.Error(message)
		c.JSON(500, models.Response{Code: 500, Message: message})
		return
	}
	// check params
	seen := make(map[string]struct{})
	for _, nodeInfo := range nodeInfoList.NodeInfos {
		if nodeInfo.Username == "" || nodeInfo.Password == "" ||
			nodeInfo.Port == 0 || nodeInfo.Address == "" {
			message := fmt.Sprintf("Invalid node params(some fields are empty): %v", nodeInfoList)
			log.Error(message)
			c.JSON(500, models.Response{Code: 500, Message: message})
			return
		}
		if _, ok := seen[nodeInfo.Address]; !ok {
			seen[nodeInfo.Address] = struct{}{}
		} else {
			message := fmt.Sprintf("Invalid node params(some fields are duplicated): %v", nodeInfoList)
			log.Error(message)
			c.JSON(500, models.Response{Code: 500, Message: message})
			return
		}
	}
	err = this.corednsService.Register(nodeInfoList)
	if err != nil {
		message := fmt.Sprintf("Register Node failed: %v", err.Error())
		log.Error(message)
		c.JSON(500, models.Response{Code: 500, Message: message})
		return
	}
	c.JSON(200, models.Response{Code: 0, Message: "Register Node success."})
	return
}

// @Summary GetNode
// @Description GetNode
// @Success 201 {object} models.Response OK
// @Failure 400 {object} models.Response Bad Request
// @Failure 401 {object} models.Response Unauthorized
// @Failure 403 {object} models.Response Forbidden
// @Failure 500 {object} models.Response Internal Server Error
// @router /node [get]
func (this *corednsController) GetNode(c *gin.Context) {
	nodeInfoList, err := this.corednsService.Get()
	if err != nil {
		message := fmt.Sprintf("Get Node failed: %v", err.Error())
		log.Error(message)
		c.JSON(500, models.Response{Code: 500, Message: message})
		return
	}
	c.JSON(200, models.Response{Code: 0, Message: "Get Node success.", Data: nodeInfoList})
	return
}

// @Summary AddDomain
// @Description AddDomain
// @Param   body   body   models.DomainInfoList true   "Domain info list you want to add to coredns."
// @Success 201 {object} models.Response OK
// @Failure 400 {object} models.Response Bad Request
// @Failure 401 {object} models.Response Unauthorized
// @Failure 403 {object} models.Response Forbidden
// @Failure 500 {object} models.Response Internal Server Error
// @router /domain [post]
func (this *corednsController) AddDomain(c *gin.Context) {
	domainInfoList := &models.DomainInfoList{}
	err := c.ShouldBindJSON(domainInfoList)
	if err != nil {
		message := fmt.Sprintf("Json Unmarshal error: %v", err.Error())
		log.Error(message)
		c.JSON(500, models.Response{Code: 500, Message: message})
		return
	}
	// check params
	seen := make(map[string]struct{})
	for _, domainInfo := range domainInfoList.DomainInfos {
		if domainInfo.Domain == "" || domainInfo.IP == "" {
			message := fmt.Sprintf("Invalid domain params(some fields are empty): %v", domainInfoList)
			log.Error(message)
			c.JSON(500, models.Response{Code: 500, Message: message})
			return
		}
		if _, ok := seen[domainInfo.Domain]; !ok {
			seen[domainInfo.Domain] = struct{}{}
		} else {
			message := fmt.Sprintf("Invalid domain params(some fields are duplicated): %v", domainInfoList)
			log.Error(message)
			c.JSON(500, models.Response{Code: 500, Message: message})
			return
		}
	}
	err = this.corednsService.AddDomain(domainInfoList)
	if err != nil {
		message := fmt.Sprintf("Add Domain failed: %v", err.Error())
		log.Error(message)
		c.JSON(500, models.Response{Code: 500, Message: message})
		return
	}
	c.JSON(200, models.Response{Code: 0, Message: "Add Domain success."})
	return
}
