# frozen_string_literal: true

class Web::Admin::Management::UsersController < Web::Admin::Management::ApplicationController
  def index
    query = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')

    @search = User.includes(language_members: :language).ransack(query)

    users = @search.result

    respond_to do |format|
      format.html do
        @users = users.page(params[:page])
      end

      format.csv do
        headers.delete('Content-Length')
        headers['Cache-Control'] = 'no-cache'
        headers['Content-Type'] = 'text/csv'
        headers['Content-Disposition'] = "attachment; filename=users-#{Time.zone.today}.csv"
        headers['X-Accel-Buffering'] = 'no'
        headers['Last-Modified'] = Time.current.httpdate

        response.status = 200

        self.response_body = csv_enumerator(users)
      end
    end
  end

  private

  def csv_enumerator(users)
    @csv_enumerator ||= Enumerator.new do |yielder|
      yielder << %w[id email stack finished_lessons]

      users.find_each do |user|
        language_members = user.language_members

        language_members.each do |language_member|
          language = language_member.language

          yielder << CSV.generate_line([
                                         user.id,
                                         user.email,
                                         language.name,
                                         user.finished_lessons_for_language(language).size
                                       ])
        end
      end
    end
  end
end
