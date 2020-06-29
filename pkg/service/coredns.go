package service

import (
	"context"
	"fmt"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/config"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/log"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/models"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/util"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CorednsService interface {
	Register(*models.NodeInfoList) error
	Get() (models.NodeInfoList, error)
	AddDomain(*models.DomainInfoList) error
}

type corednsService struct {
	NodeList *models.NodeInfoList
}

func NewCorednsService() CorednsService {
	return &corednsService{}
}

// Registry Node to dynapi
func (c *corednsService) Register(list *models.NodeInfoList) error {
	c.NodeList = list
	return nil
}

// Get NodeInfoList
func (c *corednsService) Get() (models.NodeInfoList, error) {
	if c.NodeList != nil {
		return *c.NodeList, nil
	}
	return models.NodeInfoList{}, nil
}

// Add Domain to coredns
func (c *corednsService) AddDomain(list *models.DomainInfoList) error {
	log.Debugf("Begin Add Domain: %v to Coredns ..., list")
	errMap := make(map[string]error)
	util.ParallelizeUntil(context.TODO(), 16, len(c.NodeList.NodeInfos), func(index int) {
		for _, domainInfo := range list.DomainInfos {
			// Open a ssh connection
			sshConfig := &ssh.ClientConfig{
				User: c.NodeList.NodeInfos[index].Username,
				Auth: []ssh.AuthMethod{
					ssh.Password(c.NodeList.NodeInfos[index].Password),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
			client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.NodeList.NodeInfos[index].Address, c.NodeList.NodeInfos[index].Port), sshConfig)
			if err != nil {
				log.Error("Failed to dial: ", err)
				errMap[c.NodeList.NodeInfos[index].Address] = err
				return
			}
			defer client.Close()
			// Open a SFTP session over an existing ssh connection.
			sftpClient, err := sftp.NewClient(client)
			if err != nil {
				log.Error("Failed to sftp: ", err)
				errMap[c.NodeList.NodeInfos[index].Address] = err
				return
			}
			defer sftpClient.Close()
			// Open the coredns relevant zoneFile(create it if not exist)
			zoneFilePath := filepath.Join(config.GetConfig().Coredns.ZonesDir, domainInfo.Domain)
			zoneFile, err := sftpClient.OpenFile(zoneFilePath, os.O_RDWR|os.O_CREATE)
			if err != nil {
				log.Error("Failed to open zoneFile: ", err)
				errMap[c.NodeList.NodeInfos[index].Address] = err
				return
			}
			defer zoneFile.Close()
			// Write relevant coredns zoneContent
			zoneContent := strings.ReplaceAll(util.CorednsZoneTemplate, "{ZONE}", domainInfo.Domain)
			zoneContent = strings.ReplaceAll(zoneContent, "{IP}", domainInfo.IP)
			sftpZoneBytes, err := ioutil.ReadAll(zoneFile)
			if err != nil {
				log.Error("Failed to read zoneFile: ", err)
				errMap[c.NodeList.NodeInfos[index].Address] = err
				return
			}
			if !strings.Contains(string(sftpZoneBytes), zoneContent) {
				log.Infof("Try to add domain/ip: %s/%s to zoneFile: %s(%s) ...", domainInfo.Domain, domainInfo.IP, zoneFilePath, c.NodeList.NodeInfos[index].Address)
				_, err = zoneFile.Write([]byte(zoneContent))
				if err != nil {
					log.Error("Failed to write zoneFile: ", err)
					errMap[c.NodeList.NodeInfos[index].Address] = err
					return
				}
			} else {
				log.Infof("Ignore existing domain/ip: %s/%s of zoneFile: %s(%s)", domainInfo.Domain, domainInfo.IP, zoneFilePath, c.NodeList.NodeInfos[index].Address)
			}
			// Open the coredns relevant corefile(create it if not exist)
			corefilePath := config.GetConfig().Coredns.CorefilePath
			corefile, err := sftpClient.OpenFile(corefilePath, os.O_RDWR|os.O_CREATE)
			if err != nil {
				log.Error("Failed to open corefile: ", err)
				errMap[c.NodeList.NodeInfos[index].Address] = err
				return
			}
			defer corefile.Close()
			// Write relevant coredns server block
			corefileContent := strings.ReplaceAll(util.CorednsServerBlockTemplate, "{ZONE}", domainInfo.Domain)
			corefileContent = strings.ReplaceAll(corefileContent, "{ZONESDIR}", config.GetConfig().Coredns.ZonesDir)
			sftpCorefileBytes, err := ioutil.ReadAll(corefile)
			if err != nil {
				log.Error("Failed to read corefile: ", err)
				errMap[c.NodeList.NodeInfos[index].Address] = err
				return
			}
			if !strings.Contains(string(sftpCorefileBytes), corefileContent) {
				log.Infof("Try to add server block: %s to corefile: %s(%s) ...", domainInfo.Domain, corefilePath, c.NodeList.NodeInfos[index].Address)
				_, err = corefile.Write([]byte(corefileContent))
				if err != nil {
					log.Error("Failed to write corefile: ", err)
					errMap[c.NodeList.NodeInfos[index].Address] = err
					return
				}
			} else {
				log.Infof("Ignore existing server block: %s of corefile: %s(%s)", domainInfo.Domain, corefilePath, c.NodeList.NodeInfos[index].Address)
			}

		}
	})
	if len(errMap) != 0 {
		return fmt.Errorf("Add Domain error: %v", errMap)
	}
	return nil
}
