package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Endpoint             string `json:"endpoint"`
	Region               string `json:"region"`
	Bucket               string `json:"bucket"`
	AccessKeyID          string `json:"access_key_id"`
	SecretAccessKey      string `json:"secret_access_key"`
	BasePath             string `json:"base_path"`
	ForcePathStyle       bool   `json:"force_path_style"`
	PublicBaseURL        string `json:"public_base_url"`
	PresignExpireMinutes int    `json:"presign_expire_minutes"`
}

type S3 struct {
	client  *s3.Client
	presign *s3.PresignClient
	cfg     S3Config
}

func NewS3(ctx context.Context, cfg S3Config) (*S3, error) {
	if cfg.Region == "" {
		cfg.Region = "auto"
	}
	if cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("s3 bucket and credentials are required")
	}
	opts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(cfg.Region), awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""))}
	if cfg.Endpoint != "" {
		opts = append(opts, awsconfig.WithBaseEndpoint(cfg.Endpoint))
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("load s3 config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) { o.UsePathStyle = cfg.ForcePathStyle })
	return &S3{client: client, presign: s3.NewPresignClient(client), cfg: cfg}, nil
}
func (s *S3) Kind() string { return "s3" }
func (s *S3) objectKey(key string) string {
	return strings.Trim(s.cfg.BasePath, "/") + func() string {
		if s.cfg.BasePath != "" {
			return "/"
		}
		return ""
	}() + strings.TrimLeft(key, "/")
}
func (s *S3) Put(ctx context.Context, key string, r io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(s.cfg.Bucket), Key: aws.String(s.objectKey(key)), Body: r, ContentLength: aws.Int64(size)})
	return err
}
func (s *S3) Open(ctx context.Context, key string, opt *OpenOptions) (io.ReadCloser, error) {
	in := &s3.GetObjectInput{Bucket: aws.String(s.cfg.Bucket), Key: aws.String(s.objectKey(key))}
	if opt != nil && (opt.Offset > 0 || opt.Length > 0) {
		end := ""
		if opt.Length > 0 {
			end = fmt.Sprintf("-%d", opt.Offset+opt.Length-1)
		}
		in.Range = aws.String(fmt.Sprintf("bytes=%d%s", opt.Offset, end))
	}
	out, err := s.client.GetObject(ctx, in)
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}
func (s *S3) Stat(ctx context.Context, key string) (*ObjectInfo, error) {
	out, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(s.cfg.Bucket), Key: aws.String(s.objectKey(key))})
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{Key: key, Size: aws.ToInt64(out.ContentLength), ModTime: aws.ToTime(out.LastModified)}, nil
}
func (s *S3) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.cfg.Bucket), Key: aws.String(s.objectKey(key))})
	return err
}
func (s *S3) PresignURL(ctx context.Context, key, filename string, expire time.Duration) (string, error) {
	if s.cfg.PublicBaseURL != "" {
		return strings.TrimRight(s.cfg.PublicBaseURL, "/") + "/" + strings.TrimLeft(s.objectKey(key), "/"), nil
	}
	if expire <= 0 {
		expire = time.Duration(s.cfg.PresignExpireMinutes) * time.Minute
	}
	out, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.cfg.Bucket), Key: aws.String(s.objectKey(key)), ResponseContentDisposition: aws.String("attachment; filename*=UTF-8''" + url.PathEscape(filename))}, func(o *s3.PresignOptions) { o.Expires = expire })
	if err != nil {
		return "", err
	}
	return out.URL, nil
}
