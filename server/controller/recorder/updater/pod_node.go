package updater

import (
	cloudmodel "server/controller/cloud/model"
	"server/controller/db/mysql"
	"server/controller/recorder/cache"
	"server/controller/recorder/common"
	"server/controller/recorder/db"
)

type PodNode struct {
	UpdaterBase[cloudmodel.PodNode, mysql.PodNode, *cache.PodNode]
}

func NewPodNode(wholeCache *cache.Cache, cloudData []cloudmodel.PodNode) *PodNode {
	updater := &PodNode{
		UpdaterBase[cloudmodel.PodNode, mysql.PodNode, *cache.PodNode]{
			cache:        wholeCache,
			dbOperator:   db.NewPodNode(),
			diffBaseData: wholeCache.PodNodes,
			cloudData:    cloudData,
		},
	}
	updater.dataGenerator = updater
	updater.cacheHandler = updater
	return updater
}

func (n *PodNode) getDiffBaseByCloudItem(cloudItem *cloudmodel.PodNode) (diffBase *cache.PodNode, exists bool) {
	diffBase, exists = n.diffBaseData[cloudItem.Lcuuid]
	return
}

func (n *PodNode) generateDBItemToAdd(cloudItem *cloudmodel.PodNode) (*mysql.PodNode, bool) {
	vpcID, exists := n.cache.ToolDataSet.GetVPCIDByLcuuid(cloudItem.VPCLcuuid)
	if !exists {
		log.Errorf(resourceAForResourceBNotFound(
			common.RESOURCE_TYPE_VPC_EN, cloudItem.VPCLcuuid,
			common.RESOURCE_TYPE_POD_NODE_EN, cloudItem.Lcuuid,
		))
		return nil, false
	}
	podClusterID, exists := n.cache.ToolDataSet.GetPodClusterIDByLcuuid(cloudItem.PodClusterLcuuid)
	if !exists {
		log.Errorf(resourceAForResourceBNotFound(
			common.RESOURCE_TYPE_POD_CLUSTER_EN, cloudItem.PodClusterLcuuid,
			common.RESOURCE_TYPE_POD_NODE_EN, cloudItem.Lcuuid,
		))
		return nil, false
	}
	dbItem := &mysql.PodNode{
		Name:         cloudItem.Name,
		Type:         cloudItem.Type,
		MemTotal:     cloudItem.MemTotal,
		VCPUNum:      cloudItem.VCPUNum,
		ServerType:   cloudItem.ServerType,
		State:        cloudItem.State,
		IP:           cloudItem.IP,
		PodClusterID: podClusterID,
		SubDomain:    cloudItem.SubDomainLcuuid,
		Domain:       n.cache.DomainLcuuid,
		Region:       cloudItem.RegionLcuuid,
		AZ:           cloudItem.AZLcuuid,
		VPCID:        vpcID,
	}
	dbItem.Lcuuid = cloudItem.Lcuuid
	return dbItem, true
}

func (n *PodNode) generateUpdateInfo(diffBase *cache.PodNode, cloudItem *cloudmodel.PodNode) (map[string]interface{}, bool) {
	updateInfo := make(map[string]interface{})
	if diffBase.State != cloudItem.State {
		updateInfo["state"] = cloudItem.State
	}
	if diffBase.RegionLcuuid != cloudItem.RegionLcuuid {
		updateInfo["region"] = cloudItem.RegionLcuuid
	}
	if diffBase.AZLcuuid != cloudItem.AZLcuuid {
		updateInfo["az"] = cloudItem.AZLcuuid
	}
	if diffBase.VCPUNum != cloudItem.VCPUNum {
		updateInfo["vcpu_num"] = cloudItem.VCPUNum
	}
	if diffBase.MemTotal != cloudItem.MemTotal {
		updateInfo["mem_total"] = cloudItem.MemTotal
	}

	if len(updateInfo) > 0 {
		return updateInfo, true
	}
	return nil, false
}

func (n *PodNode) addCache(dbItems []*mysql.PodNode) {
	n.cache.AddPodNodes(dbItems)
}

func (n *PodNode) updateCache(cloudItem *cloudmodel.PodNode, diffBase *cache.PodNode) {
	diffBase.Update(cloudItem)
}

func (n *PodNode) deleteCache(lcuuids []string) {
	n.cache.DeletePodNodes(lcuuids)
}
