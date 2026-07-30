package testsupport

import (
	"context"
	"fmt"

	"hexletbasics/internal/api"
)

// Client keeps happy-path handler tests concise while the generated client
// exposes every documented HTTP response as an operation-specific union.
// Unexpected variants are returned as errors; auth/error contract tests use the
// generated client or HTTP transport directly when they need to inspect them.
type Client struct {
	*api.Client
}

func success[T any](response any, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	if value, ok := response.(*T); ok {
		return value, nil
	}
	return nil, fmt.Errorf("unexpected API response %T", response)
}

func noContent[T any](response any, err error) error {
	_, err = success[T](response, err)
	return err
}

func (c *Client) AdminListBanners(ctx context.Context, params api.AdminListBannersParams) (*api.BannerPage, error) {
	return success[api.BannerPage](c.Client.AdminListBanners(ctx, params))
}

func (c *Client) AdminGetBanner(ctx context.Context, params api.AdminGetBannerParams) (*api.Banner, error) {
	return success[api.Banner](c.Client.AdminGetBanner(ctx, params))
}

func (c *Client) AdminCreateBanner(ctx context.Context, request *api.BannerInput) (*api.Banner, error) {
	return success[api.Banner](c.Client.AdminCreateBanner(ctx, request))
}

func (c *Client) AdminUpdateBanner(ctx context.Context, request *api.BannerInput, params api.AdminUpdateBannerParams) (*api.Banner, error) {
	return success[api.Banner](c.Client.AdminUpdateBanner(ctx, request, params))
}

func (c *Client) AdminDeleteBanner(ctx context.Context, params api.AdminDeleteBannerParams) error {
	return noContent[api.AdminDeleteBannerNoContent](c.Client.AdminDeleteBanner(ctx, params))
}

func (c *Client) AdminListBlogPosts(ctx context.Context, params api.AdminListBlogPostsParams) (*api.BlogPostPage, error) {
	return success[api.BlogPostPage](c.Client.AdminListBlogPosts(ctx, params))
}

func (c *Client) AdminGetBlogPost(ctx context.Context, params api.AdminGetBlogPostParams) (*api.BlogPost, error) {
	return success[api.BlogPost](c.Client.AdminGetBlogPost(ctx, params))
}

func (c *Client) AdminListCourseCategories(ctx context.Context, params api.AdminListCourseCategoriesParams) (*api.CourseCategoryPage, error) {
	return success[api.CourseCategoryPage](c.Client.AdminListCourseCategories(ctx, params))
}

func (c *Client) AdminGetCourseCategory(ctx context.Context, params api.AdminGetCourseCategoryParams) (*api.CourseCategory, error) {
	return success[api.CourseCategory](c.Client.AdminGetCourseCategory(ctx, params))
}

func (c *Client) AdminCreateCourseCategory(ctx context.Context, request *api.CourseCategoryInput) (*api.CourseCategory, error) {
	return success[api.CourseCategory](c.Client.AdminCreateCourseCategory(ctx, request))
}

func (c *Client) AdminUpdateCourseCategory(ctx context.Context, request *api.CourseCategoryInput, params api.AdminUpdateCourseCategoryParams) (*api.CourseCategory, error) {
	return success[api.CourseCategory](c.Client.AdminUpdateCourseCategory(ctx, request, params))
}

func (c *Client) AdminDeleteCourseCategory(ctx context.Context, params api.AdminDeleteCourseCategoryParams) error {
	return noContent[api.AdminDeleteCourseCategoryNoContent](c.Client.AdminDeleteCourseCategory(ctx, params))
}

func (c *Client) AdminListCourses(ctx context.Context, params api.AdminListCoursesParams) (*api.CoursePage, error) {
	return success[api.CoursePage](c.Client.AdminListCourses(ctx, params))
}

func (c *Client) AdminGetCourse(ctx context.Context, params api.AdminGetCourseParams) (*api.Course, error) {
	return success[api.Course](c.Client.AdminGetCourse(ctx, params))
}

