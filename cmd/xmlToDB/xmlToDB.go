package xmlToDB

import (
	"encoding/xml"
	"log"
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

func XMLToDB(fileName string) {
	mxFile := new(Mxfile)

	data, err := os.ReadFile("cmd/xmlToDB/xml/" + fileName + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	err = xml.Unmarshal(data, &mxFile)
	if err != nil {
		log.Fatal(err)
	}

	mxCell := mxFile.Diagram.MxGraphModel.Root.MxCell

	for i := range mxCell {
		mxCell[i].Value = utils.HtmlToString(mxCell[i].Value)
	}

	_ = database.DropTable(fileName)
	_ = database.CreateTable(fileName)
	_ = database.InsertLines(fileName, mxCell)
}
