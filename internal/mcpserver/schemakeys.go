package mcpserver

// The words of JSON Schema and the names of tool arguments that the tool
// definitions repeat. A misspelt key would be a schema nobody validates, so each
// is written once.
const (
	typeString  = "string"
	typeObject  = "object"
	typeInteger = "integer"
	typeNull    = "null"

	keyType                 = "type"
	keyEnum                 = "enum"
	keyDescription          = "description"
	keyProperties           = "properties"
	keyRequired             = "required"
	keyAdditionalProperties = "additionalProperties"
	keyItems                = "items"
	keyAnyOf                = "anyOf"
	keyMinimum              = "minimum"
	keyMinLength            = "minLength"
	keyName                 = "name"

	// locationPath is where an operation's parameter sits in the URL.
	locationPath = "path"

	keyResult     = "result"
	keyCategories = "categories"
	keyTotal      = "total"
	keyTools      = "tools"
	keyState      = "state"
	keyPosition   = "position"

	argPage       = "page"
	keyCourses    = "courses"
	argPerPage    = "perPage"
	argSlug       = "slug"
	argCourseSlug = "courseSlug"
	argLessonSlug = "lessonSlug"
)