func (c *Client) AdminCreateCourse(ctx context.Context, request *api.CourseInput) (*api.Course, error) {
	return success[api.Course](c.Client.AdminCreateCourse(ctx, request))
}

func (c *Client) AdminUpdateCourse(ctx context.Context, request *api.CourseInput, params api.AdminUpdateCourseParams) (*api.Course, error) {
	return success[api.Course](c.Client.AdminUpdateCourse(ctx, request, params))
}

func (c *Client) AdminListCourseLandingPages(ctx context.Context, params api.AdminListCourseLandingPagesParams) (*api.CourseLandingPagePage, error) {
	return success[api.CourseLandingPagePage](c.Client.AdminListCourseLandingPages(ctx, params))
}

func (c *Client) AdminGetCourseLandingPage(ctx context.Context, params api.AdminGetCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return success[api.CourseLandingPage](c.Client.AdminGetCourseLandingPage(ctx, params))
}

func (c *Client) AdminCreateCourseLandingPage(ctx context.Context, request *api.CourseLandingPageInput) (*api.CourseLandingPage, error) {
	return success[api.CourseLandingPage](c.Client.AdminCreateCourseLandingPage(ctx, request))
}

func (c *Client) AdminUpdateCourseLandingPage(ctx context.Context, request *api.CourseLandingPageInput, params api.AdminUpdateCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return success[api.CourseLandingPage](c.Client.AdminUpdateCourseLandingPage(ctx, request, params))
}

func (c *Client) AdminDeleteCourseLandingPage(ctx context.Context, params api.AdminDeleteCourseLandingPageParams) error {
	return noContent[api.AdminDeleteCourseLandingPageNoContent](c.Client.AdminDeleteCourseLandingPage(ctx, params))
}

func (c *Client) AdminListCourseLessons(ctx context.Context, params api.AdminListCourseLessonsParams) (*api.CourseLessonListItemPage, error) {
	return success[api.CourseLessonListItemPage](c.Client.AdminListCourseLessons(ctx, params))
}

func (c *Client) AdminListCourseLessonMembers(ctx context.Context, params api.AdminListCourseLessonMembersParams) (*api.CourseLessonMemberPage, error) {
	return success[api.CourseLessonMemberPage](c.Client.AdminListCourseLessonMembers(ctx, params))
}

func (c *Client) AdminListCourseLessonReviews(ctx context.Context, params api.AdminListCourseLessonReviewsParams) (*api.CourseLessonReviewPage, error) {
	return success[api.CourseLessonReviewPage](c.Client.AdminListCourseLessonReviews(ctx, params))
}

func (c *Client) AdminListReviews(ctx context.Context, params api.AdminListReviewsParams) (*api.ReviewPage, error) {
	return success[api.ReviewPage](c.Client.AdminListReviews(ctx, params))
}

func (c *Client) AdminGetReview(ctx context.Context, params api.AdminGetReviewParams) (*api.Review, error) {
	return success[api.Review](c.Client.AdminGetReview(ctx, params))
}

func (c *Client) AdminCreateReview(ctx context.Context, request *api.ReviewInput) (*api.Review, error) {
	return success[api.Review](c.Client.AdminCreateReview(ctx, request))
}

func (c *Client) AdminUpdateReview(ctx context.Context, request *api.ReviewInput, params api.AdminUpdateReviewParams) (*api.Review, error) {
	return success[api.Review](c.Client.AdminUpdateReview(ctx, request, params))
}

func (c *Client) AdminDeleteReview(ctx context.Context, params api.AdminDeleteReviewParams) error {
	return noContent[api.AdminDeleteReviewNoContent](c.Client.AdminDeleteReview(ctx, params))
}

func (c *Client) AdminListRoles(ctx context.Context, params api.AdminListRolesParams) (*api.StaffRolePage, error) {
	return success[api.StaffRolePage](c.Client.AdminListRoles(ctx, params))
}

