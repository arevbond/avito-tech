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

	returningID = returning + fieldID
)

var (
	queryAllBannersWithFilters = `
		SELECT 
			b.id, 
			b.feature_id, 
			b.content, 
			b.is_active, 
			b.created_at, 
			b.updated_at,
			ARRAY_AGG(CASE WHEN bt.tag_id != 0 THEN bt.tag_id END) FILTER (WHERE bt.banner_id = b.id) AS tag_ids
		FROM 
			banners AS b
		LEFT JOIN 
			banner_tags AS bt ON b.id = bt.banner_id
		WHERE 
			(b.feature_id = $1 OR $1 = 0)
			AND
			($2 = 0 OR EXISTS (SELECT 1 FROM banner_tags WHERE banner_id = b.id AND tag_id = $2))
		GROUP BY
			b.id
		ORDER BY 
			b.id
		OFFSET $3 LIMIT $4;`
)
