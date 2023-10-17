package common

import (
	"fmt"
)

const (
	// tcp
	tcpGatewayNodeFmt = "tcpnode:%d"

	// account
	userIdIndexKeyFmt         = "uidindex"
	accountKeyFmt             = "acc:%d"
	accountUIDKeyFmt          = "accuid:%s"
	nonceKeyFmt               = "nonce:%s"
	walletAccKeyFmt           = "walletacc:%s"
	emailSendDailyLimitKeyFmt = "esdl:%s:%s"
	emailBindCodeKeyFmt       = "ebc:%s"
	emailAcc                  = "emailacc:%s"
	emailResetPWCodeKeyFmt    = "erpwc:%s"
	emailResetPWNonceKeyFmt   = "erpwn:%s"

	// login
	tokenKeyFmt          = "token:%s"
	uidTokenKeyFmt       = "uidtoken:%d"
	deviceAccountsKeyFmt = "devaccs:%s"
	loginInfoKeyFmt      = "login:%d"

	// user
	userKeyFmt      = "user:%d"
	userStateKeyFmt = "ustate:%d"

	// item
	itemsKeyFmt       = "items:%d"
	itemsShardKeyFmt  = "items:%d"
	itemsUShardKeyFmt = "uitems:%d"

	// mail
	mailsKeyFmt = "mails:%d"

	// social
	friendsKeyFmt         = "friends:%d"
	friendAppliesKeyFmt   = "friendapplies:%d"
	friendBlackListKeyFmt = "friendbl:%d"

	// nft
	nftGraphiLimitKeyFmt             = "nftgl:%s"
	nftGraphiStartIndexKeyFmt        = "nftgsi:%s"
	nftsKey                          = "nfts:%d"
	nftsListKey                      = "nftl:%d"
	nftUIDKey                        = "nftu:%s"
	collectionGraphiStartIndexKeyFmt = "colgsi:%s"
	collectionTranVersionsKeyFmt     = "coltranvs:%s"
	collectionTranActivityKeyFmt     = "coltran:%d"
	collectionGraphiLimitKeyFmt      = "colgl:%s"
)

// Get key
// tcp
func TCPGatewayNodeKey(userId uint64) string {
	return fmt.Sprintf(tcpGatewayNodeFmt, userId)
}

// account
func UserIdIndexKey() string {
	return userIdIndexKeyFmt
}

func AccountKey(userId uint64) string {
	return fmt.Sprintf(accountKeyFmt, userId)
}

func AccountUIDKey(acc string) string {
	return fmt.Sprintf(accountUIDKeyFmt, acc)
}

func NonceKey(nonce string) string {
	return fmt.Sprintf(nonceKeyFmt, nonce)
}

func WalletAccKey(walletAddr string) string {
	return fmt.Sprintf(walletAccKeyFmt, walletAddr)
}

func EmailSendDailyLimitKey(emailAddr string, date string) string {
	return fmt.Sprintf(emailSendDailyLimitKeyFmt, emailAddr, date)
}

func EmailBindCodeKey(emailAddr string) string {
	return fmt.Sprintf(emailBindCodeKeyFmt, emailAddr)
}

func EmailAccKey(emailAddr string) string {
	return fmt.Sprintf(emailAcc, emailAddr)
}

func EmailResetPasswordValidationCodeKey(emailAddr string) string {
	return fmt.Sprintf(emailResetPWCodeKeyFmt, emailAddr)
}

func EmailResetPasswordNonceKey(emailAddr string) string {
	return fmt.Sprintf(emailResetPWNonceKeyFmt, emailAddr)
}

// login
func TokenKey(token string) string {
	return fmt.Sprintf(tokenKeyFmt, token)
}

func UIDTokenKey(userId uint64) string {
	return fmt.Sprintf(uidTokenKeyFmt, userId)
}

func DeviceAccountsKey(deviceId string) string {
	return fmt.Sprintf(deviceAccountsKeyFmt, deviceId)
}

func LoginInfoKey(userId uint64) string {
	return fmt.Sprintf(loginInfoKeyFmt, userId)
}

// user
func UserKey(userId uint64) string {
	return fmt.Sprintf(userKeyFmt, userId)
}

func UserSateKey(userId uint64) string {
	return fmt.Sprintf(userStateKeyFmt, userId)
}

// item
func ItemsKey(userId uint64) string {
	return fmt.Sprintf(itemsKeyFmt, userId)
}

func ItemsShardKeys(ids []uint32) []string {
	if len(ids) == 0 {
		return nil
	}
	ret := make([]string, 0, len(ids))
	for _, id := range ids {
		ret = append(ret, fmt.Sprintf(itemsShardKeyFmt, id))
	}
	return ret
}

func ItemsShardKey(id uint32) string {
	return fmt.Sprintf(itemsShardKeyFmt, id)
}

func ItemsShardKeysByShardCnt(cnt uint32) []string {
	if cnt == 0 {
		return nil
	}
	ret := make([]string, cnt)
	for i := uint32(0); i < cnt; i++ {
		ret[i] = fmt.Sprintf(itemsShardKeyFmt, i)
	}
	return ret
}

func UItemsShardKeys(ids []uint32) []string {
	if len(ids) == 0 {
		return nil
	}
	ret := make([]string, 0, len(ids))
	for _, id := range ids {
		ret = append(ret, fmt.Sprintf(itemsUShardKeyFmt, id))
	}
	return ret
}

func UItemsShardKey(id uint32) string {
	return fmt.Sprintf(itemsUShardKeyFmt, id)
}

func UItemsShardKeysByShardCnt(cnt uint32) []string {
	if cnt == 0 {
		return nil
	}
	ret := make([]string, cnt)
	for i := uint32(0); i < cnt; i++ {
		ret[i] = fmt.Sprintf(itemsUShardKeyFmt, i)
	}
	return ret
}

// mail
func MailsKey(userId uint64) string {
	return fmt.Sprintf(mailsKeyFmt, userId)
}

// social
func FriendsKey(userId uint64) string {
	return fmt.Sprintf(friendsKeyFmt, userId)
}

func FriendAppliesKey(userId uint64) string {
	return fmt.Sprintf(friendAppliesKeyFmt, userId)
}

func FriendBlackList(userId uint64) string {
	return fmt.Sprintf(friendBlackListKeyFmt, userId)
}

// nft
func NFTGraphiLimitKey(addr string) string {
	return fmt.Sprintf(nftGraphiLimitKeyFmt, addr)
}

func NFTGraphiStartIndexKey(addr string) string {
	return fmt.Sprintf(nftGraphiStartIndexKeyFmt, addr)
}

func NFTsKey(userId uint64) string {
	return fmt.Sprintf(nftsKey, userId)
}

func NFTsListKey(userId uint64) string {
	return fmt.Sprintf(nftsListKey, userId)
}

func NFTUIDKey(id string) string {
	return fmt.Sprintf(nftUIDKey, id)
}

func CollectionGraphiStartIndex(collectionId string) string {
	return fmt.Sprintf(collectionGraphiStartIndexKeyFmt, collectionId)
}

func CollectionTranVersions(collectionId string) string {
	return fmt.Sprintf(collectionTranVersionsKeyFmt, collectionId)
}

func CollectionTranActivity(version uint64) string {
	return fmt.Sprintf(collectionTranActivityKeyFmt, version)
}

func CollectionGraphiLimit(collectionId string) string {
	return fmt.Sprintf(collectionGraphiLimitKeyFmt, collectionId)
}