func (c *Client) AdminGetRole(ctx context.Context, params api.AdminGetRoleParams) (*api.StaffRoleDetail, error) {
	return success[api.StaffRoleDetail](c.Client.AdminGetRole(ctx, params))
}

func (c *Client) AdminCreateRole(ctx context.Context, request *api.RoleInput) (*api.StaffRoleDetail, error) {
	return success[api.StaffRoleDetail](c.Client.AdminCreateRole(ctx, request))
}

func (c *Client) AdminUpdateRole(ctx context.Context, request *api.RoleInput, params api.AdminUpdateRoleParams) (*api.StaffRoleDetail, error) {
	return success[api.StaffRoleDetail](c.Client.AdminUpdateRole(ctx, request, params))
}

func (c *Client) AdminDeleteRole(ctx context.Context, params api.AdminDeleteRoleParams) error {
	return noContent[api.AdminDeleteRoleNoContent](c.Client.AdminDeleteRole(ctx, params))
}

func (c *Client) AdminGetRolePermissions(ctx context.Context, params api.AdminGetRolePermissionsParams) (*api.StaffRoleDetail, error) {
	return success[api.StaffRoleDetail](c.Client.AdminGetRolePermissions(ctx, params))
}

func (c *Client) AdminUpdateRolePermissions(ctx context.Context, request *api.RolePermissionsInput, params api.AdminUpdateRolePermissionsParams) (*api.StaffRoleDetail, error) {
	return success[api.StaffRoleDetail](c.Client.AdminUpdateRolePermissions(ctx, request, params))
}

func (c *Client) AdminListStaffMembers(ctx context.Context, params api.AdminListStaffMembersParams) (*api.StaffMemberPage, error) {
	return success[api.StaffMemberPage](c.Client.AdminListStaffMembers(ctx, params))
}

func (c *Client) AdminGetStaffMember(ctx context.Context, params api.AdminGetStaffMemberParams) (*api.StaffMember, error) {
	return success[api.StaffMember](c.Client.AdminGetStaffMember(ctx, params))
}

func (c *Client) AdminCreateStaffMember(ctx context.Context, request *api.StaffMemberInput) (*api.StaffMember, error) {
	return success[api.StaffMember](c.Client.AdminCreateStaffMember(ctx, request))
}

func (c *Client) AdminUpdateStaffMember(ctx context.Context, request *api.StaffMemberInput, params api.AdminUpdateStaffMemberParams) (*api.StaffMember, error) {
	return success[api.StaffMember](c.Client.AdminUpdateStaffMember(ctx, request, params))
}

func (c *Client) AdminDeleteStaffMember(ctx context.Context, params api.AdminDeleteStaffMemberParams) error {
	return noContent[api.AdminDeleteStaffMemberNoContent](c.Client.AdminDeleteStaffMember(ctx, params))
}

func (c *Client) AdminListManagementUsers(ctx context.Context, params api.AdminListManagementUsersParams) (*api.UserCrudPage, error) {
	return success[api.UserCrudPage](c.Client.AdminListManagementUsers(ctx, params))
}

func (c *Client) AdminGetManagementUser(ctx context.Context, params api.AdminGetManagementUserParams) (*api.UserCrud, error) {
	return success[api.UserCrud](c.Client.AdminGetManagementUser(ctx, params))
}

func (c *Client) AdminUpdateManagementUser(ctx context.Context, request *api.UserInput, params api.AdminUpdateManagementUserParams) (*api.UserCrud, error) {
	return success[api.UserCrud](c.Client.AdminUpdateManagementUser(ctx, request, params))
}

func (c *Client) AdminListUsers(ctx context.Context, params api.AdminListUsersParams) (*api.UserCrudPage, error) {
	return success[api.UserCrudPage](c.Client.AdminListUsers(ctx, params))
}

func (c *Client) AdminSearchUsers(ctx context.Context, params api.AdminSearchUsersParams) ([]api.UserCrud, error) {
	response, err := c.Client.AdminSearchUsers(ctx, params)
	items, err := success[api.AdminSearchUsersOKApplicationJSON](response, err)
	if err != nil {
		return nil, err
	}
	return []api.UserCrud(*items), nil
}

