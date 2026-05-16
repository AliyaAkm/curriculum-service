package dictionary

import (
	"context"
	"curriculum-service/internal/domain/durationcategory"
	"curriculum-service/internal/domain/level"
	"curriculum-service/internal/domain/locale"
	"curriculum-service/internal/domain/status"
	"curriculum-service/internal/domain/tag"
	"curriculum-service/internal/domain/topic"
	"curriculum-service/internal/service/cache"
	"time"
)

type StatusRepository interface {
	GetAllStatus(ctx context.Context) ([]status.Status, error)
}

type LevelRepository interface {
	GetAllLevel(ctx context.Context) ([]level.Level, error)
}

type DurationCategoryRepository interface {
	GetAllDurationCategories(ctx context.Context) ([]durationcategory.DurationCategory, error)
}

type TopicRepository interface {
	GetAllTopics(ctx context.Context) ([]topic.Topic, error)
}

type TagRepository interface {
	GetAllTags(ctx context.Context) ([]tag.Tag, error)
}

type LocaleRepository interface {
	GetAllLocales(ctx context.Context) ([]locale.Locale, error)
}

type Status struct {
	repo  StatusRepository
	cache *cache.JSONCache
	ttl   time.Duration
}

func NewStatus(repo StatusRepository, cacheClient *cache.JSONCache, ttl time.Duration) *Status {
	return &Status{repo: repo, cache: cacheClient, ttl: ttl}
}

func (r *Status) GetAllStatus(ctx context.Context) ([]status.Status, error) {
	return cache.GetOrSet(ctx, r.cache, "dictionary:status", r.ttl, r.repo.GetAllStatus)
}

type Level struct {
	repo  LevelRepository
	cache *cache.JSONCache
	ttl   time.Duration
}

func NewLevel(repo LevelRepository, cacheClient *cache.JSONCache, ttl time.Duration) *Level {
	return &Level{repo: repo, cache: cacheClient, ttl: ttl}
}

func (r *Level) GetAllLevel(ctx context.Context) ([]level.Level, error) {
	return cache.GetOrSet(ctx, r.cache, "dictionary:level", r.ttl, r.repo.GetAllLevel)
}

type DurationCategory struct {
	repo  DurationCategoryRepository
	cache *cache.JSONCache
	ttl   time.Duration
}

func NewDurationCategory(repo DurationCategoryRepository, cacheClient *cache.JSONCache, ttl time.Duration) *DurationCategory {
	return &DurationCategory{repo: repo, cache: cacheClient, ttl: ttl}
}

func (r *DurationCategory) GetAllDurationCategories(ctx context.Context) ([]durationcategory.DurationCategory, error) {
	return cache.GetOrSet(ctx, r.cache, "dictionary:duration_category", r.ttl, r.repo.GetAllDurationCategories)
}

type Topic struct {
	repo  TopicRepository
	cache *cache.JSONCache
	ttl   time.Duration
}

func NewTopic(repo TopicRepository, cacheClient *cache.JSONCache, ttl time.Duration) *Topic {
	return &Topic{repo: repo, cache: cacheClient, ttl: ttl}
}

func (r *Topic) GetAllTopics(ctx context.Context) ([]topic.Topic, error) {
	return cache.GetOrSet(ctx, r.cache, "dictionary:topic", r.ttl, r.repo.GetAllTopics)
}

type Tag struct {
	repo  TagRepository
	cache *cache.JSONCache
	ttl   time.Duration
}

func NewTag(repo TagRepository, cacheClient *cache.JSONCache, ttl time.Duration) *Tag {
	return &Tag{repo: repo, cache: cacheClient, ttl: ttl}
}

func (r *Tag) GetAllTags(ctx context.Context) ([]tag.Tag, error) {
	return cache.GetOrSet(ctx, r.cache, "dictionary:tag", r.ttl, r.repo.GetAllTags)
}

type Locale struct {
	repo  LocaleRepository
	cache *cache.JSONCache
	ttl   time.Duration
}

func NewLocale(repo LocaleRepository, cacheClient *cache.JSONCache, ttl time.Duration) *Locale {
	return &Locale{repo: repo, cache: cacheClient, ttl: ttl}
}

func (r *Locale) GetAllLocales(ctx context.Context) ([]locale.Locale, error) {
	return cache.GetOrSet(ctx, r.cache, "dictionary:locale", r.ttl, r.repo.GetAllLocales)
}
