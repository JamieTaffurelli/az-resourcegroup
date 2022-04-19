resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.location

  lifecycle {
    prevent_destroy = true
  }

  tags     = var.tags
}

resource "azurerm_management_lock" "rg" {
  name       = "${var.resource_group_name}-CanNotDelete"
  scope      = azurerm_resource_group.rg.id
  lock_level = "CanNotDelete"
  notes      = "Prevents this resource group from being deleted"
}