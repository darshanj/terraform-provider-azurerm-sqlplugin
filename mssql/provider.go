package mssql

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"time"
	_ "github.com/denisenkom/go-mssqldb"
	"database/sql"
	"context"
	"log"
	"errors"
)

func Provider() terraform.ResourceProvider {
	dataSources := make(map[string]*schema.Resource)

	resources := map[string]*schema.Resource{
		"sqlplugin_mssql_employee": resourceArmMsSqlUser(),
	}

	p := &schema.Provider{
		Schema: make(map[string]*schema.Schema),
		DataSourcesMap: dataSources,
		ResourcesMap:   resources,
	}
	p.ConfigureFunc = providerConfigure(p)
	return p
}

func resourceArmMsSqlUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlUserCreate,
		Read:   resourceArmMsSqlUserRead,
		Update: resourceArmMsSqlUserUpdate,
		Delete: resourceArmMsSqlUserDelete,
		//Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
		//	_, err := parse.MsSqlDatabaseID(id)
		//	return err
		//}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"admin_user": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				// TODO: ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},

			"admin_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"server": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
			},
			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				// TODO: ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},
			"port": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
			},
		},
	}
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		//user := d.Get("admin_user")
		//password := d.Get("admin_password")
		//database := d.Get("database_name")
		//server := d.Get("server")
		//port := d.Get("port")
		// Build connection string
		connString := "server=localhost;user id=sa;password=yourStrong(!)Password1;port=1433;database=tempdb;"

		var err error

		// Create connection pool
		db, err := sql.Open("sqlserver", connString)
		if err != nil {
			log.Fatal("Error creating connection pool: ", err.Error())
			return nil,err
		}
		log.Print(fmt.Sprintf("Connection opened with %s!\n",connString))

		return db,nil
	}
}

func resourceArmMsSqlUserCreate(d *schema.ResourceData, meta interface{}) error {
	var err error
	db := meta.(*sql.DB)

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Print("Not Connected!\n")
		log.Fatal(err.Error())
	}
	log.Print("Connected!\n")
	// Create employee
	createID, err := CreateEmployee(db, "Jake", "United States")
	if err != nil {
		log.Fatal("Error creating Employee: ", err.Error())
		return err
	}
	log.Print("Inserted ID: %d successfully.\n", createID)
	d.Set("ID",createID)
	return nil
}

func resourceArmMsSqlUserUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error
	db := meta.(*sql.DB)

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Connected!\n")
	// Create employee
	createID, err := UpdateEmployee(db,"Jake", "United States")
	if err != nil {
		log.Fatal("Error creating Employee: ", err.Error())
		return err
	}
	log.Print("Inserted ID: %d successfully.\n", createID)
	d.Set("ID",createID)
	return nil
}
func resourceArmMsSqlUserDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	db := meta.(*sql.DB)

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Connected!\n")
	// Delete employee
	_, err = DeleteEmployee(db, "Jake", "United States")
	return err
}

func resourceArmMsSqlUserRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	db := meta.(*sql.DB)

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Connected!\n")
	// Read employees
	id, err := ReadEmployees(db,"Jake", "United States")
	if err != nil || id == -1 {
		log.Fatal("Error reading Employees: ", err.Error())
		d.SetId("")
		return err
	}

	log.Print("Read employee with id %d successfully.\n", id)
	d.Set("ID",id)
	return nil
}

// CreateEmployee inserts an employee record
func CreateEmployee(db *sql.DB, name string, location string) (int64, error) {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("CreateEmployee: db is null")
		return -1, err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := "INSERT INTO TestSchema.Employees (Name, Location) VALUES (@Name, @Location); select convert(bigint, SCOPE_IDENTITY());"

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Name", name),
		sql.Named("Location", location))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}


// DeleteEmployee deletes an employee from the database
func DeleteEmployee(db *sql.DB, name string, location string) (int64, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("DELETE FROM TestSchema.Employees WHERE Name = @Name AND Location = @Location;")

	// Execute non-query with named parameters
	result, err := db.ExecContext(ctx, tsql,
		sql.Named("Name", name),
		sql.Named("Location", location))

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

// ReadEmployees reads all employee records
func ReadEmployees(db *sql.DB, name string, location string) (int, error) {
	ctx := context.Background()
	var id = -1

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return id, err
	}


	tsql := "SELECT Id, Name, Location FROM TestSchema.Employees WHERE Name == @Name AND Location = @Location;"

	// Execute query
	rows, err := db.QueryContext(ctx, tsql,
		sql.Named("Name", name),
		sql.Named("Location", location))

	if err != nil {
		return id, err
	}

	defer rows.Close()

	// Iterate through the result set.
	for rows.Next() {
		var name, location string
		var dbId int

		// Get values from row.
		err := rows.Scan(&id, &name, &location)
		if err != nil {
			return -1, err
		}
		log.Print("ID: %d, Name: %s, Location: %s\n", dbId, name, location)
		id = dbId
		break
	}

	return id, nil
}

// UpdateEmployee updates an employee's information
func UpdateEmployee(db *sql.DB, name string, location string) (int64, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("UPDATE TestSchema.Employees SET Location = @Location WHERE Name = @Name")

	// Execute non-query with named parameters
	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("Location", location),
		sql.Named("Name", name))
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}