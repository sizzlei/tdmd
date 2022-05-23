package lib


import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	GetTable = `
		SELECT
			table_name,
			table_type,
			ENGINE,
			row_format,
			table_collation,
			table_comment 
		FROM
			information_schema.TABLES 
		WHERE
			%s
	`
	GetColumn = `
		SELECT
			concat("|",
			column_name,"|",
			column_type,"|",
			is_nullable,"|",
			IF(column_default IS NULL,"",column_default),"|",
			IF(character_set_name IS NULL,"",character_set_name),"|",
			IF(collation_name IS NULL,"",collation_name),"|",
			IF(column_key = "", "", column_key),"|",
			IF(extra = "","",extra),"|",
			IF(column_comment = "", "",replace(column_comment,"\n","")),"|")
		FROM
			information_schema.COLUMNS 
		WHERE
			table_schema ='%s'
			and table_name ='%s'
		ORDER BY
		ordinal_position;
	`
	GetIndex = `
		SELECT
			concat(
				IF(NON_UNIQUE=1,"[Normal] ","[Unique] "),
				INDEX_NAME,
				" (",GROUP_CONCAT( COLUMN_NAME ORDER BY SEQ_IN_INDEX ASC SEPARATOR ',' ),")"
			) 
		FROM
			information_schema.STATISTICS 
		WHERE
			TABLE_SCHEMA = '%s' 
		 	AND TABLE_NAME = '%s' 
			AND INDEX_NAME != 'PRIMARY' 
		GROUP BY
			TABLE_SCHEMA,
			TABLE_NAME,
			INDEX_NAME 
		ORDER BY
			INDEX_NAME
	`
	GetConstraint = `
		SELECT
			concat(
				x.constraint_name,
				" (",
				group_concat( x.column_name ),
				") <-- ",
				concat( x.referenced_table_name, '.', x.referenced_column_name ),
				"\n    - *ON Delete* : ",
				y.DELETE_RULE,
				" / *ON Update* : ",
				y.UPDATE_RULE 
			) 
		FROM
			information_schema.KEY_COLUMN_USAGE x
			INNER JOIN information_schema.REFERENTIAL_CONSTRAINTS y ON x.constraint_name = y.constraint_name 
		WHERE
			x.CONSTRAINT_SCHEMA = '%s'	
			AND x.table_name = '%s'
			AND x.constraint_name <> 'PRIMARY' 
		GROUP BY
			x.constraint_name
	`
)


func CreateDBobject(Conf Configure) (*sql.DB, error) {
	var dbObj *sql.DB 
	DSNFormat := "%s:%s@tcp(%s:%d)/information_schema"
	DSN := fmt.Sprintf(DSNFormat,Conf.User,Conf.Pass,Conf.Endpoint,Conf.Port)
	dbObj, err := sql.Open("mysql",DSN)
	if err != nil {
		return dbObj, err
	}
	return dbObj, nil
}

func MakeinCondition(Str []string) string{
	for i,v := range Str {
		Str[i] = fmt.Sprintf("'%s'",v)
	}
	z := strings.Join(Str,",")

	return z
}


