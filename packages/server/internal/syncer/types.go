package syncer

import "context"

type Artifact struct {
	ID         string
	ExternalID string
	S3Key      string
}

type ArtifactHandler interface {
	// Process should respect ctx and return the S3 key for the completed artifact.
	Process(ctx context.Context, artifact Artifact) (s3Key string, err error)
}

type ArtifactStore interface {
	// ClaimPending marks record as "processing" and MarkDone
	ClaimPending(ctx context.Context, limit int) ([]Artifact, error)
	MarkDone(ctx context.Context, id string, s3Key string) error
	MarkFailed(ctx context.Context, id string, err error) error
}

// Source is responsible for rolling discovery + persistence of metadata.
type Source[E any] interface {
	Sync(ctx context.Context) error
}
