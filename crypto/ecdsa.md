[*crypto目录*]()
------

### 目录

- [生成私钥](#生成私钥)
- [检查秘钥对](#检查秘钥对)
- [私钥到PEM格式文件](#私钥到PEM格式文件)
- [公钥与PEM格式文件](#公钥与PEM格式文件)

### 生成私钥

使用`P256`曲线生成ECDSA的私钥。

```go
func GenECDSAKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.WithMessage(err, "generate ecdsa key: %v")
	}

	if err := CheckECDSAKey(key); err != nil {
		return nil, errors.WithMessage(err, "check ecdsa key")
	}
	return key, err
}
```

[↑top](#目录)

### 检查秘钥对

使用签名和延签验证ECDSA公私钥对。

```go
func CheckECDSAKey(key *ecdsa.PrivateKey, pub *ecdsa.PublicKey) error {
	fakeMsgHash := []byte("ecdsa")

	r, s, err := ecdsa.Sign(rand.Reader, key, fakeMsgHash)
	if err != nil {
		return errors.WithMessage(err, "sign error")
	}
	if !ecdsa.Verify(pub, fakeMsgHash, r, s) {
		return errors.WithMessage(err, "ecdsa private key mismatch with public key")
	}
	return nil
}
```

[↑top](#目录)


### 私钥到PEM格式文件

为了让私钥文件能够在不同厂家、不同程序识别，让私钥通用，需要把ECDSA的私钥转换为`PKCS8`定义的格式，然后使用`ANS1`保存PKCS8格式的数据，最后为了可以以文本方式传播私钥，把ANS1的数据转换为PEM。

```go
func WriteECDSAPrivKeyPem(key *ecdsa.PrivateKey, path string) error {
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return errors.WithMessage(err, "marhshal ecdsa key")
	}

	ecf, err := os.Create(path)
	if err != nil {
		return errors.WithMessage(err, "create ecdsa file")
	}
	defer ecf.Close()

	if err := pem.Encode(ecf, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return errors.WithMessage(err, "write to private key pem file failed")
	}

	return nil
}
```

从PEM恢复ECDSA私钥。`ParsePKCS8PrivateKey`还可以用来恢复RSA的私钥。

```go
func LoadECDSAPrivKeyPem(path string) (*ecdsa.PrivateKey, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.WithMessage(err, "open private file error")
	}
	defer f.Close()

	buf := make([]byte, 1000)
	n, err := f.Read(buf)
	if err != nil {
		return nil, errors.WithMessage(err, "read private key pem failed")
	}
	buf = buf[:n]

	block, _ := pem.Decode(buf)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		errors.WithMessage(err, "parse pkc8 private key")
	}

	if k, ok := key.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("not a ecdsa private key")
	} else {
		return k, nil
	}
}
```


[↑top](#目录)

### 公钥与PEM格式文件

因为ECDSA的公钥是可以通过私钥计算的，PKCS系列中的PKCS1标准，定义了RSA的公钥传播格式，没有ECDSA的。

但存在另外一个序列化方案，`PKIX`，PKIX是支持X509的一个工作组，Go中提供了`MarshalPKIXPublicKey`能够序列化RSA、ECDSA的公钥。

```go
func WriteECDSAPubKeyPem(key *ecdsa.PublicKey, path string) error {
	// Not write public key to pem
	if path == "" {
		return nil
	}

	// ecdsa.Pub -> PKI pub(ans1 der) -> pem
	buf, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		errors.WithMessage(err, "marshal to pki public key")
	}

	ecf, err := os.Create(path)
	if err != nil {
		return errors.WithMessage(err, "create ecdsa file")
	}
	defer ecf.Close()

	if err := pem.Encode(ecf, &pem.Block{Type: "", Bytes: buf}); err != nil {
		errors.WithMessage(err, "write to public key pem file failed")
	}
	return nil
}
```

从PEM恢复公钥。

```go
func LoadECDSAPubKeyPem(path string) (*ecdsa.PublicKey, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.WithMessage(err, "open public file error")
	}
	defer f.Close()

	buf := make([]byte, 1000)
	n, err := f.Read(buf)
	if err != nil {
		return nil, errors.WithMessage(err, "read public key pem failed")
	}
	buf = buf[:n]

	block, _ := pem.Decode(buf)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		errors.WithMessage(err, "parse pkix public key")
	}

	if k, ok := key.(*ecdsa.PublicKey); !ok {
		return nil, errors.New("not a ecdsa public key")
	} else {
		return k, nil
	}
}
```

[↑top](#目录)

### 公私钥的PEM单元测试

生成ECDSA公私钥，写入PEM文件，然后再次读取，结果PASS。

```go
func TestECDSAPem(t *testing.T) {
	if err := GenECDSAKeyPem("ecdsa-priv.pem", "ecdsa-pub.pem"); err != nil {
		t.Fatal(err)
	}

	key, err := LoadECDSAPrivKeyPem("ecdsa-priv.pem")
	if err != nil {
		t.Fatal(err)
	}
	pub, err := LoadECDSAPubKeyPem("ecdsa-pub.pem")
	if  err != nil {
		t.Fatal(err)
	}

	if err := CheckECDSAKey(key, pub); err != nil {
		t.Fatal(err)
	}
}

func GenECDSAKeyPem(privPath, pubPath string) error {
	key, err := GenECDSAKey()
	if err != nil {
		return errors.WithMessage(err, "GenECDSAKey error")
	}

	if err := WriteECDSAPrivKeyPem(key, privPath); err != nil {
		return errors.WithMessage(err, "WriteECDSAPrivKeyPem")
	}

	if err := WriteECDSAPubKeyPem(&key.PublicKey, pubPath); err != nil {
		return errors.WithMessage(err, "WriteECDSAPubKeyPem")
	}


	return nil
}
```


[↑top](#目录)
