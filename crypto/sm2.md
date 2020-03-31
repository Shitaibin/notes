[*crypto目录*](https://github.com/Shitaibin/notes/tree/master/crypto)
------

使用[`gmsm`](https://github.com/tjfoc/gmsm/)进行一些国密和标准加密的对比。

### 目录

- [密钥生成与PEM](#密钥生成与PEM)
- [密钥检查](#密钥检查)
- [测试SM2使用ECDSA私钥](#测试SM2使用ECDSA私钥)
- [测试ECDSA使用SM2私钥](#测试ECDSA使用SM2私钥)

### 密钥生成与PEM

```go
func TestSm2GenerateECKey(t *testing.T) {
	key, err := sm2.GenerateKey() // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}

	CheckSMKey(t, key)

	ok, err := sm2.WritePrivateKeytoPem("priv.pem", key, nil) // 生成密钥文件
	if ok != true {
		log.Fatal(err)
	}

	pubKey, _ := key.Public().(*sm2.PublicKey)
	ok, err = sm2.WritePublicKeytoPem("pub.pem", pubKey, nil) // 生成公钥文件
	if ok != true {
		log.Fatal(err)
	}
}
```

[↑top](#目录)

### 密钥检查

使用公钥加密和私钥解密，检查公私钥对。

```go
func CheckSMKey(t *testing.T, priv *sm2.PrivateKey) {
	pub := &priv.PublicKey
	msg := []byte("123456")
	d0, err := pub.Encrypt(msg)
	if err != nil {
		t.Errorf("Error: failed to encrypt %s: %v\n", msg, err)
	}

	d1, err := priv.Decrypt(d0)
	if err != nil {
		t.Errorf("Error: failed to decrypt: %v\n", err)
	}

	if !bytes.Equal(d1, msg) {
		t.Fatalf("PrivateKey is mismatched with PublicKey")
	}
}
```

[↑top](#目录)

### 测试SM2使用ECDSA私钥

SM2和ECDSA都是基于椭圆曲线的签名算法。

对于基于椭圆曲线的非对称加密算法：`Q = d * G`，其中d是私钥，是从一定取值范围中选择出来的一个随机数，G是基点，Q是公钥。

所以，私钥本质是一个随机数，而公钥是使用椭圆曲线参数和私钥计算出来，在椭圆曲线上的点。

SM2和ECDSA的私钥随机数取值范围不同，哪个范围更大呢，范围小的私钥可以作为范围大的私钥，反过来则可能不行。

`gmsm`的实现
```go
func TestECDSA2SmPem(t *testing.T) {
	eckey, err := GenECDSAKey()
	if err != nil {
		t.Fatal(err)
	}
	CheckECDSAKey(eckey)
	if err := WriteECDSAPrivKeyPem(eckey, "ecdsa-priv.pem"); err != nil {
		t.Error(err)
	}

	smkey, err := sm2.ReadPrivateKeyFromPem("ecdsa-priv.pem", nil)
	CheckSMKey(t, smkey)

	CmpSmECDSAKey(t, smkey, eckey)
}

func CmpSmECDSAKey(t *testing.T, smkey *sm2.PrivateKey, eckey *ecdsa.PrivateKey) {
	if eckey.D.Cmp(smkey.D) == 0 {
		t.Logf("Sm's key equals ecdsa's key")
	}
	if eckey.X.Cmp(smkey.X) != 0 || eckey.Y.Cmp(smkey.Y) != 0 {
		t.Logf("Sm's public key not ecdsa's public key")
	}
	t.Logf("ECDSA key: \nd: %v \nx: %v \ny: %x",
		eckey.D.String(), eckey.X.String(), eckey.Y.String())
	t.Logf("SM key: \nd: %v \nx: %v \ny: %x",
		smkey.D.String(), smkey.X.String(), smkey.Y.String())
}
```

结果：ECDSA的私钥可以作为SM2的私钥，但2种曲线参数不同，所以计算出的公钥不同。

```
=== RUN   TestECDSA2SmPem
--- PASS: TestECDSA2SmPem (0.02s)
    core_test.go:221: Sm's key equals ecdsa's key
    core_test.go:224: Sm's public key not ecdsa's public key
    core_test.go:226: ECDSA key: 
        d: 8494798200422691620569254120599992833974414226270812979695404970688056822702 
        x: 58587531562374079467902292251685072723903412227336782448333595234457170649938 
        y: 3430353433343131343732333331373932303035383131323230333934313336353739303432303938353036353033333832363934343831313830313236313136343830313638303736313434
    core_test.go:228: SM key: 
        d: 8494798200422691620569254120599992833974414226270812979695404970688056822702 
        x: 76996726420467879635691605013885200246944046645551263951175762448586093361391 
        y: 313130363136353937323938323838353137353038323738363337313635343737343030363037303537313633323434343939383936353837313336363036383635323738303139313334353539
PASS
```

[↑top](#目录)

### 测试ECDSA使用SM2私钥



```go
func TestSm2ECDSAPem(t *testing.T) {
	smkey, err := sm2.GenerateKey() // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}

	CheckSMKey(t, smkey)

	ok, err := sm2.WritePrivateKeytoPem("priv.pem", smkey, nil) // 生成密钥文件
	if ok != true {
		log.Fatal(err)
	}

	eckey, err := LoadECDSAPrivateKeyPem("priv.pem")
	if err != nil {
		t.Fatal(err)  // panic here: not a ecdsa private key
	}

	if err := CheckECDSAKey(eckey); err != nil {
		t.Fatal(err)
	}

	CmpSmECDSAKey(t, smkey, eckey)
}
```

结果：ECDSA不能解析SM2私钥，因为Go的标准库，不识别SM2的椭圆曲线标识。并不能直接证明SM2的私钥范围，超出了ECDSA的私钥范围。

```
=== RUN   TestSm2ECDSAPem
--- FAIL: TestSm2ECDSAPem (0.01s)
    core_test.go:193: not a ecdsa private key
FAIL
```


[↑top](#目录)
