package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosMongoCollection_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_mongo_collection.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosMongoCollection_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "database_name", rn),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMCosmosMongoCollection_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_mongo_collection.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosMongoCollection_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "database_name", rn),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMCosmosMongoCollection_update(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_mongo_collection.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosMongoCollection_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "database_name", rn),
				),
			},
			{
				Config: testAccAzureRMCosmosMongoCollection_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "database_name", rn),
					//todo check set values when the SDK actually reads them
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMCosmosMongoCollection_updated(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "database_name", rn),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMCosmosMongoCollection_debug(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_mongo_collection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosMongoCollection_debug(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoCollectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMCosmosMongoCollectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).cosmosAccountsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for rn, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmos_mongo_collection" {
			continue
		}

		if err := tf.AccCheckResourceAttributes(rs.Primary.Attributes, "name", "resource_group_name", "account_name", "database_name"); err != nil {
			return fmt.Errorf("resource %s is missing an attribute: %v", rn, err)
		}
		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		database := rs.Primary.Attributes["database_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos Mongo Collection %s (account %s, database %s) still exists:\n%v", name, account, database, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Mongo Collection %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosMongoCollectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).cosmosAccountsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if err := tf.AccCheckResourceAttributes(rs.Primary.Attributes, "name", "resource_group_name", "account_name", "database_name"); err != nil {
			return fmt.Errorf("resource %s is missing an attribute: %v", resourceName, err)
		}
		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		database := rs.Primary.Attributes["database_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (account: '%s', database: %s) does not exist", name, account, database)
		}

		return nil
	}
}

func testAccAzureRMCosmosMongoCollection_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmos_mongo_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmos_mongo_database.test.account_name}"
  database_name       = "${azurerm_cosmos_mongo_database.test.name}"
}
`, testAccAzureRMCosmosMongoDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosMongoCollection_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmos_mongo_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmos_mongo_database.test.account_name}"
  database_name       = "${azurerm_cosmos_mongo_database.test.name}"

  default_ttl_seconds = 707
  shard_key           = "seven"

  indexes {
    key = "seven"
  }

  indexes {
    key    = "day"
    unique = false
  }

  indexes {
    key    = "fool"
    unique = false
  }
}
`, testAccAzureRMCosmosMongoDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosMongoCollection_updated(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmos_mongo_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmos_mongo_database.test.account_name}"
  database_name       = "${azurerm_cosmos_mongo_database.test.name}"

  default_ttl_seconds = 70707
  shard_key           = "days"

  indexes {
    key    = "seven"
    unique = false
  }

  indexes {
    key    = "day"
    unique = true
  }

  indexes {
    key = "fools"
  }
}
`, testAccAzureRMCosmosMongoDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosMongoCollection_debug(rInt int, location string) string {
	return fmt.Sprintf(`


resource "azurerm_cosmos_mongo_collection" "test" {
  name                = "seven-day-tables-cola"
  resource_group_name = "kt-cosmos-201905"
  account_name        = "kt-cosmos-mongo"
  database_name       = "SevenDayDBs22"

  default_ttl_seconds = 10000

  indexes {
    key = "seven"
    unique = false
  }

  indexes {
    key = "seven11"
  }
}
`)
}

func testAccAzureRMCosmosMongoCollection_debug2(rInt int, location string) string {
	return fmt.Sprintf(`


resource "azurerm_cosmos_mongo_collection" "test" {
  name                = "seven-day-tables-more123ugg"
  resource_group_name = "kt-cosmos-201905"
  account_name        = "kt-cosmos-mongo"
  database_name       = "SevenDayDBs"

  indexes {
    key = "seven"
	unique = false
  }
}
`)
}
