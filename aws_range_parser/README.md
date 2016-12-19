# AWS IP Ranges parser

This script provides an easy, CLI based, way to display IP ranges that AWS uses for different services. Such information usually needed to configure Firewall rules.
Instructions are available by running `aws_range_parser -h`, running `aws_range_parser` without any options, or in this document.

```
This script is used to display AWS specific IP ranges that could be used for Firewall or Security Group configurations. These ranges specify public IPs that AWS uses for a each public facing service.

Usage:
aws_range_parser [-h] [-region REGION_NAME] -service SERVICE_NAME

Service:
    Valid values: [AMAZON CLOUDFRONT EC2 ROUTE53 ROUTE53_HEALTHCHECKS]

Region:
    Valid values: [GLOBAL ap-northeast-1 ap-northeast-2 ap-south-1 ap-southeast-1 ap-southeast-2 ca-central-1 cn-north-1 eu-central-1 eu-west-1 eu-west-2 sa-east-1 us-east-1 us-east-2 us-gov-west-1 us-west-1 us-west-2]

Notes:
Please remember that some services, such as CloudFront and Route53 are Global and as such use only GLOBAL as their region. Their information can be gathered with or without specifying region name.
```
