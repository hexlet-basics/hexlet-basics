package localization

import (
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:generate go tool goi18n extract -sourceLanguage en -format json -outdir locales .

// Backend message registry. English fallback text stays next to each stable ID;
// the embedded catalogs contain the complete en/ru/es translations.
var (
	WrongCredentials = Message{
		value: i18n.Message{
			ID:    "auth.wrong_credentials",
			Other: "Wrong email or password",
		},
	}
	PasswordProcessingFailed = Message{
		value: i18n.Message{
			ID:    "auth.password_processing_failed",
			Other: "Could not process the password",
		},
	}
	EmailTaken = Message{
		value: i18n.Message{
			ID:    "auth.email_taken",
			Other: "This email is already taken",
		},
	}
	CourseNotFound = Message{
		value: i18n.Message{
			ID:    "courses.not_found",
			Other: "course not found",
		},
	}
	LessonNotFound = Message{
		value: i18n.Message{
			ID:    "lessons.not_found",
			Other: "lesson not found",
		},
	}
	FileTooLarge = Message{
		value: i18n.Message{
			ID:    "attachments.file_too_large",
			Other: "file is too large",
		},
	}
	FileRequired = Message{
		value: i18n.Message{
			ID:    "attachments.file_required",
			Other: "a file part is required",
		},
	}
	UnsupportedFileType = Message{
		value: i18n.Message{
			ID:    "attachments.unsupported_file_type",
			Other: "unsupported file type",
		},
	}
	StoreFileFailed = Message{
		value: i18n.Message{
			ID:    "attachments.store_failed",
			Other: "failed to store file",
		},
	}
	ReadFileFailed = Message{
		value: i18n.Message{
			ID:    "attachments.read_failed",
			Other: "failed to read file",
		},
	}
	ReadWebhookBodyFailed = Message{
		value: i18n.Message{
			ID:    "webhooks.read_body_failed",
			Other: "failed to read body",
		},
	}
	InvalidWebhookSignature = Message{
		value: i18n.Message{
			ID:    "webhooks.invalid_signature",
			Other: "invalid signature",
		},
	}
	InvalidWebhookPayload = Message{
		value: i18n.Message{
			ID:    "webhooks.invalid_payload",
			Other: "invalid payload",
		},
	}
	TriggerBuildFailed = Message{
		value: i18n.Message{
			ID:    "webhooks.trigger_build_failed",
			Other: "failed to trigger build",
		},
	}
	BadRequest = Message{
		value: i18n.Message{
			ID:    "http.bad_request",
			Other: "Bad Request",
		},
	}
	Unauthorized = Message{
		value: i18n.Message{
			ID:    "http.unauthorized",
			Other: "Unauthorized",
		},
	}
	Forbidden = Message{
		value: i18n.Message{
			ID:    "http.forbidden",
			Other: "Forbidden",
		},
	}
	NotFound = Message{
		value: i18n.Message{
			ID:    "http.not_found",
			Other: "Not Found",
		},
	}
	MethodNotAllowed = Message{
		value: i18n.Message{
			ID:    "http.method_not_allowed",
			Other: "Method Not Allowed",
		},
	}
	Conflict = Message{
		value: i18n.Message{
			ID:    "http.conflict",
			Other: "Conflict",
		},
	}
	RequestEntityTooLarge = Message{
		value: i18n.Message{
			ID:    "http.request_entity_too_large",
			Other: "Request Entity Too Large",
		},
	}
	UnsupportedMediaType = Message{
		value: i18n.Message{
			ID:    "http.unsupported_media_type",
			Other: "Unsupported Media Type",
		},
	}
	UnprocessableEntity = Message{
		value: i18n.Message{
			ID:    "http.unprocessable_entity",
			Other: "Unprocessable Entity",
		},
	}
	InternalServerError = Message{
		value: i18n.Message{
			ID:    "http.internal_server_error",
			Other: "Internal Server Error",
		},
	}
)

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
