package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/github"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
	"golang.org/x/sync/errgroup"
)

const (
	baseImage = "cgr.dev/chainguard/static:latest"
)

var (
	amazonKeychain authn.Keychain = authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithLogger(io.Discard)))
	azureKeychain  authn.Keychain = authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper())
	keychain                      = authn.NewMultiKeychain(
		amazonKeychain,
		authn.DefaultKeychain,
		google.Keychain,
		github.Keychain,
		azureKeychain,
	)
)

func main() {
	ctx := context.Background()
	cfg, err := loadBobYAML()
	if err != nil {
		log.Fatalf("LoadYaml: %v", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	for _, koo := range cfg.Builds {
		koo := koo
		eg.Go(func() error {
			b, err := build.NewGo(ctx, ".",
				build.WithPlatforms(cfg.DefaultPlatforms...), // only build for these platforms.
				build.WithBaseImages(func(ctx context.Context, _ string) (name.Reference, build.Result, error) {
					ref := name.MustParseReference(baseImage)
					base, err := remote.Index(ref, remote.WithContext(ctx))
					return ref, base, err
				}),
			)
			if err != nil {
				return fmt.Errorf("build error: %w", err)
			}

			c, err := build.NewCaching(b)
			if err != nil {
				log.Fatalf("NewCaching: %v", err)
			}

			r, err := c.Build(ctx, koo.Main)
			if err != nil {
				log.Fatalf("Build: %v", err)
			}

			targetRepo := filepath.Join(cfg.Docker.Registry, koo.ID)
			p, err := publish.NewDefault(targetRepo, // publish to example.registry/my-repo
				publish.WithTags([]string{Commit}),     // tag with :deadbeef
				publish.WithAuthFromKeychain(keychain), // use credentials from ~/.docker/config.json
				publish.WithNamer(func(string, string) string {
					return targetRepo
				}),
			)
			if err != nil {
				log.Fatalf("NewDefault: %v", err)
			}

			ref, err := p.Publish(ctx, r, koo.Main)
			if err != nil {
				log.Fatalf("Publish: %v", err)
			}

			fmt.Println(ref.String())
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
