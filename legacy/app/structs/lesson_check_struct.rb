# typed: strict
# frozen_string_literal: true

class LessonCheckStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :version_id, T.nilable(Integer)
  const :code, T.nilable(String)

  validates :version_id, presence: true
end
