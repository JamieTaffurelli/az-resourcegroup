resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.location

  lifecycle {
    prevent_destroy = true
  }

  tags     = var.tags
}

resource "azurerm_management_lock" "rg" {
  name       = "subscription-level"
  scope      = data.azurerm_resource_group.rg.id
  lock_level = "CanNotDelete"
  notes      = "Prevents resources in this resource group from being deleted"
}