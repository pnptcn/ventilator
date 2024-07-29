package data

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"capnproto.org/go/capnp/v3"
	"github.com/minio/minio-go/v7"
)

type ArtifactService struct {
	minioClient *minio.Client
	bucket      string
}

func NewArtifactService(minioClient *minio.Client, bucket string) *ArtifactService {
	return &ArtifactService{
		minioClient: minioClient,
		bucket:      bucket,
	}
}

func (s *ArtifactService) CreateArtifact(ctx context.Context, header Artifact_Header) (*Artifact, error) {
	var (
		artifact Artifact
		seg      *capnp.Segment
		err      error
	)

	if _, seg, err = capnp.NewMessage(capnp.SingleSegment(nil)); err != nil {
		return nil, err
	}

	if artifact, err = NewRootArtifact(seg); err != nil {
		return nil, err
	}

	if err = artifact.SetHeader(header); err != nil {
		return nil, err
	}

	// Generate a unique ID
	id := generateUniqueID()
	header.SetId(id)

	return &artifact, nil
}

func (s *ArtifactService) GetArtifact(ctx context.Context, id string) (*Artifact, error) {
	// Retrieve from MinIO
	obj, err := s.minioClient.GetObject(ctx, s.bucket, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	// Deserialize
	msg, err := capnp.NewDecoder(obj).Decode()
	if err != nil {
		return nil, err
	}

	artifact, err := ReadRootArtifact(msg)
	if err != nil {
		return nil, err
	}

	return &artifact, nil
}

func (s *ArtifactService) UpdateArtifact(ctx context.Context, artifact *Artifact) error {
	// Add new version
	versions, err := artifact.Versions()
	if err != nil {
		return err
	}
	newVersion, err := versions.AddWithCaveats(1)
	if err != nil {
		return err
	}
	newVersion.SetVersionNumber(uint32(versions.Len()))
	newVersion.SetTimestamp(time.Now().Unix())

	// Add audit entry
	header, err := artifact.Header()
	if err != nil {
		return err
	}
	auditTrail, err := header.AuditTrail()
	if err != nil {
		return err
	}
	newAudit, err := auditTrail.AddWithCaveats(1)
	if err != nil {
		return err
	}
	newAudit.SetTimestamp(time.Now().Unix())
	newAudit.SetAction("updated")

	// Store updated artifact
	return s.storeArtifact(ctx, *artifact)
}

func (s *ArtifactService) storeArtifact(ctx context.Context, artifact Artifact) error {
	// Serialize
	msg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return err
	}
	err = msg.SetRoot(Artifact_TypeID, artifact.ToPtr())
	if err != nil {
		return err
	}
	buf, err := msg.Marshal()
	if err != nil {
		return err
	}

	// Store in MinIO
	header, err := artifact.Header()
	if err != nil {
		return err
	}
	id, err := header.Id()
	if err != nil {
		return err
	}
	_, err = s.minioClient.PutObject(ctx, s.bucket, id, buf, int64(len(buf)), minio.PutObjectOptions{})
	return err
}

func generateUniqueID() string {
	hash := sha256.Sum256([]byte(time.Now().String()))
	return hex.EncodeToString(hash[:])
}
