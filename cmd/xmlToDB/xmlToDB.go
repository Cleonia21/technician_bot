package xmlToDB

import (
	"encoding/xml"
	"technician_bot/database"

	//"main/database"
	"os"
	"technician_bot/cmd/utils"
)

type Mxfile struct {
	Diagram Diagram `xml:"diagram"`
}

type Diagram struct {
	MxGraphModel MxGraphModel `xml:"mxGraphModel"`
}

type MxGraphModel struct {
	Root Root `xml:"root"`
}

type Root struct {
	MxCell []database.Line `xml:"mxCell"`
}

func FileToDB(filePath string, tableName string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return ByteToDB(data, tableName)
}

func ByteToDB(data []byte, tableName string) error {
	mxFile := new(Mxfile)

	err := xml.Unmarshal(data, &mxFile)
	if err != nil {
		return err
	}

	mxCell := mxFile.Diagram.MxGraphModel.Root.MxCell

	for i := range mxCell {
		mxCell[i].Value = utils.HtmlToString(mxCell[i].Value)
	}

	_ = database.DropTable(tableName)
	_ = database.CreateTable(tableName)
	_ = database.InsertLines(tableName, mxCell)
	return nil
}
