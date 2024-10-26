package xmlToDB

import (
	"encoding/xml"
	"technician_bot/database"

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

	err = database.DropTable(tableName)
	if err = database.CreateTable(tableName); err != nil {
		return err
	}
	if err = database.InsertLines(tableName, mxCell); err != nil {
		return err
	}
	return nil
}
