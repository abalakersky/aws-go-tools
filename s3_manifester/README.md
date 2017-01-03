# S3 Manifester binary

s3_manifester is a tool that provides ability to index all files in a specific S3 bucket.

## Options

running s3_manifester with `-h` switch provides the following usage information:

```bash
$ ./s3_manifester -h

Usage of ./s3_manifester:

  -akid string
    AWS Access Key
  -bucket string
    Bucket Name to list objects from. REQUIRED
  -creds string
    Credentials Profile to use [default "default"]
  -csv string
    Create CSV log with full output [default "no"]
  -file string
    Save output to file instead of displaying on the screen [default "yes"]
  -region string
    Region to connect to. [default "us-east-1"]
  -search string
    Search string to find in object paths
  -seckey string
    AWS Secret Access Key
```

- `-bucket` - the name of the bucket to be scanned. This information is required, and is the only required option besides credentials
- `-region` - specifies AWS region that bucket is located in. Defaults to `us-east-1` unless specified on command line
- `-file` - by default Manifester will create a log file with the list of all files in the bucket. If wanted, this option can be specified as `no` and Manifester will display output on the screen instead
- `-search` - provides ability to search the index and only display entries containing search string in its path
- `-csv` - an option to output more information in csv format (to screen or to file with .csv extension) with patch, size, and etag for each object

## Credentials

Manifester uses default shared credentials, usually configured if you have AWS CLI installed and used, [as listed here](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).
Shared Credentials file can also be manually created by following [these instructions](https://github.com/aws/aws-sdk-go/wiki/configuring-sdk#creating-the-credentials-file)

If you do have shared credentials file Manifester will default to use `"default"` profile. However, if you have more then one profile you can specify which one to use with `-creds` option on the command line

If, however, you do not have shared credentials file, you can specify API credentials on the command line using options `-akid` for Access Key and `-secKey` for Secret Access Key.

### Examples

Following are some examples of s3_manifester usage:

- create index file for bucket "Bucket-Name" in your default account

```bash
 ./s3_manifester -bucket Bucket-Name
```

- create index file for bucket "Client-Bucket" in client's account

```bash
 ./s3_manifester -bucket Client-Bucket -creds ClientProfile
```

- Display list of files that have ".pdf" in their name, on the screen, for bucket "Some-Bucket" in second client's account

```bash
 ./s3_manifester -bucket Some-Bucket -creds SecondClient -file no -search .pdf
```

- create index file of Third-Bucket for new customer

```bash
 ./s3_manifester -bucket Third-Bucket -akid My_Access_Key -secKey My_Secret_Access_Key
```

- create a CSV file with full detailed information about every file in "Client-Bucket" in client's account

```bash
 ./s3_manifester -bucket Client-Bucket -creds ClientProfile -csv yes
```
