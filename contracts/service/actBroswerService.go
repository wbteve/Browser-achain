package service

import (
	"Browser-achain/common"
	"Browser-achain/contracts/models"
	"Browser-achain/util"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"

)

type UserBalanceVo struct {
	CoinType string
	Balance  string
}

const coinType = "ACT"
const contractLength = 30
const CONTARCT_PREFIX = "CON"

// Check the balance of all currencies in the address according to the address
func QueryBalanceByAddress(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		common.WebResultFail(c)
	}

	list, err := models.ListByAddress(address)

	if err != nil {
		log.Fatal("QueryBalanceByAddress|query data ERROR:", err)
		common.WebResultFail(c)
	}

	userBalanceVoList := make([]UserBalanceVo, 0)

	for _, value := range list {
		var userBalanceVo UserBalanceVo
		userBalanceVo.CoinType = value.CoinType
		actualAmountString := util.GetActualAmount(value.Balance)
		actualAmount, _ := strconv.ParseFloat(actualAmountString, 32)

		// if coinType is not ACT and balance less than 0,ignore
		if coinType != userBalanceVo.CoinType && actualAmount <= float64(0) {
			return
		}
		// if coinType is ACT and balance less than 0, replace with 0
		if coinType == userBalanceVo.CoinType && actualAmount <= float64(0) {
			actualAmount = float64(0)
		}
		userBalanceVo.Balance = strconv.FormatFloat(actualAmount, 'E', -1, 32)
		userBalanceVoList = append(userBalanceVoList, userBalanceVo)
	}
	common.WebResultSuccess(userBalanceVoList, c)
}

// Check the address balance by the keyword
func QueryContractByKey(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	perPage, _ := strconv.Atoi(c.Param("perPage"))
	keyword := c.Query("keyword")
	log.Printf("QueryContractByKey|page=%s|perPage=%s|keyword=%s\n", page, perPage, keyword)
	if page < 1 || perPage < 1 {
		common.WebResultFail(c)
	}

	queryType := 1
	// keyword is not empty and keyword startWith CONN and the length of keyword greater than 30
	if keyword != "" && strings.Index(keyword, CONTARCT_PREFIX) == 0 && len(keyword) > 30 {
		queryType = 0
	}
	contractInfoPageVO, err := models.ListContractInfoByKey(keyword, models.Forever, page, perPage, queryType)
	if err != nil {
		common.WebResultFail(c)
	}

	contractInfoVOList := make([]models.ContractInfoVO, 0)
	tbActContractInfoList := contractInfoPageVO.ActContractInfoList
	for _, actContractInfo := range tbActContractInfoList {
		var contractInfoVO models.ContractInfoVO
		circulation := actContractInfo.Circulation
		intCirculation, _ := strconv.ParseInt(util.GetActualAmount(&circulation), 10, 64)
		contractInfoVO.Circulation = intCirculation
		contractInfoVO.CoinType = actContractInfo.CoinType
		contractInfoVO.CoinAddress = actContractInfo.OwnerAddress
		contractInfoVO.ContractName = actContractInfo.Name
		contractInfoVO.RegisterTime = actContractInfo.RegTime
		contractInfoVO.Status = actContractInfo.Status
		contractInfoVO.ContractId = actContractInfo.ContractId
		contractInfoVO.Coin = int(actContractInfo.Type)

		contractInfoVOList = append(contractInfoVOList, contractInfoVO)
	}
	contractInfoPageVO.ActContractInfoList = nil
	contractInfoPageVO.ContractInfoVOList = contractInfoVOList
	common.WebResultSuccess(contractInfoPageVO, c)
}

// Query the act address account information
func QueryAddressInfo(c *gin.Context)  {
	userActAddress := c.Param("userAddress")
	userAddressList, err := models.ListByAddressAndCoinType(userActAddress, "ACT")
	if err != nil {
		common.WebResultFail(c)
	}
	var userAddressVO models.UserAddressVO

	if len(userAddressList) > 0{
		tbUserAddress := userAddressList[0]
		actualAmount := util.GetActualAmount(tbUserAddress.Balance)
		userAddressVO.Balance = actualAmount
		userAddressVO.Address = *tbUserAddress.UserAddress
	}
	common.WebResultSuccess(userAddressVO, c)
}


