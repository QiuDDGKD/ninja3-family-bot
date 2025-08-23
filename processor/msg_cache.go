package processor

import (
	"sync"
	"time"
)

type MsgCache struct {
	cache map[string]time.Time // 存储消息ID及其过期时间
	mu    sync.RWMutex         // 读写锁保护并发访问
	ttl   time.Duration        // 消息的存活时间
}

// NewMsgCache 创建一个新的消息ID缓存器
func NewMsgCache(ttl time.Duration) *MsgCache {
	return &MsgCache{
		cache: make(map[string]time.Time),
		ttl:   ttl,
	}
}

// Add 添加消息ID到缓存，返回是否真的添加到缓存
func (mc *MsgCache) Add(msgID string) bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	_, exists := mc.cache[msgID]
	if exists {
		return false // 已存在，未添加
	}
	mc.cache[msgID] = time.Now().Add(mc.ttl)
	return true // 新添加
}

// Exists 检查消息ID是否存在
func (mc *MsgCache) Exists(msgID string) bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	expiry, exists := mc.cache[msgID]
	if !exists {
		return false
	}
	// 检查是否过期
	if time.Now().After(expiry) {
		return false
	}
	return true
}

// CleanUp 清理过期的消息ID
func (mc *MsgCache) CleanUp() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	now := time.Now()
	for msgID, expiry := range mc.cache {
		if now.After(expiry) {
			delete(mc.cache, msgID)
		}
	}
}
