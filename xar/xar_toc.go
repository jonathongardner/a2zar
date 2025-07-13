package xar

import (
	"encoding/xml"
)

type xarStyleAttr struct {
	Text  string `xml:",chardata"`
	Style string `xml:"style,attr"`
}

type xarOffsetSize struct {
	xarStyleAttr
	Offset string `xml:"offset"`
	Size   string `xml:"size"`
}

type xarSignature struct {
	xarOffsetSize
	KeyInfo struct {
		XmlNamespace string `xml:"xmlns,attr"`
		X509Data     struct {
			X509Certificate []string `xml:"X509Certificate"`
		} `xml:"X509Data"`
	} `xml:"KeyInfo"`
}

type xarToc struct {
	XMLName xml.Name `xml:"xar"`
	Toc     struct {
		CreationTime string        `xml:"creation-time"`
		Checksum     xarOffsetSize `xml:"checksum"`
		Signature    xarSignature  `xml:"signature"`
		XSignature   xarSignature  `xml:"x-signature"`
		File         []xarFile     `xml:"file"`
	} `xml:"toc"`
}

type xarFileData struct {
	Length            int64        `xml:"length"`
	Offset            int64        `xml:"offset"`
	Size              int64        `xml:"size"`
	Encoding          xarStyleAttr `xml:"encoding"`
	ExtractedChecksum xarStyleAttr `xml:"extracted-checksum"`
	ArchivedChecksum  xarStyleAttr `xml:"archived-checksum"`
}

type xarFile struct {
	Id               string      `xml:"id,attr"`
	Data             xarFileData `xml:"data"`
	Type             string      `xml:"type"`
	Name             string      `xml:"name"` // this was an array?
	Link             string      `xml:"link"`
	FinderCreateTime struct {
		NanoSeconds string `xml:"nanoseconds"`
		Time        string `xml:"time"`
	} `xml:"FinderCreateTime"`
	CTime    string    `xml:"ctime"`
	MTime    string    `xml:"mtime"`
	ATime    string    `xml:"atime"`
	Group    string    `xml:"group"`
	Gid      string    `xml:"gid"`
	User     string    `xml:"user"`
	Uid      string    `xml:"uid"`
	Mode     string    `xml:"mode"`
	DeviceNo string    `xml:"deviceno"`
	Inode    string    `xml:"inode"`
	File     []xarFile `xml:"file"`
}

func (xf *xarFile) fileInfo(dir string) []*XarFileInfo {
	order := make([]*XarFileInfo, 0, len(xf.File))
	toAdd := newFileInfo(dir, xf)
	order = append(order, toAdd)

	if xf.Type == typeDirectory {
		for _, df := range xf.File {
			order = append(order, df.fileInfo(toAdd.Path())...)
		}
	}

	return order
}
