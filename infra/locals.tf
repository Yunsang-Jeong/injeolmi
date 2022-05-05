locals {
  delimiter       = var.context.delimiter
  name_tag_prefix = join(local.delimiter, [var.service_name])
}
