# frozen_string_literal: true

class Web::LanguagesController < Web::ApplicationController
  def show
    @language = Language.find_by!(slug: params[:id])

    unless @language.current_version
      f('empty_language_current_version', type: :warning)

      redirect_to root_path
      return
    end

    if @language.progress_in_development?
      f('.typescript_description_html', type: :info, values: { language: @language.to_s.humanize, link_to_repo: ExternalLinks.source_code }, now: true)
    end

    @current_module_versions = @language.current_module_versions
                                        .includes(:module)
                                        .order(:order)
                                        .eager_load(:lesson_versions)
                                        .joins(:infos)
                                        .merge(Language::Module::Version::Info.with_locale)
                                        .merge(Language::Lesson::Version.includes(:lesson).order(:order))

    @infos_by_module = @language.current_module_infos.with_locale.index_by(&:version_id)
    @infos_by_lesson = @language.current_lesson_infos.with_locale.index_by(&:version_id)

    @finished_lessons_by_id = current_user.finished_lessons_for_language(@language).index_by(&:id)
    @language_member = @language.members.find_by(user: current_user) || Language::MemberFake.new

    @first_lesson = @language.current_lessons.ordered.first
    @next_lesson = current_user.not_finished_lessons_for_language(@language).ordered.first

    @human_language_title = [t("human_languages.#{@language}"), @language.learn_as.text].join(' ')

    if @current_module_versions.empty?
      f('warning', now: true)
      nil
    end
  end
end
