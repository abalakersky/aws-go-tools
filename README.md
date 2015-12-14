# Tools to work with AWS written in GoLang


[AWS Range Parser](src/aws_range_parser) Provides command line shortcut to pull IP ranges used by AWS for different services. Helpful for Firewall configurations

[S3 Bucket List](src/s3_list) Simple utility to list all buckets in the account

[S3 Manifester with Search](src/s3_manifester) and [Instructions](src/s3_manifester/s3_manifester.md)

[S3 Simple Object List](src/s3_list) Simple object list on a small bucket. Not useful on buckets with more then 1000 objects

[S3 Object List](src/s3_list) Manifester for large buckets. Outputs keys of every object in the bucket to STOUT and to a manifest file.
