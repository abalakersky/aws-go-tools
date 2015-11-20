// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package s3

import (
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/private/waiter"
)

func (c *S3) WaitUntilBucketExists(input *HeadBucketInput) error {
	waiterCfg := waiter.Config{
		Operation:   "HeadBucket",
		Delay:       5,
		MaxAttempts: 20,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "status",
				Argument: "",
				Expected: 200,
			},
			{
				State:    "success",
				Matcher:  "status",
				Argument: "",
				Expected: 403,
			},
			{
				State:    "retry",
				Matcher:  "status",
				Argument: "",
				Expected: 404,
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}

func (c *S3) WaitUntilBucketNotExists(input *HeadBucketInput) error {
	waiterCfg := waiter.Config{
		Operation:   "HeadBucket",
		Delay:       5,
		MaxAttempts: 20,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "status",
				Argument: "",
				Expected: 404,
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}

func (c *S3) WaitUntilObjectExists(input *HeadObjectInput) error {
	waiterCfg := waiter.Config{
		Operation:   "HeadObject",
		Delay:       5,
		MaxAttempts: 20,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "status",
				Argument: "",
				Expected: 200,
			},
			{
				State:    "retry",
				Matcher:  "status",
				Argument: "",
				Expected: 404,
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}

func (c *S3) WaitUntilObjectNotExists(input *HeadObjectInput) error {
	waiterCfg := waiter.Config{
		Operation:   "HeadObject",
		Delay:       5,
		MaxAttempts: 20,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "status",
				Argument: "",
				Expected: 404,
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}
