package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/pkg/errcode"
)

type CourseService struct {
	courseRepo *repository.CourseRepository
	catRepo    *repository.CategoryRepository
	tagRepo    *repository.TagRepository
	favRepo    *repository.FavoriteRepository
}

func NewCourseService(courseRepo *repository.CourseRepository, catRepo *repository.CategoryRepository, tagRepo *repository.TagRepository, favRepo *repository.FavoriteRepository) *CourseService {
	return &CourseService{courseRepo: courseRepo, catRepo: catRepo, tagRepo: tagRepo, favRepo: favRepo}
}

func (s *CourseService) List(ctx context.Context, req *model.CourseListRequest) (*model.PaginatedList, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 16
	}
	courses, total, err := s.courseRepo.List(ctx, req)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	items := make([]model.CourseListItem, len(courses))
	for i, c := range courses {
		items[i] = toCourseListItem(&c)
	}
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize, TotalPages: totalPages}, nil
}

func (s *CourseService) Detail(ctx context.Context, slug string, userID int64) (*model.CourseDetailResponse, error) {
	course, err := s.courseRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	prev, next, _ := s.courseRepo.GetPrevNext(ctx, course.ID, course.PublishedAt)
	related, _ := s.courseRepo.Related(ctx, course.CategoryID, course.ID, 3)

	resp := &model.CourseDetailResponse{
		ID:             course.ID,
		Title:          course.Title,
		Subtitle:       course.Subtitle,
		Slug:           course.Slug,
		Cover:          course.CoverImage,
		Content:        course.Content,
		DetailTitle:    course.DetailTitle,
		DetailSubtitle: course.DetailSubtitle,
		Price:          course.Price,
		VipPrice:       course.VipPrice,
		QualityLabel:   course.QualityLabel,
		LikeCount:      course.LikeCount,
		CommentCount:   course.CommentCount,
		PublishedAt:    course.PublishedAt,
	}

	if course.Category != nil {
		resp.Category = &model.CategoryBrief{ID: course.Category.ID, Name: course.Category.Name, Slug: course.Category.Slug}
	}
	resp.Tags = make([]model.TagBrief, len(course.Tags))
	for i, t := range course.Tags {
		resp.Tags[i] = model.TagBrief{ID: t.ID, Name: t.Name, Slug: t.Slug}
	}

	if prev != nil {
		resp.PrevCourse = &model.CourseBrief{Slug: prev.Slug, Title: prev.Title}
	}
	if next != nil {
		resp.NextCourse = &model.CourseBrief{Slug: next.Slug, Title: next.Title}
	}

	relatedItems := make([]model.CourseListItem, len(related))
	for i, c := range related {
		relatedItems[i] = toCourseListItem(&c)
	}
	resp.RelatedCourses = relatedItems

	// User access
	access := &model.UserAccess{}
	if userID > 0 {
		access.IsFavorited = s.favRepo.IsFavorited(ctx, userID, course.ID)
		access.HasPurchased = s.favRepo.HasPurchased(ctx, userID, course.ID)
		access.IsVip = s.favRepo.IsVip(ctx, userID)
		access.CanDownload = access.HasPurchased || access.IsVip
	}
	resp.UserAccess = access

	// Resources: only visible if can_download
	if access.CanDownload {
		resp.Resources = make([]model.ResourceItem, len(course.Resources))
		for i, r := range course.Resources {
			resp.Resources[i] = model.ResourceItem{ID: r.ID, Name: r.Name, URL: r.URL, Password: r.Password}
		}
	} else {
		resp.Resources = []model.ResourceItem{}
	}

	return resp, nil
}

func (s *CourseService) Latest(ctx context.Context) ([]model.CourseListItem, error) {
	courses, err := s.courseRepo.Latest(ctx, 8)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	items := make([]model.CourseListItem, len(courses))
	for i, c := range courses {
		items[i] = toCourseListItem(&c)
	}
	return items, nil
}

