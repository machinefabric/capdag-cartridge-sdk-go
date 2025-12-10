// Package metadata provides consolidated metadata structures for document processing
package metadata

import (
	"time"
)

// FileMetadata represents consolidated metadata for any file/document type
type FileMetadata struct {
	// Basic file info
	FilePath      string `json:"file_path"`
	FileSizeBytes uint64 `json:"file_size_bytes"`
	ContentLength uint64 `json:"content_length"`
	MimeType      *string `json:"mime_type,omitempty"`
	Encoding      *string `json:"encoding,omitempty"`

	// Common document metadata
	Title      *string `json:"title,omitempty"`
	Subject    *string `json:"subject,omitempty"`
	Identifier *string `json:"identifier,omitempty"`
	Language   *string `json:"language,omitempty"`
	Creator    *string `json:"creator,omitempty"`
	Producer   *string `json:"producer,omitempty"`

	// People and keywords
	Authors      []string `json:"authors"`
	Contributors []string `json:"contributors"`
	Keywords     []string `json:"keywords"`

	// Dates (ISO format recommended)
	CreationDate     *time.Time `json:"creation_date,omitempty"`
	ModificationDate *time.Time `json:"modification_date,omitempty"`

	// Counts
	PageCount      *uint64 `json:"page_count,omitempty"`
	ChapterCount   *uint64 `json:"chapter_count,omitempty"`
	WordCount      *uint64 `json:"word_count,omitempty"`
	CharacterCount *uint64 `json:"character_count,omitempty"`

	// Document classification
	DocumentType  string  `json:"document_type"`
	FormatVersion *string `json:"format_version,omitempty"`

	// Extended/format-specific metadata (arbitrary JSON values)
	ExtendedMetadata map[string]interface{} `json:"extended_metadata"`

	// PDF-specific
	PDFVersion      *string `json:"pdf_version,omitempty"`
	HasForms        bool    `json:"has_forms"`
	IsEncrypted     bool    `json:"is_encrypted"`
	AttachmentCount uint64  `json:"attachment_count"`
	IsLinearized    bool    `json:"is_linearized"`

	// EPUB-specific
	EPUBVersion     *string    `json:"epub_version,omitempty"`
	Publisher       *string    `json:"publisher,omitempty"`
	PublicationDate *time.Time `json:"publication_date,omitempty"`
	Rights          *string    `json:"rights,omitempty"`
	HasDRM          bool       `json:"has_drm"`
	ThumbnailPath   *string    `json:"thumbnail_path,omitempty"`
}

// NewFileMetadata creates a new FileMetadata with required fields
func NewFileMetadata(filePath, documentType string, fileSizeBytes uint64) *FileMetadata {
	return &FileMetadata{
		FilePath:         filePath,
		DocumentType:     documentType,
		FileSizeBytes:    fileSizeBytes,
		ContentLength:    0,
		Authors:          make([]string, 0),
		Contributors:     make([]string, 0),
		Keywords:         make([]string, 0),
		ExtendedMetadata: make(map[string]interface{}),
		HasForms:         false,
		IsEncrypted:      false,
		AttachmentCount:  0,
		IsLinearized:     false,
		HasDRM:           false,
	}
}

// NewMinimalFileMetadata creates minimal metadata with just basic file info
func NewMinimalFileMetadata(filePath string, fileSizeBytes uint64, mimeType *string, documentType string) *FileMetadata {
	metadata := NewFileMetadata(filePath, documentType, fileSizeBytes)
	metadata.MimeType = mimeType
	return metadata
}

// AddAuthor adds an author to the authors list
func (m *FileMetadata) AddAuthor(author string) {
	m.Authors = append(m.Authors, author)
}

// AddContributor adds a contributor to the contributors list
func (m *FileMetadata) AddContributor(contributor string) {
	m.Contributors = append(m.Contributors, contributor)
}

// AddKeyword adds a keyword to the keywords list
func (m *FileMetadata) AddKeyword(keyword string) {
	m.Keywords = append(m.Keywords, keyword)
}

// SetExtended sets an extended metadata value
func (m *FileMetadata) SetExtended(key string, value interface{}) {
	m.ExtendedMetadata[key] = value
}

// GetExtended gets an extended metadata value
func (m *FileMetadata) GetExtended(key string) (interface{}, bool) {
	value, exists := m.ExtendedMetadata[key]
	return value, exists
}

// PrimaryAuthor returns the first author if any
func (m *FileMetadata) PrimaryAuthor() *string {
	if len(m.Authors) > 0 {
		return &m.Authors[0]
	}
	return nil
}

// IsEmpty checks if metadata is essentially empty (no meaningful fields)
func (m *FileMetadata) IsEmpty() bool {
	return m.Title == nil &&
		len(m.Authors) == 0 &&
		m.Subject == nil &&
		m.Creator == nil &&
		m.Producer == nil
}

// ToPrompt creates a short prompt-style summary of key metadata fields
func (m *FileMetadata) ToPrompt() string {
	var parts []string

	// Filename only (extract from path)
	// Simple implementation - in production might want proper path parsing
	if m.FilePath != "" {
		parts = append(parts, "File: "+m.FilePath)
	}

	if m.Title != nil && *m.Title != "" {
		title := *m.Title
		if len(title) > 100 {
			title = title[:97] + "..."
		}
		parts = append(parts, "Title: "+title)
	}

	if len(m.Authors) > 0 {
		authorsStr := ""
		for i, author := range m.Authors {
			if i > 0 {
				authorsStr += ", "
			}
			authorsStr += author
		}
		if len(authorsStr) > 150 {
			authorsStr = authorsStr[:147] + "..."
		}
		parts = append(parts, "Authors: "+authorsStr)
	}

	if len(m.Contributors) > 0 {
		contribStr := ""
		for i, contrib := range m.Contributors {
			if i > 0 {
				contribStr += ", "
			}
			contribStr += contrib
		}
		if len(contribStr) > 150 {
			contribStr = contribStr[:147] + "..."
		}
		parts = append(parts, "Contributors: "+contribStr)
	}

	if len(m.Keywords) > 0 {
		keywordsLimited := m.Keywords
		if len(keywordsLimited) > 5 {
			keywordsLimited = keywordsLimited[:5]
		}
		keywordsStr := ""
		for i, keyword := range keywordsLimited {
			if i > 0 {
				keywordsStr += ", "
			}
			keywordsStr += keyword
		}
		parts = append(parts, "Keywords: "+keywordsStr)
	}

	if m.CreationDate != nil {
		// Display just YYYY-MM-DD if possible
		dateStr := m.CreationDate.Format("2006-01-02")
		parts = append(parts, "Created: "+dateStr)
	}

	if len(parts) == 0 {
		return ""
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += "\n"
		}
		result += part
	}

	if len(result) > 300 {
		return result[:297] + "..."
	}
	return result
}

