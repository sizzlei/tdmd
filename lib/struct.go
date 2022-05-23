package lib

type Configure struct {
	Endpoint 		string 
	Port			int
	User 			string 
	Pass			string
}

type Export struct {
	DB				[]string
	DBcnt 			int
	Table			[]string
	FilePath		string
}

type Tableinfo struct{
	TableName 		string 
	TableType 		string
	Engine 			string 
	RowFormat		string
	Collation 		string
	Comment 		string
	Columns 		[]string
	Indexes			[]string
	Constraints		[]string
}