func (s *CourseService) ToggleLike(ctx context.Context, userID, courseID int64) (*model.LikeResponse, error) {
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	liked, err := s.favRepo.Toggle(ctx, userID, courseID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	delta := 1
	if !liked {
		delta = -1
	}
	s.courseRepo.IncrementLikeCount(ctx, courseID, delta)
	return &model.LikeResponse{Liked: liked, LikeCount: course.LikeCount + delta}, nil
}

func (s *CourseService) Search(ctx context.Context, req *model.SearchRequest) (*model.PaginatedList, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 16
	}
	// Trim and validate keyword
	req.Q = strings.TrimSpace(req.Q)
	if req.Q == "" {
		return &model.PaginatedList{Items: []model.CourseListItem{}, Total: 0, Page: req.Page, PageSize: req.PageSize, TotalPages: 0}, nil
	}
	courses, total, err := s.courseRepo.Search(ctx, req)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	items := make([]model.CourseListItem, len(courses))
	for i, c := range courses {
		items[i] = toCourseListItem(&c)
	}
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize, TotalPages: totalPages}, nil
}

func (s *CourseService) Categories(ctx context.Context) ([]model.Category, error) {
	return s.catRepo.List(ctx)
}

func (s *CourseService) Tags(ctx context.Context) ([]model.Tag, error) {
	return s.tagRepo.List(ctx)
}

// Admin methods

func (s *CourseService) AdminList(ctx context.Context, page, pageSize, categoryID int, status, keyword string) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	courses, total, err := s.courseRepo.AdminList(ctx, page, pageSize, categoryID, status, keyword)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: courses, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *CourseService) AdminCreate(ctx context.Context, req *model.AdminCourseRequest) (*model.Course, error) {
	now := time.Now()
	course := &model.Course{
		Title:          req.Title,
		Subtitle:       req.Subtitle,
		Slug:           req.Slug,
		CategoryID:     req.CategoryID,
		QualityLabel:   req.QualityLabel,
		CoverImage:     req.CoverImage,
		Content:        req.Content,
		DetailTitle:    req.DetailTitle,
		DetailSubtitle: req.DetailSubtitle,
		Price:          req.Price,
		VipPrice:       req.VipPrice,
		Status:         "draft",
		PublishedAt:    &now,
	}
	if err := s.courseRepo.Create(ctx, course); err != nil {
		return nil, errcode.ErrInternal
	}
	if len(req.TagIDs) > 0 {
		s.courseRepo.ReplaceTags(ctx, course.ID, req.TagIDs)
	}
	if len(req.Resources) > 0 {
		resources := make([]model.CourseResource, len(req.Resources))
		for i, r := range req.Resources {
			resources[i] = model.CourseResource{Name: r.Name, URL: r.URL, Password: r.Password, SortOrder: r.SortOrder}
		}
		s.courseRepo.ReplaceResources(ctx, course.ID, resources)
	}
	return course, nil
}

func (s *CourseService) AdminUpdate(ctx context.Context, id int64, req *model.AdminCourseRequest) (*model.Course, error) {
	course, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	course.Title = req.Title
	course.Subtitle = req.Subtitle
	course.Slug = req.Slug
	course.CategoryID = req.CategoryID
	course.QualityLabel = req.QualityLabel
	course.CoverImage = req.CoverImage
	course.Content = req.Content
	course.DetailTitle = req.DetailTitle
	course.DetailSubtitle = req.DetailSubtitle
	course.Price = req.Price
	course.VipPrice = req.VipPrice

	if err := s.courseRepo.UpdateFields(ctx, course); err != nil {
		return nil, errcode.ErrInternal
	}
	if req.TagIDs != nil {
		s.courseRepo.ReplaceTags(ctx, id, req.TagIDs)
	}
	if req.Resources != nil {
		resources := make([]model.CourseResource, len(req.Resources))
		for i, r := range req.Resources {
			resources[i] = model.CourseResource{Name: r.Name, URL: r.URL, Password: r.Password, SortOrder: r.SortOrder}
		}
		s.courseRepo.ReplaceResources(ctx, id, resources)
	}
	// Reload with associations
	course, _ = s.courseRepo.FindByID(ctx, id)
	return course, nil
}

