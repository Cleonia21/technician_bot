package xmlToDB

import (
	"encoding/xml"
	"log"
	"technician_bot/cmd/db"

	//"main/db"
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
	MxCell []db.Line `xml:"mxCell"`
}

func XMLToDB(fileName string, dataBase *db.Data) {
	mxFile := new(Mxfile)

	data, err := os.ReadFile("cmd/xmlToDB/xml/" + fileName + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	err = xml.Unmarshal(data, &mxFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("file success read")

	mxCell := mxFile.Diagram.MxGraphModel.Root.MxCell

	for i := range mxCell {
		mxCell[i].Value = utils.HtmlToString(mxCell[i].Value)
	}

	//database.DB.Db.AutoMigrate(&Line{})
	//
	//database.DB.Db.Create(mxCell[0])
	//
	//time.Sleep(10 * time.Minute)
	dataBase.Exec(fileName, mxCell)
}
