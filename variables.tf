variable "resource_group_name" {
  type        = string
  description = "Name of Resource Group to deploy"
}

variable "location" {
  type        = string
  description = "Location of the Resource Group"
}

variable "tags" {
  type        = map(string)
  description = "Tags of the Resource Group"
}
