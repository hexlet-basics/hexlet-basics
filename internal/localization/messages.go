package localization

import "net/http"

// Backend message registry. English fallback text stays next to each stable ID;
// the embedded catalogs contain the complete en/ru/es translations.
var (
	WrongCredentials = Message{
		id:      "auth.wrong_credentials",
		english: "Wrong email or password",
	}
	PasswordProcessingFailed = Message{
		id:      "auth.password_processing_failed",
		english: "Could not process the password",
	}
	EmailTaken = Message{
		id:      "auth.email_taken",
		english: "This email is already taken",
	}
	CourseNotFound = Message{
		id:      "courses.not_found",
		english: "course not found",
	}
	FileTooLarge = Message{
		id:      "attachments.file_too_large",
		english: "file is too large",
	}
	FileRequired = Message{
		id:      "attachments.file_required",
		english: "a file part is required",
	}
	UnsupportedFileType = Message{
		id:      "attachments.unsupported_file_type",
		english: "unsupported file type",
	}
	StoreFileFailed = Message{
		id:      "attachments.store_failed",
		english: "failed to store file",
	}
	ReadFileFailed = Message{
		id:      "attachments.read_failed",
		english: "failed to read file",
	}
	ReadWebhookBodyFailed = Message{
		id:      "webhooks.read_body_failed",
		english: "failed to read body",
	}
	InvalidWebhookSignature = Message{
		id:      "webhooks.invalid_signature",
		english: "invalid signature",
	}
	InvalidWebhookPayload = Message{
		id:      "webhooks.invalid_payload",
		english: "invalid payload",
	}
	TriggerBuildFailed = Message{
		id:      "webhooks.trigger_build_failed",
		english: "failed to trigger build",
	}
	BadRequest = Message{
		id:      "http.bad_request",
		english: http.StatusText(http.StatusBadRequest),
	}
	Unauthorized = Message{
		id:      "http.unauthorized",
		english: http.StatusText(http.StatusUnauthorized),
	}
	Forbidden = Message{
		id:      "http.forbidden",
		english: http.StatusText(http.StatusForbidden),
	}
	NotFound = Message{
		id:      "http.not_found",
		english: http.StatusText(http.StatusNotFound),
	}
	MethodNotAllowed = Message{
		id:      "http.method_not_allowed",
		english: http.StatusText(http.StatusMethodNotAllowed),
	}
	Conflict = Message{
		id:      "http.conflict",
		english: http.StatusText(http.StatusConflict),
	}
	RequestEntityTooLarge = Message{
		id:      "http.request_entity_too_large",
		english: http.StatusText(http.StatusRequestEntityTooLarge),
	}
	UnsupportedMediaType = Message{
		id:      "http.unsupported_media_type",
		english: http.StatusText(http.StatusUnsupportedMediaType),
	}
	UnprocessableEntity = Message{
		id:      "http.unprocessable_entity",
		english: http.StatusText(http.StatusUnprocessableEntity),
	}
	InternalServerError = Message{
		id:      "http.internal_server_error",
		english: http.StatusText(http.StatusInternalServerError),
	}
)

var allMessages = []Message{
	WrongCredentials,
	PasswordProcessingFailed,
	EmailTaken,
	CourseNotFound,
	FileTooLarge,
	FileRequired,
	UnsupportedFileType,
	StoreFileFailed,
	ReadFileFailed,
	ReadWebhookBodyFailed,
	InvalidWebhookSignature,
	InvalidWebhookPayload,
	TriggerBuildFailed,
	BadRequest,
	Unauthorized,
	Forbidden,
	NotFound,
	MethodNotAllowed,
	Conflict,
	RequestEntityTooLarge,
	UnsupportedMediaType,
	UnprocessableEntity,
	InternalServerError,
}

var statusMessages = map[int]Message{
	http.StatusBadRequest:            BadRequest,
	http.StatusUnauthorized:          Unauthorized,
	http.StatusForbidden:             Forbidden,
	http.StatusNotFound:              NotFound,
	http.StatusMethodNotAllowed:      MethodNotAllowed,
	http.StatusConflict:              Conflict,
	http.StatusRequestEntityTooLarge: RequestEntityTooLarge,
	http.StatusUnsupportedMediaType:  UnsupportedMediaType,
	http.StatusUnprocessableEntity:   UnprocessableEntity,
	http.StatusInternalServerError:   InternalServerError,
}
