package notice

import (
	"context"
	"database/sql"
	"errors"

	"aiub-companion/internal/database"
	"aiub-companion/internal/database/sqlc"
)

type Repository interface {
	WithTx(ctx context.Context, fn func(Repository) error) error
	GetNotices(ctx context.Context, f Filter) ([]Notice, error)
	GetNoticeByID(ctx context.Context, id string) (*Notice, error)
	GetNoticeAttachments(ctx context.Context, noticeID string) ([]Attachment, error)
	GetLatestNoticeSourceOrder(ctx context.Context) (int64, error)
	UpsertNotice(ctx context.Context, n Notice, sourceOrder int64) (bool, error)
	UpdateNotice(ctx context.Context, n Notice) error
	UpdateNoticeDetails(ctx context.Context, id string, fullTitle, content string) error
	UpsertNoticeAttachment(ctx context.Context, noticeID string, att Attachment) error
	SetPinState(ctx context.Context, id string, pinned bool) error
	SetReadState(ctx context.Context, id string, read bool) error
	ClearNotices(ctx context.Context) error
}

type repository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *repository) WithTx(ctx context.Context, fn func(Repository) error) error {
	return database.RunInTx(ctx, r.db, func(tx *sql.Tx) error {
		return fn(&repository{queries: r.queries.WithTx(tx), db: r.db})
	})
}

func (r *repository) GetNotices(ctx context.Context, f Filter) ([]Notice, error) {
	rows, err := r.queries.ListNoticesWithState(ctx, sqlc.ListNoticesWithStateParams{
		Category: database.StringOrNull(f.Category),
		Urgent:   database.BoolOrNull(f.Urgent),
		Pinned:   database.BoolOrNull(f.Pinned),
		Unread:   database.BoolOrNull(f.Unread),
	})
	if err != nil {
		return nil, err
	}

	notices := make([]Notice, len(rows))
	for i := range rows {
		notices[i] = Notice{
			ID:         rows[i].ID,
			Title:      rows[i].Title,
			FullTitle:  rows[i].FullTitle.String,
			Summary:    rows[i].Summary.String,
			Content:    rows[i].Content.String,
			PostedDate: rows[i].PostedDate,
			Category:   rows[i].Category,
			IsCached:   rows[i].IsCached == 1,
			IsUrgent:   rows[i].IsUrgent == 1,
			IsPinned:   rows[i].IsPinned == 1,
			IsRead:     rows[i].IsRead == 1,
		}
	}
	return notices, nil
}

func (r *repository) GetNoticeByID(ctx context.Context, id string) (*Notice, error) {
	row, err := r.queries.GetNoticeWithState(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Notice{
		ID:         row.ID,
		Title:      row.Title,
		FullTitle:  row.FullTitle.String,
		Summary:    row.Summary.String,
		Content:    row.Content.String,
		PostedDate: row.PostedDate,
		Category:   row.Category,
		IsCached:   row.IsCached == 1,
		IsPinned:   row.IsPinned == 1,
		IsRead:     row.IsRead == 1,
	}, nil
}

func (r *repository) GetNoticeAttachments(ctx context.Context, noticeID string) ([]Attachment, error) {
	rows, err := r.queries.GetNoticeAttachments(ctx, noticeID)
	if err != nil {
		return nil, err
	}

	attachments := make([]Attachment, len(rows))
	for i := range rows {
		attachments[i] = Attachment{
			ID:        rows[i].ID,
			Label:     rows[i].Label,
			URL:       rows[i].Url,
			LocalPath: rows[i].LocalPath.String,
		}
	}
	return attachments, nil
}

func (r *repository) GetLatestNoticeSourceOrder(ctx context.Context) (int64, error) {
	info, err := r.queries.GetLatestNoticeInfo(ctx)
	if err != nil {
		return 0, err
	}
	return info.SourceOrder, nil
}

func (r *repository) UpsertNotice(ctx context.Context, n Notice, sourceOrder int64) (bool, error) {
	isNew, err := r.queries.UpsertNotice(ctx, sqlc.UpsertNoticeParams{
		ID:          n.ID,
		Title:       n.Title,
		Summary:     database.StringOrNull(n.Summary),
		PostedDate:  n.PostedDate,
		Category:    n.Category,
		IsUrgent:    database.BoolToInt64(n.IsUrgent),
		SourceOrder: sourceOrder,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isNew, nil
}

func (r *repository) UpdateNotice(ctx context.Context, n Notice) error {
	return r.queries.UpdateNotice(ctx, sqlc.UpdateNoticeParams{
		ID:         n.ID,
		Title:      n.Title,
		Summary:    database.StringOrNull(n.Summary),
		PostedDate: n.PostedDate,
	})
}

func (r *repository) UpdateNoticeDetails(ctx context.Context, id string, fullTitle, content string) error {
	return r.queries.UpdateNoticeDetails(ctx, sqlc.UpdateNoticeDetailsParams{
		ID:        id,
		FullTitle: database.StringOrNull(fullTitle),
		Content:   database.StringOrNull(content),
	})
}

func (r *repository) UpsertNoticeAttachment(ctx context.Context, noticeID string, att Attachment) error {
	return r.queries.UpsertNoticeAttachment(ctx, sqlc.UpsertNoticeAttachmentParams{
		ID:       att.ID,
		NoticeID: noticeID,
		Url:      att.URL,
		Label:    att.Label,
	})
}

func (r *repository) SetPinState(ctx context.Context, id string, pinned bool) error {
	return r.queries.SetPinState(ctx, sqlc.SetPinStateParams{
		NoticeID: id,
		IsPinned: database.BoolToInt64(pinned),
	})
}

func (r *repository) SetReadState(ctx context.Context, id string, read bool) error {
	return r.queries.SetReadState(ctx, sqlc.SetReadStateParams{
		NoticeID: id,
		IsRead:   database.BoolToInt64(read),
	})
}

func (r *repository) ClearNotices(ctx context.Context) error {
	return r.queries.ClearNotices(ctx)
}
