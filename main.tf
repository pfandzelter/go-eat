data "aws_region" "current" {}

resource "aws_dynamodb_table" "go-eat-table" {
  name           = "go-eat"
  hash_key       = "canteen"
  range_key = "date"
  billing_mode = "PROVISIONED"
  write_capacity = 1
  read_capacity = 1

  attribute {
    name = "canteen"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }
}

resource "aws_lambda_function" "go-eat" {
  function_name    = "go-eat"
  filename         = "go-eat.zip"
  handler          = "go-eat"
  source_code_hash = filebase64sha256("go-eat.zip")
  role             = aws_iam_role.go-eat-role.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 20

  environment {
    variables = {
      DYNAMODB_TABLE = aws_dynamodb_table.go-eat-table.name,
      DYNAMODB_REGION = data.aws_region.current.name
    }
  }
}

resource "aws_iam_role" "go-eat-role" {
  name               = "go-eat"
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": {
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    },
    "Effect": "Allow"
  }
}
POLICY
}

resource "aws_iam_role_policy_attachment" "go-eat-basic-exec-role" {
  role       = aws_iam_role.go-eat-role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "go-eat-lambda_logging" {
  name = "go-eat-lambda_logging"
  path = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "go-eat-dynamo" {
  name = "go-eat-dynamo"
  path = "/"
  description = "IAM policy for DynamoDB access from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Stmt1582485790003",
      "Action": [
        "dynamodb:PutItem"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:dynamodb:*:*:*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "go-eat-lambda_logs" {
  role = aws_iam_role.go-eat-role.name
  policy_arn = aws_iam_policy.go-eat-lambda_logging.arn
}

resource "aws_iam_role_policy_attachment" "go-eat-dynamo" {
  role = aws_iam_role.go-eat-role.name
  policy_arn = aws_iam_policy.go-eat-dynamo.arn
}

# we want to run this on weekdays between 7am and 4pm, every full hour
resource "aws_cloudwatch_event_rule" "go-eat-cron" {
  name                = "go-eat-cron"
  schedule_expression = "cron(0 10 ? * 2-6 *)"
}

resource "aws_cloudwatch_event_target" "go-eat-lambda" {
  target_id = "runLambda"
  rule      = aws_cloudwatch_event_rule.go-eat-cron.name
  arn       = aws_lambda_function.go-eat.arn
}

resource "aws_lambda_permission" "go-eat-cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.go-eat.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.go-eat-cron.arn
}
