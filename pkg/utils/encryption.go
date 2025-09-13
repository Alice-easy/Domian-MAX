package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptionService 加密服务
type EncryptionService struct {
	key []byte
}

// NewEncryptionService 创建加密服务
func NewEncryptionService(key string) (*EncryptionService, error) {
	if len(key) != 32 {
		return nil, errors.New("密钥长度必须为32字节")
	}
	
	return &EncryptionService{
		key: []byte(key),
	}, nil
}

// Encrypt 加密数据
func (e *EncryptionService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// 编码为base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密数据
func (e *EncryptionService) Decrypt(ciphertext string) (string, error) {
	// 解码base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文太短")
	}

	// 分离nonce和密文
	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	
	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptMap 加密映射中的敏感字段
func (e *EncryptionService) EncryptMap(data map[string]interface{}, sensitiveFields []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	for key, value := range data {
		// 检查是否为敏感字段
		isSensitive := false
		for _, field := range sensitiveFields {
			if key == field {
				isSensitive = true
				break
			}
		}
		
		if isSensitive {
			if strValue, ok := value.(string); ok {
				encrypted, err := e.Encrypt(strValue)
				if err != nil {
					return nil, err
				}
				result[key] = encrypted
			} else {
				result[key] = value
			}
		} else {
			result[key] = value
		}
	}
	
	return result, nil
}

// DecryptMap 解密映射中的敏感字段
func (e *EncryptionService) DecryptMap(data map[string]interface{}, sensitiveFields []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	for key, value := range data {
		// 检查是否为敏感字段
		isSensitive := false
		for _, field := range sensitiveFields {
			if key == field {
				isSensitive = true
				break
			}
		}
		
		if isSensitive {
			if strValue, ok := value.(string); ok {
				decrypted, err := e.Decrypt(strValue)
				if err != nil {
					return nil, err
				}
				result[key] = decrypted
			} else {
				result[key] = value
			}
		} else {
			result[key] = value
		}
	}
	
	return result, nil
}