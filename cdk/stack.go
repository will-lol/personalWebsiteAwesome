package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"fmt"
)

type CommandHook struct{}

func (t CommandHook) BeforeBundling(inputDir *string, outputDir *string) *[]*string {
	command1 := "./prebuild.sh"
	return &[]*string{&command1}
}
func (t CommandHook) AfterBundling(inputDir *string, outputDir *string) *[]*string {
	return &[]*string{}
}

func NewWebsiteStack(scope constructs.Construct, id string, props awscdk.StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props)

	subscriptions := awsdynamodb.NewTable(stack, jsii.String("subscriptions"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("Endpoint"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	metadata := awsdynamodb.NewTable(stack, jsii.String("metadata"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("ID"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	secretArn := "arn:aws:secretsmanager:ap-southeast-2:301436506805:secret:website/vapid-keys-WIfxAO"
	a := awssecretsmanager.Secret_FromSecretCompleteArn(stack, jsii.String("SecretFromCompleteArn"), jsii.String(secretArn))

	f := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("handler"), &awscdklambdagoalpha.GoFunctionProps{
		Entry: jsii.String("../"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			CommandHooks: &CommandHook{},
		},
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"SUBSCRIPTIONS_TABLE_NAME": subscriptions.TableName(),
			"METADATA_TABLE_NAME":      metadata.TableName(),
			"SECRET_ARN":								&secretArn,
		},
		ParamsAndSecrets: awslambda.ParamsAndSecretsLayerVersion_FromVersion(awslambda.ParamsAndSecretsVersions_V1_0_103, &awslambda.ParamsAndSecretsOptions{}),
	})

	a.GrantRead(f, nil)

	// Add a Function URL.
	lambda_url := f.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})
	awscdk.NewCfnOutput(stack, jsii.String("lambdaFunctionUrl"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("lambdaFunctionUrl"),
		Value:      lambda_url.Url(),
	})

	metadata.GrantFullAccess(f)
	subscriptions.GrantFullAccess(f)

	assetsBucket := awss3.NewBucket(stack, jsii.String("assets"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		EnforceSSL:        jsii.Bool(true),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		Versioned:         jsii.Bool(false),
	})
	// Allow CloudFront to read from the bucket.
	cfOAI := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("cfnOriginAccessIdentity"), &awscloudfront.OriginAccessIdentityProps{})
	cfs := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{})
	cfs.AddActions(jsii.String("s3:GetBucket*"))
	cfs.AddActions(jsii.String("s3:GetObject*"))
	cfs.AddActions(jsii.String("s3:List*"))
	cfs.AddResources(assetsBucket.BucketArn())
	cfs.AddResources(jsii.String(fmt.Sprintf("%v/*", *assetsBucket.BucketArn())))
	cfs.AddCanonicalUserPrincipal(cfOAI.CloudFrontOriginAccessIdentityS3CanonicalUserId())
	assetsBucket.AddToResourcePolicy(cfs)

	// Add a CloudFront distribution to route between the public directory and the Lambda function URL.
	lambdaURLDomain := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), lambda_url.Url(), nil))
	lambdaOrigin := awscloudfrontorigins.NewHttpOrigin(lambdaURLDomain, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	cf := awscloudfront.NewDistribution(stack, jsii.String("customerFacing"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
			Origin:               lambdaOrigin,
			CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
			OriginRequestPolicy:  awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
			CachePolicy:          awscloudfront.CachePolicy_CACHING_DISABLED(),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
	})

	// Add /assets* to the distribution backed by S3.
	assetsOrigin := awscloudfrontorigins.NewS3Origin(assetsBucket, &awscloudfrontorigins.S3OriginProps{
		// Get content from the / directory in the bucket.
		OriginPath:           jsii.String("/"),
		OriginAccessIdentity: cfOAI,
	})
	cf.AddBehavior(jsii.String("/assets*"), assetsOrigin, nil)

	// Export the domain.
	awscdk.NewCfnOutput(stack, jsii.String("cloudFrontDomain"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("cloudfrontDomain"),
		Value:      cf.DomainName(),
	})

	// Deploy the contents of the ./assets directory to the S3 bucket.
	awss3deployment.NewBucketDeployment(stack, jsii.String("assetsDeployment"), &awss3deployment.BucketDeploymentProps{
		DestinationBucket: assetsBucket,
		Sources: &[]awss3deployment.ISource{
			awss3deployment.Source_Asset(jsii.String("../assets"), nil),
		},
		DestinationKeyPrefix: jsii.String("assets"),
		Distribution:         cf,
		DistributionPaths:    jsii.Strings("/assets*"),
	})

	return stack
}

func main() {
	defer jsii.Close()
	app := awscdk.NewApp(nil)
	NewWebsiteStack(app, "WebsiteStack", awscdk.StackProps{
		Env: &awscdk.Environment{
			Region: jsii.String("ap-southeast-2"),
		},
	})
	app.Synth(nil)
}
