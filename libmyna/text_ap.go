package libmyna

import (
	"errors"
	"fmt"
	"github.com/jpki/myna/asn1"
	"strconv"
)

type TextAP struct {
	reader *Reader
}

type TextAttrs struct {
	Header  []byte `asn1:"private,tag:33"`
	Name    string `asn1:"private,tag:34,utf8"`
	Address string `asn1:"private,tag:35,utf8"`
	Birth   string `asn1:"private,tag:36"`
	Sex     string `asn1:"private,tag:37"`
}

type TextSignature struct {
	MyNumDigest []byte `asn1:"private,tag:49"`
	AttrsDigest []byte `asn1:"private,tag:50"`
	Signature   []byte `asn1:"private,tag:51"`
}

type TextCertificate struct {
	Raw []byte `asn1:"application,tag:78"`
}

type TextBasicInfo struct {
	APInfo []byte `asn1:"private,tag:65"`
	KeyID  []byte `asn1:"private,tag:66"`
}

func (self *TextAP) LookupPin() (int, error) {
	err := self.reader.SelectEF("0011") // 券面事項入力補助用PIN
	if err != nil {
		return 0, err
	}
	count := self.reader.LookupPin()
	return count, nil
}

func (self *TextAP) VerifyPin(pin string) error {
	err := self.reader.SelectEF("0011")
	if err != nil {
		return err
	}
	err = self.reader.Verify(pin)
	return err
}

func (self *TextAP) LookupPinA() (int, error) {
	err := self.reader.SelectEF("0014")
	if err != nil {
		return 0, err
	}
	count := self.reader.LookupPin()
	return count, nil
}

func (self *TextAP) VerifyPinA(pin string) error {
	err := self.reader.SelectEF("0014")
	if err != nil {
		return err
	}
	err = self.reader.Verify(pin)
	return err
}

func (self *TextAP) LookupPinB() (int, error) {
	err := self.reader.SelectEF("0015")
	if err != nil {
		return 0, err
	}
	count := self.reader.LookupPin()
	return count, nil
}

func (self *TextAP) VerifyPinB(pin string) error {
	err := self.reader.SelectEF("0015")
	if err != nil {
		return err
	}
	err = self.reader.Verify(pin)
	return err
}

func (self *TextAP) ReadMyNumber() (string, error) {
	err := self.reader.SelectEF("0001")
	if err != nil {
		return "", err
	}
	data := self.reader.ReadBinary(17)
	var mynumber asn1.RawValue
	_, err = asn1.UnmarshalWithParams(data, &mynumber, "private,tag:16")
	if err != nil {
		return "", err
	}
	return string(mynumber.Bytes), nil
}

func (self *TextAP) ReadAttributes() (*TextAttrs, error) {
	err := self.reader.SelectEF("0002")
	if err != nil {
		return nil, err
	}

	data := self.reader.ReadBinary(7)
	if len(data) != 7 {
		return nil, errors.New("Error at ReadBinary()")
	}

	parser := ASN1PartialParser{}
	err = parser.Parse(data)
	if err != nil {
		return nil, err
	}
	data = self.reader.ReadBinary(parser.GetSize())
	var attrs TextAttrs
	_, err = asn1.UnmarshalWithParams(data, &attrs, "private,tag:32")
	if err != nil {
		return nil, err
	}
	return &attrs, nil
}

func (self *TextAP) ReadSignature() (*TextSignature, error) {
	err := self.reader.SelectEF("0003")
	if err != nil {
		return nil, err
	}
	data := self.reader.ReadBinary(336)
	if len(data) != 336 {
		return nil, errors.New("Error at ReadBinary()")
	}
	var signature TextSignature
	_, err = asn1.UnmarshalWithParams(data, &signature, "private,tag:48")
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (self *TextAP) ReadCertificate() (*TextCertificate, error) {
	err := self.reader.SelectEF("0004")
	if err != nil {
		return nil, err
	}
	data := self.reader.ReadBinary(568)
	if len(data) != 568 {
		return nil, errors.New("Error at ReadBinary()")
	}
	var certificate TextCertificate
	_, err = asn1.UnmarshalWithParams(data, &certificate, "application,tag:33")
	if err != nil {
		return nil, err
	}
	return &certificate, nil
}

func (self *TextAP) ReadBasicInfo() (*TextBasicInfo, error) {
	err := self.reader.SelectEF("0005")
	if err != nil {
		return nil, err
	}
	data := self.reader.ReadBinary(256)
	if len(data) != 256 {
		return nil, errors.New("Error at ReadBinary()")
	}
	var basicInfo TextBasicInfo
	_, err = asn1.UnmarshalWithParams(data, &basicInfo, "private,tag:64")
	if err != nil {
		return nil, err
	}

	if len(basicInfo.APInfo) != 4 {
		return nil, errors.New("invalid APInfo length")
	}
	if len(basicInfo.KeyID) != 16 {
		return nil, errors.New("invalid KeyID length")
	}
	return &basicInfo, nil
}

// ヘッダーをHEX文字列に変換
func (self *TextAttrs) HeaderString() string {
	return fmt.Sprintf("% X", self.Header)
}

// ISO5218コードから日本語文字列に変換
func (self *TextAttrs) SexString() string {
	n, err := strconv.Atoi(self.Sex)
	if err != nil {
		return "エラー"
	}
	switch n {
	case 1:
		return "男性"
	case 2:
		return "女性"
	case 9:
		return "適用不能"
	default:
		return "不明"
	}
}
