resource "azurerm_mssql_employee" "d" {
  admin_user = "sa"
  admin_password = "yourStrong(!)Password1"
  server = "localhost"
  port = "1433"
  database_name = "tempdb"
}