func (c *Client) AdminGetUser(ctx context.Context, params api.AdminGetUserParams) (*api.UserCrud, error) {
	return success[api.UserCrud](c.Client.AdminGetUser(ctx, params))
}

func (c *Client) AdminCreateUser(ctx context.Context, request *api.UserInput) (*api.UserCrud, error) {
	return success[api.UserCrud](c.Client.AdminCreateUser(ctx, request))
}

func (c *Client) AdminUpdateUser(ctx context.Context, request *api.UserInput, params api.AdminUpdateUserParams) (*api.UserCrud, error) {
	return success[api.UserCrud](c.Client.AdminUpdateUser(ctx, request, params))
}

func (c *Client) AdminDeleteUser(ctx context.Context, params api.AdminDeleteUserParams) error {
	return noContent[api.AdminDeleteUserNoContent](c.Client.AdminDeleteUser(ctx, params))
}

func (c *Client) AdminListCategoryQnaItems(ctx context.Context, params api.AdminListCategoryQnaItemsParams) ([]api.QnaItem, error) {
	response, err := c.Client.AdminListCategoryQnaItems(ctx, params)
	items, err := success[api.AdminListCategoryQnaItemsOKApplicationJSON](response, err)
	if err != nil {
		return nil, err
	}
	return []api.QnaItem(*items), nil
}

func (c *Client) AdminCreateCategoryQnaItem(ctx context.Context, request *api.QnaItemInput, params api.AdminCreateCategoryQnaItemParams) (*api.QnaItem, error) {
	return success[api.QnaItem](c.Client.AdminCreateCategoryQnaItem(ctx, request, params))
}

func (c *Client) AdminUpdateCategoryQnaItem(ctx context.Context, request *api.QnaItemInput, params api.AdminUpdateCategoryQnaItemParams) (*api.QnaItem, error) {
	return success[api.QnaItem](c.Client.AdminUpdateCategoryQnaItem(ctx, request, params))
}

func (c *Client) AdminDeleteCategoryQnaItem(ctx context.Context, params api.AdminDeleteCategoryQnaItemParams) error {
	return noContent[api.AdminDeleteCategoryQnaItemNoContent](c.Client.AdminDeleteCategoryQnaItem(ctx, params))
}

func (c *Client) AdminListLandingPageQnaItems(ctx context.Context, params api.AdminListLandingPageQnaItemsParams) ([]api.QnaItem, error) {
	response, err := c.Client.AdminListLandingPageQnaItems(ctx, params)
	items, err := success[api.AdminListLandingPageQnaItemsOKApplicationJSON](response, err)
	if err != nil {
		return nil, err
	}
	return []api.QnaItem(*items), nil
}

func (c *Client) AdminCreateLandingPageQnaItem(ctx context.Context, request *api.QnaItemInput, params api.AdminCreateLandingPageQnaItemParams) (*api.QnaItem, error) {
	return success[api.QnaItem](c.Client.AdminCreateLandingPageQnaItem(ctx, request, params))
}

func (c *Client) AdminUpdateLandingPageQnaItem(ctx context.Context, request *api.QnaItemInput, params api.AdminUpdateLandingPageQnaItemParams) (*api.QnaItem, error) {
	return success[api.QnaItem](c.Client.AdminUpdateLandingPageQnaItem(ctx, request, params))
}

func (c *Client) AdminDeleteLandingPageQnaItem(ctx context.Context, params api.AdminDeleteLandingPageQnaItemParams) error {
	return noContent[api.AdminDeleteLandingPageQnaItemNoContent](c.Client.AdminDeleteLandingPageQnaItem(ctx, params))
}

func (c *Client) AdminListLeads(ctx context.Context, params api.AdminListLeadsParams) (*api.LeadPage, error) {
	return success[api.LeadPage](c.Client.AdminListLeads(ctx, params))
}