func GetDefinition(Conf Configure, Export Export)  map[string][]Tableinfo{
	dbObj, err := CreateDBobject(Conf)
	if err != nil {
		log.Errorf("Failed to Create Database Object. %s",err)
		os.Exit(1)
	}

	defer dbObj.Close()

	var dbSchemata map[string][]Tableinfo
	dbSchemata = make(map[string][]Tableinfo)

	for _, v := range Export.DB {
		// Make Where Condition
		var dbQueries string
		if Export.DBcnt == 1 && len(Export.Table) >= 1  {
			tIn := MakeinCondition(Export.Table)
			dbQueries = fmt.Sprintf(GetTable,fmt.Sprintf(" table_schema = '%s' and table_name in (%s)",v,tIn))
		} else {
			dbQueries = fmt.Sprintf(GetTable,fmt.Sprintf(" table_schema = '%s'",v))
		}

		// Get Table Metadata
		dbMeta, err := dbObj.Query(dbQueries)
		if err != nil {
			log.Errorf("Failed to Get %s Database Metadata. %s",v,err)
			os.Exit(1)
		}

		for dbMeta.Next() {
			var Table Tableinfo
			err := dbMeta.Scan(&Table.TableName,&Table.TableType,&Table.Engine,&Table.RowFormat,&Table.Collation,&Table.Comment)
			if err != nil {
				log.Errorf("Failed to Get Table Metadata. %s",v,err)
				os.Exit(1)
			}

			// Get Column Metadata
			colMeta, err := dbObj.Query(fmt.Sprintf(GetColumn,v,Table.TableName))
			if err != nil{
				log.Errorf("Failed to Get %s Column Metadata. %s",Table.TableName,err)
			}
			for colMeta.Next() {
				var Column string
				err := colMeta.Scan(&Column)
				if err != nil {
					log.Errorf("Failed to Scan %s Column Metadata. %s",Table.TableName,err)
					os.Exit(1)
				}
				Table.Columns = append(Table.Columns,Column)
			}
			

			// Get Secondary Index
			indexMeta, err := dbObj.Query(fmt.Sprintf(GetIndex,v,Table.TableName))
			if err != nil{
				log.Errorf("Failed to Get %s Index Metadata. %s",Table.TableName,err)
			}
			for indexMeta.Next() {
				var Index string 
				err := indexMeta.Scan(&Index)
				if err != nil{
					log.Errorf("Failed to Scan %s Index Metadata. %s",Table.TableName,err)
				}
				Table.Indexes = append(Table.Indexes,Index)
			}

			// Get Constraint
			constMeta, err := dbObj.Query(fmt.Sprintf(GetConstraint,v,Table.TableName))
			if err != nil{
				log.Errorf("Failed to Get %s Constraint Metadata. %s",Table.TableName,err)
			}
			for constMeta.Next() {
				var Constraint string
				err := constMeta.Scan(&Constraint)
				if err != nil{
					log.Errorf("Failed to Scan %s Constraint Metadata. %s",Table.TableName,err)
				}
				Table.Constraints = append(Table.Constraints,Constraint)
			}
			dbSchemata[v] = append(dbSchemata[v],Table)
		}
	}
	return dbSchemata
}

func MakeMarkdown(Path string,Def map[string][]Tableinfo) {
	now := time.Now()
	checkTime := now.Format("2006-01-02")
	
	for k, v := range Def {
		fName := fmt.Sprintf("%s/%s.log",Path,k)
		mdName := fmt.Sprintf("%s/%s.md",Path,k)

		// Open New File
		_ = os.Remove(mdName)
		f, err := os.OpenFile(fName,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Errorf("Failed to Create %s File",k)
		}
		// Database Name
		Writefile(f,fmt.Sprintf("%s \n",k))
		Writefile(f,"=============\n")
		Writefile(f,fmt.Sprintf("**Last Update** : %s\n",checkTime))
		
		// Title
		Writefile(f,"## Table List\n")
		for _, t := range v {
			if t.Comment == "" {
				Writefile(f,fmt.Sprintf("- [%s](#%s)\n ",t.TableName,strings.ToLower(t.TableName)))
			} else {
				Writefile(f,fmt.Sprintf("- [%s (%s)](#%s)\n ",t.TableName,t.Comment,strings.ToLower(t.TableName)))
			}
		}

		// Table Data
		for _, tb := range v {
			Writefile(f,fmt.Sprintf("## %s\n",strings.ToLower(tb.TableName)))
			Writefile(f,"**Information**\n")
			Writefile(f,"|Table type|Engine|Row format|Collation|Comment|\n")
			Writefile(f,"|---|---|---|---|---|\n")
			Writefile(f,fmt.Sprintf("|%s|%s|%s|%s|%s|\n\n",
				tb.TableType,
				tb.Engine,
				tb.RowFormat,
				tb.Collation,
				tb.Comment,
			))

			// Write Column
			Writefile(f,"**Columns**\n")
			Writefile(f,"|Name|Type|Nullable|Default|Charset|Collation|Key|Extra|Comment|\n")
			Writefile(f,"|---|---|---|---|---|---|---|---|---|\n")
			for _, c := range tb.Columns {
				Writefile(f,fmt.Sprintf("%s\n",c))
			}
			
			// Write Index
			if len(tb.Indexes) > 0 {
				Writefile(f,"\n**Index**\n")
				for _, i := range tb.Indexes {
					Writefile(f,fmt.Sprintf("- %s\n",i))
				}
			}
			// Write Constraint
			if len(tb.Constraints) > 0 {
				Writefile(f,"\n**Constraint**\n")
				for _, c := range tb.Constraints {
					Writefile(f,fmt.Sprintf("- %s\n",c))
				}
			}
			Writefile(f,"\n\n")
		}
		// File Close
		if err := f.Close(); err != nil {
			panic(err)
		}

		// log --> md
		_ = os.Rename(fName,mdName)
	}
	

	
}

func Writefile(File *os.File, Data string) {
	_, _ = File.Write([]byte(Data))
}