func (s *CourseService) AdminDelete(ctx context.Context, id int64) error {
	_, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.courseRepo.Delete(ctx, id)
}

func (s *CourseService) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.courseRepo.UpdateStatus(ctx, id, status)
}

// Admin Category methods

func (s *CourseService) AdminCategoryList(ctx context.Context) ([]model.Category, error) {
	return s.catRepo.AdminList(ctx)
}

func (s *CourseService) AdminCategoryCreate(ctx context.Context, req *model.AdminCategoryRequest) (*model.Category, error) {
	cat := &model.Category{
		Name:      req.Name,
		Slug:      req.Slug,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		IsActive:  true,
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}
	if err := s.catRepo.Create(ctx, cat); err != nil {
		return nil, errcode.ErrInternal
	}
	return cat, nil
}

func (s *CourseService) AdminCategoryUpdate(ctx context.Context, id int, req *model.AdminCategoryRequest) (*model.Category, error) {
	cat, err := s.catRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	cat.Name = req.Name
	cat.Slug = req.Slug
	cat.ParentID = req.ParentID
	cat.SortOrder = req.SortOrder
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}
	if err := s.catRepo.Update(ctx, cat); err != nil {
		return nil, errcode.ErrInternal
	}
	return cat, nil
}

func (s *CourseService) AdminCategoryDelete(ctx context.Context, id int) error {
	_, err := s.catRepo.FindByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.catRepo.Delete(ctx, id)
}

func (s *CourseService) AdminCategoryUpdateStatus(ctx context.Context, id int, isActive bool) error {
	cat, err := s.catRepo.FindByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	cat.IsActive = isActive
	return s.catRepo.Update(ctx, cat)
}

func (s *CourseService) AdminCourseBatch(ctx context.Context, req *model.AdminCourseBatchRequest) error {
	for _, id := range req.IDs {
		switch req.Action {
		case "publish":
			s.courseRepo.UpdateStatus(ctx, id, "published")
		case "offline":
			s.courseRepo.UpdateStatus(ctx, id, "draft")
		case "delete":
			s.courseRepo.Delete(ctx, id)
		}
	}
	return nil
}

// Admin Tag methods

func (s *CourseService) AdminTagList(ctx context.Context) ([]model.Tag, error) {
	return s.tagRepo.List(ctx)
}

func (s *CourseService) AdminTagCreate(ctx context.Context, req *model.AdminTagRequest) (*model.Tag, error) {
	tag := &model.Tag{Name: req.Name, Slug: req.Slug}
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, errcode.ErrInternal
	}
	return tag, nil
}

func (s *CourseService) AdminTagUpdate(ctx context.Context, id int, req *model.AdminTagRequest) (*model.Tag, error) {
	tag, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	tag.Name = req.Name
	tag.Slug = req.Slug
	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, errcode.ErrInternal
	}
	return tag, nil
}

func (s *CourseService) AdminTagDelete(ctx context.Context, id int) error {
	_, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.tagRepo.Delete(ctx, id)
}

func toCourseListItem(c *model.Course) model.CourseListItem {
	item := model.CourseListItem{
		ID:          c.ID,
		Title:       c.Title,
		Subtitle:    c.Subtitle,
		Slug:        c.Slug,
		Cover:       c.CoverImage,
		Price:       c.Price,
		VipPrice:    c.VipPrice,
		PublishedAt: c.PublishedAt,
	}
	if c.PublishedAt != nil {
		item.RelativeTime = relativeTime(*c.PublishedAt)
	}
	if c.Category != nil {
		item.Category = &model.CategoryBrief{ID: c.Category.ID, Name: c.Category.Name, Slug: c.Category.Slug}
	}
	return item
}

func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Hour:
		return "刚刚"
	case d < 24*time.Hour:
		return fmt.Sprintf("%d小时前", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%d天前", int(d.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
}
