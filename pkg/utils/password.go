package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordService 密码服务
type PasswordService struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// NewPasswordService 创建密码服务
func NewPasswordService() *PasswordService {
	return &PasswordService{
		memory:      64 * 1024, // 64 MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

// HashPassword 哈希密码
func (p *PasswordService) HashPassword(password string) (string, error) {
	salt, err := p.generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// 编码为base64字符串
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword 验证密码
func (p *PasswordService) VerifyPassword(password, encodedHash string) (bool, error) {
	salt, hash, err := p.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// generateRandomBytes 生成随机字节
func (p *PasswordService) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// decodeHash 解码哈希字符串
func (p *PasswordService) decodeHash(encodedHash string) (salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, errors.New("无效的哈希格式")
	}

	if parts[1] != "argon2id" {
		return nil, nil, errors.New("不兼容的哈希算法")
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, err
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, err
	}

	return salt, hash, nil
}

// GenerateRandomPassword 生成随机密码
func (p *PasswordService) GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	
	for i := range b {
		randomBytes := make([]byte, 1)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return "", err
		}
		b[i] = charset[randomBytes[0]%byte(len(charset))]
	}
	
	return string(b), nil
}