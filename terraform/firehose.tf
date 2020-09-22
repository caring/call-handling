module "firehose" {
  source = "git@github.com:caring/tf-modules.git//aws/firehose_monitoring?ref=v1.6.0"
  
  service_name              = local.service_name
  task_role_name            = aws_iam_role.ecs_task_role.name

  tags = local.tags
}
