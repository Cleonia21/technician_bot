package xmlToDB

import (
	"encoding/xml"
	"log"
	"main/db"
	"main/utils"
	"os"
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
	MxCell []db.Line `xml:"mxCell"`
}

func XMLToDB(fileName string, dataBase *db.Data) {
	mxFile := new(Mxfile)

	data, err := os.ReadFile("xmlToDB/xml/" + fileName + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	err = xml.Unmarshal(data, &mxFile)
	if err != nil {
		log.Fatal(err)
	}

	mxCell := mxFile.Diagram.MxGraphModel.Root.MxCell

	for i, _ := range mxCell {
		mxCell[i].Value = utils.HtmlToString(mxCell[i].Value)
	}

	dataBase.Exec(fileName, mxCell)
}
