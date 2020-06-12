package cache

type ServentCache interface {
	Set(key string, value []byte)
	Remove(key string)
	Exists(key string) bool
	Get(key string) []byte
	GetAll() map[string][]byte
}

type ServentCacheHelper struct {
	PeerList       ServentCache
	DescriptorList ServentCache
	PeerFilesList  ServentCache
	FindFilesList  ServentCache
}

func NewServentCacheHelper() *ServentCacheHelper {
	return &ServentCacheHelper{
		PeerList:       NewPeerList(),
		DescriptorList: NewDescriptorList(),
		PeerFilesList:  NewPeerFilesList(),
		FindFilesList:  NewFindFilesList(),
	}
}
