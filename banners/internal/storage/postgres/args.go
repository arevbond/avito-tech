package postgres

const (
	featureTable      = "features"
	bannerTable       = "banners"
	tagTable          = "tags"
	bannerToTagsTable = "banner_tags"
)

const (
	returning       = "RETURNING "
	separator       = ","
	suffixDoNothing = "ON CONFLICT DO NOTHING"
)

const (
	fieldID   = "id"
	fieldName = "name"

	fieldContent  = "content"
	fieldIsActive = "is_active"

	fieldCreatedAt = "created_at"
	fieldUpdatedAt = "updated_at"

	fieldBannerID  = "banner_id"
	fieldFeatureID = "feature_id"
	fieldTagID     = "tag_id"
)

var (
	featureFields = []string{
		fieldID, fieldName,
	}
	tagFields = []string{
		fieldID, fieldName,
	}
	bannerToTagsFields = []string{
		fieldBannerID, fieldTagID,
	}
	bannerFields = []string{
		fieldID, fieldFeatureID, fieldContent, fieldIsActive, fieldCreatedAt, fieldUpdatedAt,
	}
)
