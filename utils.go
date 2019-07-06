package goengine

import (
  "crypto/md5"
  "encoding/hex"
  "github.com/satori/go.uuid"
)

func GenerateSid() string {
  md5Gen := md5.New()
  uuid_v1, _ := uuid.NewV1()
  uuid_v4, _ := uuid.NewV4()

  md5Gen.Write([]byte(uuid_v1.String() + "-" + uuid_v4.String()))
  cipherStr := md5Gen.Sum(nil)
  return hex.EncodeToString(cipherStr)
}
